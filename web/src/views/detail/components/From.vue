<template>
  <div>
    <div class="flex mt-6 mb-6" id="from">
      <img src="@/assets/img/detail/web_link.png" class="w-4 h-4 mt-1"/>
      <div class="text-default text-size16 ml-3">{{ $t('detail.from') }}</div>
      <div class="flex-1"></div>
    </div>
    <div class="mt-3">
      <div
          class="mb-3 text-size14 text-[#00000080] cursor-pointer w-fit"
          v-for="(item,index) in referList"
          :key="index"
          @click="toDetail(item)">
        <el-tooltip placement="right" :hide-after="0" :show-after="100">
          <template v-slot:content>
            <p class="max-w-[300px] max-h-[1000px] overflow-y-scroll">{{ item.content }}</p>
          </template>
          <div class="border-b-transparent border-b border-solid hover:text-default hover:border-b-default">
            {{ index + 1 }}. {{ item.title }}
          </div>
        </el-tooltip>
      </div>
    </div>
  </div>
</template>

<script setup>
import {defineProps, defineEmits} from 'vue'
import {useRouter} from 'vue-router'
import router from "@/router";

const route = useRouter()

const props = defineProps({
  showFrom: {
    type: Boolean,
    default: true
  },
  referList: {
    type: Array,
    default: () => []
  },
  kdbId: {
    type: String,
    default: '0'
  },
  authType: {
    type: Number,
    default: -1
  }
})

const emits = defineEmits([
  'showFrom'
])

const setShow = () => {
  emits('showFrom', !props.showFrom)
}
const toDetail = (item) => {
  if (props.kdbId === '0') {
    window.open(item.url, '_blank')
  } else {
    const {href} = router.resolve({
      path: "/online-file",
      query: {
        url: item.url
      },
    })
    if (props.authType === 3 || props.authType === 9) {
      window.open(href, '_blank')
    }
  }
}
</script>

<style scoped lang="less">
.tip {
  visibility: hidden;

}

.item:hover .tip {
  visibility: visible
}

</style>