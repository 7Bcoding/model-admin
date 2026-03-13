create table llm_admin.users
(
    id            int auto_increment
        primary key,
    username      varchar(255)               not null,
    account_name  varchar(255)               null,
    password_hash varchar(255)               not null,
    role          varchar(50) default 'user' not null,
    created_at    datetime(3)                null,
    updated_at    datetime(3)                null,
    constraint uni_users_username
        unique (username)
);

