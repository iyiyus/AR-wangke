-- 盖章订单表
CREATE TABLE IF NOT EXISTS `qingka_wangke_stamp` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL COMMENT '用户ID',
  `user` varchar(100) NOT NULL COMMENT '用户名',
  `company` varchar(200) NOT NULL COMMENT '公司名称',
  `file` varchar(500) NOT NULL COMMENT '文件路径',
  `color` varchar(20) NOT NULL COMMENT '颜色：red/blue/black',
  `size` varchar(50) NOT NULL COMMENT '规格',
  `side` varchar(10) NOT NULL COMMENT '单双面：single/double',
  `copies` int(11) NOT NULL COMMENT '份数',
  `note` text COMMENT '备注',
  `status` varchar(20) NOT NULL DEFAULT '待处理' COMMENT '状态：待处理/进行中/已完成/已取消',
  `addtime` datetime NOT NULL COMMENT '添加时间',
  `complete` datetime DEFAULT NULL COMMENT '完成时间',
  PRIMARY KEY (`id`),
  KEY `idx_uid` (`uid`),
  KEY `idx_status` (`status`),
  KEY `idx_addtime` (`addtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='盖章订单表';

-- 病历订单表
CREATE TABLE IF NOT EXISTS `qingka_wangke_medical` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL COMMENT '用户ID',
  `user` varchar(100) NOT NULL COMMENT '用户名',
  `hospital` varchar(200) NOT NULL COMMENT '医院名称',
  `department` varchar(100) NOT NULL COMMENT '科室',
  `type` varchar(50) NOT NULL COMMENT '病历类型',
  `files` text NOT NULL COMMENT '文件路径（多个用逗号分隔）',
  `note` text COMMENT '备注',
  `status` varchar(20) NOT NULL DEFAULT '待处理' COMMENT '状态：待处理/进行中/已完成/已取消',
  `addtime` datetime NOT NULL COMMENT '添加时间',
  `complete` datetime DEFAULT NULL COMMENT '完成时间',
  PRIMARY KEY (`id`),
  KEY `idx_uid` (`uid`),
  KEY `idx_status` (`status`),
  KEY `idx_addtime` (`addtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='病历订单表';