CREATE TABLE if not exists `xorm_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键自增id',
  `uid` int(11) NOT NULL COMMENT '用户id',
  `keywords` varchar(255) COMMENT '索引词',
  `degree` char(2) NOT NULL COMMENT '学历',
  `gender` char(1) NOT NULL COMMENT '性别',
  `city` char(2) NOT NULL COMMENT '城市',
  create_time datetime NOT NULL comment '创建时间',
  update_time datetime NOT NULL comment '最后修改时间',
  delete_time datetime comment '软删除时间',
  version int not null default 0 comment '每次更新时版本加1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户信息表';