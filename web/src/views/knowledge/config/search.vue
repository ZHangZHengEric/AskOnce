<template>
  <div class="h-full w-full overflow-y-scroll scrollbar-hidden pt-20 flex justify-center">
    <div class="w-11/12 lg:w-1/2">
      <div class="text-center text-size32 font-[600]">AskOnce</div>
      <el-input
          size="large"
          class="rounded-xl mt-10"
          @keyup.enter="toSearch"
          v-model="data.question"
          :placeholder="$t('app.placeholder')">
        <template #prefix>
          <el-icon class="el-input__icon text-themeBlue">
            <search/>
          </el-icon>
        </template>
      </el-input>
      <Loading v-show="data.loading" class="mt-4"></Loading>
      <div class="text-size20 text-default mt-6" v-show="data.list.length">搜索结果</div>
      <div class="list mt-4">
        <div v-for="item in data.list" :key="item.searchContent"
             class="bg-[#FFFFFF66] mb-4 rounded-xl item py-4 px-7 text-default">
          <div class="text-size16 font-normal">{{ item.dataName }}</div>
          <div class="text-size14 text-normal mt-5">
            <TextClamp :text='item.searchContent' :max-lines='2' auto-resize class="relative leading-6">
              <template #after="{ toggle, expanded, clamped }">
                <!--                <button v-if="expanded " class="text-themeBlue ml-4 absolute bottom-0 right-2" @click="toggle">-->
                <button v-if="expanded " class="text-themeBlue ml-4" @click="toggle">
                  收起
                  <el-icon class="ml-1">
                    <ArrowUp/>
                  </el-icon>
                </button>
                <button v-if="clamped " class="text-themeBlue ml-4" @click="toggle">
                  展开
                  <el-icon class="ml-1 mt-0.5">
                    <ArrowDown/>
                  </el-icon>
                </button>
              </template>
            </TextClamp>
          </div>
        </div>
      </div>
      <div class="h-4"></div>
    </div>
  </div>
</template>
<script setup>
import {reactive} from 'vue'
import {useRoute} from 'vue-router'
import {knowledgeSearch} from "@/http/api/knowledge";
import TextClamp from 'vue3-text-clamp';
import Loading from '@/components/Loading/index.vue'
import {ArrowRight} from "@element-plus/icons-vue";

const route = useRoute()
const data = reactive({
  list: [],
  question: '',
  loading: false
})

const toSearch = () => {
  if (!data.question) {
    return
  }
  data.loading = true
  data.list = []
  knowledgeSearch({query: data.question, kdbId: parseInt(route.query.id)}).then(res => {
    data.list = res.data.list
    data.loading = false
  })
}
</script>


<style scoped lang="less">
.item {
  box-shadow: 0px 0px 4px 0px #00000014;

}

</style>