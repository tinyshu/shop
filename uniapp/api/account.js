import request from "@/utils/request"

// 获取用户账户信息
export const getAccountInfo = (groupId, refs) => {
    const data = {
        ID: groupId,
    }
    return request({
        url: `/account/findUserAccount`,
        method: 'GET',
        loading: true,
        data
    },refs)
}

