CREATE TABLE users
(
    id bigserial not null unique,
    userid int not null unique,
    fname varchar(250),
    lname varchar(250),
    username varchar(250) not null unique,
    languages varchar(250),
    chatid bigint not null 
);

CREATE TABLE messages
(
    id serial not null unique,
    mowner int references users(userid) on delete cascade,
    mdate varchar(200) not null,
	inputdata varchar not null,
    chatid bigint not null,
    chattitle varchar(250) DEFAULT 'private_mess'
);

CREATE TABLE admins
(
    userid int not null unique,
    username varchar(250) not null unique,
    masteradmin boolean
);