<template>
  <el-breadcrumb :separator-icon="ArrowRight" class="absolute top-4 left-12">
    <el-breadcrumb-item :to="{ path: '/knowledge-manage' }">{{ $t('knowledge.knowledgeManage') }}</el-breadcrumb-item>
    <el-breadcrumb-item :to="{ path: `/knowledge-config/detail`,query:{id:route.query.id,type:route.query.type} }">
      {{ $t('knowledge.detail') }}
    </el-breadcrumb-item>
    <el-breadcrumb-item>{{ $t('knowledge.add') }}</el-breadcrumb-item>
  </el-breadcrumb>
  <div class="pl-10 pt-4 pr-16 h-full">
    <div class="text-size22 text-default font-black">{{ $t('knowledge.knowledgeBase') }}</div>
    <div class="mt-6 text-default text-size16">
      设置数据库信息
    </div>
    <el-form :model="formData" class="w-[50%] mt-5">
      <el-form-item label="数据库类型">
        <el-select>
          <el-option
              v-for="option in data.options"
              :key="option.value"
              :value="option.value"
              :label="option.label"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="数据库地址">
        <el-input></el-input>
      </el-form-item>
      <el-form-item label="数据库用户">
        <el-input></el-input>
      </el-form-item>
      <el-form-item label="数据库密码">
        <el-input></el-input>
      </el-form-item>
      <el-form-item label="数据库名称">
        <el-input></el-input>
      </el-form-item>
      <el-form-item label="数据库备注">
        <el-input type="textarea" :autosize="{ minRows: 3, maxRows: 5 }"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button>确认</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import {reactive} from 'vue'
import {ArrowRight} from "@element-plus/icons-vue";
import {useRoute} from 'vue-router'

const route = useRoute()

const data = reactive({
  text: '',
  title: '',
  options: [{
    label: 'mysql',
    value: 'mysql'
  }]
})
const formData = reactive({})
</script>


<style scoped lang="less">
.selected {
  color: #6793fb;
  border-bottom: 2px solid #6793fb;
}

.text-input {
  box-shadow: 0px 0px 10px 0px #00000014;
}

.name {
  overflow: hidden; //超出的文本隐藏
  text-overflow: ellipsis; //溢出用省略号显示
  white-space: nowrap; //溢出不换行
}
</style>