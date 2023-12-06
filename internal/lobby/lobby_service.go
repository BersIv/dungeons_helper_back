package lobby

import (
	"context"
	"dungeons_helper_server/internal/character"
	"errors"
	"time"
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

func (s *service) CreateLobby(c context.Context, lobby *CreateLobbyReq) (*CreateLobbyRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	res, err := s.Repository.CreateLobby(ctx, lobby)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (s *service) GetAllLobby(c context.Context) ([]GetLobbyRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	res, err := s.Repository.GetAllLobby(ctx)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (s *service) JoinLobby(c context.Context, req *JoinLobbyReq) ([]character.Character, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	password, err := s.Repository.GetLobbyById(ctx, req.IdLobby)
	if err != nil {
		return nil, err
	}

	if password.Password != req.LobbyPassword {
		return nil, errors.New("Wrong lobby password")
	}

	res, err := s.Repository.JoinLobby(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, err
}
