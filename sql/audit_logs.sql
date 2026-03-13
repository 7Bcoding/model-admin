create table llm_admin.audit_logs
(
    id           bigint auto_increment
        primary key,
    operator     varchar(255)                        not null comment '''操作者''',
    request_url  text                                not null comment '''请求URL''',
    method       varchar(10)                         not null comment '''HTTP方法''',
    action       varchar(255)                        not null comment '''操作类型''',
    target       varchar(255)                        null,
    detail       text                                null,
    result       text                                not null comment '''操作结果''',
    status       bigint                              not null comment '''HTTP状态码''',
    request_body text                                null comment '''请求体''',
    created_at   timestamp default CURRENT_TIMESTAMP not null comment '创建时间'
)
    comment '审计日志表' collate = utf8mb4_unicode_ci;

create index idx_action
    on llm_admin.audit_logs (action);

create index idx_created_at
    on llm_admin.audit_logs (created_at);

create index idx_operator
    on llm_admin.audit_logs (operator);




