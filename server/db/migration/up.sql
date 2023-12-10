create table image
(
    id    int auto_increment
        primary key,
    image blob not null
);

create table account
(
    id       int auto_increment
        primary key,
    email    varchar(30) not null unique,
    password varchar(60) null,
    nickname varchar(16) not null,
    idAvatar int         not null,
    constraint account_image_id_fk
        foreign key (idAvatar) references image (id)
);

create table stats
(
    id           int auto_increment
        primary key,
    strength     int not null,
    dexterity    int not null,
    constitution int not null,
    intelligence int not null,
    wisdom       int not null,
    charisma     int not null
);

create table races
(
    id       int auto_increment
        primary key,
    raceName varchar(20) not null unique
);

create table subrace
(
    id          int auto_increment
        primary key,
    subraceName varchar(20) not null unique,
    idRace      int         not null,
    idStats     int         not null,
    constraint races_stats_id_fk
        foreign key (idStats) references stats (id),
    constraint races_race_id_fk
        foreign key (idRace) references races (id)
);

create table skills
(
    id        int auto_increment
        primary key,
    skillName varchar(20) not null unique,
    idStats   int         not null,
    constraint skills_stats_id_fk
        foreign key (idStats) references stats (id)
);

create table class
(
    id        int auto_increment
        primary key,
    className varchar(20) not null unique
);

create table alignment
(
    id            int auto_increment
        primary key,
    alignmentName varchar(20) not null unique
);

create table characters
(
    id            int auto_increment
        primary key,
    hp            int         not null,
    exp           int         not null,
    idAvatar      int         not null,
    charName      varchar(16) not null,
    sex           boolean     not null,
    weight        int,
    height        int,
    idClass       int         not null,
    idRace        int         not null,
    idSubrace     int         not null,
    idStats       int         not null,
    addLanguage   varchar(20),
    idAlignment   int,
    ideals        varchar(255),
    weaknesses    varchar(255),
    traits        varchar(255),
    allies        varchar(255),
    organizations varchar(255),
    enemies       varchar(255),
    story         varchar(1000),
    goals         varchar(1000),
    treasures     varchar(255),
    notes         varchar(1000),
    constraint characters_images_id_fk
        foreign key (idAvatar) references image (id)
            ON DELETE CASCADE,
    constraint characters_class_id_fk
        foreign key (idClass) references class (id),
    constraint characters_race_id_fk
        foreign key (idRace) references races (id),
    constraint characters_subrace_id_fk
        foreign key (idSubrace) references subrace (id),
    constraint characters_stats_id_fk
        foreign key (idStats) references stats (id)
            ON DELETE CASCADE,
    constraint characters_alignment_id_fk
        foreign key (idAlignment) references alignment (id)
);

create table charSkills
(
    id      int auto_increment
        primary key,
    idSkill int not null,
    idChar  int not null,
    constraint charSkills_char_id_fk
        foreign key (idChar) references characters (id),
    constraint charSkills_skills_id_fk
        foreign key (idSkill) references skills (id),
    constraint unique_char_skill
        unique (idChar, idSkill)
);

create table accChar
(
    id        int auto_increment
        primary key,
    act       boolean not null,
    idAccount int     not null,
    idChar    int     not null unique,
    constraint accChar_account_id_fk
        foreign key (idAccount) references account (id),
    constraint accChar_characters_id_fk
        foreign key (idChar) references characters (id)
);

CREATE TRIGGER unique_act_true
    BEFORE INSERT
    ON accChar
    FOR EACH ROW
BEGIN
    DECLARE count_act_true INT;
    SELECT COUNT(*) INTO count_act_true FROM accChar WHERE idAccount = NEW.idAccount AND act = true;
    IF count_act_true > 0 AND NEW.act = true THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Only one record with act = true allowed per account';
    END IF;
END;

create table lobby
(
    id            int auto_increment
        primary key,
    lobbyMasterId int,
    lobbyName     varchar(16) unique,
    lobbyPassword varchar(60),
    amount        int
);

