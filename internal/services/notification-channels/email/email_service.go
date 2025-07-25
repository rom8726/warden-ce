//nolint:gosec,gocyclo // it's ok
package email

import (
	"bytes"
	"context"
	"crypto/tls"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-mail/mail"

	"github.com/rom8726/warden/internal/domain"
)

const (
	timeout = 30 * time.Second
)

//go:embed templates/issues_summary_email.tmpl
var issuesSummaryEmailTemplate string

//go:embed templates/issue_email.tmpl
var issueEmailTemplate string

//go:embed templates/reset_password_email.tmpl
var resetPasswordEmailTemplate string

//go:embed templates/2fa_code_email.tmpl
var twoFACodeEmailTemplate string

//go:embed templates/user_notification_email.tmpl
var userNotificationEmailTemplate string

type Service struct {
	cfg          *Config
	teamsRepo    TeamsRepository
	usersRepo    UsersRepository
	projectsRepo ProjectsRepository

	sendEmailFunc func(ctx context.Context, toEmails []string, subject, body string) error
}

type Config struct {
	SMTPHost      string
	Username      string
	Password      string
	CertFile      string
	KeyFile       string
	AllowInsecure bool
	UseTLS        bool

	BaseURL string
	From    string
	LogoURL string
}

func New(
	cfg *Config,
	teamsRepo TeamsRepository,
	usersRepo UsersRepository,
	projectsRepo ProjectsRepository,
) *Service {
	service := &Service{
		cfg:          cfg,
		teamsRepo:    teamsRepo,
		usersRepo:    usersRepo,
		projectsRepo: projectsRepo,
	}
	service.sendEmailFunc = service.SendEmail

	return service
}

func (s *Service) Type() domain.NotificationType {
	return domain.NotificationTypeEmail
}

func (s *Service) Send(
	ctx context.Context,
	issue *domain.Issue,
	project *domain.Project,
	configData json.RawMessage,
	isRegress bool,
) error {
	emails, err := s.getMembersEmails(ctx, project)
	if err != nil {
		return fmt.Errorf("get members emails: %w", err)
	}

	var emailCfg EmailConfig
	if err := json.Unmarshal(configData, &emailCfg); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	emails = append(emails, emailCfg.EmailTo)

	newOrRegress := "new"
	if isRegress {
		newOrRegress = "regress"
	}

	subject := fmt.Sprintf("[%s][%s] Issue #%d from project %q: %s",
		issue.Level, newOrRegress, issue.ID, project.Name, issue.Title)

	body, err := renderSingleIssueEmailBody(issue, project, s.cfg.BaseURL, isRegress)
	if err != nil {
		return fmt.Errorf("render body: %w", err)
	}

	return s.SendEmail(ctx, emails, subject, body)
}

func (s *Service) SendResetPasswordEmail(ctx context.Context, email, token string) error {
	slog.Debug("sending reset password email", "base_url", s.cfg.BaseURL)

	tpl, err := template.New("reset_password").Parse(resetPasswordEmailTemplate)
	if err != nil {
		slog.Error("failed to parse reset password template", "error", err)

		return fmt.Errorf("parse template: %w", err)
	}

	resetLink := s.cfg.BaseURL + "/reset-password?token=" + token
	slog.Debug("generated reset link", "reset_link", resetLink)

	renderData := struct {
		ResetLink string
	}{
		ResetLink: resetLink,
	}

	var body bytes.Buffer
	if err := tpl.Execute(&body, renderData); err != nil {
		slog.Error("failed to execute reset password template", "error", err)

		return fmt.Errorf("execute template: %w", err)
	}

	err = s.SendEmail(ctx, []string{email}, "Warden: Reset Your Password", body.String())
	if err != nil {
		slog.Error("failed to send reset password email", "error", err)

		return err
	}

	slog.Info("reset password email sent successfully")

	return nil
}

