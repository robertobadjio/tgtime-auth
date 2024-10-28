create table "user"
(
    id                 bigserial
        primary key,
    name               varchar(255)                               not null,
    email              varchar(255)                               not null
        constraint users_email_unique
            unique,
    mac_address        macaddr                                    not null
        constraint users_mac_address_unique
            unique,
    telegram_id        varchar(255)
        constraint users_telegram_id_unique
            unique,
    email_verified_at  timestamp(0),
    password           varchar(255)                               not null,
    remember_token     varchar(100),
    device_mac_address varchar(17)
        constraint users_device_mac_address_unique
            unique,
    created_at         timestamp(0),
    updated_at         timestamp(0),
    deleted            boolean      default false                 not null,
    department_id      integer      default 0                     not null,
    role               varchar(255) default ''::character varying not null,
    surname            varchar(255) default ''::character varying not null,
    lastname           varchar(255) default ''::character varying not null,
    birth_date         date                                       not null,
    position           varchar(255) default ''::character varying
);