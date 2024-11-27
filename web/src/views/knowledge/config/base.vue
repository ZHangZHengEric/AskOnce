<template>
  <div class="w-full pt-4">
    <el-form
        :model="formData"
        label-position="left"
        label-suffix=":"
        class="flex flex-col h-full"
        label-width="100px">
      <el-form-item :label="$t('knowledge.name')">
        <el-input v-model="formData.name"></el-input>
      </el-form-item>
      <el-form-item :label="$t('knowledge.introduction')">
        <el-input type="textarea" rows="5" v-model="formData.intro"></el-input>
      </el-form-item>
      <div id="image" class="flex-1">
        <el-form-item :label="$t('knowledge.setCover')" class="overflow-hidden">
          <div class="flex bg-white rounded-xl p-4 h-full w-full overflow-hidden">
            <img :src="formData.cover" :style="{height:data.imageHeight - 60 +'px'}" class="rounded-xl">
            <div class="flex-1 ml-8" :style="{height:data.imageHeight - 80 +'px'}">
              <div class="top flex text-size16 font-normal">
                <div class="px-4 mr-4 cursor-pointer"
                     v-for="(item,index) in typeList"
                     :key="index"
                     :class="{selected:data.type===index}"
                     @click="selectType(index)">
                  {{ item.name }}
                </div>
                <div class="flex-1"></div>
                <div class="flex cursor-pointer" @click="choiceRandom">
                  <svg-icon icon-class="knowledge_random" class="w-4 h-4 mt-2"></svg-icon>
                  <div class="ml-2 text-size16 font-normal text-default">{{ $t('knowledge.randomCover') }}</div>
                </div>
              </div>
              <div class="overflow-y-scroll h-full pr-4 mb-10 pb-6 mt-4">
                <div class="mb-10">
                  <el-row :gutter="32">
                    <el-col :xs="12" :sm="12" :md="6" :lg="6" :xl="4"
                            class="mb-8"
                            v-for="(item,index) in data.coverList"
                            :key="item.url"
                            @click="itemCLick(item,index)">
                      <img
                          v-lazy="item.url"
                          class="w-full cursor-pointer rounded-xl border-2 border-solid border-transparent"
                          :class="{'select-cover':data.selectCoverIndex===index}">
                    </el-col>
                  </el-row>
                </div>
              </div>
            </div>
          </div>
        </el-form-item>
      </div>
      <div class="flex">
        <div class="flex-1"></div>
        <div class="px-8 bg-default rounded-full py-1 text-white cursor-pointer" @click="toSave">{{ $t('app.save') }}
        </div>
      </div>
    </el-form>
  </div>
</template>
<script setup>
import {reactive, onMounted, toRefs, computed} from 'vue'
import {knowledgeCovers, knowledgeDetail, knowledgeUpdate} from "@/http/api/knowledge";
import {useRoute, useRouter} from 'vue-router'
import {ElMessage} from "element-plus";
import {useI18n} from 'vue-i18n'
import {cloneDeep} from 'lodash'

const route = useRoute()
const router = useRouter()
const {t} = useI18n()

const data = reactive({
  imageHeight: 0,
  selectCoverIndex: -1,
  formData: {},
  type: 0,
  coverAllList: [],
  coverList: []
})

const typeList = computed(() => {
  return [
    {
      name: t('knowledge.all'),
      type: 'all'
    },
    {
      name: t('knowledge.colour'),
      type: 'color'
    },
    {
      name: t('knowledge.bussiness'),
      type: 'office'
    },
    {
      name: t('knowledge.technology'),
      type: 'science'
    },
    {
      name: t('knowledge.scenery'),
      type: 'scenery'
    }
  ]
})

const itemCLick = (item, index) => {
  data.selectCoverIndex = index
  data.formData.cover = item.url
  data.formData.coverId = item.id
}

const choiceRandom = () => {
  data.selectCoverIndex = parseInt(Math.random() * data.coverList.length)
  data.formData.cover = data.coverList[data.selectCoverIndex].url
  data.formData.coverId = data.coverList[data.selectCoverIndex].url.id
}

const selectType = (type) => {
  data.type = type
  data.selectCoverIndex = -1
  if (type === 0) {
    data.coverList = cloneDeep(data.coverAllList)
  } else {
    data.coverList = data.coverAllList.filter(item => item.type === typeList.value[type].type)
  }
}
const loadData = () => {
  knowledgeDetail({
    kdbId: route.query.id
  }).then(res => {
    data.formData = res.data
  })
}

const toSave = () => {
  knowledgeUpdate(data.formData).then(res => {
    ElMessage.success('修改成功')
    router.back()
  })
}

const loadCovers = () => {
  knowledgeCovers({type: ''}).then(res => {
    data.coverAllList = res.data.list
    data.coverList = cloneDeep(data.coverAllList)
  })
}

onMounted(() => {
  data.imageHeight = document.getElementById('image').clientHeight
  loadData()
  loadCovers()
})


const {formData} = toRefs(data)
</script>


<style scoped lang="less">

/deep/ .el-textarea__inner {
  resize: none;
}

.selected {
  color: white;
  background: #7269FB;
  border-radius: 100px;
}

.select-cover {
  border: 2px solid #7269FB;
  border-radius: 12px;
}

</style>


