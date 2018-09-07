CREATE TABLE `com_table_mapping` (
  `id` char(32) NOT NULL COMMENT '主键',
  `aliasName` varchar(50) NOT NULL DEFAULT '' COMMENT '数据库别名',
  `log_table_name` varchar(50) NOT NULL DEFAULT '' COMMENT '表名称',
  `field_name` varchar(50) NOT NULL DEFAULT '' COMMENT '字段名称',
  `field_type` varchar(16) NOT NULL DEFAULT '' COMMENT '字段类型',
  `field_title` varchar(50) NOT NULL DEFAULT '' COMMENT '字段标题',
  `field_sort` tinyint(4) NOT NULL DEFAULT '0' COMMENT '字段排序',
  `order_by` varchar(4) DEFAULT NULL COMMENT '排序（ASC:升序,DESC:降序)',
  `is_show` char(1) NOT NULL DEFAULT '' COMMENT '是否显示',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，0正常 1禁用',
  `is_primary` tinyint(4) unsigned zerofill NOT NULL COMMENT '是否主键【0:否，1:是】',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tab_field` (`aliasName`,`log_table_name`,`field_name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;




CREATE TABLE `com_biz_log` (
  `log_id` varchar(32) NOT NULL COMMENT '日志表id，uuid',
  `user_id` varchar(32) DEFAULT NULL COMMENT '用户id,记录操作用户',
  `module_name` varchar(225) DEFAULT NULL COMMENT '模块名称',
  `create_time` datetime DEFAULT NULL COMMENT '操作时间',
  `class_name` varchar(225) DEFAULT NULL COMMENT '类名称',
  `method_name` varchar(225) DEFAULT NULL COMMENT '方法名称',
  `params` longtext COMMENT '传入参数',
  `ip` varchar(225) DEFAULT NULL COMMENT '操作ip',
  `commemts` longtext COMMENT '备注',
  `status` int(11) DEFAULT NULL COMMENT '状态',
  PRIMARY KEY (`log_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='日志表';



CREATE TABLE `com_user` (
  `id` char(32) NOT NULL,
  `user_name` varchar(20) NOT NULL DEFAULT '' COMMENT '用户名',
  `email` varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱',
  `password` char(32) NOT NULL DEFAULT '' COMMENT '密码',
  `salt` char(10) NOT NULL DEFAULT '' COMMENT '密码盐',
  `last_login` int(11) NOT NULL DEFAULT '0' COMMENT '最后登录时间',
  `last_ip` char(15) NOT NULL DEFAULT '' COMMENT '最后登录IP',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，0正常 -1禁用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_name` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



