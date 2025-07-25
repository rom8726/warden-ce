package users

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/pkg/crypt"
)

const issuerName = "Warden"

var twoFACodeStore = struct {
	sync.Mutex
	codes map[domain.UserID]twoFACodeEntry
}{codes: make(map[domain.UserID]twoFACodeEntry)}

type twoFACodeEntry struct {
	Code      string
	Action    string
	ExpiresAt time.Time
}

// In-memory store for 2FA session IDs.
type twoFASessionEntry struct {
	UserID    domain.UserID
	Username  string
	CreatedAt time.Time
}

var twoFASessionStore = struct {
	sync.Mutex
	sessions map[string]twoFASessionEntry
}{sessions: make(map[string]twoFASessionEntry)}

func generate2FACode() string {
	if env, ok := os.LookupEnv("WARDEN_ENVIRONMENT"); ok && env == "test" { // TODO: refactor
		return "654321"
	}

	otp := ""
	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			otp += "0"
		} else {
			otp += strconv.FormatInt(num.Int64(), 10)
		}
	}

	return otp
}

func store2FACode(userID domain.UserID, code, action string, ttl time.Duration) {
	twoFACodeStore.Lock()
	defer twoFACodeStore.Unlock()
	twoFACodeStore.codes[userID] = twoFACodeEntry{
		Code:      code,
		Action:    action,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func validate2FACode(userID domain.UserID, code, action string) bool {
	twoFACodeStore.Lock()
	defer twoFACodeStore.Unlock()
	entry, ok := twoFACodeStore.codes[userID]
	if !ok || entry.Action != action || entry.Code != code || time.Now().After(entry.ExpiresAt) {
		return false
	}
	delete(twoFACodeStore.codes, userID)

	return true
}

func (s *UsersService) Setup2FA(ctx context.Context, userID domain.UserID) (secret, qrURL, qrImage string, err error) {
	user, err := s.usersRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", "", fmt.Errorf("get user: %w", err)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuerName,
		AccountName: user.Email,
	})
	if err != nil {
		return "", "", "", fmt.Errorf("generate totp secret: %w", err)
	}

	secretStr := key.Secret()
	qrURL = key.URL()

	encKey := []byte(s.tokenizer.SecretKey())
	encSecret, err := crypt.EncryptAESGCM([]byte(secretStr), encKey)
	if err != nil {
		return "", "", "", fmt.Errorf("encrypt secret: %w", err)
	}

	encSecretB64 := base64.StdEncoding.EncodeToString(encSecret)

	if err := s.usersRepo.Update2FA(ctx, userID, false, encSecretB64, nil); err != nil {
		return "", "", "", fmt.Errorf("save user: %w", err)
	}

	qrPNG, err := qrcode.Encode(qrURL, qrcode.Medium, 256)
	if err != nil {
		return "", "", "", fmt.Errorf("generate qr: %w", err)
	}
	qrImage = base64.StdEncoding.EncodeToString(qrPNG)

	return secretStr, qrURL, qrImage, nil
}

// Confirm2FA enables 2FA for the user after validating the provided TOTP code.
func (s *UsersService) Confirm2FA(ctx context.Context, userID domain.UserID, code string) error {
	if s.twoFARateLimiter.IsBlocked(userID) {
		return domain.ErrTooMany2FAAttempts
	}

	user, err := s.usersRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	if user.TwoFASecret == "" {
		return errors.New("2FA secret not set")
	}

	encKey := []byte(s.tokenizer.SecretKey())
	encSecret, err := base64.StdEncoding.DecodeString(user.TwoFASecret)
	if err != nil {
		return fmt.Errorf("decode secret: %w", err)
	}
	plainSecret, err := crypt.DecryptAESGCM(encSecret, encKey)
	if err != nil {
		return fmt.Errorf("decrypt secret: %w", err)
	}

	valid := totp.Validate(code, string(plainSecret))
	if !valid {
		_, blocked := s.twoFARateLimiter.Inc(userID)
		if blocked {
			return domain.ErrTooMany2FAAttempts
		}

		return domain.ErrInvalid2FACode
	}

	s.twoFARateLimiter.Reset(userID)

	now := time.Now().UTC()
	if err := s.usersRepo.Update2FA(ctx, userID, true, user.TwoFASecret, &now); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

// Send2FACode Call this to initiate 2FA disable/reset: generates and sends code.
func (s *UsersService) Send2FACode(ctx context.Context, userID domain.UserID, action string) error {
	user, err := s.usersRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	code := generate2FACode()
	store2FACode(userID, code, action, 15*time.Minute)

	return s.emailer.Send2FACodeEmail(ctx, user.Email, code, action)
}

// Disable2FA disables 2FA for the user after validating the email code.
func (s *UsersService) Disable2FA(ctx context.Context, userID domain.UserID, emailCode string) error {
	if !validate2FACode(userID, emailCode, "disable") {
		return domain.ErrInvalidEmailCode
	}

	if err := s.usersRepo.Update2FA(ctx, userID, false, "", nil); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

func (s *UsersService) Reset2FA(
	ctx context.Context,
	userID domain.UserID,
	emailCode string,
) (secret, qrURL, qrImage string, err error) {
	user, err := s.usersRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", "", fmt.Errorf("get user: %w", err)
	}
	if !validate2FACode(userID, emailCode, "reset") {
		return "", "", "", errors.New("invalid or expired email code")
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuerName,
		AccountName: user.Email,
	})
	if err != nil {
		return "", "", "", fmt.Errorf("generate totp secret: %w", err)
	}
	secretStr := key.Secret()
	qrURL = key.URL()
	encKey := []byte(s.tokenizer.SecretKey())
	encSecret, err := crypt.EncryptAESGCM([]byte(secretStr), encKey)
	if err != nil {
		return "", "", "", fmt.Errorf("encrypt secret: %w", err)
	}
	encSecretB64 := base64.StdEncoding.EncodeToString(encSecret)
	if err := s.usersRepo.Update2FA(ctx, userID, false, encSecretB64, nil); err != nil {
		return "", "", "", fmt.Errorf("save user: %w", err)
	}
	qrPNG, err := qrcode.Encode(qrURL, qrcode.Medium, 256)
	if err != nil {
		return "", "", "", fmt.Errorf("generate qr: %w", err)
	}
	qrImage = base64.StdEncoding.EncodeToString(qrPNG)

	return secretStr, qrURL, qrImage, nil
}

