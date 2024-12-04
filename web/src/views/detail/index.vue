<template>
  <div class="relative h-full ">
    <div class="h-full overflow-y-scroll scrollbar-hidden" ref="refScroll">
      <div class="w-full flex m-auto relative">
        <div class="pl-3 w-3/5 pdf" id="pdf"
             :class="{'w-full':!data.showRight}">
          <question
              :question="data.question"
              :processList="data.processList"
              :loading="data.loading"
              :scrollTop="data.scrollTop"
              @showEditQuestion="showEditQuestion"/>
          <web-progress
              class="mb-16"
              v-if="data.loading"
              :list='progressList'
              :loadingText="data.loadingText"
              :progress="data.progress"
              :analyseing="data.analyseing"
              :detail="data.detail"
              :type="data.type"/>
          <div class="detail mb-10 relative" v-show="data.detail">
            <div class="flex">
              <div class="flex-1 flex text-size18 pt-4 font-[500] leading-10 text-default">
                <div :class="{selected:data.searchType===0}">
                  {{ data.kdbName ? data.kdbName : $t('home.internet') }}
                </div>
              </div>
              <div class="flex leading-10 mt-4 cursor-pointer"
                   v-show="data.referList&&data.referList.length"
                   @click="toFrom">
                <img src="@/assets/img/detail/web_link.png"
                     class="w-4 h-4 mt-3 mr-2"/>
                <div class="text-size16 text-color333 font-[400]">
                  {{ data.referList.length }}
                </div>
              </div>
            </div>
            <model-answer
                class="mt-4"
                :content="data.detail"
                :loading="data.loading"/>
            <div class="flex mt-6 mb-6 text-normal text-size16 font-[400] ">
              <div class="flex cursor-pointer" @click="toShare">
                <img src="@/assets/img/detail/web_share.png" class="w-4 h-4 mt-1">
                <div class="ml-2">{{ $t('detail.share') }}</div>
              </div>

              <div v-if="data.type==='simple'" class="flex ml-4 cursor-pointer" @click="toSearch">
                <img src="@/assets/img/detail/web_deep.png" class="w-4 h-4 mt-1">
                <div class="ml-2">{{ $t('home.complex') }}</div>
              </div>
              <el-tooltip
                  class="box-item"
                  effect="dark"
                  v-if="data.type==='complex'"
                  content="研发中..."
                  placement="bottom">
                <div class="flex ml-4 cursor-pointer">
                  <img src="@/assets/img/detail/web_deep.png" class="w-4 h-4 mt-1">
                  <div class="ml-2">{{ $t('home.research') }}</div>
                </div>
              </el-tooltip>
              <div class="flex ml-4 cursor-pointer" @click="exportPdf">
                <img src="@/assets/img/detail/web_export.png" class="w-4 h-4 mt-1">
                <div class="ml-2">{{ $t('detail.export') }}</div>
              </div>
              <div class="flex ml-4 cursor-pointer" @click="refFeedBack.showDialog = true">
                <img v-if="!data.unlike" src="@/assets/img/detail/no_like_no.png" class="w-4 h-4 mt-1">
                <img v-else src="@/assets/img/detail/no_like_select.png" class="w-4 h-4 mt-1">
                <div class="ml-2">{{ $t('detail.noLike') }}</div>
              </div>
              <div class="flex-1"></div>
            </div>
            <div class="border-b border-solid  border-[#DCDFE6]"></div>
            <dependent-event v-if="data.eventsInfo.length" :eventsInfo="data.eventsInfo"></dependent-event>
            <organization v-if="data.orgsInfo.length" :orgsInfo="data.orgsInfo"></organization>
            <related-personnel v-if="data.peopleInfo.length" :peopleInfo="data.peopleInfo"></related-personnel>
            <from
                v-if="data.referList.length"
                :show-from="data.showFrom"
                :kdbId="data.kdbId"
                :authType="data.authType"
                @showFrom="(value)=>{data.showFrom=value}"
                :referList="data.referList"></from>
          </div>
        </div>
        <div class="w-2/5" v-show="data.showRight">
          <out-line :showOutLine="data.showOutLine" :outLineList="data.outLineList" @hide="data.showRight=false"/>
        </div>
      </div>
    </div>
    <div class="absolute right-0 top-[20%] border border-b pl-3"
         style="z-index: 10000"
         v-show="!data.showRight&&!data.loading"
         @click="data.showRight=true">
      <svg-icon icon-class="detail_show_menu" class="w-8 h-8 cursor-pointer"></svg-icon>
    </div>
    <feedback ref="refFeedBack" :sessionId="route.params.id" @setNoLike="data.unlike = true"></feedback>
    <edit-question ref="refEditQuestion" @send="send" :name="data.kdbName" :loading="data.loading"></edit-question>
  </div>
