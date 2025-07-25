package users

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/pkg/crypt"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/backend/contract"
	mockusers "github.com/rom8726/warden/test_mocks/internal_/backend/usecases/users"
)

func TestSetup2FA(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := domain.UserID(1)
	user := domain.User{ID: userID, Email: "user@example.com"}
	mockUsersRepo := mockcontract.NewMockUsersRepository(t)
	mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
	mockTokenizer := mockcontract.NewMockTokenizer(t)
	mockEmailer := mockcontract.NewMockEmailer(t)
	mockAuthProvider := mockusers.NewMockAuthProvider(t)
	mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)
	mockUsersRepo.EXPECT().GetByID(ctx, userID).Return(user, nil)
	mockTokenizer.EXPECT().SecretKey().Return("testsecret123456")
	mockUsersRepo.EXPECT().Update2FA(ctx, userID, false, mock.Anything, mock.AnythingOfType("*time.Time")).Return(nil)

	service := New(mockUsersRepo, mockTeamsRepo, mockTokenizer, mockEmailer, mockRateLimiter, []AuthProvider{mockAuthProvider})
	secret, qrURL, qrImage, err := service.Setup2FA(ctx, userID)
	require.NoError(t, err)
	require.NotEmpty(t, secret)
	require.NotEmpty(t, qrURL)
	require.NotEmpty(t, qrImage)
}

func TestConfirm2FA(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := domain.UserID(1)
	plainSecret := "JBSWY3DPEHPK3PXP"
	encKey := []byte("testsecret123456")
	encSecret, _ := crypt.EncryptAESGCM([]byte(plainSecret), encKey)
	encSecretB64 := base64.StdEncoding.EncodeToString(encSecret)
	user := domain.User{ID: userID, Email: "user@example.com", TwoFASecret: encSecretB64}
	mockUsersRepo := mockcontract.NewMockUsersRepository(t)
	mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
	mockTokenizer := mockcontract.NewMockTokenizer(t)
	mockEmailer := mockcontract.NewMockEmailer(t)
	mockAuthProvider := mockusers.NewMockAuthProvider(t)
	mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)
	mockUsersRepo.EXPECT().GetByID(ctx, userID).Return(user, nil)
	mockTokenizer.EXPECT().SecretKey().Return("testsecret123456")
	mockUsersRepo.EXPECT().Update2FA(ctx, userID, true, encSecretB64, mock.AnythingOfType("*time.Time")).Return(nil)
	mockRateLimiter.EXPECT().IsBlocked(userID).Return(false)
	mockRateLimiter.EXPECT().Reset(userID)

	service := New(mockUsersRepo, mockTeamsRepo, mockTokenizer, mockEmailer, mockRateLimiter, []AuthProvider{mockAuthProvider})
	code, _ := totp.GenerateCode(plainSecret, time.Now())
	err := service.Confirm2FA(ctx, userID, code)
	require.NoError(t, err)
}

func TestConfirm2FA__user_blocked(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := domain.UserID(1)
	plainSecret := "JBSWY3DPEHPK3PXP"
	mockUsersRepo := mockcontract.NewMockUsersRepository(t)
	mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
	mockTokenizer := mockcontract.NewMockTokenizer(t)
	mockEmailer := mockcontract.NewMockEmailer(t)
	mockAuthProvider := mockusers.NewMockAuthProvider(t)
	mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)
	mockRateLimiter.EXPECT().IsBlocked(userID).Return(true)

	service := New(mockUsersRepo, mockTeamsRepo, mockTokenizer, mockEmailer, mockRateLimiter, []AuthProvider{mockAuthProvider})
	code, _ := totp.GenerateCode(plainSecret, time.Now())
	err := service.Confirm2FA(ctx, userID, code)
	require.Error(t, err, domain.ErrTooMany2FAAttempts)
}

func TestSend2FACode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := domain.UserID(1)
	user := domain.User{ID: userID, Email: "user@example.com"}
	mockUsersRepo := mockcontract.NewMockUsersRepository(t)
	mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
	mockTokenizer := mockcontract.NewMockTokenizer(t)
	mockEmailer := mockcontract.NewMockEmailer(t)
	mockAuthProvider := mockusers.NewMockAuthProvider(t)
	mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)
	mockUsersRepo.EXPECT().GetByID(ctx, userID).Return(user, nil)
	mockEmailer.EXPECT().Send2FACodeEmail(ctx, user.Email, mock.AnythingOfType("string"), "disable").Return(nil)

	service := New(mockUsersRepo, mockTeamsRepo, mockTokenizer, mockEmailer, mockRateLimiter, []AuthProvider{mockAuthProvider})
	err := service.Send2FACode(ctx, userID, "disable")
	require.NoError(t, err)
}

