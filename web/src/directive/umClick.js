/**
 * create by Administrator on 2023/11/17 14:03
 * 类注释：v-um-click=['','','','']
 * 备注：
 */
import {trackEvent} from '@/utils/UmUtils'

export default {
    bind(el, binding) {
        el.addEventListener('click', () => {
            // window._czc && window._czc.push(['_trackEvent', ...binding.value])
            trackEvent(...binding.value)
        }, false)
    },
    unbind(el) {
        el.removeEventListener('click', () => {
            console.log('remove')
        })
    }
}
