<template>
  <div>
    <div class="mt-8">
      <el-icon class="mr-2 icon-loading w-4 h-4">
        <Loading/>
      </el-icon>
      {{ loadingText }}
    </div>
    <div class="relative ml-5 mt-8 w-11/12 ">
      <div class="mb-5 h-1 rounded-full bg-gray-200 flex-1 mt-[20px]">
        <div class="h-1 rounded-full bg-themeRed" :style="{width:props.progress+'%'}"></div>
      </div>
      <div class="absolute z-10 -top-2" :style="{left:getIndex(index)+'%'}" v-for="(item,index) in list" :key="item">
        <svg-icon :icon-class="getIconClass(index)" class="w-5 h-5 bg-theme rounded-full"></svg-icon>
        <div class="-ml-5 mt-3 w-40" :class="{'text-themeRed':props.progress>getIndex(index)}">{{ item }}</div>
      </div>
    </div>
    <div class="mt-16" v-if="analyseing">
      <div class="flex text-size18 text-color333">
        <svg-icon icon-class="analysis" class="w-5 h-5 mt-1 mr-2"></svg-icon>
        <span class="flex">{{ $t('detail.analyse') }}<span class="dot ml-1">...</span></span>
      </div>
      <div class="bg-[#F5F5F5] p-4 mt-3 rounded-xl" v-show="!detail">
        <el-skeleton :rows="4"/>
      </div>
    </div>

  </div>
</template>

<script setup>
import {defineProps} from 'vue'
import SvgIcon from "@/components/SvgIcon/index.vue";


const props = defineProps({
  analyseing: {
    type: Boolean,
    default: false
  },
  detail: {
    type: String,
    default: ""
  },
  loadingText: {
    type: String,
    default: ''
  },
  progress: {
    type: Number,
    default: 0
  },
  type: {
    type: String,
    default: "simple"
  },
  list: {
    type: Array,
    default: () => {
      return []
    }
  }
})

const getIndex = (index) => {
  const length = props.list.length - 1
  return parseInt(100 / length * index)
}

const getIconClass = (index) => {
  const length = props.list.length - 1
  const aa = parseInt(100 / length * (index + 1))
  if (props.progress > getIndex(index) && props.progress <= aa) {
    return 'web_loading'
  } else if (props.progress >= aa) {
    return 'web_ok'
  } else {
    return "web_waiting"
  }
}


</script>

<style scoped lang="less">
.icon-loading {
  animation: loading 1s linear infinite;
}

@keyframes loading {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.dot {
  /*让点垂直居中*/
  height: 1em;
  line-height: 1;
  /*让点垂直排列*/
  display: flex;
  flex-direction: column;
  /*溢出部分的点隐藏*/
  overflow: hidden;
}

.dot::before {
  /*三行三种点，需要搭配white-space:pre使用才能识别\A字符*/
  content: '...\A..\A.';
  white-space: pre-wrap;
  animation: dot 1s infinite step-end; /*step-end确保一次动画结束后直接跳到下一帧而没有过渡*/
}

@keyframes dot {
  33% {
    transform: translateY(-2em);
  }
  66% {
    transform: translateY(-1em);
  }
}

</style>