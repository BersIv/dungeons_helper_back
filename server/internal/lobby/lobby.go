package lobby

import (
	"context"
	"dungeons_helper/server/internal/character"
)

type Lobby struct {
	Id            int64  `json:"id"`
	LobbyMasterID int64  `json:"lobbyMasterID"`
	LobbyName     string `json:"lobbyName"`
	LobbyPassword string `json:"lobbyPassword"`
	Amount        int64  `json:"amount"`
}

type CreateLobbyReq struct {
	LobbyMasterID int64  `json:"lobbyMasterID"`
	LobbyName     string `json:"lobbyName"`
	LobbyPassword string `json:"lobbyPassword"`
	Amount        int64  `json:"amount"`
}

type CreateLobbyRes struct {
	Id int64 `json:"id"`
}

type GetLobbyRes struct {
	Id          int64  `json:"id"`
	LobbyName   string `json:"lobbyName"`
	LobbyMaster string `json:"lobbyMaster"`
}

type JoinLobbyReq struct {
	IdLobby       int64  `json:"idLobby"`
	IdAcc         int64  `json:"idAcc"`
	LobbyPassword string `json:"lobbyPassword"`
}

type GetLobbyByIdRes struct {
	Password string `json:"password"`
}

type Repository interface {
	CreateLobby(ctx context.Context, lobby *CreateLobbyReq) (*CreateLobbyRes, error)
	GetAllLobby(ctx context.Context) ([]GetLobbyRes, error)
	GetLobbyById(ctx context.Context, id int64) (*GetLobbyByIdRes, error)
	JoinLobby(ctx context.Context, req *JoinLobbyReq) ([]character.Character, error)
}

type Service interface {
	CreateLobby(c context.Context, lobby *CreateLobbyReq) (*CreateLobbyRes, error)
	GetAllLobby(c context.Context) ([]GetLobbyRes, error)
	JoinLobby(c context.Context, req *JoinLobbyReq) ([]character.Character, error)
}
