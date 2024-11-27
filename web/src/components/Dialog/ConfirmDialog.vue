/**
* create by wanghu on 2024/2/29 13:48
* 类注释：
* 备注：
*/
<template>
  <el-dialog v-model="data.showDialog" class="rounded-xl  w-1/2  xl:w-1/3">
    <template v-slot:header>
      <div class="text-default text-size20 font-normal">
        {{ $t('app.systemPrompt') }}
      </div>
    </template>
    <div class="flex justify-center items-center">
      <div>
        <div class="flex text-default text-size18 font-normal mt-6">
          <svg-icon icon-class="waring" class="w-6 h-6"/>
          <div class="ml-2">{{ tip }}</div>
        </div>
        <div class="text-size16 mt-6 mb-10">{{ detail }}</div>
      </div>
    </div>
    <div class="w-full flex justify-center items-center">
      <div class="flex">
        <div
            class="px-8 py-1 border border-solid cursor-pointer border-[#969AAE] text-[#969AAE] rounded-full"
            @click="cancel">
          {{ $t('app.cancel') }}
        </div>
        <div
            class="px-8 py-1 ml-4 rounded-full cursor-pointer bg-default text-white"
            @click="confirm">
          {{ $t('app.confirm') }}
        </div>
      </div>

    </div>
  </el-dialog>
</template>

<script setup>
import {reactive, defineExpose, defineProps, defineEmits, toRefs} from 'vue'
import SvgIcon from "@/components/SvgIcon/index.vue";

const data = reactive({
  showDialog: false,
})

const props = defineProps({
  message: {
    type: Object,
    default: () => {
    }

  },
  tip: String,
  detail: String,
})

const emit = defineEmits(["confirm"])

const confirm = () => {
  data.showDialog = false
  emit("confirm", props.message)
}

const cancel = () => {
  data.showDialog = false
}

const {showDialog} = toRefs(data)

defineExpose({
  showDialog
})
</script>

<style lang="less" rel="stylesheet/less">
</style>