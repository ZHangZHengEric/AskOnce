<template>
  <el-breadcrumb :separator-icon="ArrowRight" class="absolute top-4 left-12">
    <el-breadcrumb-item>{{ $t('knowledge.knowledgeManage') }}</el-breadcrumb-item>
  </el-breadcrumb>
  <div class="w-full h-full flex flex-wrap mt-4">
    <div class="flex w-full">
      <div class="text-default text-size22 font-[700]">{{ $t('knowledge.knowledgeBase') }}</div>
      <svg-icon icon-class="knowledge_all" class="w-6 h-6 mt-1 ml-4"></svg-icon>
      <div class="flex-1"></div>
      <el-input class="w-60 rounded-xl mr-7" :placeholder="$t('app.placeholder')" v-model="data.queryName"
                @keyup.enter="loadData(true)">
        <template #prefix>
          <el-icon class="el-input__icon text-themeBlue">
            <search/>
          </el-icon>
        </template>
      </el-input>
      <button @click="router.push({path:'/knowledge-create'})"
              class="text-default px-3 rounded-full border border-solid border-default">
        {{ $t('knowledge.create') }}
      </button>
    </div>
    <div v-infinite-scroll="loadData"
         :infinite-scroll-disabled="disabled"
         class="w-full h-full mt-10 overflow-y-scroll">
      <el-row>
        <el-col
            v-for="(item,index) in data.list"
            :key="item.id"
            :lg="4"
            :xl="4"
            :sm="6"
            class="mb-8">
          <div class="relative h-full cursor-pointer border-solid border-2 border-transparent mr-4"
               @click="toDetail(item)"
               :class="{selected:data.selectIndex===index}"
               @mouseleave="item.showDetail=false">
            <img v-lazy="item.cover" class="w-full rounded-xl"/>
            <div class="absolute top-10 w-full text-center text-[#697FB7] text-size20 font-bold px-4"
                 :class="{'text-white':item.defaultColor}">
              {{ item.name }}
            </div>
            <div class="absolute flex top-0 left-0 bg-themeRed rounded-tl-xl text-white px-4 py-1 rounded-br-3xl"
                 :class="{'bg-default':item.type===0}"
                 @click.stop="showDetail(item)">
              <div>
                {{ item.type === 0 ? $t('knowledge.private') : $t('knowledge.company') }}
              </div>
              <svg-icon icon-class="down" class="w-3 h-3 mt-2 ml-2"></svg-icon>
            </div>
            <div class="absolute top-9 left-1 bg-white rounded" v-if="item.showDetail">
              <div class="text-default text-size14 py-1 px-3 border-b border-solid border-[#DBDAF9] flex">
                <svg-icon icon-class="detail" class="w-3 h-3 mr-2 mt-1"></svg-icon>
                {{ $t('knowledge.detail') }}
              </div>
              <div class="text-size12 font-normal text-normal leading-6 px-3 mb-2">
                <div>
                  {{ $t('knowledge.founder') }}:
                </div>
                <div>
                  {{ item.creator ? item.creator : "admin" }}
                </div>
                <div>
                  {{ $t('knowledge.createTime') }}:
                </div>
                <div>
                  {{ item.createTime.split(" ")[0] }}
                </div>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import {reactive, computed, ref} from 'vue'
import {ArrowRight} from '@element-plus/icons-vue'
import {useRouter} from 'vue-router'
import {
  knowledgeList,
} from "@/http/api/knowledge";
import {useI18n} from 'vue-i18n'

const {t} = useI18n()
const disabled = computed(() => data.loading || data.noMore)
const router = useRouter()
const data = reactive({
  selectIndex: -1,
  selectItem: {},
  queryName: '',
  pageNo: 1,
  pageSize: 12,
  list: [],
  total: 0,
  noMore: false,
  loading: false,
  showConfig: -1,
})

const showDetail = (item) => {
  data.list.forEach(tem => {
    if (item.id === tem.id) {
      tem.showDetail = !tem.showDetail
    } else {
      tem.showDetail = false
    }
  })
}

const loadData = (refresh) => {
  data.selectIndex = -1
  data.selectItem = {}
  data.loading = true
  if (refresh) {
    data.pageNo = 1
  }
  knowledgeList({
    queryName: data.queryName,
    pageNo: data.pageNo,
    pageSize: data.pageSize
  }).then(res => {
    data.loading = false
    data.total = res.data.total
    if (data.pageNo === 1)
      data.list = res.data.list
    else {
      data.list = [...data.list, ...res.data.list]
    }
    data.noMore = data.list.length >= data.total
    data.pageNo++
  }).catch(err => {
    data.loading = false
  })
}

const toDetail = (item) => {
  router.push({path: '/knowledge-config/detail', query: {id: item.id, type: item.type, dataSource: item.dataSource}})
}

</script>

<style scoped lang="less">
.selected {
  border: 2px solid #7269FB;;
  border-radius: 15px;
}

.bg-default {
  background: #0A2540;
}
</style>