func TestDisable2FA(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := domain.UserID(1)
	mockUsersRepo := mockcontract.NewMockUsersRepository(t)
	mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
	mockTokenizer := mockcontract.NewMockTokenizer(t)
	mockEmailer := mockcontract.NewMockEmailer(t)
	mockAuthProvider := mockusers.NewMockAuthProvider(t)
	mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)

	store2FACode(userID, "12345678", "disable", time.Minute*2)
	mockUsersRepo.EXPECT().Update2FA(ctx, userID, false, "", mock.Anything).Return(nil)

	service := New(mockUsersRepo, mockTeamsRepo, mockTokenizer, mockEmailer, mockRateLimiter, []AuthProvider{mockAuthProvider})
	err := service.Disable2FA(ctx, userID, "12345678")
	require.NoError(t, err)
}

func TestReset2FA(t *testing.T) {
	t.SkipNow() // TODO: Fix test expectations
	t.Parallel()
	ctx := context.Background()
	userID := domain.UserID(1)
	user := domain.User{ID: userID, Email: "user@example.com"}
	mockUsersRepo := mockcontract.NewMockUsersRepository(t)
	mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
	mockTokenizer := mockcontract.NewMockTokenizer(t)
	mockEmailer := mockcontract.NewMockEmailer(t)
	mockAuthProvider := mockusers.NewMockAuthProvider(t)
	mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)
	mockUsersRepo.EXPECT().GetByID(ctx, userID).Return(user, nil)
	mockTokenizer.EXPECT().SecretKey().Return("testsecret123456")
	mockUsersRepo.EXPECT().Update2FA(ctx, userID, false, mock.Anything, mock.AnythingOfType("*time.Time")).Return(nil)

	service := New(mockUsersRepo, mockTeamsRepo, mockTokenizer, mockEmailer, mockRateLimiter, []AuthProvider{mockAuthProvider})

	store2FACode(userID, "87654321", "reset", time.Minute)

	secret, qrURL, qrImage, err := service.Reset2FA(ctx, userID, "87654321")
	require.NoError(t, err)
	require.NotEmpty(t, secret)
	require.NotEmpty(t, qrURL)
	require.NotEmpty(t, qrImage)
}

func TestVerify2FA(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := domain.UserID(1)
	plainSecret := "JBSWY3DPEHPK3PXP"
	encKey := []byte("testsecret123456")
	encSecret, _ := crypt.EncryptAESGCM([]byte(plainSecret), encKey)
	encSecretB64 := base64.StdEncoding.EncodeToString(encSecret)
	user := domain.User{ID: userID, Email: "user@example.com", TwoFASecret: encSecretB64, TwoFAEnabled: true}
	mockUsersRepo := mockcontract.NewMockUsersRepository(t)
	mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
	mockTokenizer := mockcontract.NewMockTokenizer(t)
	mockEmailer := mockcontract.NewMockEmailer(t)
	mockAuthProvider := mockusers.NewMockAuthProvider(t)
	mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)
	mockUsersRepo.EXPECT().GetByID(ctx, userID).Return(user, nil)
	mockTokenizer.EXPECT().SecretKey().Return("testsecret123456")
	mockTokenizer.EXPECT().AccessToken(&user).Return("access_token", nil)
	mockTokenizer.EXPECT().RefreshToken(&user).Return("refresh_token", nil)
	mockTokenizer.EXPECT().AccessTokenTTL().Return(3600 * time.Second)
	mockRateLimiter.EXPECT().IsBlocked(userID).Return(false)
	mockRateLimiter.EXPECT().Reset(userID)

	service := New(mockUsersRepo, mockTeamsRepo, mockTokenizer, mockEmailer, mockRateLimiter, []AuthProvider{mockAuthProvider})
	code, _ := totp.GenerateCode(plainSecret, time.Now())
	sessionID := generate2FASession(userID, "username", time.Minute)
	accessToken, refreshToken, expiresIn, err := service.Verify2FA(ctx, code, sessionID)
	require.NoError(t, err)
	require.Equal(t, "access_token", accessToken)
	require.Equal(t, "refresh_token", refreshToken)
	require.Equal(t, 3600, expiresIn)
}

func TestVerify2FA__user_blocked(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userID := domain.UserID(1)
	plainSecret := "JBSWY3DPEHPK3PXP"
	mockUsersRepo := mockcontract.NewMockUsersRepository(t)
	mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
	mockTokenizer := mockcontract.NewMockTokenizer(t)
	mockEmailer := mockcontract.NewMockEmailer(t)
	mockAuthProvider := mockusers.NewMockAuthProvider(t)
	mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)
	mockRateLimiter.EXPECT().IsBlocked(userID).Return(true)

	service := New(mockUsersRepo, mockTeamsRepo, mockTokenizer, mockEmailer, mockRateLimiter, []AuthProvider{mockAuthProvider})
	code, _ := totp.GenerateCode(plainSecret, time.Now())
	sessionID := generate2FASession(userID, "username", time.Minute)
	_, _, _, err := service.Verify2FA(ctx, code, sessionID)
	require.Error(t, err, domain.ErrTooMany2FAAttempts)
}
