package lobby

import (
	"dungeons_helper/util"
	"encoding/json"
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

func (h *Handler) CreateLobby(w http.ResponseWriter, r *http.Request) {
	var lobbyReq CreateLobbyReq
	err := json.NewDecoder(r.Body).Decode(&lobbyReq)
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

	//idAcc, err := util.GetIdFromToken(r)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusUnauthorized)
	//	return
	//}

	lobbyReq = CreateLobbyReq{LobbyMasterID: lobbyReq.LobbyMasterID, LobbyPassword: lobbyReq.LobbyPassword, LobbyName: lobbyReq.LobbyName, Amount: lobbyReq.Amount}
	ctx := r.Context()
	res, err := h.Service.CreateLobby(ctx, &lobbyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func (h *Handler) GetAllLobby(w http.ResponseWriter, r *http.Request) {
	_, err := util.GetIdFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	res, err := h.Service.GetAllLobby(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)

}

func (h *Handler) JoinLobby(w http.ResponseWriter, r *http.Request) {
	idAcc, err := util.GetIdFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var lobbyReq JoinLobbyReq
	err = json.NewDecoder(r.Body).Decode(&lobbyReq)
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

	lobbyReq = JoinLobbyReq{IdLobby: lobbyReq.IdLobby, LobbyPassword: lobbyReq.LobbyPassword, IdAcc: idAcc}
	ctx := r.Context()
	res, err := h.Service.JoinLobby(ctx, &lobbyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}
