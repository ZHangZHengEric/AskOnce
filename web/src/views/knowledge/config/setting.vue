<template>
  <div class="w-full h-full">
    <div class="text-size20 font-bold text-default mt-8">{{ $t('knowledge.knowledgeMange') }}</div>
    <div class="bg-[#E6A23C40] flex px-3 py-1.5 text-size14 text-normal w-full font-normal rounded mt-8">
      <svg-icon icon-class="setting_waring" class="w-4 h-4 mr-2 mt-0.5"></svg-icon>
      {{ $t('knowledge.delKnowTip') }}
    </div>

    <div
        class="bg-default text-white mt-16 cursor-pointer rounded-full px-12  py-1.5 w-fit flex justify-center items-center"
        @click="showConfirm">
      {{ $t('knowledge.delKnowLedge') }}
    </div>
  </div>
  <confirm-dialog ref="refConfirm" :tip="delTip" :detail="delDetail" @confirm="confirm">
  </confirm-dialog>
</template>
<script setup>
import {knowledgeDel, knowledgeUserDeleteSelf} from "@/http/api/knowledge";
import ConfirmDialog from "@/components/Dialog/ConfirmDialog.vue";
import {computed, ref, reactive} from "vue";
import {useI18n} from 'vue-i18n'
import {useRouter, useRoute} from 'vue-router'
import {useKnowledgeStore} from '@/store'

const knowledgeStore = useKnowledgeStore()
const {t} = useI18n()
const refConfirm = ref()
const route = useRoute()
const router = useRouter()

const data = reactive({})

const delTip = computed(() => {
  return t('knowledge.knowDelTip')
})
const delDetail = computed(() => {
  return t('knowledge.knowDelDetail')
})

const showConfirm = () => {
  refConfirm.value.showDialog = true
}
const confirm = () => {
  if (knowledgeStore.authType === 9) {
    knowledgeDel({
      kdbId: parseInt(route.query.id)
    }).then(res => {
      router.push({path: '/knowledge-manage'})
    })
  } else {
    knowledgeUserDeleteSelf({
      kdbId: parseInt(route.query.id)
    }).then(res => {
      router.push({path: '/knowledge-manage'})
    })
  }
}
</script>


<style scoped lang="less">

</style>