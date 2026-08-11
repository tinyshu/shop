import request from "@/utils/request"

export const getWeChatOpenIdByCode = (data) => {
  return request({
    url: '/wechat/code2Session',
    data,
  })
}

export const wxLogin = (data) => {
  return request({
    url: '/base/loginWx',
    method: 'POST',
    loading: true,
    data,
  })
}
