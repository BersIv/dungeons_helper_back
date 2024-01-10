drop table dnd.accchar;

drop table dnd.acclobby;

drop table dnd.account;

drop table dnd.charskills;

drop table dnd.characters;

drop table dnd.image;

drop table dnd.subrace;

drop table dnd.races;

drop table dnd.skills;

drop table dnd.stats;

drop table dnd.class;

drop table dnd.alignment;

drop table dnd.lobby;

SELECT l.id, l.lobbyName, count(ac.idAcc) FROM lobby l
                                                   LEFT JOIN acclobby ac on l.id = ac.idLobby
                                                   LEFT JOIN account a on ac.idAcc = a.id
GROUP BY l.id, l.lobbyName