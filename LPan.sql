create table user
(
    user_name varchar(32) default 'LPanUser' not null,
    user_id   bigint auto_increment
        primary key,
    user_mail varchar(32)                    not null
);

