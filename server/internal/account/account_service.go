package account

import (
	"context"
	util2 "dungeons_helper/server/util"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
	"time"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
	}
}

func (s *service) CreateAccount(c context.Context, req *CreateAccountReq) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := util2.HashPassword(req.Password)
	if err != nil {
		return err
	}

	account := &Account{
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Nickname,
		IdAvatar: req.IdAvatar,
	}

	err = s.Repository.CreateAccount(ctx, account)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Login(c context.Context, req *LoginAccountReq) (*LoginAccountRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	account, err := s.Repository.GetAccountByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	err = util2.CheckPassword(req.Password, account.Password)
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, util2.MyJWTClaims{
		Id:       account.Id,
		Nickname: account.Nickname,
		IdAvatar: account.IdAvatar,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(account.Id)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	secretKey := os.Getenv("SECRET_KEY")
	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}
	return &LoginAccountRes{accessToken: ss, Id: account.Id, Email: account.Email, Nickname: account.Nickname, IdAvatar: account.IdAvatar}, nil
}

func (s *service) RestorePassword(c context.Context, email string) (string, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	account, err := s.Repository.GetAccountByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	tempPassword := util2.GeneratePassword()
	hashedPassword, err := util2.HashPassword(tempPassword)
	if err != nil {
		return "", err
	}

	account.Password = hashedPassword
	if err := s.Repository.UpdatePassword(ctx, account); err != nil {
		return "", err
	}

	//TODO: Send email with new password and delete tempPassword from return
	return tempPassword, nil
}

func (s *service) UpdateNickname(c context.Context, req *UpdateNicknameReq) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	account, err := s.Repository.GetAccountById(ctx, req.Id)
	if err != nil {
		return err
	}

	account.Nickname = req.Nickname
	err = s.Repository.UpdateNickname(ctx, account)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) UpdatePassword(c context.Context, req *UpdatePasswordReq) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	account, err := s.Repository.GetAccountById(ctx, req.Id)
	if err != nil {
		return err
	}

	err = util2.CheckPassword(req.OldPassword, account.Password)
	if err != nil {
		return err
	}

	hashedPassword, err := util2.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	account.Password = hashedPassword
	err = s.Repository.UpdatePassword(ctx, account)
	if err != nil {
		return err
	}

	return nil
}
