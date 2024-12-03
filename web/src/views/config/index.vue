<template>
  <div class="mt-4">
    <div class="text-size20 text-themeBlue px-2 w-fit border-b-2 border-solid border-themeRed">
      {{ $t('config.modelSetting') }}
    </div>
    <el-form :model="formData"
             label-position="left"
             label-suffix=":"
             size="large"
             label-width="80px"
             class="w-1/2 xl:w-1/3  mt-10">
      <el-form-item :label="$t('config.language')">
        <el-select v-model="formData.language">
          <el-option v-for="item in data.language" :key="item.name"
                     :label="locale==='zh-cn'?item.name:item.enName" :value="item.value"></el-option>
        </el-select>
      </el-form-item>
<!--      <el-form-item :label="$t('config.model')">-->
<!--        <el-select v-model="formData.modelType">-->
<!--          <el-option v-for="item in data.models" :key="item.name" :label="item.name" :value="item.value"></el-option>-->
<!--        </el-select>-->
<!--      </el-form-item>-->
      <div class="text-center  ml-10 mt-12 flex justify-center items-center">
        <button type="button" class="px-16 w-fit bg-default text-white rounded-full py-1.5 cursor-pointer"
                @click="toSave">{{ $t('app.save') }}
        </button>
      </div>
    </el-form>
  </div>
</template>
<script setup>
import {reactive, toRefs, onMounted} from 'vue'
import {configDetail, configDict, configSave} from "@/http/api/config";
import {ElMessage} from "element-plus";
import {useI18n} from 'vue-i18n'

const {locale} = useI18n()

const data = reactive({
  formData: {
    language: 'zh-cn',
    modelType: 'Atom-13B-Chat',
  },
  models: [],
  language: []
})

onMounted(() => {
  loadData()
})

const loadData = () => {
  configDict({}).then(res => {
    data.models = res.data.models
    data.language = res.data.language
    configDetail({}).then(res => {
      data.formData = res.data
    })
  })
}

const toSave = () => {
  configSave(data.formData).then(res => {
    ElMessage.success('保存成功')
  })
}


const {formData} = toRefs(data)

</script>

<style scoped lang="less">

</style>