import request from "@/utils/request"

export const getBrandListAll = (data) => {
    return request({
        url: '/brand/getBrandListAll',
        method: 'GET',
        data
    })
}

export const getBrandListByCategoryId = (data) => {
    return request({
        url: '/brand/getBrandListByCategoryId',
        method: 'GET',
        data
    })
}
