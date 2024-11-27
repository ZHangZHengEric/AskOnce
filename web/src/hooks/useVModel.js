/**
 * create by zhangxiang on 2023-07-11 13:36
 * 类注释：
 * 备注：
 */

import {computed} from 'vue'

/**
 *  v-model 绑定对象
 *  const emit = defineEmits(['update:' + propName])
 * @param props
 * @param propName
 * @param emit
 * @returns {WritableComputedRef<*>|boolean|*}
 */
export const useVModel = (props, propName, emit) => {
    return computed({
        get() {
            return new Proxy(props[propName], {
                set(obj, name, val) {
                    emit('update:' + propName, {
                        ...obj,
                        [name]: val
                    })
                    return true
                }
            })
        },
        set(val) {
            emit('update:' + propName, val)
        }
    })
}