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
      {{ data.fileType === 0 ? $t('knowledge.addTextTip') : $t('knowledge.addFileTip') }}
    </div>
    <div class=" flex text-size18 mt-4 font-[500] leading-10 border-b border-solid cursor-pointer border-[#F0F0F0]">
      <div :class="{selected:data.fileType===0}" @click="data.fileType=0"> {{ $t('knowledge.addText') }}</div>
      <div class="ml-6" :class="{selected:data.fileType===1}" @click="data.fileType=1">{{
          $t('knowledge.addFile')
        }}
      </div>
    </div>
    <div v-show="data.fileType===0">
      <div class="text-size20 text-default font-medium mt-10">{{ $t('knowledge.textName') }}</div>
      <el-input class="py-1 mt-4 " v-model="data.title" size="large" :maxlength="200" show-word-limit/>
      <div class="text-size20 text-default font-medium mt-10">{{ $t('knowledge.textDetail') }}</div>
      <el-input type="textarea" class="mt-4 rounded-xl" :rows="8" v-model="data.text" :maxlength="10000"
                show-word-limit/>
      <button class="bg-default text-white px-20 py-2 rounded-full mt-10" @click="addData">{{
          $t('knowledge.uploadData')
        }}
      </button>
    </div>

    <div v-show="data.fileType===1">
      <div class="text-size20 text-default font-medium mt-10">{{ $t('knowledge.upHere') }}</div>
      <el-upload
          class="mt-10 rounded-xl bg-white"
          drag
          :disabled="data.uploading"
          multiple
          action="#"
          ref="refUpload"
          accept=".pdf, .docx, .doc, .txt"
          :auto-upload="false"
          :show-file-list="false"
          v-model:file-list="data.fileList"
      >
        <el-icon class="el-icon--upload">
          <upload-filled/>
        </el-icon>
        <div class="el-upload__text">
          Drop file here or <em>click to upload</em>
        </div>
      </el-upload>
      <div class="mt-4">
        <div v-for="(item,index) in data.fileList" :key="item.uid"
             class="text-size14 text-default font-normal mb-4 flex w-2/3">
          <svg-icon icon-class="file_logo" class="w-4 h-4"></svg-icon>
          <el-tooltip
              class="box-item"
              effect="dark"
              :content="item.name"
              placement="bottom"
          >
            <div class="ml-3 w-2/3 name">{{ item.name }}</div>
          </el-tooltip>
          <div class="ml-10 text-themeBlue flex-1">{{ item.upLoading }}</div>
          <div class="cursor-pointer" v-if="!data.uploading" @click="data.fileList.splice(index,1)">x</div>
        </div>
      </div>
      <button class="bg-default text-white px-20 py-2 rounded-full mt-10" @click="uploadFile">
        {{ $t('knowledge.uploadData') }}
      </button>
    </div>
  </div>
</template>
<script setup>
import {reactive, ref} from 'vue'
import {ArrowRight} from "@element-plus/icons-vue";
import {useRoute} from 'vue-router'
import {ElMessage} from "element-plus";
import {knowledgeDataAdd} from "@/http/api/knowledge";
import SvgIcon from "@/components/SvgIcon/index.vue";

const route = useRoute()
const refUpload = ref()

const data = reactive({
  fileType: 0,
  fileList: [],
  text: '',
  title: '',
  uploading: false
})

const addData = () => {
  const formData = new FormData()
  formData.append('type', 'text')
  formData.append('kdbId', route.query.id)
  formData.append('text', data.text)
  formData.append('title', data.title)
  if (!data.text || !data.title) {
    ElMessage.error('标题内容不能为空')
    return
  }
  knowledgeDataAdd(formData, {
    headers: {'Content-Type': 'multipart/form-data'}
  }).then(res => {
    data.title = ''
    data.text = ''
    ElMessage.success('上传成功')
  })
}

const getName = (item) => {
  return item.name.slice(0, 20)
}

const uploadFile = () => {
  if (data.uploading) {
    return
  }
  if (data.fileList.length) {
    uploadToServe(data.fileList[0].raw, 0)
    data.uploading = true
  }
}

const uploadSuccess = (index) => {
  if (index === data.fileList.length - 1) {
    refUpload.value.clearFiles()
    ElMessage.success('上传完成')
    data.uploading = false
  } else {
    index++
    uploadToServe(data.fileList[index].raw, index)
  }
  return index
}
// 单文件上传
const uploadToServe = (file, index) => {
  if (file.size > 1024 * 1024 * 200) {
    ElMessage.error('单文件只支持小于200M')
    if (data.fileList.length === 1) {
      refUpload.value.clearFiles()
      ElMessage.success('上传完成')
      data.uploading = false
    } else {
      index = uploadSuccess(index)
    }
  } else {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('type', 'file')
    formData.append('kdbId', route.query.id)
    knowledgeDataAdd(formData, {
      headers: {'Content-Type': 'multipart/form-data'},
      onUploadProgress: (progressEvent) => {
        const complete = (((progressEvent.loaded / progressEvent.total) * 100) | 0);
        data.fileList[index].upLoading = complete === 100 ? '上传完成' : complete + '%'
      }
    }).then(res => {
      data.fileList[index].upLoading = '上传完成'
      index = uploadSuccess(index)
    })
  }
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
