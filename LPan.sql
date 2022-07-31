create table private
(
    user_id     int                     not null,
    file_name   varchar(32)             not null,
    file_id     int                     not null,
    deleted     datetime                null comment '删除的时间',
    father_path varchar(32) default '/' not null,
    share       tinyint     default 0   not null comment '表示这条记录是否为分享的',
    expr_time   time                    null comment '分享文件过期时间'
);

create table public_file
(
    file_name varchar(32)  not null,
    file_id   bigint auto_increment
        primary key,
    hash      varchar(255) not null,
    size      bigint       not null
);

create table url
(
    origin varchar(255) not null,
    sha1   varchar(255) not null,
    constraint url_sha1_uindex
        unique (sha1)
);

create table user
(
    user_name varchar(32) default 'LPanUser' not null,
    user_id   bigint auto_increment
        primary key,
    user_mail varchar(32)                    not null,
    vip       tinyint     default 0          not null
);

