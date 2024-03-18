CREATE TABLE `id_segment_tab`
(
    `id`           bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '',
    `segment_type` int(11) unsigned NOT NULL DEFAULT 0 COMMENT '',
    `max_id`       bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '目前该业务的最大id',
    `step`         int(10) unsigned NOT NULL DEFAULT 0 COMMENT '每次获取的步长',
    `ctime`        int(10) unsigned NOT NULL DEFAULT 0 COMMENT '',
    `mtime`        int(10) unsigned NOT NULL DEFAULT 0 COMMENT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_segment_type` (`segment_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='id_segment_tab';