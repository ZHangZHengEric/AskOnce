<template>
  <el-breadcrumb :separator-icon="ArrowRight" class="absolute top-4 left-12">
    <el-breadcrumb-item :to="{ path: '/knowledge-manage' }">{{ $t('knowledge.knowledgeManage') }}</el-breadcrumb-item>
    <el-breadcrumb-item>{{ getTipText() }}</el-breadcrumb-item>
  </el-breadcrumb>
  <div class="w-full flex flex-col h-full overflow-hidden">
    <div class="flex text-size18 pt-4 font-[500] leading-10 border-b border-solid cursor-pointer border-[#F0F0F0]">
      <div :class="{selected:route.path==='/knowledge-config/detail'}"
           @click="toChild('/knowledge-config/detail')">
        {{ $t('knowledge.detail') }}
      </div>
      <div class="ml-6"
           :class="{selected:route.path==='/knowledge-config/search'}"
           @click="toChild('/knowledge-config/search')">
        {{ $t('knowledge.search') }}
      </div>
      <div v-if="data.authType===9"
           class="ml-6"
           :class="{selected:route.path==='/knowledge-config/base'}"
           @click="toChild('/knowledge-config/base')">
        {{ $t('knowledge.basicInformation') }}
      </div>
      <div v-if="data.authType===9"
           class="ml-6"
           :class="{selected:route.path==='/knowledge-config/member'}"
           @click="toChild('/knowledge-config/member')">
        {{ $t('knowledge.memberSettings') }}
      </div>
      <div v-if="data.showSetting" class="ml-6"
           :class="{selected:route.path==='/knowledge-config/setting'}"
           @click="toChild('/knowledge-config/setting')">
        {{ $t('knowledge.setting') }}
      </div>
    </div>
    <router-view class="h-full flex-1"></router-view>
  </div>
</template>
<script setup>
import {onMounted, reactive} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {ArrowRight} from "@element-plus/icons-vue";
import {knowledgeAuth} from "@/http/api/knowledge";
import {useI18n} from 'vue-i18n'
import {useKnowledgeStore} from '@/store'

const knowLedgeStore = useKnowledgeStore()

const {t} = useI18n()
const route = useRoute();
const router = useRouter()
const data = reactive({
  authType: '',
  showSetting: false
})
onMounted(() => {
  knowledgeAuth({
    kdbId: parseInt(route.query.id)
  }).then(res => {
    data.authType = res.data.authType
    knowLedgeStore.setAuthType(data.authType)
    data.showSetting = route.query.type === '0'
  })
})

const getTipText = () => {
  switch (route.name) {
    case "knowledgeDetail":
      return t('knowledge.detail')
    case 'knowledgeConfigBase':
      return t('knowledge.basicInformation')
    case 'knowledgeConfigMember':
      return t('knowledge.memberSettings')
    case 'knowledgeSearch':
      return t('knowledge.search')
    case 'knowledgeSetting':
      return t('knowledge.setting')
  }
}

const toChild = (url) => {
  router.push({path: url, query: {id: route.query.id, type: route.query.type, dataSource: route.query.dataSource}})
}
</script>

<style scoped lang="less">
.selected {
  color: #7269FB;
  border-bottom: 2px solid #7269FB;
}
</style>