<template>
  <div
      class="h-full w-full overflow-y-scroll scrollbar-hidden pt-20 flex justify-center">
    <div class="w-11/12 lg:w-4/5 xl:3/5 2xl:1/2">
      <div class="text-center text-size32 font-[600]">You Only Ask Once</div>
      <div
          class="w-full lg:w-3/4 mt-10 relative visible h-52 border-solid border-1 border-[#DCDFE680] rounded-2xl bg-white m-auto">
        <textarea
            @keydown="toSend" v-model="data.question"
            class="h-40 w-full rounded-2xl pl-4 resize-none overflow-scroll scrollbar-hidden outline-none
             text-size14 lg:text-size16 xl:text-18 2xl:text-20 pt-2"
            :placeholder="$t('home.inputTip')"/>
        <div class="text-size15 text-default flex">
          <div class="relative flex cursor-pointer visible" :class="{'flex-1':data.searchType===0}"
               @click="data.showPopOne =!data.showPopOne ">
            <svg-icon icon-class="home_point_one" class="w-4 h-4 mt-2.5 ml-4 align-middle"></svg-icon>
            <div class="ml-2 my-1.5">{{ typeList[data.searchType] }}</div>
            <el-icon class="ml-2 mt-2" size="16px">
              <ArrowDown/>
            </el-icon>
            <el-collapse-transition>
              <div v-show="data.showPopOne"
                   class="absolute top-10 w-fit p-2.5 bg-white pop leading-8 z-100 rounded-xl">
                <div class="px-4 mb-2 last:mb-0" v-for="(item,index) in typeList" :key="index"
                     :class="{'search-type-select':data.searchType===index}"
                     @click="setBase(index)">
                  {{ item }}
                </div>
              </div>
            </el-collapse-transition>
          </div>
          <div class="relative flex  ml-4 visible w-fit flex-1" v-show="data.searchType===1">
            <svg-icon
                icon-class="home_point_two"
                @click="showKnowDialog"
                class="w-4 h-4 mt-2.5 align-middle cursor-pointer"></svg-icon>
            <div class="ml-2 my-1.5 w-fit cursor-pointer" @click="showKnowDialog">
              {{ !data.selectKnowledge.kdbId ? $t('home.choice') : data.selectKnowledge.kdbName }}
            </div>
          </div>
          <svg-icon icon-class="home_send" @click="toDetail(data.question)" class="w-8 h-8 mr-4 cursor-pointer"/>
        </div>
      </div>
      <div class="mt-12 shadow-a  bg-white flex w-fit m-auto rounded-full text-default text-size16">
        <div class="pt-2 pb-2 pl-6 pr-6 cursor-pointer no-select"
             :class="{'select-type':data.type==='simple'}"
             @click="data.type='simple'">
          {{ $t('home.simple') }}
        </div>
        <div class="pt-2 pb-2 pl-6 pr-6 cursor-pointer no-select"
             @click="data.type='complex'"
             :class="{'select-type':data.type==='complex'}">
          {{ $t('home.complex') }}
        </div>
        <div class="pt-2 pb-2 pl-6 pr-6 flex cursor-pointer no-select"
             @click="data.type='research'"
             :class="{'select-type':data.type==='research'}">
          <svg-icon :icon-class="data.type==='research'?'study-white':'study'" class="w-4 h-4 mt-1 mr-2"></svg-icon>
          {{ $t('home.research') }}
        </div>
      </div>
      <div class="flex mt-12 m-auto w-11/12 flex-wrap justify-center items-center">
        <div v-for="item in data.caseList"
             :key="item"
             class="item cursor-pointer text-normal bg-white mr-5 border-solid border border-[#DCDFE6]  text-size16 flex rounded-full mb-3 py-1.5 px-4"
             @click="toDetail(item)">
          {{ item }}
          <svg-icon icon-class="home_share" class="w-4 h-4 mt-1 ml-2"></svg-icon>
        </div>
      </div>
    </div>
    <select-knowledge ref="refSelectKnowledge" @selectKnowLedge="selectKnowLedge"></select-knowledge>
  </div>
</template>
<script setup>
import {reactive, onMounted, computed, onUnmounted, ref} from 'vue'
import {useRouter} from 'vue-router'
import {aiSearchCase, aiSearchSession} from "@/http/api/aisearch";
import {useI18n} from 'vue-i18n'
import { getLocal, setLocal} from '@/utils/tools'
import SelectKnowledge from "./components/SelectKnowledge.vue";

const {t} = useI18n()
const refSelectKnowledge = ref()
const router = useRouter()
const data = reactive({
  question: '',
  type: 'simple',
  caseList: [],
  searchType: 0,
  selectKnowledge: {},
  showPopOne: false,
  knowledgeList: []
})

onMounted(() => {
  const searchType = getLocal('searchType')
  if (searchType)
    data.searchType = parseInt(searchType)
  const selectKnowledge = getLocal('selectKnowledge')
  if (selectKnowledge && data.searchType === 1) {
    data.selectKnowledge = JSON.parse(selectKnowledge)
    getCase(data.selectKnowledge.kdbId)
  } else {
    getCase(0)
  }
})

onUnmounted(() => {
  setLocal('searchType', data.searchType)
  setLocal('selectKnowledge', JSON.stringify(data.selectKnowledge))
})

const typeList = computed(() => {
  return [t('home.internet'), t('home.knowledge')]
})

const showKnowDialog = () => {
  refSelectKnowledge.value.showDialog = true
}


const setBase = (index) => {
  data.searchType = index;
  if (data.searchType === 0) {
    getCase(0)
  } else {
    if (data.selectKnowledge && data.selectKnowledge.kdbId)
      getCase(data.selectKnowledge.kdbId)
    else {
      data.caseList = []
    }
  }
}

const getCase = (kdbId) => {
  data.caseList = []
  aiSearchCase({
    kdbId: kdbId
  }).then(res => {
    data.caseList = res.data.cases
  })
}


const selectKnowLedge = (item, index) => {
  data.selectKnowledge = item
  getCase(data.selectKnowledge.kdbId)
}
const toSend = (event) => {
  if (!event.shiftKey && event.keyCode === 13) {
    event.cancelBubble = true; //ie阻止冒泡行为
    event.stopPropagation();//Firefox阻止冒泡行为
    event.preventDefault(); //取消事件的默认动作*换行
    if (event.srcElement.value) {
      toDetail(data.question)
    }
  }
}

const toDetail = async (value) => {
  if (!value) {
    return
  }
  if (data.searchType === 1 && !data.selectKnowledge.kdbId) {
    return
  }
  const res = await aiSearchSession({
    question: value
  })
  router.push({
    path: `/detail/${res.data.sessionId}`,
    query: {
      question: value,
      type: data.type,
      searchEngine: data.searchType === 0 ? 'web' : 'kdb',
      kdbId: data.searchType === 0 ? 0 : data.selectKnowledge.kdbId,
      kdbName: data.searchType === 0 ? '' : data.selectKnowledge.kdbName,
    }
  })
}

</script>


<style scoped lang="less">


.no-select {
  margin: 4px 6px;
}

.select-type {
  background: #0A2540;;
  color: white;
  margin: 4px 6px;
  border-radius: 100px;
}

.shadow-a {
  box-shadow: 0px 0px 6px 0px #0000001A;
}

.pop {
  box-shadow: 0px 0px 4px 0px #0000001A;
}

.item {
  border: 1px solid #DFDFDF
}

.search-type-select {
  background: #7269FB;
  color: white;
  border-radius: 100px;
}

</style>