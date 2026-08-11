import { getUser, setUser } from '@/store/storage'
import { getUserAuditStatus } from '@/api/user'

/**
 * 刷新审核通过标记（信任服务端生效态；免审时 getAuditStatus 返回 1）。
 * @param {object} vm Vue 实例，写入 vm.isAudit
 * @returns {Promise<boolean>}
 */
export function refreshUserAudit(vm) {
	const user = getUser()
	if (!user) {
		if (vm) vm.isAudit = false
		return Promise.resolve(false)
	}
	if (user.auditStatus === 1) {
		if (vm) vm.isAudit = true
		return Promise.resolve(true)
	}
	return getUserAuditStatus().then(res => {
		const ok = !!(res.data && res.data.auditStatus === 1)
		if (ok) {
			user.auditStatus = 1
			setUser(user)
		}
		if (vm) vm.isAudit = ok
		return ok
	}).catch(() => {
		if (vm) vm.isAudit = false
		return false
	})
}
