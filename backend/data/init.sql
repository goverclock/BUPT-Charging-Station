drop table users;

create table users (
  id         serial primary key,
  uuid       varchar(64) not null,
  name       varchar(255) not null unique,
  password   varchar(255) not null,
  isadmin    boolean,
  balance    float,
  batteryCapacity float
);