func (s *Service) SendUnresolvedIssuesSummaryEmail(ctx context.Context, issues []domain.IssueExtended) error {
	projectsMap := make(map[domain.ProjectID]domain.Project)
	teamProjectMap := make(map[domain.TeamID]domain.ProjectID)
	issuesNewByProjectMap := make(map[domain.ProjectID][]domain.IssueExtended)
	issuesRegressByProjectMap := make(map[domain.ProjectID][]domain.IssueExtended)

	for i := range issues {
		issue := issues[i]

		_, exists := projectsMap[issue.ProjectID]
		if !exists {
			project, err := s.projectsRepo.GetByID(ctx, issue.ProjectID)
			if err != nil {
				return fmt.Errorf("get project: %w", err)
			}

			projectsMap[issue.ProjectID] = project
			issuesNewByProjectMap[issue.ProjectID] = make([]domain.IssueExtended, 0)
			issuesRegressByProjectMap[issue.ProjectID] = make([]domain.IssueExtended, 0)

			teamID := domain.TeamIDCommon
			if project.TeamID != nil {
				teamID = *project.TeamID
			}
			teamProjectMap[teamID] = project.ID
		}

		if issue.ResolvedAt == nil {
			issuesNewByProjectMap[issue.ProjectID] = append(issuesNewByProjectMap[issue.ProjectID], issue)
		} else {
			issuesRegressByProjectMap[issue.ProjectID] = append(issuesRegressByProjectMap[issue.ProjectID], issue)
		}
	}

	teamIDs := make([]domain.TeamID, 0, len(teamProjectMap))
	for teamID := range teamProjectMap {
		teamIDs = append(teamIDs, teamID)
	}

	userIDs, err := s.teamsRepo.GetUniqueUserIDsByTeamIDs(ctx, teamIDs)
	if err != nil {
		return fmt.Errorf("get teams: %w", err)
	}

	usersList, err := s.usersRepo.FetchByIDs(ctx, userIDs)
	if err != nil {
		return fmt.Errorf("fetch users: %w", err)
	}

	teamsByUserIDMap, err := s.teamsRepo.GetTeamsByUserIDs(ctx, userIDs)
	if err != nil {
		return fmt.Errorf("get teams by user ids map: %w", err)
	}

	userProjectIDsMap := make(map[domain.UserID][]domain.ProjectID)
	for userID, teams := range teamsByUserIDMap {
		userProjectIDs := make([]domain.ProjectID, 0, len(teamIDs))
		for _, team := range teams {
			projectID, exists := teamProjectMap[team.ID]
			if exists {
				userProjectIDs = append(userProjectIDs, projectID)
			}
		}
		userProjectIDsMap[userID] = userProjectIDs
	}

	const maxWorkers = 10

	// Prepare data for sending emails
	emailsToSend := make([]emailData, 0, len(usersList))

	for _, user := range usersList {
		type projectRenderData struct {
			ProjectName   string
			ProjectID     domain.ProjectID
			NewIssues     []domain.IssueExtended
			RegressIssues []domain.IssueExtended
		}
		userProjects := make([]projectRenderData, 0)

		for _, projectID := range userProjectIDsMap[user.ID] {
			newIssues := issuesNewByProjectMap[projectID]
			regressIssues := issuesRegressByProjectMap[projectID]

			if len(newIssues) > 0 || len(regressIssues) > 0 {
				project := projectsMap[projectID]
				userProjects = append(userProjects, projectRenderData{
					ProjectName:   project.Name,
					ProjectID:     project.ID,
					NewIssues:     newIssues,
					RegressIssues: regressIssues,
				})
			}
		}

		if len(userProjects) == 0 {
			continue
		}

		type renderDataType struct {
			BaseURL  string
			Name     string
			Projects []projectRenderData
		}

		renderData := renderDataType{
			BaseURL:  s.cfg.BaseURL,
			Name:     user.Username,
			Projects: userProjects,
		}

		tpl, err := template.New("unresolved_issues_summary").Parse(issuesSummaryEmailTemplate)
		if err != nil {
			return fmt.Errorf("parse template for user %s: %w", user.Username, err)
		}

		var body bytes.Buffer
		if err := tpl.Execute(&body, renderData); err != nil {
			slog.Error("execute template failed", "error", err, "user", user.Username)

			return fmt.Errorf("execute template for user %s: %w", user.Username, err)
		}

		emailsToSend = append(emailsToSend, emailData{
			toEmails: []string{user.Email},
			subject:  "Warden: unresolved issues summary",
			body:     body.String(),
		})
	}

	// Send emails in parallel
	return s.sendEmailsParallel(ctx, maxWorkers, emailsToSend)
}

