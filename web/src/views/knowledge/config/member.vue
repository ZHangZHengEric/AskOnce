<template>
  <div class="w-full flex flex-col">
    <div class="rounded-full bg-white mt-8 border-solid border border-[#DCDFE6]
    flex w-fit px-4 py-2 leading-5 cursor-pointer text-default">
      <div class="pr-2 border-r border-solid border-[#DCDFE6]"
           :class="{'text-themeBlue':data.type===3}"
           @click="setType(3)">
        {{ $t('knowledge.editableMember') }}
      </div>
      <div class="ml-2"
           :class="{'text-themeBlue':data.type===2}"
           @click="setType(2)">
        {{ $t('knowledge.readableMember') }}
      </div>
    </div>
    <div class="w-full flex mt-8">
      <el-input class="rounded-xl w-1/3" :placeholder="$t('app.placeholder')" size="large" v-model="data.queryName">
        <template #prefix>
          <el-icon class="el-input__icon">
            <search/>
          </el-icon>
        </template>
      </el-input>
      <div class="flex-1"></div>
      <div v-if="false" class="bg-default text-white cursor-pointer rounded-full px-10 flex justify-center items-center"
           @click="toAdd">
        <svg-icon icon-class="add" class="w-4 h-4 mr-2"></svg-icon>
        {{ $t('knowledge.addMember') }}
      </div>
      <div class="bg-default text-white cursor-pointer rounded-full px-10 flex justify-center items-center"
           @click="toCreateShare">
        <svg-icon icon-class="add" class="w-4 h-4 mr-2"></svg-icon>
        {{ $t('knowledge.createShare') }}
      </div>
    </div>
    <div class="flex-1 mt-8 rounded-xl w-full overflow-y-scroll bg-white border border-solid border-[#DBDAF9]">
      <div class="h-full" v-if="data.list.length">
        <div v-for="(item,index) in data.list" :key="index"
             class="flex mx-4 py-3 border-b border-b-[#DBDAF9] border-solid last:border-b-transparent">
          <img src="@/assets/img/detail/detail_member.png" class="w-12 h-12 ">
          <div class="ml-4 flex-1 text-size16 text-default">
            <div class="">{{ item.userName }}</div>
            <div class="">{{ item.joinTime }}</div>
          </div>
          <el-icon @click="deletePerson(item,index)" class="mt-4 mr-1 cursor-pointer">
            <Delete/>
          </el-icon>
        </div>
      </div>
      <el-empty
          v-else
          :description="$t('knowledge.noData')"
          :image-size="300"
          :image="require('@/assets/img/config/no_data.png')"
      />
    </div>
  </div>
  <select-person
      ref="refSelectPerson"
      :id="route.query.id"
      :type="data.type"
      @load-data="loadData"
      :list="data.list"></select-person>
</template>
<script setup>
import {reactive, ref, onMounted} from 'vue'
import SelectPerson from "@/views/knowledge/componetns/SelectPerson.vue";
import SvgIcon from "@/components/SvgIcon/index.vue";
import {knowledgeShareCodeGen, knowledgeUserDelete, knowledgeUserList} from "@/http/api/knowledge";
import {useRoute} from 'vue-router'
import {copyToClip} from "@/utils/tools";
import {ElMessage} from 'element-plus'

const route = useRoute()

const data = reactive({
  type: 3,
  list: [],
  queryName: ''
})

const refSelectPerson = ref()

onMounted(() => {
  loadData()
})

const loadData = () => {
  knowledgeUserList({
    kdbId: parseInt(route.query.id),
    authType: data.type,
    pageNo: 1,
    pageSize: 1000,
    queryName: data.queryName
  }).then(res => {
    data.list = res.data.list
  })
}

const toCreateShare = () => {
  knowledgeShareCodeGen({
    kdbId: parseInt(route.query.id),
    authType: data.type,
  }).then(res => {
    const shareCode = res.data.shareCode
    copyToClip(`${window.location.origin}/knowledge-share?code=${shareCode}`)
    ElMessage.success('已复制到剪切板')
  })
}

const deletePerson = (item, index) => {
  knowledgeUserDelete({
    kdbId: parseInt(route.query.id),
    authType: data.type,
    userIds: [item.userId]
  }).then(res => {
    data.list.splice(index, 1)
  })
}

const setType = (type) => {
  data.type = type
  data.queryName = ''
  loadData()
}

const toAdd = () => {
  refSelectPerson.value.showDialog = true
}

</script>
<style scoped lang="less">
/deep/ .el-input__wrapper {
  border-radius: 8px;
}
</style>