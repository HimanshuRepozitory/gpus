drop table if exists gpus;

create table gpus(
      id serial primary key,
      username varchar primary key not null,
      comment varchar not null
);
