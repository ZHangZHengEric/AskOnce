<template>
  <el-dialog
      v-model="data.showDialog"
      width="600px"
      draggable
      :style="{padding:0}"
      class="rounded-xl"
      custom-class="my-dialog">
    <div class="feed-bg px-6 py-6">
      <div class="text-size22 font-normal text-default">
        您的反馈将帮助我们优化进步
      </div>

      <div class="mt-6">
        <el-checkbox-group v-model="data.select">
          <el-row>
            <el-col v-for="item in data.list" :span="8" :key="item.label">
              <el-checkbox :label="item.label" :value="item.label" class="mb-3"
                           style="--el-checkbox-checked-text-color: #7269FB;
                           --el-checkbox-checked-input-border-color:#7269FB;
                            --el-checkbox-checked-bg-color:#7269FB"/>
            </el-col>
          </el-row>
        </el-checkbox-group>
      </div>
      <div class="mt-2">
        <el-input v-model="data.other" resize="none" type="textarea" :rows="5" :maxlength="512" show-word-limit></el-input>
      </div>
      <div class="mt-6 flex justify-center">
        <div class="px-8 py-1.5 bg-default text-white rounded-full cursor-pointer" @click="commit">提交</div>
      </div>
    </div>

  </el-dialog>
</template>

<script setup>
import {reactive, defineExpose, toRefs, onMounted, defineProps, defineEmits} from 'vue'
import {unlike} from "@/http/api/aisearch";
import {ElMessage} from "element-plus";

const data = reactive({
  showDialog: false,
  other: '',
  select: [],
  list: [
    {label: '来源与问题不相关', value: ''},
    {label: '回答格式混乱', value: ''},
    {label: '答非所问', value: ''},
    {label: '内容不完整', value: ''},
    {label: '内容重复', value: ''},
    {label: '内容存在幻觉', value: ''},
  ]
})
const {showDialog} = toRefs(data)
const emits = defineEmits(['setNoLike'])
const props = defineProps({
  sessionId: {
    type: String,
    default: ''
  }
})

onMounted(() => {
})

const commit = () => {
  let dataselect = []
  if (data.other) {
    dataselect = [...data.select, data.other]
  } else {
    dataselect = data.select
  }
  unlike({
    sessionId: props.sessionId,
    reasons: dataselect
  }).then(res => {
    ElMessage.success('已反馈')
    data.unlike = true
    data.showDialog = false
    emits('setNoLike')
  })
}

defineExpose({
  showDialog
})
</script>


<style scoped lang="less">
/deep/ .el-dialog {
  padding: 0 !important;
}

.my-dialog {
  padding: 0 !important;
}

.feed-bg {
  background-image: url("@/assets/img/detail/feedback_bg.png");
  background-size: 100% 100%;
  background-repeat: no-repeat;
}

</style>