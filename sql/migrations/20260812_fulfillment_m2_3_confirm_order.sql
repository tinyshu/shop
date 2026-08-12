-- M2-3 / FUL-02: register confirmOrder API + casbin for roles that already have cancelOrder
-- Idempotent by path+method / ptype+v0+v1+v2.

INSERT INTO `sys_apis` (`created_at`, `updated_at`, `deleted_at`, `path`, `description`, `api_group`, `method`)
SELECT NOW(3), NOW(3), NULL, '/order/confirmOrder', '确认收货', 'order', 'POST'
FROM DUAL
WHERE NOT EXISTS (
  SELECT 1 FROM `sys_apis`
  WHERE `path` = '/order/confirmOrder' AND `method` = 'POST' AND `deleted_at` IS NULL
);

INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`)
SELECT 'p', r.`v0`, '/order/confirmOrder', 'POST', '', '', ''
FROM (
  SELECT DISTINCT `v0` FROM `casbin_rule`
  WHERE `ptype` = 'p' AND `v1` = '/order/cancelOrder' AND `v2` = 'POST'
) r
WHERE NOT EXISTS (
  SELECT 1 FROM `casbin_rule` c
  WHERE c.`ptype` = 'p' AND c.`v0` = r.`v0`
    AND c.`v1` = '/order/confirmOrder' AND c.`v2` = 'POST'
);
