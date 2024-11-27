<template>
  <el-dialog v-model="data.showDialog" width="30%" draggable custom-class="rounded-xl" class="rounded-xl">
    <template #header>
      <div class="text-size20 font-normal text-default">{{$t('knowledge.add')}} {{ props.type === 3 ? $t('knowledge.editableMember') : $t('knowledge.readableMember')}}</div>
    </template>
    <div class="mt-10 mb-4 text-size16 text-normal">
      <el-input v-model="data.queryName" @keyup.enter="loadData">
        <template #prefix>
          <el-icon class="el-input__icon text-themeBlue">
            <search/>
          </el-icon>
        </template>
      </el-input>
      <div class="text-size16 mt-4 ">{{$t('knowledge.hasSelect')}}：{{ selectNum }}</div>
      <div class="h-40 overflow-y-scroll mt-4">
        <div v-for="(item,index) in data.list" :key="index" class="flex mb-2">
          <el-checkbox v-model="item.select"></el-checkbox>
          <div class="ml-2 pt-1">{{ item.userName }}</div>
        </div>
      </div>
      <div class="flex items-center justify-center">
        <div class="flex">
          <button class="border border-solid px-4 py-1 border-[969AAE] rounded-full" @click="data.showDialog=false">
            {{$t('app.cancel')}}
          </button>
          <button class="ml-4 bg-default text-white px-4 py-1 rounded-full" @click="toSave">  {{$t('app.confirm')}}</button>
        </div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import {reactive, defineExpose, toRefs, onMounted, watch, computed, defineProps, defineEmits} from 'vue'
import {knowledgeUserAdd, knowledgeUserQuery} from "@/http/api/knowledge";
import {ElMessage} from "element-plus";

const data = reactive({
  showDialog: false,
  queryName: '',
  list: []
})
const {showDialog} = toRefs(data)

watch(showDialog, (val) => {
  if (val) {
    loadData()
  }
})

const emits = defineEmits(['loadData'])

const props = defineProps({
  id: {
    type: String,
    default: '-1'
  },
  type: {
    type: Number,
    default: -1
  },
  list: {
    type: Array,
    default: () => {
      return []
    }
  }
})

onMounted(() => {
})

const selectNum = computed(() => {
  const list = data.list.filter(item => item.select)
  return list.length
})

const toSave = () => {
  const list = data.list.filter(item => item.select).map(item => {
    return item.userId
  })
  if (list.length === 0) {
    return
  }
  knowledgeUserAdd({
    kdbId: parseInt(props.id),
    authType: props.type,
    userIds: list
  }).then(res => {
    data.list = data.list.map(item => {
      item.select = false
      return item
    })
    data.showDialog = false
    emits("loadData")
    ElMessage.success('添加成功')
  })
}

const loadData = () => {
  knowledgeUserQuery({
    queryName: data.queryName
  }).then(res => {
    data.list = res.data.list.map(item => {
      item.select = props.list.some(obj1 => obj1.userId === item.userId)
      return item
    })
  })
}


defineExpose({
  showDialog
})
</script>


<style scoped lang="less">
/deep/ .el-dialog__header {
  border-bottom: 1px solid #DBDAF9;
  padding: 10px 0;
}

</style>