package account

import (
	"context"
	"dungeons_helper_server/util"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

const (
	secretKey = "secret"
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

func (s *service) CreateAccount(c context.Context, req *CreateAccountReq) (*CreateAccountRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	account := &Account{
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Nickname,
		IdAvatar: req.IdAvatar,
	}

	r, err := s.Repository.CreateAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	res := &CreateAccountRes{
		Id:       r.Id,
		Email:    r.Email,
		Nickname: r.Nickname,
		IdAvatar: r.IdAvatar,
	}

	return res, nil
}

type MyJWTClaims struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
	IdAvatar int64  `json:"idAvatar"`
	jwt.RegisteredClaims
}

func (s *service) Login(c context.Context, req *LoginAccountReq) (*LoginAccountRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	time.Sleep(5 * time.Second)

	account, err := s.Repository.GetAccountByEmail(ctx, req.Email)
	if err != nil {
		return &LoginAccountRes{}, err
	}

	err = util.CheckPassword(req.Password, account.Password)
	if err != nil {
		return &LoginAccountRes{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		Id:       account.Id,
		Nickname: account.Nickname,
		IdAvatar: account.IdAvatar,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(account.Id)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return &LoginAccountRes{}, err
	}
	return &LoginAccountRes{accessToken: ss, Id: account.Id, Nickname: account.Nickname, IdAvatar: account.IdAvatar}, nil
}

func (s *service) RestorePassword(c context.Context, email string) (string, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	account, err := s.Repository.GetAccountByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	tempPassword := util.GeneratePassword()
	hashedPassword, err := util.HashPassword(tempPassword)
	if err != nil {
		return "", err
	}

	account.Password = hashedPassword
	if err := s.Repository.UpdatePassword(c, account); err != nil {
		return "", err
	}

	//TODO: Send email with new password

	return tempPassword, nil
}

func (s *service) UpdateNickname(c context.Context, id int64, newNickname string) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	account, err := s.Repository.GetAccountById(ctx, id)
	if err != nil {
		return err
	}

	account.Nickname = newNickname
	err = s.Repository.UpdateNickname(ctx, account)
	if err != nil {
		return err
	}

	return nil
}
