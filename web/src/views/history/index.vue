<template>
  <div class="pt-4">
    <div class="flex">
      <div class="text-size22 font-[900] text-default flex-1">{{ $t('history.history') }}</div>
      <el-input class="w-60 rounded-xl mr-7" :placeholder="$t('app.placeholder')" v-model="data.queryName"
                @keyup.enter="loadByKeyword">
        <template #prefix>
          <el-icon class="el-input__icon text-themeBlue">
            <search/>
          </el-icon>
        </template>
      </el-input>
      <el-radio-group style="--el-color-primary:#7269FB" v-model="data.queryType" @change="loadByKeyword">
        <el-radio-button value="">{{ $t('knowledge.all') }}</el-radio-button>
        <el-radio-button value="simple">{{ $t('home.simple') }}</el-radio-button>
        <el-radio-button value="complex">{{ $t('home.complex') }}</el-radio-button>
      </el-radio-group>
    </div>
    <el-table
        class="mt-10 rounded-xl"
        :data="data.tableData"
        @row-click="toDetail"
        :header-cell-style="{'background-color':'#E0E0E033'}">
      <el-table-column :label="$t('history.createTime')" prop="createTime"></el-table-column>
      <el-table-column :label="$t('history.knowledgeBase')" prop="kdbName">
        <template #default="scope">
          <span>{{ scope.row.kdbName ? scope.row.kdbName : $t('home.internet') }} </span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('history.keyWord')" prop="question"></el-table-column>
      <el-table-column :label="$t('history.agentName')">
        <template #default="scope">
          <span>{{ scope.row.askType ? askType[scope.row.askType] : '暂无' }} </span>
        </template>
      </el-table-column>
    </el-table>
    <div class="flex">
      <div class="flex-1"></div>
      <el-pagination
          small
          background
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
</template>
<script setup>
import {reactive, onMounted, computed} from 'vue'
import {historyAsk} from "@/http/api/aisearch";
import {useRouter} from 'vue-router'
import {useI18n} from 'vue-i18n'

const {t} = useI18n()

const router = useRouter()
const data = reactive({
  tableData: [],
  pageNo: 1,
  pageSize: 10,
  total: 0,
  queryName: '',
  queryType: ''
})

const askType = computed(() => {
  return {
    simple: t('home.simple'),
    complex: t('home.complex'),
    research: t('home.research')
  }
})

const loadByKeyword = () => {
  data.pageNo = 1
  loadData()
}

onMounted(() => {
  loadData()
})

const toDetail = (row) => {
  router.push({
    path: `/detail/${row.sessionId}`,
    query: {
      kdbName: row.kdbName,
      type: row.askType,
      question: row.question,
      kdbId: row.kdbId,
      searchEngine: row.kdbId === 0 ? 'web' : 'kdb',
    }
  })
}

const loadData = () => {
  historyAsk({
    pageNo: data.pageNo,
    pageSize: data.pageSize,
    query: data.queryName,
    queryType: data.queryType
  }).then(res => {
    data.tableData = res.data.list
    data.total = res.data.total
  })
}

</script>


<style scoped lang="less">
:deep(.el-pagination.is-background .el-pager li.is-active ) {
  background-color: white;;
  border: 1px solid rgba(114, 105, 251, 1);
  color: rgba(114, 105, 251, 1);
}

</style>

