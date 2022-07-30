create table private
(
    user_id     int                     not null,
    file_name   varchar(32)             not null,
    file_id     int                     not null,
    deleted     datetime                null comment '删除的时间',
    father_path varchar(32) default '/' not null
);

create table public_file
(
    file_name varchar(32) not null,
    file_id   bigint auto_increment
        primary key
);

create table user
(
    user_name varchar(32) default 'LPanUser' not null,
    user_id   bigint auto_increment
        primary key,
    user_mail varchar(32)                    not null
);

