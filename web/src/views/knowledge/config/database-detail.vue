<template>
  <el-breadcrumb :separator-icon="ArrowRight" class="absolute top-4 left-12">
    <el-breadcrumb-item :to="{ path: '/knowledge-manage' }">{{ $t('knowledge.knowledgeManage') }}</el-breadcrumb-item>
    <el-breadcrumb-item
        :to="{ path: `/knowledge-config/detail`,query:{id:route.query.id,type:route.query.type,dataSource:'database'} }">
      {{ $t('knowledge.detail') }}
    </el-breadcrumb-item>
    <el-breadcrumb-item>{{ $t('knowledge.dbDetail') }}</el-breadcrumb-item>
  </el-breadcrumb>
  <el-table :data="data.tableData"
            style="width: 100%"
            show-overflow-tooltip
            class="mt-4 rounded-xl"
            :header-cell-style="{'background-color':'#E0E0E033'}">
    <el-table-column label="表名" prop="eventName"></el-table-column>
    <el-table-column label="类型" prop="eventName"></el-table-column>
    <el-table-column :label="$t('knowledge.operate')" width="150">
      <template #default="scope">
        <el-button type="text" @click="()=>{console.log(scope.row)}" ·>
          查看表结构
        </el-button>
      </template>
    </el-table-column>
  </el-table>

</template>

<script setup>
import {useRoute} from 'vue-router'
import {onMounted, reactive} from "vue"

const data = reactive({
  tableData: []
})
const route = useRoute();
const loadData = () => {

}

onMounted(() => {
  loadData()
})

</script>


<style scoped lang="less">

</style>