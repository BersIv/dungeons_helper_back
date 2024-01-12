package account

import (
	"context"
	"dungeons_helper/internal/image"
	"dungeons_helper/util"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return err
	}

	account := &CreateAccountReq{
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
	}

	err = s.Repository.CreateAccount(ctx, account)
	if err != nil {
		return err
	}

	err = sendWelcomeEmail(req.Email)
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

	err = util.CheckPassword(req.Password, account.Password)
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, util.MyJWTClaims{
		Id:       account.Id,
		Nickname: account.Nickname,
		Avatar:   image.Image{},
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
	return &LoginAccountRes{accessToken: ss, Id: account.Id, Email: account.Email, Nickname: account.Nickname, Avatar: account.Avatar}, nil
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

	err = util.CheckPassword(req.OldPassword, account.Password)
	if err != nil {
		return err
	}

	hashedPassword, err := util.HashPassword(req.NewPassword)
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

func sendWelcomeEmail(toEmail string) error {
	// Параметры для почтового сервера mail.ru
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpUsername := "ivanbers1998@gmail.com"
	smtpPassword := "tajh doie vrtv azfj"

	// Сообщение для отправки
	subject := "Добро пожаловать!"
	body := "Спасибо за регистрацию!"

	// Формирование заголовков сообщения
	headers := make(map[string]string)
	headers["From"] = smtpUsername
	headers["To"] = toEmail
	headers["Subject"] = subject

	// Формирование тела сообщения
	message := ""
	for key, value := range headers {
		message += key + ": " + value + "\r\n"
	}
	message += "\r\n" + body

	// Настройка параметров аутентификации
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// Формирование адреса сервера
	serverAddr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	fmt.Println(serverAddr)
	fmt.Println(auth)
	fmt.Println(smtpUsername)
	// Отправка письма
	err := smtp.SendMail(serverAddr, auth, smtpUsername, []string{toEmail}, []byte(message))
	if err != nil {
		return err
	}

	return nil
}
