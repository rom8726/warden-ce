//nolint:errcheck,gocyclo,gosec,nestif // need refactoring
package installer

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	_ "embed"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"os"
	"os/user"
	"strings"
	"text/template"
	"time"

	"github.com/rom8726/warden/internal/version"
)

// ErrContextCancelled is returned when the context is cancelled.
var ErrContextCancelled = errors.New("installation cancelled by user")

//go:embed docker-compose.yml
var dockerCompose string

//go:embed platform.env.tmpl
var platformEnv string

//go:embed config.env.tmpl
var configEnv string

//go:embed Makefile
var makefile string

var DockerRegistry string

// Config holds all the configuration data collected from user input.
type Config struct {
	// Domain for the platform
	Domain string

	// SSL certificate settings
	HasExistingSSLCert bool

	// SMTP server details
	MailerAddr     string
	MailerUser     string
	MailerPassword string
	MailerFrom     string
	MailerUseTLS   bool

	// Generated passwords
	PGPassword       string
	CHPassword       string
	SecretKey        string
	JWTSecretKey     string
	AdminTmpPassword string

	// Admin settings
	AdminEmail string

	// Other settings
	FrontendURL     string
	SecureLinkMD5   string
	PlatformVersion string
}

type App struct {
	config Config
	reader *bufio.Reader
}

func New() *App {
	return &App{
		reader: bufio.NewReader(os.Stdin),
	}
}

// Run executes the installer.
func (a *App) Run(ctx context.Context) error {
	// Check if running as root
	if err := a.checkRoot(); err != nil {
		return err
	}

	// Welcome message
	a.printWelcome()

	// Inform about installation directory
	fmt.Println("The platform will be installed in the /opt/warden directory.")
	fmt.Println()

	// Ask for confirmation to proceed
	confirmed, err := a.confirmContinueWithContext(ctx)
	if err != nil {
		if errors.Is(err, ErrContextCancelled) {
			fmt.Println("Installation cancelled by user.")

			return err
		}

		return err
	}
	if !confirmed {
		fmt.Println("Installation cancelled.")

		return nil
	}

	// Collect user input
	if err := a.collectUserInputWithContext(ctx); err != nil {
		if errors.Is(err, ErrContextCancelled) {
			fmt.Println("Installation cancelled by user.")
		}

		return err
	}

	// Generate passwords and other required values
	a.generateSecrets()

	// Check context before continuing
	select {
	case <-ctx.Done():
		fmt.Println("Installation cancelled by user.")

		return ErrContextCancelled
	default:
	}

	// Create required directories
	if err := a.createDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Check context before continuing
	select {
	case <-ctx.Done():
		fmt.Println("Installation cancelled by user.")

		return ErrContextCancelled
	default:
	}

	// Check context before continuing with templates
	select {
	case <-ctx.Done():
		fmt.Println("Installation cancelled by user.")

		return ErrContextCancelled
	default:
	}

	// Render templates
	if err := a.renderTemplates(); err != nil {
		return fmt.Errorf("failed to render templates: %w", err)
	}

	// Handle SSL certificate based on user's choice
	if !a.config.HasExistingSSLCert {
		// Generate self-signed SSL certificate
		if err := a.generateSSLCertificate(); err != nil {
			return fmt.Errorf("failed to generate SSL certificate: %w", err)
		}
	}

	fmt.Println("\nInstallation completed successfully!")
	fmt.Println("The platform has been installed in /opt/warden")
	fmt.Println("You can check and modify settings in /opt/warden/platform.env and /opt/warden/config.env")
	fmt.Println("A Makefile has been created in /opt/warden with commands for starting and stopping the platform")

	// Display admin login information
	fmt.Println("\nADMIN LOGIN INFORMATION:")
	fmt.Println("  - Email:", a.config.AdminEmail)
	fmt.Println("  - Temporary Password:", a.config.AdminTmpPassword)
	fmt.Println("Please use these credentials to log in to the platform. You will be prompted to change the password on first login.") //nolint:lll // it's ok here

	// Remind user about SSL certificate if they have an existing one
	if a.config.HasExistingSSLCert {
		fmt.Println("\nREMINDER: Don't forget to place your SSL certificate and key files at:")
		fmt.Println("  - Certificate: /opt/warden/nginx/ssl/nginx_cert.pem")
		fmt.Println("  - Key: /opt/warden/nginx/ssl/nginx_key.pem")
	}

	// Remind user about mailer TLS certificates if they chose to use TLS for email
	if a.config.MailerUseTLS {
		fmt.Println("\nREMINDER: Don't forget to place your email TLS certificate and key files at:")
		fmt.Println("  - Certificate: /opt/warden/secrets/mailer_cert.pem")
		fmt.Println("  - Key: /opt/warden/secrets/mailer_key.pem")
	}

	return nil
}

