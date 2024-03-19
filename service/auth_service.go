package service

import (
	"context"
	"errors"
	"github.com/luquxSentinel/spacedrive/storage"
	"github.com/luquxSentinel/spacedrive/tokens"
	"github.com/luquxSentinel/spacedrive/types"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"

	"github.com/google/uuid"
)

type AuthService interface {
	CreateUser(ctx context.Context, data *types.CreateUserData) error
}

type authService struct {
	storage storage.Storage
}


func NewAuthService() *authService {
	return &authService{}
}


func (s *authService) CreateUser(ctx context.Context, data *types.CreateUserData) error {
	// TODO: create a new user
	
	// TODO: check if no user is associated with account
	emailCount, err := s.storage.CountEmail(ctx, data.Email)
	if err != nil {
		return err
	}

	if emailCount > 0 {
		return errors.New("email already in use")
	}
	// TODO: new user from data
	newUser := new(types.User)

	// generate new user id
	newUser.UID = uuid.NewString()
	newUser.Email = data.Email
	newUser.FirstName = data.FirstName
	newUser.LastName = data.LastName

	// hash password
	newUser.Password, err = HashPassword(data.Password)
	if err != nil {
		return err
	}

	// set created at time
	newUser.CreatedAt =  time.Now().Local()

	// persist user in database
	return  s.storage.CreateUser(ctx, newUser)
}

func (s *authService) LoginUser (ctx context.Context, data *types.LoginData) (*types.User, string, error) {
//	TODO: login user logic


	// TODO: fetch user by email
	user, err := s.storage.GetUserWithEmail(ctx, data.Email)
	if err != nil {
		log.Printf("failed to fetch user with email. error : %v", err)
		return nil, "", errors.New("wrong email or password")
	}

	// TODO: verify user password
	if err := verifyPassword(user.Password, data.Password); err != nil {
		return nil, "", errors.New("wrong email or password")
	}

	// TODO: generate jwt
	signedToken, err := tokens.GenerateJWT(user.UID, user.Email)
	if err != nil {
		log.Panic(err)
		return nil, "", err
	}

	// return user, jwt, error
	return user, signedToken, nil
}


func (s *authService) DeleteUser() {
//	TODO: delete user account logic
}

func (s *authService) UpdateUser() {
//	TODO: update user logic
}


func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(b), err
}


func verifyPassword(foundPassword, givenPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(foundPassword), []byte(givenPassword))
}