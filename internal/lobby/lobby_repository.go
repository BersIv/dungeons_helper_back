package lobby

import (
	"context"
	"dungeons_helper_server/db"
	"dungeons_helper_server/internal/character"
)

type repository struct {
	db db.DatabaseTX
}

func NewRepository(db db.DatabaseTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateLobby(ctx context.Context, lobby *CreateLobbyReq) (*CreateLobbyRes, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
			return
		}
		err = tx.Commit()
	}()

	query := `INSERT INTO lobby(lobbyMasterId, lobbyName, 
                  lobbyPassword, amount) VALUES (?, ?, ?, ?)`
	res, err := tx.ExecContext(ctx, query, lobby.LobbyMasterID, lobby.LobbyName,
		lobby.LobbyPassword, lobby.Amount)
	if err != nil {
		return nil, err
	}

	var lobbyRes CreateLobbyRes
	lobbyRes.Id, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	query = `INSERT INTO acclobby(idAcc, idLobby) 
				VALUES (?, ?)`
	_, err = tx.ExecContext(ctx, query, &lobby.LobbyMasterID, lobbyRes.Id)
	if err != nil {
		return nil, err
	}

	return &lobbyRes, err
}

func (r *repository) GetAllLobby(ctx context.Context) ([]GetLobbyRes, error) {
	var lobbyList []GetLobbyRes

	query := `SELECT l.id, l.lobbyName, a.nickname FROM lobby l 
    LEFT JOIN acclobby ac on l.id = ac.idLobby 
    LEFT JOIN account a on ac.idAcc = a.id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var lobby GetLobbyRes
		err := rows.Scan(&lobby.Id, &lobby.LobbyName, &lobby.LobbyMaster)
		if err != nil {
			return nil, err
		}
		lobbyList = append(lobbyList, lobby)
	}

	return lobbyList, nil
}

func (r *repository) GetLobbyById(ctx context.Context, id int64) (*GetLobbyByIdRes, error) {
	var res GetLobbyByIdRes
	query := "SELECT lobbyPassword FROM lobby WHERE id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&res.Password)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *repository) JoinLobby(ctx context.Context, req *JoinLobbyReq) ([]character.Character, error) {
	query := "SELECT id FROM lobby WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, req.IdLobby)
	if err != nil {
		return nil, err
	}

	query = "INSERT INTO acclobby(idAcc, idLobby) VALUES (?, ?)"
	_, err = r.db.ExecContext(ctx, query, req.IdAcc, req.IdLobby)
	if err != nil {
		return nil, err
	}

	query = "SELECT c.hp, c.exp, c.charName, c.sex, c.weight, c.height, c.addLanguage, " +
		"c.ideals, c.weaknesses, c.traits, c.allies, c.organizations, c.enemies, c.story, " +
		"c.goals, c.treasures, c.notes, cl.className, r.raceName, s.subraceName, st.strength, " +
		"st.dexterity, st.constitution, st.intelligence, st.wisdom, st.charisma, " +
		"GROUP_CONCAT(sk.skillName SEPARATOR ', ') " +
		"AS characterSkills, a.alignmentName, i.image FROM characters c " +
		"JOIN accChar ac ON c.id = ac.idChar JOIN class cl ON c.idClass = cl.id " +
		"JOIN races r on c.idRace = r.id JOIN subrace s on c.idSubrace = s.id " +
		"JOIN stats st on c.idStats = st.id JOIN charSkills cs on c.id = cs.idChar " +
		"JOIN skills sk on cs.idSkill = sk.id JOIN alignment a on c.idAlignment = a.id " +
		"JOIN image i ON c.idAvatar = i.id LEFT JOIN acclobby al on ac.id = al.idAcc " +
		"WHERE al.idLobby = ? AND ac.act = true GROUP BY c.id"

	rows, err := r.db.QueryContext(ctx, query, req.IdLobby)
	if err != nil {
		return nil, err
	}

	var chars []character.Character
	for rows.Next() {
		var char character.Character
		err := rows.Scan(&char.Hp, &char.Exp, &char.CharName, &char.Sex, &char.Weight, &char.Height, &char.AddLanguage,
			&char.Ideals, &char.Weaknesses, &char.Traits, &char.Allies, &char.Organizations, &char.Enemies,
			&char.Story, &char.Goals, &char.Treasures, &char.Notes, &char.Class.ClassName, &char.Race.RaceName,
			&char.Subrace.SubraceName, &char.Stats.Strength, &char.Stats.Dexterity, &char.Stats.Constitution,
			&char.Stats.Intelligence, &char.Stats.Wisdom, &char.Stats.Charisma, &char.CharacterSkills,
			&char.Alignment.AlignmentName, &char.Avatar.Image)
		if err != nil {
			return nil, err
		}
		chars = append(chars, char)
	}

	return chars, nil
}
