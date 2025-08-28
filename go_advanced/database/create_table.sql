use test;

create table if not exists distributed_lock(
    id int auto_increment comment '自增id',
    lock_name char(30) not null comment '锁的名称',
    update_time datetime default current_timestamp on update current_timestamp comment '最后修改时间',
    primary key (id),
    unique key idx_lock (lock_name)
)default charset=utf8mb4 comment '分布式锁';