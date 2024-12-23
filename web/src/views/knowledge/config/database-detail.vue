<template>
  <el-breadcrumb :separator-icon="ArrowRight" class="absolute top-4 left-12">
    <el-breadcrumb-item :to="{ path: '/knowledge-manage' }">{{ $t('knowledge.knowledgeManage') }}</el-breadcrumb-item>
    <el-breadcrumb-item
        :to="{ path: `/knowledge-config/detail`,query:{id:route.query.id,type:route.query.type,dataSource:'database'} }">
      {{ $t('knowledge.detail') }}
    </el-breadcrumb-item>
    <el-breadcrumb-item>{{ $t('knowledge.dbDetail') }}</el-breadcrumb-item>
  </el-breadcrumb>
  <div class="pl-10 pt-4 h-full">
    <div class="mt-6 text-default text-size16">
      数据库信息
    </div>
    <el-table
        :data="data.tableData"
        style="width: 100%"
        show-overflow-tooltip
        row-key="table_name"
        class="mt-4 rounded-xl"
        :header-cell-style="{'background-color':'#E0E0E033'}">
      <el-table-column label="表名" prop="table_name"></el-table-column>
      <el-table-column label="列名" prop="column_name"></el-table-column>
      <el-table-column label="列类型" prop="column_type"></el-table-column>
      <el-table-column label="列描述" prop="column_comment"></el-table-column>
    </el-table>
  </div>

</template>

<script setup>
import {useRoute, useRouter} from 'vue-router'
import {onMounted, reactive} from "vue"
import {tableInfo} from "@/http/api/knowledge";
import {ArrowRight} from "@element-plus/icons-vue";

const data = reactive({
  tableData: []
})
const route = useRoute();
const router = useRouter();
const loadData = () => {
  tableInfo({kdbId: route.query.id, dataId: route.query.databaseId}).then(res => {
    data.tableData = res.data.dbSchema.map(item => {
      return {
        table_name: item.table_name,
        children: item.column_infos.map(column => {
          return {
            column_name: column.column_name,
            // table_name: column.column_name,
            column_type: column.column_type,
            column_comment: column.column_comment
          }
        }),

      }
    })
    console.log(data.tableData)
  })
}

onMounted(() => {
  loadData()
})

</script>


<style scoped lang="less">

</style>