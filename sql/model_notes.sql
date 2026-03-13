create table llm_admin.model_notes
(
    id               int auto_increment
        primary key,
    model_name       varchar(200)  not null,
    user_id          int           null,
    note             varchar(1000) null,
    created_at       datetime      null,
    updated_at       datetime      null,
    open_chat_id     varchar(500)  not null,
    inference_engine varchar(255)  not null,
    model_metric     varchar(255)  not null,
    constraint model_notes_model_name_uindex
        unique (model_name)
);

