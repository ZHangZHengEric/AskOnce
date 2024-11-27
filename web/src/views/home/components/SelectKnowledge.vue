<template>
  <el-dialog
      v-model="data.showDialog"
      width="600px"
      draggable
      :style="{padding:0}"
      class="rounded-xl"
      custom-class="my-dialog">
    <div class="text-size18 px-3 py-3 font-normal text-default">
      {{ $t('knowledge.knowledgeBase') }}
    </div>
    <div class="border-b border-solid border-b-[#d8d8f9]"></div>
    <div class="px-10 mt-8 pb-10">
      <el-input class="rounded-xl mr-7"
                :placeholder="$t('app.placeholder')"
                v-model="data.queryName"
                @keyup.enter="load">
        <template #prefix>
          <el-icon class="el-input__icon text-themeBlue">
            <search/>
          </el-icon>
        </template>
      </el-input>
      <div class="flex mt-7 text-size14 text-default font-normal">
        <div class="py-0.5">排序依据：</div>
        <div class="ml-2 cursor-pointer px-3 py-0.5"
             v-for="(item,index) in data.typeList"
             :key="index"
             :class="{select:data.selectIndex===index}"
             @click="setSelect(index)">
          {{ item.name }}
        </div>
      </div>
      <el-row class="mt-7">
        <el-col :span="12" v-for="item in data.knowledgeList" :key="item.kdbId">
          <div class="item flex mb-5 cursor-pointer" @click="itemClick(index,item)">
            <img src="@/assets/img/home/home_knowledge.png" class="w-11 h-11">
            <div class="text-default text-size14 font-normal ml-3">
              <div> {{ item.kdbName }}</div>
              <div class="mt-1">{{item.createTime.split(' ')[0]}}</div>
            </div>
          </div>
        </el-col>
      </el-row>
      <div class="flex justify-center">
        <el-pagination
            small
            background
            style="--el-color-primary:#0A2540"
            layout="prev, pager, next"
            :total="data.total"
            @currentChange="val => {
            data.pageNo=val
            loadData()
          }"
            class="mt-4"
        />
      </div>
    </div>
  </el-dialog>
</template>
<script setup>
import {defineExpose, onMounted, reactive, toRefs, defineEmits} from "vue";
import {kdbList} from "@/http/api/knowledge";

const data = reactive({
  showDialog: false,
  queryName: '',
  total: 0,
  pageNo: 1,
  pageSize: 10,
  selectIndex: 0,
  typeList: [{name: '从新到旧', orderType: 1}, {name: '从旧到新', orderType: 2}, {name: '上次打开时间', orderType: 3}],
  knowledgeList: [],
})

onMounted(() => {
  loadData()
})

const {showDialog} = toRefs(data)

const emit = defineEmits([
  'selectKnowLedge'
])
defineExpose({
  showDialog
})

const itemClick = (index, item) => {
  emit('selectKnowLedge', item, index)
  data.showDialog = false
}


const setSelect = (index) => {
  data.pageNo = 1;
  data.selectIndex = index;
  loadData()
}

const load = () => {
  data.pageNo = 1;
  loadData()
}


const loadData = () => {
  kdbList({
    pageNo: data.pageNo,
    pageSize: data.pageSize,
    query: data.queryName,
    orderType: data.typeList[data.selectIndex].orderType
  }).then(res => {
    data.knowledgeList = res.data.list
  })
}
</script>


<style scoped lang="less">
.select {
  background: #0A2540;
  color: white;
  border-radius: 30px;

}
</style>