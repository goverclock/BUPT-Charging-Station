drop table users;
drop table reports;

create table users (
  id         serial primary key,
  uuid       varchar(64) not null,
  name       varchar(255) not null unique,
  password   varchar(255) not null,
  isadmin    boolean,
  balance    float,
  batteryCapacity float
);

create table reports (
  id                      serial primary key,
  num                     bigint,
  charge_id               bigint,
  charge_mode             bigint,
  username                varchar(255),
  user_id                 bigint,
  request_charge_amount   float,
  real_charge_amount      float,
  charge_time             bigint,
  charge_fee              float,
  service_fee             float,
  tot_fee                 float,
  step                    bigint,
  queue_number            varchar(64),
  subtime                 bigint,
  inlinetime              bigint,
  calltime                bigint,
  charge_start_time       bigint,
  charge_end_time         bigint,
  terminate_flag          boolean,
  terminate_time          bigint,
  failed_flag             boolean,
  failed_msg              varchar(256)
)