func (s *UsersService) Verify2FA(
	ctx context.Context,
	code, sessionID string,
) (accessToken, refreshToken string, expiresIn int, err error) {
	session, ok := get2FASession(sessionID)
	if !ok {
		return "", "", 0, domain.ErrInvalidToken
	}

	userID := session.UserID
	if s.twoFARateLimiter.IsBlocked(userID) {
		return "", "", 0, domain.ErrTooMany2FAAttempts
	}

	delete2FASession(sessionID)

	user, err := s.usersRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return "", "", 0, fmt.Errorf("get user: %w", err)
	}

	if !user.TwoFAEnabled || user.TwoFASecret == "" {
		return "", "", 0, errors.New("2FA is not enabled")
	}

	encKey := []byte(s.tokenizer.SecretKey())
	encSecret, err := base64.StdEncoding.DecodeString(user.TwoFASecret)
	if err != nil {
		return "", "", 0, fmt.Errorf("decode secret: %w", err)
	}
	plainSecret, err := crypt.DecryptAESGCM(encSecret, encKey)
	if err != nil {
		return "", "", 0, fmt.Errorf("decrypt secret: %w", err)
	}

	valid := totp.Validate(code, string(plainSecret))
	if !valid {
		_, blocked := s.twoFARateLimiter.Inc(userID)
		if blocked {
			return "", "", 0, domain.ErrTooMany2FAAttempts
		}

		return "", "", 0, domain.ErrInvalid2FACode
	}

	s.twoFARateLimiter.Reset(userID)

	accessToken, err = s.tokenizer.AccessToken(&user)
	if err != nil {
		return "", "", 0, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err = s.tokenizer.RefreshToken(&user)
	if err != nil {
		return "", "", 0, fmt.Errorf("generate refresh token: %w", err)
	}

	expiresIn = int(s.tokenizer.AccessTokenTTL().Seconds())

	return accessToken, refreshToken, expiresIn, nil
}

func generate2FASession(userID domain.UserID, username string, ttl time.Duration) string {
	sessionID := uuid.NewString()
	twoFASessionStore.Lock()
	twoFASessionStore.sessions[sessionID] = twoFASessionEntry{
		UserID:    userID,
		Username:  username,
		CreatedAt: time.Now(),
	}
	twoFASessionStore.Unlock()
	// Очистка устаревших сессий (простая, неэффективная, но для in-memory ок)
	go func() {
		time.Sleep(ttl)
		twoFASessionStore.Lock()
		delete(twoFASessionStore.sessions, sessionID)
		twoFASessionStore.Unlock()
	}()

	return sessionID
}

func get2FASession(sessionID string) (twoFASessionEntry, bool) {
	twoFASessionStore.Lock()
	entry, ok := twoFASessionStore.sessions[sessionID]
	twoFASessionStore.Unlock()

	return entry, ok
}

func delete2FASession(sessionID string) {
	twoFASessionStore.Lock()
	delete(twoFASessionStore.sessions, sessionID)
	twoFASessionStore.Unlock()
}