</template>
<script setup>
import {reactive, onMounted, ref, computed, watch, onUnmounted} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import WebProgress from "./components/WebProgress.vue";
import Feedback from "./components/Feedback.vue";
import Question from "./components/Question.vue";
import From from "./components/From.vue";
import ModelAnswer from "./components/ModelAnswer.vue";
import OutLine from "./components/OutLine.vue";
import DependentEvent from "./components/DependentEvent.vue";
import Organization from "./components/Organization.vue";
import RelatedPersonnel from "./components/RelatedPersonnel.vue";
import EditQuestion from "./components/EditQuestion.vue";
import {
  aiSearchHis,
  aiSearchOutLine,
  aiSearchRefer,
  aiSearchSession,
  relation,
  searchProcess
} from "@/http/api/aisearch";
import {copyToClip} from "@/utils/tools";
import {ElMessage} from "element-plus";
import {getPdf} from "@/utils/exportPdf";
import {useI18n} from 'vue-i18n'
import {isEqual, uniqWith} from 'lodash'
import {knowledgeAuth} from "@/http/api/knowledge";

const {t} = useI18n()
const refFeedBack = ref()
const route = useRoute()
const router = useRouter()
const refEditQuestion = ref('')
const refScroll = ref()
const data = reactive({
  showFrom: true,
  loading: true,
  searchType: 0,
  showRight: true,
  question: '',
  detail: '',
  referList: [],
  peopleInfo: [],
  eventsInfo: [],
  orgsInfo: [],
  showOutLine: false,
  progress: 0,
  timer: null,
  sessionId: '',
  loadingText: '搜索中...',
  analyseing: false,
  type: 'simple',
  outLineList: [],
  kdbName: '',
  kdbId: '0',
  unlike: false,
  authType: -1,
  processList: [],
  scrollTop: 0,
})

onMounted(() => {
  loadData()
  if (refScroll.value)
    refScroll.value.addEventListener("scroll", handleScroll)
})


onUnmounted(() => {
  if (refScroll.value)
    refScroll.value.removeEventListener("scroll", handleScroll)
})

const handleScroll = () => {
  data.scrollTop = refScroll.value.scrollTop;
}

watch(() => route.params.id, (val) => {
  initData()
  loadData()
})

const initData = () => {
  data.showFrom = true
  data.loading = true
  data.showRight = true
  data.searchType = 0
  data.question = ''
  data.detail = ''
  data.referList = []
  data.showOutLine = false
  data.progress = 0
  data.timer = null
  data.sessionId = ''
  data.loadingText = '搜索中...'
  data.analyseing = false
  data.outLineList = []
  data.kdbName = ''
  data.kdbId = '0'
  data.unlike = false
  data.authType = -1
  data.peopleInfo = []
  data.eventsInfo = []
  data.orgsInfo = []
  data.processList = []
}

const progressList = computed(() => {
  if (data.type === 'simple') {
    return [data.kdbId !== "0" ? data.kdbName : t('detail.allNet'), t('detail.organizeAnswer'), t('detail.complete')]
  } else {
    return [
      t('detail.problemAnalysis'),
      data.kdbId !== "0" ? data.kdbName : t('detail.allNet'),
      t('detail.organizeAnswer'),
      t('detail.complete')
    ]
  }
})

const loadData = async () => {
  data.question = route.query.question
  data.sessionId = route.params.id
  data.type = route.query.type
  data.kdbName = route.query.kdbName
  data.kdbId = route.query.kdbId
  if (data.kdbId !== '0') {
    loadAuth()
  }
  const res = await aiSearchHis({
    sessionId: data.sessionId
  })
  if (res.data.question) {
    data.loading = false;
    if (!res.data.result) {
      ElMessage.warning('后台生成中,请稍后刷新')
      return
    }
    data.showOutLine = true
    data.question = res.data.question
    data.unlike = res.data.unlike
    loadOutLine()
    loadRelation()
    loadProcessList()
    const resultRefers = res.data.resultRefers === "null" ? '[]' : res.data.resultRefers
    const result = res.data.result
    aiSearchRefer({
      sessionId: data.sessionId
    }).then(res => {
      data.referList = res.data.list ? res.data.list : []
      data.detail = setData(resultRefers, result)
    })
  } else {
    await search(data.question)
  }
}

