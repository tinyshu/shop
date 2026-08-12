-- M0-3 industry-config: admin menu「功能开关」
-- Safe to re-run: menu/authority inserts are idempotent by menu name.

INSERT INTO `sys_base_menus` (
  `created_at`, `updated_at`, `deleted_at`,
  `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`,
  `active_name`, `keep_alive`, `default_menu`, `title`, `icon`, `close_tab`
)
SELECT
  NOW(3), NOW(3), NULL,
  0, '3', 'featureConfig', 'featureConfig', 0,
  'view/superAdmin/featureConfig/featureConfig.vue', 1,
  '', 0, 0, '功能开关', 'switch', 0
FROM DUAL
WHERE NOT EXISTS (
  SELECT 1 FROM `sys_base_menus`
  WHERE `name` = 'featureConfig' AND `deleted_at` IS NULL
);

INSERT INTO `sys_authority_menus` (`sys_base_menu_id`, `sys_authority_authority_id`)
SELECT m.`id`, 888
FROM `sys_base_menus` m
WHERE m.`name` = 'featureConfig'
  AND m.`deleted_at` IS NULL
  AND NOT EXISTS (
    SELECT 1 FROM `sys_authority_menus` am
    WHERE am.`sys_base_menu_id` = m.`id`
      AND am.`sys_authority_authority_id` = 888
  );
