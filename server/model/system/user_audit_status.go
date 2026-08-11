package system

// 用户审核状态（sys_users.audit_status）
const (
	AuditStatusNew      int8 = 0 // 新注册未填写信息
	AuditStatusPassed   int8 = 1 // 已通过
	AuditStatusPending  int8 = 2 // 已填写未审核
	AuditStatusChanging int8 = 3 // 修改信息待审核
	AuditStatusRejected int8 = 4 // 已拒绝
)
