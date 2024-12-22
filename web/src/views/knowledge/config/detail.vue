<template>
  <div class="pt-4 flex h-full">
    <div class="w-full h-full">
      <div class="w-full flex">
        <el-input
            class="rounded-xl w-1/3"
            v-model="data.queryName"
            @keyup.enter="()=>{
              data.pageNo=1
              loadData()
            }"
            :placeholder="$t('app.placeholder')">
          <template #prefix>
            <el-icon class="el-input__icon text-themeBlue">
              <search/>
            </el-icon>
          </template>
        </el-input>
        <div class="flex-1"></div>
        <button class="bg-default text-white rounded-full px-10" @click="toAdd">{{ $t('knowledge.add') }}</button>
      </div>
      <el-table
          v-loading="data.loading"
          class="mt-4 rounded-xl"
          :data="data.tableData"
          style="width: 100%"
          show-overflow-tooltip
          :header-cell-style="{'background-color':'#E0E0E033'}">
        <template v-if="route.query.dataSource === 'database'">
          <el-table-column :label="$t('knowledge.dbName')" prop="dataName"></el-table-column>
        </template>
        <template v-else>
          <el-table-column :label="$t('knowledge.dataType')" prop="dataSuffix"></el-table-column>
          <el-table-column :label="$t('knowledge.fileName')" prop="dataName"></el-table-column>
        </template>
        <el-table-column :label="$t('knowledge.createAt')" prop="createTime"></el-table-column>
        <el-table-column :label="$t('knowledge.status')">
          <template #default="scope">
            <span>{{ getStatus(scope.row) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('knowledge.operate')" width="150">
          <template #default="scope">
            <div class="flex">
              <svg-icon
                  icon-class="know_del"
                  class="w-4 h-4 cursor-pointer"
                  @click="del(scope.$index, scope.row)">
              </svg-icon>
              <el-icon class="ml-2 cursor-pointer" v-if="canAdd" @click="downLoad(scope.$index, scope.row)">
                <Download/>
              </el-icon>
              <el-icon @click="toDetail(scope.$index, scope.row)" class="ml-2 cursor-pointer" v-if="canAdd">
                <el-icon>
                  <View/>
                </el-icon>
              </el-icon>
              <el-tooltip :content="$t('knowledge.rebuild')" v-if="scope.row.status===2">
                <el-icon class="ml-2 cursor-pointer text-themeBlue text-size18" @click="build(scope.$index, scope.row)">
                  <Tools/>
                </el-icon>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
      </el-table>
      <div class="flex">
        <div class="flex-1"></div>
        <el-pagination
            small
            background
            @currentChange="(page)=>{
              data.pageNo =page
              loadData()
            }"
            layout="prev, pager, next"
            :total="data.total"
            class="mt-4"
        />
      </div>
    </div>
  </div>
  <confirm-dialog ref="refConfirm" :tip="delTip" :message="data.delMessage" :detail="delDetail" @confirm="confirm">
  </confirm-dialog>
</template>
<script setup>
import {reactive, onMounted, computed, ref} from 'vue'
import {useRouter, useRoute} from 'vue-router'
import {knowledgeDataDel, knowledgeDataList, knowledgeDataRedo} from "@/http/api/knowledge";
import {ElMessage} from "element-plus";
import ConfirmDialog from "@/components/Dialog/ConfirmDialog.vue";
import {useI18n} from 'vue-i18n'
import axios from "axios";
import {useKnowledgeStore} from '@/store'

const knowledgeStore = useKnowledgeStore()
const {t} = useI18n()
const router = useRouter()
const route = useRoute()
const refConfirm = ref()

const data = reactive({
  tableData: [],
  pageNo: 1,
  pageSize: 10,
  total: 0,
  queryName: '',
  canAdd: false,
  delMessage: {},
  loading: false
})

onMounted(() => {
  loadData()
})

const downLoad = (index, item) => {
  if (route.query.dataSource === 'database'){
    return
  }
  const url = item.dataPath
  axios({
    method: 'get',
    url: url,
    responseType: 'blob'
  }).then(res => {
    const blob = new Blob([res.data])
    if (window.navigator.msSaveOrOpenBlob) {
      navigator.msSaveBlob(blob, item.dataName)
    } else {
      const link = document.createElement('a')
      const body = document.querySelector('body')
      link.href = window.URL.createObjectURL(blob)
      link.download = item.dataName
      link.style.display = 'none'
      body.appendChild(link)
      link.click()
      body.removeChild(link)
      window.URL.revokeObjectURL(link.href)
    }
  })
}

const toDetail = (index, item) => {
  if (route.query.dataSource === 'database'){
    router.push({
      name:'databaseDetail',
      query: {
        id: route.query.id,
        type: route.query.type,
        dataSource: route.query.dataSource,
        databaseId: item.id
      }
    })
    return
  }
  const {href} = router.resolve({
    path: "/online-file",
    query: {
      url: item.dataPath
    },
  })
  window.open(href, '_blank')
}

const delTip = computed(() => {
  return t('knowledge.knowDelTip')
})
const delDetail = computed(() => {
  return t('knowledge.knowDelDetail')
})

const getStatus = (row) => {
  switch (row.status) {
    case 0:
      return '等待中'
    case 1:
      return '正在构建到知识库'
    case 2:
      return '失败'
    case 9:
      return '成功'
  }
}

const loadData = () => {
  knowledgeDataList({
    kdbId: parseInt(route.query.id),
    pageNo: data.pageNo,
    pageSize: data.pageSize,
    queryName: data.queryName
  }).then(res => {
    data.tableData = res.data.list
    data.total = res.data.total
  })
}

const canAdd = computed(() => {
  return knowledgeStore.authType === 9 || knowledgeStore.authType === 3
})

const toAdd = () => {
  if (canAdd.value) {
    if (route.query.dataSource === 'database') {
      router.push({path: '/database-add', query: {id: route.query.id, type: route.query.type}})
    } else {
      router.push({path: '/knowledge-add', query: {id: route.query.id, type: route.query.type}})
    }
  } else {
    ElMessage.error("您没有权限")
  }
}


const del = (index, row) => {
  data.delMessage = {...row, index: index}
  refConfirm.value.showDialog = true
}

const confirm = (message) => {
  data.loading = true
  knowledgeDataDel({
    kdbId: parseInt(route.query.id),
    dataId: message.id
  }).then(res => {
    data.tableData.splice(message.index, 1)
    data.loading = false
  }).catch(err => {
    data.loading = false
  })
}
const build = (index, row) => {
  knowledgeDataRedo({
    kdbId: parseInt(route.query.id),
    dataId: row.id
  }).then(res => {
    loadData()
    ElMessage.success('重新构建中')
  })
}

</script>


<style scoped lang="less">
:deep(.el-input__wrapper) {
  border-radius: 8px !important;
}

</style>