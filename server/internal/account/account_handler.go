package account

import (
	"context"
	"dungeons_helper/server/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"time"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var account CreateAccountReq
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}(r.Body)

	ctx := r.Context()
	err = h.Service.CreateAccount(ctx, &account)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		ok := errors.As(err, &mysqlErr)
		if ok && mysqlErr.Number == 1062 {
			http.Error(w, err.Error(), http.StatusConflict)
		} else if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "Wrong password or email", http.StatusRequestTimeout)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var account LoginAccountReq
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}(r.Body)

	ctx := r.Context()
	res, err := h.Service.Login(ctx, &account)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			http.Error(w, "Wrong password or email", http.StatusUnauthorized)
		} else if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "Wrong password or email", http.StatusRequestTimeout)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    res.accessToken,
		Expires:  time.Now().Add(time.Hour),
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}

func (h *Handler) Logout(w http.ResponseWriter, _ *http.Request) {
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "logout successful"}`))
}

func (h *Handler) RestorePassword(w http.ResponseWriter, r *http.Request) {
	var restoreReq struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&restoreReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}(r.Body)

	if restoreReq.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	pwd, err := h.Service.RestorePassword(ctx, restoreReq.Email)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "Wrong password or email", http.StatusRequestTimeout)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	response := fmt.Sprintf(`{"message": "Password reset instructions sent to your email. Temporary password: %s"}`, pwd)
	_, _ = w.Write([]byte(response))
}

func (h *Handler) UpdateNickname(w http.ResponseWriter, r *http.Request) {
	var req UpdateNicknameReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}(r.Body)

	id, err := util.GetIdFromHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	req = UpdateNicknameReq{Id: id, Nickname: req.Nickname}
	ctx := r.Context()
	err = h.Service.UpdateNickname(ctx, &req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "Wrong password or email", http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "Nickname updated successfully"}`))
}

func (h *Handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var req UpdatePasswordReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}(r.Body)

	id, err := util.GetIdFromHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	req = UpdatePasswordReq{Id: id, OldPassword: req.OldPassword, NewPassword: req.NewPassword}
	ctx := r.Context()
	err = h.Service.UpdatePassword(ctx, &req)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			http.Error(w, "Wrong old password", http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)

}