const loadProcessList = () => {
  searchProcess({
    sessionId: data.sessionId
  }).then(res => {
    const processList = res.data.list
    const mergedData = [];
    const proData = {
      analyze: "问题分析",
      webSearch: "全网搜索",
      vdbSearch: "知识库搜索",
      summary: "整理答案",
      finish: "回答完成"
    }
    for (let item in proData) {
      const list = processList.filter(pro => pro.stageType === item)
      if (list.length) {
        const proItem = {name: proData[item], list: list}
        mergedData.push(proItem)
      }
    }
    data.processList = mergedData
  })
}

const loadAuth = () => {
  knowledgeAuth({kdbId: parseInt(data.kdbId)}).then(res => {
    data.authType = res.data.authType
  })
}

const setData = (resultRefers, result) => {
  let resData = ''
  const referList = uniqWith(resultRefers ? JSON.parse(resultRefers)
      .filter(item => item.refers.length)
      .map(item => {
        let start = item.start
        const arr = extractMarkdownSymbols(result.substring(start, start + 8))
        if (arr.length) {
          const number = (arr[0].type === 'list' ? 3 : 1)
          const length = arr[0].length
          start = start + length + number
        }
        return {
          number: item.numberIndex,
          start: start,
          end: item.end,
          refers: item.refers
        }
      }) : [], isEqual)
  if (referList.length > 0) {
    const arr = result.split('')
    referList.forEach((item, index) => {
      const aa = item.refers.map(tem => {
        return `<span><a target="_blank" ${getHref(tem.index)}>${tem.index + 1}<span>${getContent(tem)}</span></a></span>`
      })
      arr.splice(item.start + index * 3, 0, '[[')
      arr.splice(item.number + index * 3 + 1, 0, aa.join(''));
      arr.splice(item.number + index * 3 + 2, 0, ']]');
    })
    console.log(arr)
    resData = arr.join('')
  } else {
    return result
  }
  return resData
}

const getContent = (tem) => {
  return data.referList[tem.index].content
      .substring(tem.referStart, tem.referEnd)
      .replace(/\n/g, "")
      .replace(/\r/g, "")
}

