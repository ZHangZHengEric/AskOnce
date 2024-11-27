<template>
  <el-dialog
      append-to-body
      :style="{borderRadius:'.6rem'}"
      v-model="data.showDialog"
      class="w-11/12 lg:w-2/5">
    <div
        class="w-full relative h-52 overflow-scroll scrollbar-hidden flex border-solid border-1 rounded-md bg-white m-auto shadow-i">
        <textarea @keydown="toSend"
                  v-model="data.inputValue"
                  class="flex-1 pl-2 h-40 resize-none bg-white scrollbar-hidden border-none
                  focus:outline-none text-size14 lg:text-size16 xl:text-18 2xl:text-20 pt-2"/>
      <div class="flex rounded-full pl-4 pr-4 pt-1.5 pb-1.5 absolute left-4 bottom-4 bg-themeRed cursor-pointer">
        <span class="ml-2 text-white">{{ name ? name : $t('home.internet') }}</span>
      </div>
      <svg-icon icon-class="home_send"
                @click="toSearch"
                class="absolute right-4 bottom-4 w-8 h-8 cursor-pointer"/>
    </div>
  </el-dialog>
</template>
<script setup>
import {defineExpose, reactive, toRefs, defineEmits, defineProps} from 'vue'
import SvgIcon from "@/components/SvgIcon/index.vue";
import {useRouter} from 'vue-router'

const emits = defineEmits(['send'])
const router = useRouter()

const data = reactive({
  showDialog: false,
  inputValue: ''
})

const props = defineProps({
  name: {
    type: String,
    default: ""
  },
  loading: {
    type: Boolean,
    default: false
  }

})

const {showDialog, inputValue} = toRefs(data)

const toSend = (event) => {
  if (props.loading) {
    return
  }
  if (!event.shiftKey && event.keyCode === 13) {
    event.cancelBubble = true; //ie阻止冒泡行为
    event.stopPropagation();//Firefox阻止冒泡行为
    event.preventDefault(); //取消事件的默认动作*换行
    if (event.srcElement.value) {
      data.showDialog = false
      emits('send', data.inputValue)
    }
  }
}

const toSearch = () => {
  data.showDialog = false
  emits('send', data.inputValue)
}

defineExpose({
  showDialog,
  inputValue
})

</script>


<style scoped lang="less">

</style>