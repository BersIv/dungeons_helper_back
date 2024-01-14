package account

import (
	"context"
	"dungeons_helper/internal/image"
)

type Account struct {
	Id       int64       `json:"id"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Nickname string      `json:"nickname"`
	Avatar   image.Image `json:"avatar"`
}

type IdAccountReq struct {
	IdAcc int64 `json:"idAcc"`
}

type CreateAccountReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type LoginAccountReq struct {
	Email    string `json:"Email"`
	Password string `json:"password"`
}

type LoginAccountRes struct {
	accessToken string
	Id          int64       `json:"id"`
	Email       string      `json:"email"`
	Nickname    string      `json:"nickname"`
	Avatar      image.Image `json:"IdAvatar"`
	Password    string      `json:"password"`
}

type UpdateNicknameReq struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
}

type UpdatePasswordReq struct {
	Id          int64  `json:"id"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type GoogleAcc struct {
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

type Token struct {
	Token string `json:"token"`
}

type Repository interface {
	CreateAccount(ctx context.Context, account *CreateAccountReq) error
	GetAccountByEmail(ctx context.Context, email string) (*LoginAccountRes, error)
	GetAccountById(ctx context.Context, id int64) (*LoginAccountRes, error)
	UpdatePassword(ctx context.Context, account *LoginAccountRes) error
	UpdateNickname(ctx context.Context, account *LoginAccountRes) error
}

type Service interface {
	CreateAccount(c context.Context, req *CreateAccountReq) error
	Login(c context.Context, req *LoginAccountReq) (*LoginAccountRes, error)
	RestorePassword(c context.Context, email string) error
	UpdateNickname(c context.Context, req *UpdateNicknameReq) error
	UpdatePassword(c context.Context, req *UpdatePasswordReq) error
	GoogleAuth(c context.Context, req *GoogleAcc) (*LoginAccountRes, error)
}
