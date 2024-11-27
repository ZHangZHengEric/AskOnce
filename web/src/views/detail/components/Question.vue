<template>
  <div class="relative">
    <div class="fixed top-2" v-if="scrollTop>60">
      <div class="text-size32 text-default font-[500]" :class="{'mr-10':processList.length}">
        {{ question }}
      </div>
    </div>
    <div class="text-size32 text-default font-[500]" :class="{'mr-10':processList.length}">
      {{ question }}
      <img @click="showEditQuestion"
           src="@/assets/img/detail/edit_question.png"
           class="w-8 h-8 ml-2 cursor-pointer inline-block align-middle"/>
    </div>
    <el-popover
        placement="bottom-end"
        trigger="hover"
        width="340"
        :popper-style="{'border-radius': '8px','box-shadow':'none'}"
    >
      <el-timeline class="rounded-2xl pt-5 pl-5" popper-class="rounded-2xl">
        <el-timeline-item
            v-for="(activity, index) in processList"
            :key="index"
            center
        >
          <template #dot>
            <svg-icon icon-class="point" class-name="w-4 h-4"></svg-icon>
          </template>
          <div class="">
            <div class="text-size14 text-default font-normal">{{ activity.name }}</div>
            <div v-if="activity.list.length===1">
              <div class="flex text-size12 mt-3">
                <svg-icon icon-class="detail_time" class="w-3.5 h-3.5 mr-2"></svg-icon>
                {{ dayjs(activity.list[0].time).format('YYYY-MM-DD HH:mm:ss') }}
              </div>
            </div>
            <div v-else class="text-size12 mt-2 text-normal" v-for="(tem,i) in activity.list" :key="i">
              <div v-if="i>0" class="flex">
                <div>{{ $t('detail.step') }}{{ i }}:</div>
                <div class="ml-1 flex-1">
                  {{ tem.content }}
                </div>
                <div>
                  <svg-icon icon-class="detail_mtime" class="w-3.5 h-3.5 ml-2"></svg-icon>
                </div>
                <div class="ml-1">
                  {{ tem.time - activity.list[i - 1].time }} {{ $t('detail.millisecond') }}
                </div>
              </div>
            </div>
          </div>
        </el-timeline-item>
      </el-timeline>
      <template #reference>
        <div v-show="processList.length"
             class="cursor-pointer absolute right-0 top-0 w-10 mt-2 text-size16 text-themeBlue flex">
          <svg-icon icon-class="more" class-name="w-8 h-8 ml-3 mt-1"></svg-icon>
        </div>
      </template>
    </el-popover>
  </div>
</template>
<script setup>
import {defineProps, defineEmits} from 'vue'
import dayjs from "dayjs";
import SvgIcon from "@/components/SvgIcon/index.vue";

const props = defineProps({
  question: {
    type: String,
    default: ""
  },
  scrollTop: {
    type: Number,
    default: 0
  },
  loading: {
    type: Boolean,
    default: false
  },
  processList: {
    type: Array,
    default: () => []
  }
})

const emits = defineEmits([
  'showEditQuestion'
])


const showEditQuestion = () => {
  emits('showEditQuestion')
}

</script>

<style scoped lang="less">
/deep/ .el-timeline-item__tail {
  border-color: rgba(10, 37, 64, 0.5);
  border-width: 1px;
  margin-left: 3px;
}

/deep/ .el-timeline-item__node {
  background-image: url("@/assets/img/detail/share_img.png");
}

/deep/ .el-popover.el-popper {
  border-radius: 15px !important;

}

</style>