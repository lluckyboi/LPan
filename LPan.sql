create table private
(
    user_id      int                     not null,
    file_name    varchar(32)             not null,
    file_id      int                     not null,
    deleted      datetime                null comment '删除的时间',
    father_path  varchar(32) default '/' not null,
    real_deleted tinyint     default 0   not null comment '真正删除，无法找回 该记录会在下一次清理时被删除'
);

create table public_file
(
    file_name varchar(32)  not null,
    file_id   bigint auto_increment
        primary key,
    hash      varchar(255) not null,
    size      bigint       not null
);

create table user
(
    user_name varchar(32) default 'LPanUser' not null,
    user_id   bigint auto_increment
        primary key,
    user_mail varchar(32)                    not null,
    vip       tinyint     default 0          not null
);