create table accLobby
(
    id      int auto_increment
        primary key,
    idAcc   int unique,
    idLobby int,
    constraint accLobby_lobby_id_fk
        foreign key (idLobby) references lobby (id)
            ON DELETE CASCADE,
    constraint accLobby_account_id_fk
        foreign key (idAcc) references account (id)
);

INSERT INTO image(image)
    VALUE (1);
INSERT INTO image(image)
    VALUE (2);
INSERT INTO image(image)
    VALUE (3);
INSERT INTO image(image)
    VALUE (4);
INSERT INTO image(image)
    VALUE (5);

INSERT INTO account(email, password, nickname, idAvatar)
VALUES (1, 1, 1, 1);
INSERT INTO account(email, password, nickname, idAvatar)
VALUES (2, 2, 2, 2);

INSERT INTO stats(strength, dexterity, constitution, intelligence, wisdom, charisma)
VALUES (1, 1, 1, 1, 1, 1);
INSERT INTO stats(strength, dexterity, constitution, intelligence, wisdom, charisma)
VALUES (2, 2, 2, 2, 2, 2);
INSERT INTO stats(strength, dexterity, constitution, intelligence, wisdom, charisma)
VALUES (3, 3, 3, 3, 3, 3);
INSERT INTO stats(strength, dexterity, constitution, intelligence, wisdom, charisma)
VALUES (4, 4, 4, 4, 4, 4);
INSERT INTO stats(strength, dexterity, constitution, intelligence, wisdom, charisma)
VALUES (5, 5, 5, 5, 5, 5);

INSERT INTO races(raceName)
VALUES ('Эльф');
INSERT INTO races(raceName)
VALUES ('Дворф');

INSERT INTO subrace(subraceName, idStats, idRace)
VALUES ('Полу', 1, 1);
INSERT INTO subrace(subraceName, idStats, idRace)
VALUES ('Фул', 2, 2);

INSERT INTO skills(skillName, idStats)
VALUES ('Анализ', 1);
INSERT INTO skills(skillName, idStats)
VALUES ('Акробатика', 2);

INSERT INTO class(className)
VALUES ('Воин');
INSERT INTO class(className)
VALUES ('Бард');

INSERT INTO alignment(alignmentName)
    VALUE ('Законопослушный');
INSERT INTO alignment(alignmentName)
    VALUE ('Нейтральный');

INSERT INTO characters (hp, exp, idAvatar, charName, sex, weight, height, idClass, idRace, idSubrace, idStats,
                        addLanguage, idAlignment, ideals, weaknesses, traits, allies, organizations, enemies,
                        story, goals, treasures, notes)
    VALUE (100, 0, 3, 'Влад', true, 100, 170, 1, 1, 1, 3, 'Японский', 1,
           'Нет', 'Нет', 'Нет', 'Нет', 'Нет', 'Нет', 'Нет', 'Нет', 'Нет', 'Нет');

INSERT INTO characters (hp, exp, idAvatar, charName, sex, weight, height, idClass, idRace, idSubrace, idStats,
                        addLanguage, idAlignment, ideals, weaknesses, traits, allies, organizations, enemies,
                        story, goals, treasures, notes)
    VALUE (200, 0, 4, 'Игорь', false, 50, 180, 2, 2, 2, 4, 'Японский', 2,
           'Да', 'Да', 'Да', 'Да', 'Да', 'Да', 'Да', 'Да', 'Да', 'Да');

INSERT INTO charSkills (idSkill, idChar)
    VALUE (1, 1);
INSERT INTO charSkills (idSkill, idChar)
    VALUE (2, 1);
INSERT INTO charSkills (idSkill, idChar)
    VALUE (1, 2);

INSERT INTO accChar (act, idAccount, idChar)
    VALUE (true, 1, 1);
INSERT INTO accChar (act, idAccount, idChar)
    VALUE (true, 2, 2);