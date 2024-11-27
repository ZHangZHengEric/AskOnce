<template>
  <el-container class="share h-screen">
    <el-header class="flex ">
      <div class="ml-6 mt-8 text-default text-size32 font-[400] logo" @click="router.push('/')">AskOnce</div>
    </el-header>
    <el-main class="h-full w-full flex items-center justify-center text-default">
      <div class="p-20 bg-white text-center detail rounded-xl share-dia-bg">
        <div class="flex justify-center items-center">
          <img src="@/assets/img/detail/share_img.png" class="w-12 h-12 text-center">
        </div>
        <div class="text-size22 text-default mt-6">{{ data.info.creator }} {{ $t('share.acceptJoin') }}</div>
        <div class="text-size22 text-themeBlue mt-3">{{ data.info.kdbName }}</div>
        <div @click="sure" class="mt-6 rounded-full bg-default px-20 cursor-pointer text-white py-1.5 ">
          {{ $t('share.accept') }}
        </div>
      </div>
    </el-main>
  </el-container>

</template>
<script setup>
import {onMounted, reactive} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {knowledgeShareCodeInfo, knowledgeShareCodeVerify} from "@/http/api/knowledge";

onMounted(() => {
  loadData()
})

const data = reactive({
  info: {}
})

const route = useRoute()
const router = useRouter()

const loadData = () => {
  knowledgeShareCodeInfo({
    shareCode: route.query.code
  }).then(res => {
    data.info = res.data

  })
}

const sure = () => {
  knowledgeShareCodeVerify({
    shareCode: route.query.code
  }).then(res => {
    router.push({path: '/knowledge-manage'})
  })
}
</script>


<style scoped lang="less">

.share {
  background-image: url("@/assets/img/layout/main.png");
  background-size: 100% 100%;
  background-repeat: no-repeat;

  .share-dia-bg {
    background-image: url("@/assets/img/detail/share_bg.png");
    background-size: 100% 100%;
    background-repeat: no-repeat;
  }

  .detail {
    border: 3px solid;
    border-image-source: linear-gradient(156.12deg, #FFFFFF 6.68%, rgba(255, 255, 255, 0) 84.65%);
  }
}

.logo {
  font-family: 'AlfaSlabOne'
}

</style>