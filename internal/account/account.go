package account

import "context"

type Account struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	IdAvatar int64  `json:"IdAvatar"`
}

type CreateAccountReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	IdAvatar int64  `json:"IdAvatar"`
}

type CreateAccountRes struct {
	Id       int64  `json:"id"`
	Email    string `json:"Email"`
	Nickname string `json:"nickname"`
	IdAvatar int64  `json:"IdAvatar"`
}

type LoginAccountReq struct {
	Email    string `json:"Email"`
	Password string `json:"password"`
}

type LoginAccountRes struct {
	accessToken string
	Id          int64  `json:"id"`
	Nickname    string `json:"nickname"`
	IdAvatar    int64  `json:"IdAvatar"`
}

type updateReq struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
}

type Repository interface {
	CreateAccount(ctx context.Context, account *Account) (*Account, error)
	GetAccountByEmail(ctx context.Context, email string) (*Account, error)
	GetAccountById(ctx context.Context, id int64) (*Account, error)
	UpdatePassword(ctx context.Context, account *Account) error
	UpdateNickname(ctx context.Context, account *Account) error
}

type Service interface {
	CreateAccount(ctx context.Context, req *CreateAccountReq) (*CreateAccountRes, error)
	Login(c context.Context, req *LoginAccountReq) (*LoginAccountRes, error)
	RestorePassword(c context.Context, email string) (string, error)
	UpdateNickname(ctx context.Context, id int64, newNickname string) error
}
