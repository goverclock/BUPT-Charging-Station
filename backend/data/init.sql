drop table sessions;
drop table users;
drop table stations;
drop table cars;

create table users (
  id         serial primary key,
  name       varchar(255) not null unique,
  password   varchar(255) not null,
  balance    float,
  batteryCapacity float
);

create table sessions (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  name      varchar(255),
  user_id    integer references users(id)
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

-- 2 Fast, 3 Slow
INSERT INTO stations (mode) VALUES ('Fast');
INSERT INTO stations (mode) VALUES ('Fast');
INSERT INTO stations (mode) VALUES ('Slow');
INSERT INTO stations (mode) VALUES ('Slow');
INSERT INTO stations (mode) VALUES ('Slow');

-- examples

-- user must be created with Create, because data is Encrypted
-- INSERT INTO users (name, password, balance, batteryCapacity) VALUES ('an', 'ap', 100, 233);


