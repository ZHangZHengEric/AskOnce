<template>
  <el-breadcrumb :separator-icon="ArrowRight" class="absolute top-4 left-12">
    <el-breadcrumb-item :to="{ path: '/knowledge-manage' }">{{ $t('knowledge.knowledgeManage') }}</el-breadcrumb-item>
    <el-breadcrumb-item
        :to="{ path: `/knowledge-config/detail`,query:{id:route.query.id,type:route.query.type,dataSource:'database'} }">
      {{ $t('knowledge.detail') }}
    </el-breadcrumb-item>
    <el-breadcrumb-item>{{ $t('knowledge.add') }}</el-breadcrumb-item>
  </el-breadcrumb>
  <div class="pl-10 pt-4 pr-16 h-full">
    <div class="text-size22 text-default font-black">{{ $t('knowledge.knowledgeBase') }}</div>
    <div class="mt-6 text-default text-size16">
      设置数据库信息
    </div>
    <el-form
        ref="refForm"
        :model="formData"
        label-position="right"
        label-width="100px"
        class="w-[50%] mt-5"
        :rules="rules">
      <el-form-item label="数据库类型" prop="dbType">
        <el-select v-model="formData.dbType">
          <el-option
              v-for="option in data.options"
              :key="option.value"
              :value="option.value"
              :label="option.label"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="数据库地址" prop="dbHost">
        <el-input v-model="formData.dbHost"></el-input>
      </el-form-item>
      <el-form-item label="数据库端口" prop="dbPort">
        <el-input v-model="formData.dbPort"></el-input>
      </el-form-item>
      <el-form-item label="数据库名称" prop="dbName">
        <el-input v-model="formData.dbName"></el-input>
      </el-form-item>
      <el-form-item label="数据库用户" prop="dbUser">
        <el-input v-model="formData.dbUser"></el-input>
      </el-form-item>
      <el-form-item label="数据库密码" prop="dbPwd">
        <el-input v-model="formData.dbPwd"></el-input>
      </el-form-item>
      <el-form-item label="数据库备注">
        <el-input v-model="formData.dbComment" type="textarea" :autosize="{ minRows: 3, maxRows: 5 }"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button @click="sure">确认</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import {reactive, ref} from 'vue'
import {ArrowRight} from "@element-plus/icons-vue";
import {useRoute, useRouter} from 'vue-router'
import {ElMessage} from "element-plus";
import {knowledgeDataAdd} from "@/http/api/knowledge";

const route = useRoute()
const router = useRouter()
const refForm = ref()

const rules = reactive({
  dbType: [{required: true, message: '请输入数据库类型', trigger: 'change'}],
  dbHost: [{required: true, message: '请输入数据库地址', trigger: 'blur'}],
  dbPort: [{required: true, message: '请输入数据库端口', trigger: 'blur'}],
  dbName: [{required: true, message: '请输入数据库名称', trigger: 'blur'}],
  dbUser: [{required: true, message: '请输入数据库用户', trigger: 'blur'}],
  dbPwd: [{required: true, message: '请输入数据库密码', trigger: 'blur'}],
})

const data = reactive({
  text: '',
  title: '',
  options: [{
    label: 'mysql',
    value: 'mysql'
  }, {
    label: 'postgresql',
    value: 'postgresql'
  }]
})
const formData = reactive({})
const sure = () => {
  refForm.value.validate(isOk => {
    console.log(isOk)
    if (isOk) {
      const postData = new FormData()
      postData.append('type', 'database')
      postData.append('kdbId', route.query.id)
      postData.append('dbType', formData.dbType)
      postData.append('dbHost', formData.dbHost)
      postData.append('dbPort', formData.dbPort)
      postData.append('dbName', formData.dbName)
      postData.append('dbUser', formData.dbUser)
      postData.append('dbPwd', formData.dbPwd)
      postData.append('dbComment', formData.dbComment)
      knowledgeDataAdd(postData, {
        headers: {'Content-Type': 'multipart/form-data'}
      }).then(res => {
        ElMessage.success('新增成功')
        router.back()
      })
    }
  })
}

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