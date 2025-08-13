CREATE TABLE `user` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键自增id',
  `uid` int NOT NULL COMMENT '用户id',
  `keywords` varchar(255) DEFAULT NULL COMMENT '索引词',
  `degree` char(2) NOT NULL COMMENT '学历',
  `gender` char(1) NOT NULL COMMENT '性别',
  `city` char(2) NOT NULL COMMENT '城市',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uid` (`uid`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户信息表';

CREATE TABLE `login` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` longtext,
  `password` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='登录信息表';