// SendEmail builds a MIME message and sends it via SMTP.
//
//nolint:gosec,nestif // it's ok here
func (s *Service) SendEmail(ctx context.Context, toEmails []string, subject, bodyHTML string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	slog.Debug("starting email send", "to", strings.Join(toEmails, ", "), "subject", subject, "smtp_host", s.cfg.SMTPHost)

	from := s.cfg.From
	if from == "" {
		from = s.cfg.Username
	}

	// --- Build message ------------------------------------------------------
	msg := mail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", toEmails...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html; charset=UTF-8", bodyHTML)

	// --- Dialer -------------------------------------------------------------
	host, portStr, err := net.SplitHostPort(s.cfg.SMTPHost)
	if err != nil {
		slog.Error("invalid smtp host configuration", "smtp_host", s.cfg.SMTPHost, "error", err)

		return fmt.Errorf("invalid smtp host: %w", err)
	}
	port, _ := strconv.Atoi(portStr)

	slog.Debug("creating SMTP dialer", "host", host, "port", port, "username", s.cfg.Username)
	dialer := mail.NewDialer(host, port, s.cfg.Username, s.cfg.Password)
	dialer.Timeout = timeout

	if s.cfg.UseTLS {
		var certs []tls.Certificate
		if s.cfg.CertFile != "" {
			cert, err := tls.LoadX509KeyPair(s.cfg.CertFile, s.cfg.KeyFile)
			if err != nil {
				slog.Error("failed to load TLS certificate", "cert_file", s.cfg.CertFile,
					"key_file", s.cfg.KeyFile, "error", err)

				return fmt.Errorf("load TLS key pair: %w", err)
			}
			certs = []tls.Certificate{cert}
		}

		dialer.TLSConfig = &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: s.cfg.AllowInsecure,
			Certificates:       certs,
		}
		slog.Debug("using TLS for SMTP connection", "host", host, "port", port)
	} else {
		// For MailHog, explicitly set TLS config to nil to avoid any TLS attempts
		if port == 1025 {
			dialer.TLSConfig = nil
			slog.Warn("explicitly set TLS config to nil for MailHog")
		}
		slog.Debug("using unencrypted SMTP connection", "host", host, "port", port)
	}

	// --- Send with context --------------------------------------------------
	slog.Debug("attempting to send email", "host", host, "port", port, "use_tls", s.cfg.UseTLS)

	// Special handling for MailHog to avoid TLS issues
	if port == 1025 {
		slog.Debug("using direct SMTP for MailHog")
		err = s.sendEmailDirectSMTP(ctx, host, port, s.cfg.Username, from, toEmails, subject, bodyHTML)
	} else {
		errCh := make(chan error, 1)
		go func() {
			errCh <- dialer.DialAndSend(msg)
		}()

		select {
		case <-ctx.Done():
			slog.Error("email send timeout", "to", strings.Join(toEmails, ", "), "subject", subject,
				"error", ctx.Err())

			return fmt.Errorf("send mail: %w", ctx.Err())
		case err = <-errCh:
		}
	}

	if err != nil {
		slog.Error("failed to send email", "to", strings.Join(toEmails, ", "), "subject", subject,
			"smtp_host", s.cfg.SMTPHost, "host", host, "port", port, "use_tls", s.cfg.UseTLS, "error", err)

		return fmt.Errorf("send mail: %w", err)
	}

	slog.Debug("email sent successfully", "to", strings.Join(toEmails, ", "), "subject", subject)

	return nil
}

func (s *Service) getMembersEmails(ctx context.Context, project *domain.Project) ([]string, error) {
	if project.TeamID == nil {
		slog.Warn("project has no team", "project_id", project.ID)

		return nil, nil
	}

	members, err := s.teamsRepo.GetMembers(ctx, *project.TeamID)
	if err != nil {
		return nil, fmt.Errorf("get members: %w", err)
	}

	userIDs := make([]domain.UserID, 0, len(members))
	for _, member := range members {
		userIDs = append(userIDs, member.UserID)
	}

	users, err := s.usersRepo.FetchByIDs(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("fetch users: %w", err)
	}

	emails := make([]string, 0, len(users))
	for _, user := range users {
		emails = append(emails, user.Email)
	}

	return emails, nil
}

func renderSingleIssueEmailBody(
	issue *domain.Issue,
	project *domain.Project,
	baseURL string,
	isRegress bool,
) (string, error) {
	tpl, err := template.New("issue_email").Parse(issueEmailTemplate)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	renderData := struct {
		ID          uint
		ProjectName string
		Title       string
		Level       string
		Status      string
		FirstSeen   string
		LastSeen    string
		Platform    string
		TotalEvents uint
		BaseURL     string
		ProjectID   uint
		IsRegress   bool
	}{
		ID:          uint(issue.ID),
		ProjectName: project.Name,
		Title:       issue.Title,
		Level:       string(issue.Level),
		Status:      string(issue.Status),
		FirstSeen:   issue.FirstSeen.Format(time.RFC3339),
		LastSeen:    issue.LastSeen.Format(time.RFC3339),
		Platform:    issue.Platform,
		TotalEvents: issue.TotalEvents,
		BaseURL:     baseURL,
		ProjectID:   uint(project.ID),
		IsRegress:   isRegress,
	}

	var body bytes.Buffer
	if err := tpl.Execute(&body, renderData); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return body.String(), nil
}

