package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"drivex/internal/models"
	"drivex/internal/repository"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
	sessionRepo *repository.SessionRepository
	jwtSecret string 
}

func NewAuthService(
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		sessionRepo: sessionRepo,
		jwtSecret :jwtSecret,
	}
}


func (s *AuthService) Register(
	ctx context.Context,
	name string,
	email string,
	password string,
) (*models.User, error) {

	if name == "" {
		return nil, errors.New("name is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	if password == "" {
		return nil, errors.New("password is required")
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.Create(
		ctx,
		name,
		email,
		string(passwordHash),
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func ( s *AuthService) Login ( ctx context.Context , email string , password string) ( *models.User  , string , string , error){
	user , err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", "", errors.New("invalid email or password")
	}



	if user.PasswordHash == nil {
		return nil, "", "", errors.New("please login using OAuth")
	}

	err = bcrypt.CompareHashAndPassword(
[]byte(*user.PasswordHash),
		[]byte(password),

	)


	if err != nil {
		return nil, "", "", errors.New("invalid email or password")
	}


	accessToken , err := s.generateAccessToken ( user)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, "", "", err
	}


	refreshTokenHash := hashToken(refreshToken)
	err = s.sessionRepo.Create(
    ctx,
    user.ID,
    refreshTokenHash,
    time.Now().Add(7*24*time.Hour),
)
	return user, accessToken, refreshToken, nil
}


func (s *AuthService) generateAccessToken(
	user *models.User,
) (string, error) {

	claims := jwt.MapClaims{
		"sub": user.ID,
		"email": user.Email,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(s.jwtSecret))
}




func generateRefreshToken() (string, error) {

	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
func hashToken(token string) string {
    hash := sha256.Sum256([]byte(token))
    return hex.EncodeToString(hash[:])
}

func (s *AuthService) Refresh(
	ctx context.Context,
	refreshToken string,
) (string, error) {

	if refreshToken == "" {
		return "", errors.New("refresh token required")
	}

	// Refresh token ka hash banao
	refreshTokenHash := hashToken(refreshToken)

	// Session se user ID + expiry nikalo
	userID, expiresAt, err :=
		s.sessionRepo.GetRefreshToken(
			ctx,
			refreshTokenHash,
		)

	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Refresh token expire ho gaya?
	if time.Now().After(expiresAt) {
		return "", errors.New("refresh token expired")
	}

	// User ko DB se fetch karo
	user, err := s.userRepo.GetByID(ctx, userID)

	if err != nil {
		return "", errors.New("user not found")
	}

	// New access token
	accessToken, err := s.generateAccessToken(user)

	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *AuthService) Logout(
	ctx context.Context,
	refreshToken string,
) error {

	if refreshToken == "" {
		return errors.New("refresh token required")
	}

	refreshTokenHash := hashToken(refreshToken)

	return s.sessionRepo.DeleteByRefreshTokenHash(
		ctx,
		refreshTokenHash,
	)
}