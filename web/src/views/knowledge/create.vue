<template>
  <el-breadcrumb :separator-icon="ArrowRight" class="absolute top-4 left-12">
    <el-breadcrumb-item :to="{ path: '/knowledge-manage' }">{{ $t('knowledge.knowledgeManage') }}</el-breadcrumb-item>
    <el-breadcrumb-item>{{ $t('knowledge.create') }}</el-breadcrumb-item>
  </el-breadcrumb>
  <el-form
      :model="formData"
      label-position="left"
      label-suffix=":"
      label-width="120px"
      size="large"
      :rules="rules"
      ref="refElForm"
      hide-required-asterisk
      class="w-11/12 lg:w-8/12  xl:w-6/12 pl-10 pt-4">
    <el-form-item :label="$t('knowledge.type')" prop="type">
      <el-select placeholder="Document" v-model="formData.type">
        <el-option label="doc" value="doc"></el-option>
        <el-option label="database" value="database"></el-option>
      </el-select>
    </el-form-item>
    <el-form-item :label="$t('knowledge.name')" prop="name">
      <el-input v-model="formData.name"></el-input>
    </el-form-item>
    <el-form-item :label="$t('knowledge.language')" prop="language">
      <el-select placeholder="Chinese" v-model="formData.language">
        <el-option label="中文" value="zh-cn"></el-option>
        <el-option label="English" value="en-us"></el-option>
      </el-select>
    </el-form-item>
    <el-form-item :label="$t('knowledge.introduction')" prop="intro">
      <el-input type="textarea" rows="6" v-model="formData.intro">
      </el-input>
    </el-form-item>
    <div class="text-center mt-20 flex justify-center ">
      <div class="px-16 w-fit bg-default text-white py-1 rounded-full cursor-pointer" @click="addKnowledge">
        {{ $t('app.save') }}
      </div>
    </div>
  </el-form>
</template>
<script setup>
import {reactive, ref} from "vue"
import {knowledgeAdd} from "@/http/api/knowledge";
import {useRouter} from 'vue-router'
import {ElMessage} from "element-plus";
import {ArrowRight} from "@element-plus/icons-vue";

const refElForm = ref()
const router = useRouter()


const rules = reactive({
  name: [
    {required: true, message: '请输入名称', trigger: 'blur'},
  ],
  intro: [
    {required: true, message: '请输入描述信息', trigger: 'blur'},
  ]
})
const formData = reactive({
  type: 'doc',
  language: 'zh-cn'
})
const addKnowledge = () => {
  refElForm.value.validate(isOk => {
    if (isOk) {
      knowledgeAdd(formData).then(res => {
        ElMessage.success('新增成功')
        router.back()
      })
    }
  })

}

</script>


<style scoped lang="less">

</style>