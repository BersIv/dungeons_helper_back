package character

import (
	"context"
	"dungeons_helper/internal/alignment"
	"dungeons_helper/internal/class"
	"dungeons_helper/internal/image"
	"dungeons_helper/internal/races"
	"dungeons_helper/internal/skills"
	"dungeons_helper/internal/stats"
	"dungeons_helper/internal/subraces"
)

type Character struct {
	Id              int64               `json:"id"`
	Hp              int64               `json:"hp"`
	Exp             int64               `json:"exp"`
	Avatar          image.Image         `json:"Avatar"`
	CharName        string              `json:"charName"`
	Sex             bool                `json:"sex"`
	Weight          int64               `json:"weight"`
	Height          int64               `json:"height"`
	Class           class.Class         `json:"class"`
	Race            races.Races         `json:"race"`
	Subrace         subraces.Subraces   `json:"subrace"`
	Stats           stats.Stats         `json:"stats"`
	AddLanguage     string              `json:"addLanguage"`
	CharacterSkills string              `json:"characterSkills"`
	Alignment       alignment.Alignment `json:"alignment"`
	Ideals          string              `json:"ideals"`
	Weaknesses      string              `json:"weaknesses"`
	Traits          string              `json:"traits"`
	Allies          string              `json:"allies"`
	Organizations   string              `json:"organizations"`
	Enemies         string              `json:"enemies"`
	Story           string              `json:"story"`
	Goals           string              `json:"goals"`
	Treasures       string              `json:"treasures"`
	Notes           string              `json:"notes"`
}

type CreateCharacterReq struct {
	IdAcc int64 `json:"idAcc"`
	//Id              int64               `json:"id"`
	Hp              int64               `json:"hp"`
	Exp             int64               `json:"exp"`
	Avatar          image.Image         `json:"Avatar"`
	CharName        string              `json:"charName"`
	Sex             bool                `json:"sex"`
	Weight          int64               `json:"weight"`
	Height          int64               `json:"height"`
	Class           class.Class         `json:"class"`
	Race            races.Races         `json:"race"`
	Subrace         subraces.Subraces   `json:"subrace"`
	Stats           stats.Stats         `json:"stats"`
	AddLanguage     string              `json:"addLanguage"`
	CharacterSkills []skills.Skills     `json:"characterSkills"`
	Alignment       alignment.Alignment `json:"alignment"`
	Ideals          string              `json:"ideals"`
	Weaknesses      string              `json:"weaknesses"`
	Traits          string              `json:"traits"`
	Allies          string              `json:"allies"`
	Organizations   string              `json:"organizations"`
	Enemies         string              `json:"enemies"`
	Story           string              `json:"story"`
	Goals           string              `json:"goals"`
	Treasures       string              `json:"treasures"`
	Notes           string              `json:"notes"`
}

type GetCharacterReq struct {
	Id int64 `json:"id"`
}

type Repository interface {
	GetAllCharactersByAccId(ctx context.Context, idAcc int64) ([]Character, error)
	GetCharacterById(ctx context.Context, id int64) (*Character, error)
	CreateCharacter(ctx context.Context, character *CreateCharacterReq, idAcc int64) error
	UpdateCharacterHpById(ctx context.Context, id int64, hp int64) error
}

type Service interface {
	GetAllCharactersByAccId(c context.Context, idAcc int64) ([]Character, error)
	GetCharacterById(c context.Context, id int64) (*Character, error)
	CreateCharacter(c context.Context, character *CreateCharacterReq, idAcc int64) error
}