function extractMarkdownSymbols(text) {
  const patterns = {
    header: /^(#{1,6}) /gm,
    // bold: /(\*\*|__).*?\1/g,
    // italic: /(\*|_).*?\1/g,
    // link: /(\[.*?\]\()(.*?\))/g,
    // image: /(!\[.*?\]\()(.*?\))/g,
    blockquote: /^\s*>+\s*(.*)/gm,
    inlineCode: /(`)[^`]+\1/g,
    codeBlock: /(```)[\s\S]*?\1/g,
    list: /^[ \t]*([-*+]|(\d+\.)) +.*$/gm,
    // listw:/^[ \t]*[-*+][ \t]+.*/gm,
    // listy:/^[ \t]*\d+\.[ \t]+.*/gm
  };
  const result = [];
  for (const [type, pattern] of Object.entries(patterns)) {
    const matches = text.matchAll(pattern);
    for (const match of matches) {
      result.push({
        type: type,
        start: match.index,
        length: match[1].length
      });
    }
  }
  return result;
}

const getHref = (index) => {
  if (index > data.referList.length - 1) {
    return ''
  }
  return route.query.searchEngine === 'web'
      ? `href='${data.referList[index].url}'`
      // : (data.authType === 3 || data.authType === 9 ? `href='https://askonce.atomecho.cn/online-file?url=${data.referList[index].url}'` : "")
      : (data.authType === 3 || data.authType === 9 ? `href='${data.referList[index].url}'` : "")
}


const send = async (value) => {
  const res = await aiSearchSession({
    question: data.inputValue
  })
  await router.push({
    path: `/detail/${res.data.sessionId}`,
    query: {
      question: value,
      type: data.type,
      kdbName: data.kdbName,
      kdbId: data.kdbId,
      searchEngine: route.query.searchEngine,
    }
  })
}

const toSearch = async () => {
  const res = await aiSearchSession({
    question: data.inputValue
  })
  router.push({
    path: `/detail/${res.data.sessionId}`,
    query: {
      question: data.question,
      type: 'complex',
      kdbName: data.kdbName,
      kdbId: data.kdbId,
      searchEngine: route.query.searchEngine,
    }
  })
}

const showEditQuestion = () => {
  refEditQuestion.value.inputValue = data.question
  refEditQuestion.value.showDialog = true;
}

const search = async (question) => {
  if (data.timer) {
    clearInterval(data.timer)
  }
  data.timer = setInterval(() => {
    if (data.progress > 66) {
      if (data.progress >= 99) {
        data.progress = 99
      } else {
        data.progress += 0.1
      }
    } else {
      data.progress += 0.1
    }
  }, 100)
  data.loading = true
  data.analyseing = true
  const response = await fetch('/serverApi/askonce/search/ask', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    dataType: "text/event-stream",
    body: JSON.stringify({
      question,
      sessionId: data.sessionId,
      type: data.type,
      searchEngine: route.query.searchEngine,
      kdbId: data.kdbId ? parseInt(data.kdbId) : 0
    }),
  });

  if (!response.ok) {
    ElMessage.error("接口错误")
    data.progress = 100
    data.loading = false
    clearInterval(data.timer)
    console.log('err-->', response)
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  const reader = response.body.getReader();
  let decoder = new TextDecoder();
  let result = true;
  let referText = ''
  let referList = ''
  let partialLine = "";
  let jsonLine = ''

  while (result) {
    const {done, value} = await reader.read();
    if (done) {
      console.log("Stream ended");
      data.loading = false
      result = false;
      clearInterval(data.timer)
      loadProcessList()
      break;
    }
    const decodedText = decoder.decode(value, {stream: true});
    const chunk = partialLine + decodedText;
    const newLines = chunk.split(/\r?\n/);
    partialLine = newLines.pop() ?? "";
    for (const line of newLines) {
      if (line.startsWith("data:")) {
        jsonLine = line.substring(5);
        try {
          const value = JSON.parse(jsonLine)
          if (value.message && value.code !== 200) {
            ElMessage.error(value.message)
            clearInterval(data.timer)
            data.loading = false
            return
          }
          switch (value.stage) {
            case 'analyze':
              data.progress = 1
              data.loadingText = value.text ? value.text : t('detail.problemAnalysis')
              break
            case "search":
              data.progress = 35
              data.loadingText = value.text ? value.text : (data.kdbId !== "0" ? data.kdbName : t('detail.allNet'))
              break
            case "generate":
              data.progress = 68
              data.loadingText = value.text ? value.text : t('detail.organizeAnswer')
              loadRefer()
              break
            case "appendText":
              data.analyseing = false
              data.detail += value.text
              referText += value.text
              break
            case "outline":
              loadOutLine()
              break
            case "relation":
            case 'relate':
              loadRelation();
              break
            case "refer":
              referList = value.text
              data.detail = setData(value.text, referText)
              break
              //重新调用引用
            case "refreshSearch":
              loadRefer()
              break
            case "complete":
              data.progress = 100
              referText = value.text
              data.loading = false
              console.log('-------', 'finish')
              clearInterval(data.timer)
              data.detail = setData(referList, referText)
              break
            default:
              break
          }
        } catch (e) {
          console.log('----err---', jsonLine)
        }
      }
    }

  }
}

const loadOutLine = () => {
  aiSearchOutLine({sessionId: data.sessionId}).then(res => {
    data.showOutLine = true
    data.outLineList = res.data.list
  })
}

const loadRelation = () => {
  relation({sessionId: data.sessionId}).then(res => {
    data.orgsInfo = res.data.orgsInfo ?? []
    data.eventsInfo = res.data.eventsInfo ?? []
    data.peopleInfo = res.data.peopleInfo ?? []
  })
}

const loadRefer = () => {
  aiSearchRefer({
    sessionId: data.sessionId
  }).then(res => {
    data.referList = res.data.list ? res.data.list : []
  })
}

const toFrom = () => {
  document.getElementById('from').scrollIntoView({behavior: 'smooth'})
}

const toShare = () => {
  copyToClip(window.location.href)
  ElMessage.success('分享链接已成功复制到剪贴板')
}

const exportPdf = () => {
  ElMessage.success('导出中，请稍后...')
  console.log(document.querySelector('#pdf'))
  getPdf(data.question, document.querySelector('#pdf'))
}

</script>
<style scoped lang="less">
:deep(.el-skeleton__paragraph) {
  background: #D9D9D9;
  height: 10px;
}

:deep(.el-skeleton__item) {
  background: #D9D9D9;
  height: 10px;
}

:deep(.el-skeleton__p.is-first) {
  width: 100%;
}

.selected {
  border-bottom: 2px solid #7269FB;
  color: #7269FB;
}
</style>