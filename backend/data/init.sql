drop table sessions;
drop table users;

create table users (
  id         serial primary key,
  uuid       varchar(64) not null,
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

-- create table cars (
--   id        serial primary key,
--   ownedby   varchar(64) not null,
--   stage     varchar(64) not null,
--   qid       varchar(64)
-- );

-- -- examples

-- -- user must be created with Create, because data is Encrypted
-- -- INSERT INTO users (name, password, balance, batteryCapacity) VALUES ('an', 'ap', 100, 233);

-- INSERT INTO cars (ownedby, stage, qid) VALUES ('g', 'Waiting', 'F1');

