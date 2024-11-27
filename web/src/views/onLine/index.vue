<template>
  <div>
    <VuePdfEmbed v-if="data.fileType==='pdf'" :source="data.pdfSource"/>
    <div
        v-if="data.fileType==='txt'"
        class=" w-5/6 xl:w-1/2 m-auto mt-10"
    >
      {{ data.txtDetail }}
    </div>
    <div
        v-if="data.fileType==='docx'"
        ref="refDocx"
    />

    <iframe
        v-if="data.fileType==='doc'"
        :src="data.docUrl"
        class="w-screen h-screen"
        frameborder="0"
    ></iframe>
  </div>
</template>

<script setup>
import {onMounted, reactive, ref} from 'vue'
import VuePdfEmbed from 'vue-pdf-embed'
import 'vue-pdf-embed/dist/style/index.css'
import {renderAsync} from "docx-preview";  // 引入异步渲染方法
import {useRoute} from 'vue-router'
import axios from "axios";

const data = reactive({
  url: '',
  pdfSource: {},
  fileType: '',
  txtDetail: "",
  docUrl: ""
})

const refDocx = ref()

const route = useRoute()

onMounted(() => {
  data.url = route.query.url
  if (!data.url) {
    return
  }
  data.fileType = getFileType(data.url)
  showFile()
})

const showFile = () => {
  switch (data.fileType) {
    case 'pdf':
      data.pdfSource = {
        url: data.url,
        cMapUrl: 'https://cdn.jsdelivr.net/npm/pdfjs-dist@2.5.207/cmaps/',
        cMapPacked: true
      }
      break
    case "txt":
      axios({
        method: 'get',
        url: data.url,
        responseType: 'blob'
      }).then(res => {
        txt2utf8(res.data, (text) => {
          data.txtDetail = text
        })
      })
      break
    case "docx":
      axios({
        method: 'get',
        url: data.url,
        responseType: 'blob'
      }).then(res => {
        const blob = new Blob([res.data])
        renderAsync(blob, refDocx.value)
      })
      break
    case "doc":
      data.docUrl = `https://view.officeapps.live.com/op/embed.aspx?src=${data.url}`
      break
  }
}

const txt2utf8 = (file, callback) => {
  const reader = new FileReader()
  // readAsText 可以 指定编码格式 将文件提取成 纯文本
  reader.readAsText(file)
  reader.onload = e => {
    const txtString = e.target.result
    // utf-8 的 中文编码 正则表达式
    const patrn = /[\u4E00-\u9FA5]|[\uFE30-\uFFA0]/gi
    // 两个格式的英文编码一样，所以纯英文文件也当成乱码再处理一次
    if (!patrn.exec(txtString)) {
      const reader_gb2312 = new FileReader()
      reader_gb2312.readAsText(file, 'gb2312')
      reader_gb2312.onload = e2 => {
        callback && callback(e2.target.result)
      }
    } else {
      callback && callback(txtString)
    }
  }
}


const getFileType = (val) => {
  const fileName = val.lastIndexOf('.') // 取到文件名开始到最后一个点的长度
  const fileNameLength = val.length // 取到文件名长度
  const fileFormat = val.substring(fileName + 1, fileNameLength) // 截
  if (fileFormat.toLowerCase() === 'pdf') {
    return 'pdf'
  } else if (fileFormat.toLowerCase() === 'docx') {
    return 'docx'
  } else if (fileFormat.toLowerCase() === 'html') {
    return 'html'
  } else if (fileFormat.toLowerCase() === 'doc') {
    return 'doc'
  } else if (fileFormat.toLowerCase() === 'xlsx' || fileFormat.toLowerCase() === 'xls') {
    return 'xlsx'
  } else if (fileFormat.toLowerCase() === 'png' || fileFormat.toLowerCase() === 'jpg' || fileFormat.toLowerCase() === 'jpeg' || fileFormat.toLowerCase() === 'bmp') {
    return 'img'
  } else if (fileFormat.toLowerCase() === 'txt' || fileFormat.toLowerCase() === 'msg' || fileFormat.toLowerCase() === 'eml') {
    return 'txt'
  } else if (fileFormat.toLowerCase() === 'ppt' || fileFormat.toLowerCase() === 'pptx') {
    return 'pptx'
  }
  return 'other'
}


</script>


<style scoped lang="less">

</style>