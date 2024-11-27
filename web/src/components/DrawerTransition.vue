/**
* create by zhangxiang on 2024/7/1  14:47
* 类注释：
* 备注：
*/
<template>
  <el-drawer-transition>
    <slot/>
  </el-drawer-transition>
</template>

<script>
import {
  addClass,
  removeClass
} from 'element-plus/lib/utils/dom'

export default {
  name: 'DrawerTransition',
  components: {
    'el-drawer-transition': {
      functional: true,
      render(createElement, context) {
        const data = {
          on: {
            beforeEnter(el) {
              addClass(el, 'drawer-transition')
              if (!el.dataset) el.dataset = {}

              el.dataset.oldPaddingLeft = el.style.paddingLeft
              el.dataset.oldPaddingRight = el.style.paddingRight

              el.style.width = '0'
              el.style.paddingLeft = 0
              el.style.paddingRight = 0
            },
            enter(el) {
              el.dataset.oldOverflow = el.style.overflow
              if (el.scrollWidth !== 0) {
                el.style.width = el.scrollWidth + 'px'
                el.style.paddingLeft = el.dataset.oldPaddingLeft
                el.style.paddingRight = el.dataset.oldPaddingRight
              } else {
                el.style.width = ''
                el.style.paddingLeft = el.dataset.oldPaddingLeft
                el.style.paddingRight = el.dataset.oldPaddingRight
              }
              el.style.overflow = 'hidden'
            },
            afterEnter(el) {
              removeClass(el, 'drawer-transition')
              el.style.width = ''
              el.style.overflow = el.dataset.oldOverflow
            },
            beforeLeave(el) {
              if (!el.dataset) el.dataset = {}
              el.dataset.oldPaddingLeft = el.style.paddingLeft
              el.dataset.oldPaddingRight = el.style.paddingRight
              el.dataset.oldOverflow = el.style.overflow

              el.style.width = el.scrollWidth + 'px'
              el.style.overflow = 'hidden'
            },
            leave(el) {
              if (el.scrollWidth !== 0) {
                addClass(el, 'drawer-transition')
                el.style.width = 0
                el.style.paddingLeft = 0
                el.style.paddingRight = 0
              }
            },
            afterLeave(el) {
              removeClass(el, 'drawer-transition')
              el.style.width = ''
              el.style.overflow = el.dataset.oldOverflow
              el.style.paddingLeft = el.dataset.oldPaddingLeft
              el.style.paddingRight = el.dataset.oldPaddingRight
            }
          }
        }
        return createElement('transition', data, context.children)
      }
    }
  }
}
</script>

<style scoped>
.drawer-transition {
  transition: 0.3s width ease-in-out, 0.3s padding-left ease-in-out, 0.3s padding-right ease-in-out;
}
</style>
