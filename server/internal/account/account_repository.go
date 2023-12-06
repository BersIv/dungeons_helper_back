package account

import (
	"context"
	"dungeons_helper/server/db"
)

type repository struct {
	db db.DatabaseTX
}

func NewRepository(db db.DatabaseTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateAccount(ctx context.Context, account *Account) error {
	query := "INSERT INTO account(email, password, nickname, idAvatar) VALUES (?, ?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, account.Email, account.Password, account.Nickname, account.IdAvatar)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetAccountByEmail(ctx context.Context, email string) (*Account, error) {
	account := Account{}

	query := "SELECT id, email, password, nickname, idAvatar FROM account WHERE email = ?"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&account.Id, &account.Email, &account.Password, &account.Nickname, &account.IdAvatar)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *repository) GetAccountById(ctx context.Context, id int64) (*Account, error) {
	account := Account{}

	query := "SELECT id, email, password, nickname, idAvatar FROM account WHERE id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&account.Id, &account.Email, &account.Password, &account.Nickname, &account.IdAvatar)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *repository) UpdatePassword(ctx context.Context, account *Account) error {
	query := "UPDATE account SET password = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, account.Password, account.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateNickname(ctx context.Context, account *Account) error {
	query := "UPDATE account SET nickname = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, account.Nickname, account.Id)
	if err != nil {
		return err
	}

	return nil
}
