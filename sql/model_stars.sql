create table llm_admin.model_stars
(
    id         int auto_increment
        primary key,
    model_name varchar(200) not null,
    user_id    int          not null,
    created_at datetime     null,
    constraint unique_user_model_star
        unique (model_name, user_id),
    constraint model_stars_ibfk_1
        foreign key (user_id) references llm_admin.users (id)
);

create index user_id
    on llm_admin.model_stars (user_id);

