package account

import (
	"context"
	"dungeons_helper/db"
)

type repository struct {
	db db.DatabaseTX
}

func NewRepository(db db.DatabaseTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateAccount(ctx context.Context, req *CreateAccountReq) error {
	query := "INSERT INTO image(image) VALUES (?)"
	result, err := r.db.ExecContext(ctx, query, req.Avatar)
	if err != nil {
		return err
	}
	idImage, err := result.LastInsertId()
	if err != nil {
		return err
	}

	query = "INSERT INTO account(email, password, nickname, idAvatar) VALUES (?, ?, ?, ?)"
	_, err = r.db.ExecContext(ctx, query, req.Email, req.Password, req.Nickname, idImage)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetAccountByEmail(ctx context.Context, email string) (*LoginAccountRes, error) {
	res := LoginAccountRes{}
	query := `SELECT a.id, email, password, nickname, i.image FROM account a 
				LEFT JOIN image i ON a.idAvatar = i.id 
				WHERE email = ?`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&res.Id, &res.Email, &res.Password, &res.Nickname, &res.Avatar.Image)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *repository) GetAccountById(ctx context.Context, id int64) (*LoginAccountRes, error) {
	account := LoginAccountRes{}

	query := `SELECT a.id, email, password, nickname, i.image FROM account a 
				LEFT JOIN image i ON a.idAvatar = i.id 
				WHERE a.id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&account.Id, &account.Email, &account.Password, &account.Nickname, &account.Avatar)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *repository) UpdatePassword(ctx context.Context, account *LoginAccountRes) error {
	query := "UPDATE account SET password = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, account.Password, account.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateNickname(ctx context.Context, account *LoginAccountRes) error {
	query := "UPDATE account SET nickname = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, account.Nickname, account.Id)
	if err != nil {
		return err
	}

	return nil
}
