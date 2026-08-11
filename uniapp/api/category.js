import request from "@/utils/request"

export const getCategoryListAll = (data) => {
    return request({
    url: '/category/getCategoryListAll',
    method: 'GET',
    data
  })
}

export const getHomeCategoryList = () => {
    return request({
        url: '/category/getCategoryList',
        method: 'GET',
        data: {
            page: 1,
            pageSize: 4,
            isFirst: 1
        }
    })
}
