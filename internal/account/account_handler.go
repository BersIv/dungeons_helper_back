package account

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
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
	res, err := h.Service.CreateAccount(ctx, &account)
	if err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if ok && mysqlErr.Number == 1062 {
			http.Error(w, err.Error(), http.StatusConflict)
		} else if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "Wrong password or email", http.StatusRequestTimeout)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	jsonResponse, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(jsonResponse)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var account LoginAccountReq
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
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
		MaxAge:   60 * 60 * 24,
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
	var req UpdateReq

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

	ctx := r.Context()
	err := h.Service.UpdateNickname(ctx, req.Id, req.Nickname)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "Wrong password or email", http.StatusRequestTimeout)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "Nickname updated successfully"}`))
}
