package character

import (
	"context"
	"database/sql"
	"dungeons_helper_server/db"
)

type repository struct {
	db db.DatabaseTX
}

func NewRepository(db db.DatabaseTX) Repository {
	return &repository{db: db}
}

func (r *repository) GetAllCharactersByAccId(ctx context.Context, idAcc int64) ([]Character, error) {
	var chars []Character

	query := "SELECT c.hp, c.exp, c.charName, c.sex, c.weight, c.height, c.addLanguage, " +
		"c.ideals, c.weaknesses, c.traits, c.allies, c.organizations, c.enemies, c.story, " +
		"c.goals, c.treasures, c.notes, cl.className, r.raceName, s.subraceName, st.strength, " +
		"st.dexterity, st.constitution, st.intelligence, st.wisdom, st.charisma, " +
		"GROUP_CONCAT(sk.skillName SEPARATOR ', ') " +
		"AS characterSkills, a.alignmentName, i.image FROM characters c " +
		"JOIN accChar ac ON c.id = ac.idChar JOIN class cl ON c.idClass = cl.id " +
		"JOIN races r on c.idRace = r.id JOIN subrace s on c.idSubrace = s.id " +
		"JOIN stats st on c.idStats = st.id JOIN charSkills cs on c.id = cs.idChar " +
		"JOIN skills sk on cs.idSkill = sk.id JOIN alignment a on c.idAlignment = a.id " +
		"JOIN image i ON c.idAvatar = i.id WHERE ac.idAccount = ? GROUP BY c.id"

	rows, err := r.db.QueryContext(ctx, query, idAcc)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var char Character
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

func (r *repository) GetCharacterById(ctx context.Context, id int64) (*Character, error) {
	char := Character{}
	var imageBytes []byte

	query := "SELECT c.hp, c.exp, c.charName, c.sex, c.weight, c.height, c.addLanguage, " +
		"c.ideals, c.weaknesses, c.traits, c.allies, c.organizations, c.enemies, c.story, " +
		"c.goals, c.treasures, c.notes, cl.className, r.raceName, s.subraceName, st.strength, " +
		"st.dexterity, st.constitution, st.intelligence, st.wisdom, st.charisma, " +
		"GROUP_CONCAT(sk.skillName SEPARATOR ', ') " +
		"AS characterSkills, a.alignmentName, i.image FROM characters c " +
		"JOIN accChar ac ON c.id = ac.idChar JOIN class cl ON c.idClass = cl.id " +
		"JOIN races r on c.idRace = r.id JOIN subrace s on c.idSubrace = s.id " +
		"JOIN stats st on c.idStats = st.id JOIN charSkills cs on c.id = cs.idChar " +
		"JOIN skills sk on cs.idSkill = sk.id JOIN alignment a on c.idAlignment = a.id " +
		"JOIN image i ON c.idAvatar = i.id WHERE c.id = ? GROUP BY c.id"

	err := r.db.QueryRowContext(ctx, query, id).Scan(&char.Hp, &char.Exp, &char.CharName, &char.Sex, &char.Weight, &char.Height, &char.AddLanguage,
		&char.Ideals, &char.Weaknesses, &char.Traits, &char.Allies, &char.Organizations, &char.Enemies,
		&char.Story, &char.Goals, &char.Treasures, &char.Notes, &char.Class.ClassName, &char.Race.RaceName,
		&char.Subrace.SubraceName, &char.Stats.Strength, &char.Stats.Dexterity, &char.Stats.Constitution,
		&char.Stats.Intelligence, &char.Stats.Wisdom, &char.Stats.Charisma, &char.CharacterSkills,
		&char.Alignment.AlignmentName, &imageBytes)
	if err != nil {
		return nil, err
	}

	return &char, nil
}

func (r *repository) CreateCharacter(ctx context.Context, char *CreateCharacterReq, idAcc int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
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

	query := `INSERT INTO image(image) VALUE (?)`
	result, err := tx.ExecContext(ctx, query, char.Avatar.Image)
	if err != nil {
		return err
	}
	imageId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	query = `INSERT INTO stats(strength, dexterity, constitution, intelligence, wisdom, charisma) VALUES (?, ?, ?, ?, ?, ?)`
	result, err = tx.ExecContext(ctx, query, char.Stats.Strength, char.Stats.Dexterity, char.Stats.Constitution, char.Stats.Intelligence, char.Stats.Wisdom, char.Stats.Charisma)
	if err != nil {
		return err
	}
	statsId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	query = `INSERT INTO characters(hp, exp, idAvatar, charName, sex, weight, height, idClass, 
                       idRace, idSubrace, idStats, addLanguage, idAlignment, ideals, weaknesses, 
                       traits, allies, organizations, enemies, story, goals, treasures, notes)  
					   VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err = tx.ExecContext(ctx, query, char.Hp, char.Exp, imageId, char.CharName, char.Sex, char.Weight, char.Height,
		char.Class.Id, char.Race.Id, char.Subrace.Id, statsId, char.AddLanguage, char.Alignment.Id, char.Ideals,
		char.Weaknesses, char.Traits, char.Allies, char.Organizations, char.Enemies, char.Story, char.Goals, char.Treasures, char.Notes)
	if err != nil {
		return err
	}
	charId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	query = "INSERT INTO charskills(idSkill, idChar) VALUES (?, ?)"
	for _, skill := range char.CharacterSkills {
		_, err = tx.ExecContext(ctx, query, skill.Id, charId)
		if err != nil {
			return err
		}
	}

	query = `UPDATE accchar SET act = false WHERE idAccount = ?`
	_, err = tx.ExecContext(ctx, query, idAcc)
	if err != nil {
		return err
	}

	query = `INSERT INTO accchar(act, idAccount, idChar) VALUES (?, ?, ?)`
	_, err = tx.ExecContext(ctx, query, true, idAcc, charId)
	if err != nil {
		return err
	}

	return nil
}
