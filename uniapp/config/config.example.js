// 复制本文件为 config.js 后按环境修改（config.js 可被本地覆盖；仓库保留可运行的默认 localhost）
export default {
    // #ifdef H5
    // baseUrl: '/api',
    // #endif
    // #ifndef H5
    baseUrl: 'http://localhost:48888',
    // 生产：baseUrl: 'https://your-domain.com/goapi',
    // #endif
    phone: '' // 演示用联系电话，可删
}
