-- M0-1 industry-config: expand sys_config.name + seed feature flags
-- Safe to re-run: inserts are idempotent by name.

ALTER TABLE `sys_config`
  MODIFY COLUMN `name` varchar(64) NOT NULL COMMENT '配置参数的 key';

INSERT INTO `sys_config` (`created_at`, `updated_at`, `deleted_at`, `name`, `value`, `group_type`, `desc`, `status`)
SELECT NOW(3), NOW(3), NULL, 'feature.user_audit', '0', 'feature', '用户审核开关(1开/0关，默认关=B2C)', 1
FROM DUAL
WHERE NOT EXISTS (
  SELECT 1 FROM `sys_config` WHERE `name` = 'feature.user_audit' AND `deleted_at` IS NULL
);

INSERT INTO `sys_config` (`created_at`, `updated_at`, `deleted_at`, `name`, `value`, `group_type`, `desc`, `status`)
SELECT NOW(3), NOW(3), NULL, 'feature.settle_month', '0', 'feature', '月结能力开关(1开/0关，默认关)', 1
FROM DUAL
WHERE NOT EXISTS (
  SELECT 1 FROM `sys_config` WHERE `name` = 'feature.settle_month' AND `deleted_at` IS NULL
);

INSERT INTO `sys_config` (`created_at`, `updated_at`, `deleted_at`, `name`, `value`, `group_type`, `desc`, `status`)
SELECT NOW(3), NOW(3), NULL, 'feature.courier_mode', 'courier', 'feature', '物流模式 delivery=城配 courier=快递(默认)', 1
FROM DUAL
WHERE NOT EXISTS (
  SELECT 1 FROM `sys_config` WHERE `name` = 'feature.courier_mode' AND `deleted_at` IS NULL
);