// checkRoot verifies that the installer is running as root.
func (a *App) checkRoot() error {
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	if currentUser.Uid != "0" {
		return errors.New("installer must be run as root (sudo)")
	}

	return nil
}

// printWelcome displays the welcome message.
func (a *App) printWelcome() {
	fmt.Println("=================================================")
	fmt.Println("       Welcome to Warden Platform Installer      ")
	fmt.Println("=================================================")
	fmt.Println("This installer will set up the Warden platform on your system.")
	fmt.Println("It will create necessary directories and configuration files.")
	fmt.Println()
}

// confirmContinueWithContext asks the user for confirmation to proceed with context awareness.
func (a *App) confirmContinueWithContext(ctx context.Context) (bool, error) {
	return a.readYesNoWithContext(ctx, "Do you want to continue with the installation?")
}

// readInputWithContext reads a line of input from the user with context awareness.
func (a *App) readInputWithContext(ctx context.Context, prompt string) (string, error) {
	// Check if context is cancelled
	select {
	case <-ctx.Done():
		return "", ErrContextCancelled
	default:
		// Continue with input
	}

	fmt.Print(prompt)

	// Create a channel for the input result
	resultCh := make(chan struct {
		input string
		err   error
	})

	// Read input in a goroutine
	go func() {
		input, err := a.reader.ReadString('\n')
		resultCh <- struct {
			input string
			err   error
		}{strings.TrimSpace(input), err}
	}()

	// Wait for either input or context cancellation
	select {
	case <-ctx.Done():
		fmt.Println()

		return "", ErrContextCancelled
	case result := <-resultCh:
		return result.input, result.err
	}
}

// readYesNoWithContext reads a yes/no answer from the user with context awareness.
func (a *App) readYesNoWithContext(ctx context.Context, prompt string) (bool, error) {
	input, err := a.readInputWithContext(ctx, prompt+" (y/n): ")
	if err != nil {
		return false, err
	}
	input = strings.ToLower(input)

	return input == "y" || input == "yes", nil
}

