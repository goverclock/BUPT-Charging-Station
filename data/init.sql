drop table sessions;
drop table users;
drop table stations;
drop table cars;

create table users (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  email      varchar(255) not null unique,
  password   varchar(255) not null,
  created_at timestamp not null   
);

create table sessions (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  email      varchar(255),
  user_id    integer references users(id),
  created_at timestamp not null   
);

create table stations (
  id        serial primary key,
  mode      varchar(64) not null,
  usedby    varchar(64),
  slot1     varchar(64),
  slot2     varchar(64)
  -- ok        boolean  -- is the station functional?
);

create table cars (
  id        serial primary key,
  ownedby   varchar(64) not null,
  stage     varchar(64) not null,
  qid       varchar(64)
);
