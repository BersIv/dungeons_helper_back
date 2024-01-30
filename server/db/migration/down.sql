drop table dnd.accChar;

drop table dnd.accLobby;

drop table dnd.account;

drop table dnd.charSkills;

drop table dnd.characters;

drop table dnd.image;

drop table dnd.subrace;

drop table dnd.races;

drop table dnd.skills;

drop table dnd.stats;

drop table dnd.class;

drop table dnd.alignment;

drop table dnd.lobby;

SELECT a.id, email, password, nickname, image FROM account a 
LEFT JOIN image i ON a.idAvatar = i.id
WHERE email = 1;

SELECT a.id, email, password, nickname, i.image FROM account a 
LEFT JOIN image i ON a.idAvatar = i.id 
WHERE email = 1;

SELECT a.id, email, password, nickname, i.image FROM account a 
				LEFT JOIN image i ON a.idAvatar = i.id 
				WHERE email = 1;

use dnd;
SELECT id, subraceName, idStats FROM subrace WHERE idRace = 1;