func (a *App) collectUserInputWithContext(ctx context.Context) error {
	fmt.Println("\n=== Platform Configuration ===")

	// Get admin email
	for {
		adminEmail, err := a.readInputWithContext(ctx, "Enter administrator email: ")
		if err != nil {
			return fmt.Errorf("failed to read admin email: %w", err)
		}
		if adminEmail == "" {
			fmt.Println("Error: administrator email cannot be empty. Please try again.")

			continue
		}
		// Simple validation to check if it looks like an email
		if !strings.Contains(adminEmail, "@") {
			fmt.Println("Error: please enter a valid email address. Please try again.")

			continue
		}
		a.config.AdminEmail = adminEmail

		break
	}

	// Get domain
	for {
		domain, err := a.readInputWithContext(ctx, "Enter domain for the platform: ")
		if err != nil {
			return fmt.Errorf("failed to read domain: %w", err)
		}
		if domain == "" {
			fmt.Println("Error: domain cannot be empty. Please try again.")

			continue
		}
		a.config.Domain = domain
		a.config.FrontendURL = "https://" + domain

		break
	}

	// Ask about SSL certificate
	hasExistingSSLCert, err := a.readYesNoWithContext(ctx, "Do you have an existing SSL certificate for this domain?")
	if err != nil {
		return fmt.Errorf("failed to read SSL certificate option: %w", err)
	}
	a.config.HasExistingSSLCert = hasExistingSSLCert

	if hasExistingSSLCert {
		fmt.Println("\nYou will need to place your SSL certificate and key files at:")
		fmt.Println("  - Certificate: /opt/warden/nginx/ssl/nginx_cert.pem")
		fmt.Println("  - Key: /opt/warden/nginx/ssl/nginx_key.pem")
		fmt.Println("You will be reminded about this at the end of installation.")
	} else {
		fmt.Println("\nA self-signed SSL certificate will be generated for you at the end of installation.")
	}

	fmt.Println("\n=== SMTP Server Configuration ===")

	// Get SMTP server details
	for {
		mailerAddr, err := a.readInputWithContext(ctx, "Enter SMTP server address (including port): ")
		if err != nil {
			return fmt.Errorf("failed to read SMTP address: %w", err)
		}
		if mailerAddr == "" {
			fmt.Println("Error: SMTP address cannot be empty. Please try again.")

			continue
		}
		a.config.MailerAddr = mailerAddr

		break
	}

	// SMTP user
	mailerUser, err := a.readInputWithContext(ctx, "Enter SMTP user: ")
	if err != nil {
		return fmt.Errorf("failed to read SMTP user: %w", err)
	}
	a.config.MailerUser = mailerUser

	// SMTP password
	mailerPassword, err := a.readInputWithContext(ctx, "Enter SMTP password: ")
	if err != nil {
		return fmt.Errorf("failed to read SMTP password: %w", err)
	}
	a.config.MailerPassword = mailerPassword

	// Email from
	for {
		mailerFrom, err := a.readInputWithContext(ctx, "Enter email address for sending emails (from): ")
		if err != nil {
			return fmt.Errorf("failed to read email from: %w", err)
		}
		if mailerFrom == "" {
			fmt.Println("Error: email from cannot be empty. Please try again.")

			continue
		}
		a.config.MailerFrom = mailerFrom

		break
	}

	// TLS option
	mailerUseTLS, err := a.readYesNoWithContext(ctx, "Use TLS for SMTP connection?")
	if err != nil {
		return fmt.Errorf("failed to read TLS option: %w", err)
	}
	a.config.MailerUseTLS = mailerUseTLS

	return nil
}

// generateSecrets creates random passwords and keys.
func (a *App) generateSecrets() {
	// Set platform version from version package
	a.config.PlatformVersion = version.Version

	// Set empty string for SecureLinkMD5
	a.config.SecureLinkMD5 = ""

	// Generate PostgreSQL password (12 characters, alphanumeric)
	a.config.PGPassword = generateRandomString(12)

	// Generate ClickHouse password (12 characters, alphanumeric)
	a.config.CHPassword = generateRandomString(12)

	// Generate SecretKey (32 characters, alphanumeric)
	a.config.SecretKey = generateRandomString(32)

	// Generate JWTSecretKey (32 characters, alphanumeric)
	a.config.JWTSecretKey = generateRandomString(32)

	// Generate admin temporary password (12 characters, alphanumeric)
	a.config.AdminTmpPassword = generateRandomString(12)

	fmt.Println("\nGenerated secure passwords and keys.")
}

// generateRandomString creates a random alphanumeric string of the specified length.
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[mathrand.Intn(len(charset))]
	}

	return string(result)
}