// sendEmailsParallel sends emails in parallel with a limit on the number of workers.
func (s *Service) sendEmailsParallel(
	ctx context.Context,
	maxWorkers int,
	emails []emailData,
) error {
	if len(emails) == 0 {
		return nil
	}

	// Create a channel to limit the number of concurrent sends
	semaphore := make(chan struct{}, maxWorkers)

	// Create a WaitGroup to wait for all sends to complete
	var wg sync.WaitGroup

	// Channel to collect errors
	errorChan := make(chan error, len(emails))

	// Function to send a single email
	sendSingleEmail := func(email emailData) {
		defer wg.Done()

		// Get a slot for sending
		semaphore <- struct{}{}
		defer func() { <-semaphore }()

		err := s.sendEmailFunc(ctx, email.toEmails, email.subject, email.body)
		if err != nil {
			errorChan <- fmt.Errorf("send email to %s: %w", strings.Join(email.toEmails, ", "), err)

			return
		}
	}

	// Start sending emails in goroutines
	for _, email := range emails {
		wg.Add(1)
		go sendSingleEmail(email)
	}

	// Wait for all sends to complete
	wg.Wait()
	close(errorChan)

	// Check for errors
	for err := range errorChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// emailData contains data for sending a single email.
type emailData struct {
	toEmails []string
	subject  string
	body     string
}

// Send2FACodeEmail sends a 2FA confirmation code for a specific action (disable/reset).
func (s *Service) Send2FACodeEmail(ctx context.Context, email, code, action string) error {
	subject := "Warden: 2FA confirmation code"
	var actionText string
	switch action {
	case "disable":
		actionText = "to disable two-factor authentication"
	case "reset":
		actionText = "to reset two-factor authentication"
	default:
		actionText = "for your action"
	}

	tmpl, err := template.New("2fa_code_email").Parse(twoFACodeEmailTemplate)
	if err != nil {
		return err
	}
	var body bytes.Buffer
	err = tmpl.Execute(&body, struct {
		Code       string
		ActionText string
	}{
		Code:       code,
		ActionText: actionText,
	})
	if err != nil {
		return err
	}

	return s.SendEmail(ctx, []string{email}, subject, body.String())
}

// sendEmailDirectSMTP sends an email directly using net/smtp.
func (s *Service) sendEmailDirectSMTP(
	_ context.Context,
	host string,
	port int,
	username, from string,
	toEmails []string,
	subject, bodyHTML string,
) error {
	slog.Debug("starting direct SMTP send", "host", host, "port", port, "username", username)

	// Create email message
	msg := []byte("To: " + strings.Join(toEmails, ",") + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + bodyHTML + "\r\n")

	// Connect to SMTP server
	conn, err := smtp.Dial(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		slog.Error("failed to connect to SMTP server", "host", host, "port", port, "error", err)

		return err
	}
	defer func() { _ = conn.Close() }()

	// Say hello
	if err := conn.Hello("localhost"); err != nil {
		slog.Error("failed to say hello to SMTP server", "host", host, "port", port, "error", err)

		return err
	}

	// Set sender
	if err := conn.Mail(from); err != nil {
		slog.Error("failed to set from address", "host", host, "port", port, "error", err)

		return err
	}

	// Set recipients
	for _, to := range toEmails {
		if err := conn.Rcpt(to); err != nil {
			slog.Error("failed to set to address", "host", host, "port", port, "to", to, "error", err)

			return err
		}
	}

	// Send data
	writeCloser, err := conn.Data()
	if err != nil {
		slog.Error("failed to open data connection", "host", host, "port", port, "error", err)

		return err
	}
	defer func() { _ = writeCloser.Close() }()

	if _, err := writeCloser.Write(msg); err != nil {
		slog.Error("failed to write email data", "host", host, "port", port, "error", err)

		return err
	}

	slog.Debug("email sent successfully via direct SMTP", "host", host, "port", port)

	return nil
}

func (s *Service) SendUserNotificationEmail(
	ctx context.Context,
	toEmail string,
	notifType domain.UserNotificationType,
	content domain.UserNotificationContent,
) error {
	tpl, err := template.New("user_notification").Parse(userNotificationEmailTemplate)
	if err != nil {
		slog.Error("failed to parse user notification template", "error", err)

		return fmt.Errorf("parse template: %w", err)
	}

	subject := "Warden: New notification"
	switch notifType {
	case domain.UserNotificationTypeTeamAdded:
		subject = "Warden: You've been added to a team"
	case domain.UserNotificationTypeTeamRemoved:
		subject = "Warden: You've been removed from a team"
	case domain.UserNotificationTypeRoleChanged:
		subject = "Warden: Your team role has been changed"
	case domain.UserNotificationTypeIssueRegression:
		subject = "Warden: Issue regression detected in project"
	}

	var body bytes.Buffer
	if err := tpl.ExecuteTemplate(&body, "UserNotificationEmail", content); err != nil {
		slog.Error("failed to execute user notification template", "error", err)

		return fmt.Errorf("execute template: %w", err)
	}

	return s.SendEmail(ctx, []string{toEmail}, subject, body.String())
}
