<template>
  <div class="mt-6">
    <div class="flex">
      <svg-icon icon-class="event" class="w-5 h-4 mt-2"></svg-icon>
      <div class="text-size20 ml-3 font-bold text-default flex-1">{{ $t('detail.eventDependent') }}</div>
      <svg-icon icon-class="detail_hide" class="w-6 h-6 cursor-pointer mt-2 " @click="show=!show" :class="{'rotate-180':!show}"></svg-icon>
    </div>
    <el-collapse-transition>
      <el-table
          :data="eventsInfo"
          v-if="show"
          size="large"
          :header-cell-style="{'background-color':'#E0E0E033'}"
          class="rounded-xl mt-6 my-table">
        <el-table-column :label="$t('detail.eventName')" prop="eventName"></el-table-column>
        <el-table-column :label="$t('detail.eventDate')" prop="eventDate"></el-table-column>
        <el-table-column :label="$t('detail.eventDescribe')" prop="eventDescribe"></el-table-column>
      </el-table>
    </el-collapse-transition>
  </div>
</template>

<script setup>
import {defineProps, ref} from 'vue'
import SvgIcon from "@/components/SvgIcon/index.vue";

const props = defineProps({
  eventsInfo: {
    type: Array,
    default: () => []
  }
})

const show = ref(true)

</script>

<style scoped lang="less">
.my-table {
  box-shadow: 0px 0px 10px 0px rgba(10, 37, 64, 0.06);
}
</style>