// createDirectories creates the required directories for the installation.
func (a *App) createDirectories() error {
	// Define the directories to create
	directories := []string{
		"/opt/warden",
		"/opt/warden/nginx/ssl",
		"/opt/warden/secrets",
	}

	// Create each directory
	for _, dir := range directories {
		fmt.Printf("Creating directory: %s\n", dir)
		if err := os.MkdirAll(dir, 0o600); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// renderTemplates generates configuration files from templates.
func (a *App) renderTemplates() error {
	// Define template data for platform.env
	platformData := struct {
		DockerRegistry  string
		Domain          string
		PGPassword      string
		CHPassword      string
		SecureLinkMD5   string
		PlatformVersion string
	}{
		DockerRegistry:  DockerRegistry,
		Domain:          a.config.Domain,
		PGPassword:      a.config.PGPassword,
		CHPassword:      a.config.CHPassword,
		SecureLinkMD5:   a.config.SecureLinkMD5,
		PlatformVersion: a.config.PlatformVersion,
	}

	// Define template data for config.env
	configData := struct {
		FrontendURL      string
		SecretKey        string
		JWTSecretKey     string
		PGPassword       string
		CHPassword       string
		MailerAddr       string
		MailerUser       string
		MailerPassword   string
		MailerFrom       string
		MailerUserTLS    bool
		AdminEmail       string
		AdminTmpPassword string
	}{
		FrontendURL:      a.config.FrontendURL,
		SecretKey:        a.config.SecretKey,
		JWTSecretKey:     a.config.JWTSecretKey,
		PGPassword:       a.config.PGPassword,
		CHPassword:       a.config.CHPassword,
		MailerAddr:       a.config.MailerAddr,
		MailerUser:       a.config.MailerUser,
		MailerPassword:   a.config.MailerPassword,
		MailerFrom:       a.config.MailerFrom,
		MailerUserTLS:    a.config.MailerUseTLS,
		AdminEmail:       a.config.AdminEmail,
		AdminTmpPassword: a.config.AdminTmpPassword,
	}

	// Render platform.env
	platformTmpl, err := template.New("platform.env").Parse(platformEnv)
	if err != nil {
		return fmt.Errorf("failed to parse platform.env template: %w", err)
	}

	platformFile, err := os.Create("/opt/warden/platform.env")
	if err != nil {
		return fmt.Errorf("failed to create platform.env file: %w", err)
	}
	defer platformFile.Close()

	if err := platformTmpl.Execute(platformFile, platformData); err != nil {
		return fmt.Errorf("failed to render platform.env template: %w", err)
	}

	fmt.Println("Created /opt/warden/platform.env")

	// Render config.env
	configTmpl, err := template.New("config.env").Parse(configEnv)
	if err != nil {
		return fmt.Errorf("failed to parse config.env template: %w", err)
	}

	configFile, err := os.Create("/opt/warden/config.env")
	if err != nil {
		return fmt.Errorf("failed to create config.env file: %w", err)
	}
	defer configFile.Close()

	if err := configTmpl.Execute(configFile, configData); err != nil {
		return fmt.Errorf("failed to render config.env template: %w", err)
	}

	fmt.Println("Created /opt/warden/config.env")

	// Write docker-compose.yml
	dockerComposeFile, err := os.Create("/opt/warden/docker-compose.yml")
	if err != nil {
		return fmt.Errorf("failed to create docker-compose.yml file: %w", err)
	}
	defer dockerComposeFile.Close()

	if _, err := dockerComposeFile.WriteString(dockerCompose); err != nil {
		return fmt.Errorf("failed to write docker-compose.yml file: %w", err)
	}

	fmt.Println("Created /opt/warden/docker-compose.yml")

	// Write Makefile
	makefileFile, err := os.Create("/opt/warden/Makefile")
	if err != nil {
		return fmt.Errorf("failed to create Makefile file: %w", err)
	}
	defer makefileFile.Close()

	if _, err := makefileFile.WriteString(makefile); err != nil {
		return fmt.Errorf("failed to write Makefile file: %w", err)
	}

	fmt.Println("Created /opt/warden/Makefile")

	return nil
}

// generateSSLCertificate generates a self-signed SSL certificate for the domain.
func (a *App) generateSSLCertificate() error {
	const (
		keyPath  = "/opt/warden/nginx/ssl/nginx_key.pem"
		certPath = "/opt/warden/nginx/ssl/nginx_cert.pem"
	)

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("SSL key generation failed: %w", err)
	}

	serialNumber, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		return fmt.Errorf("SSL serial number generation failed: %w", err)
	}

	// Use the domain from config for the certificate
	certTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   a.config.Domain,
			Organization: []string{"Warden"},
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().AddDate(10, 0, 0), // Valid for 10 years
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &key.PublicKey, key)
	if err != nil {
		return fmt.Errorf("SSL certificate creation failed: %w", err)
	}

	keyFile, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("open SSL key file: %w", err)
	}
	err = pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	if err != nil {
		_ = keyFile.Close()

		return fmt.Errorf("write SSL key: %w", err)
	}
	_ = keyFile.Close()

	certFile, err := os.OpenFile(certPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("open SSL cert file: %w", err)
	}
	if err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		_ = certFile.Close()

		return fmt.Errorf("write SSL cert: %w", err)
	}
	_ = certFile.Close()

	fmt.Println("Generated self-signed SSL certificate for", a.config.Domain)

	return nil
}
