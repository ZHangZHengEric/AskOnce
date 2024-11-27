/**
 * create by zhangxiang on 2023-06-10 12:45
 * 类注释：
 * 备注：
 */

import {createPinia} from 'pinia'

export const store = createPinia()

export function setupStore(app) {
    app.use(store)
}

export * from './modules'
