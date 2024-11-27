/**
* create by Administrator on 2023/10/13 9:13
* 类注释：
* 备注：
*/
<template>
  <div class="relative">
    <mavon-editor
        class="text-default relative"
        ref="refEditor"
        v-model="detail"
        v-show="detail"
        previewBackground="transparent"
        :boxShadow="false"
        defaultOpen="preview"
        :subfield="false"
        :transition="false"
        :toolbarsFlag="false"/>
  </div>
</template>

<script setup>
import {defineProps, computed, ref} from 'vue'

const refEditor = ref(null)
import 'mavon-editor/dist/css/index.css'
import {mavonEditor} from 'mavon-editor'

const props = defineProps({
  content: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  }
})


const detail = computed(() => {
  return props.content
//   return `| Syntax      | Description  |
// | ----- | ----- |
// | Header      | Title       |
// | Paragraph   | Text        |`
})
</script>

<style scoped lang="less">

:deep(.v-note-wrapper) {
  border: none !important;
  min-height: 0;
  min-width: 30px;
  background: transparent;
}

:deep(.v-note-show) {
  position: relative !important;
}

/deep/ .v-note-wrapper .v-note-panel .v-note-show .v-show-content, .v-note-wrapper .v-note-panel .v-note-show .v-show-content-html {
  padding: 0 !important;
  line-height: 1.8 !important;
  overflow: hidden;
}


/deep/ .markdown-body blockquote, .markdown-body dl, .markdown-body ol, .markdown-body pre, .markdown-body table, .markdown-body ul {
  margin-bottom: 0px !important;
}

/deep/ .markdown-body p {
  margin-bottom: 10px !important;
}

/deep/ ol, ul, menu {
  list-style: auto !important;
}

/deep/ .markdown-body span a {
  position: relative;
  overflow: visible; /* 让内容超出父元素布局 */
  clip: auto;
  display: inline-block;
  text-decoration: none;
  color: #7269FB;
  line-height: 20px;
  width: 20px;
  margin-right: 10px;
  text-align: center;
  font-size: 12px;
}

/deep/ .markdown-body span a:hover {
  color: white;
  background: #7269FB;
  border-radius: 50%; /* 将伪元素的边框半径设置为50%，使背景呈现圆形 */
}

/deep/ .markdown-body span a::before {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  width: 20px;
  height: 20px;
  background-color: #7269FB1A;
  border-radius: 50%; /* 将伪元素的边框半径设置为50%，使背景呈现圆形 */
  z-index: -1; /* 将伪元素放在最底层，使其在内容之后 */
}

/deep/ .markdown-body span:hover {
  border-bottom: 1px dashed #7269FB;
}

/deep/ .markdown-body span a span {
  visibility: hidden;
  max-width: 300px;
  min-width: 250px;
  max-height: 500px;
  overflow-y: scroll;
  background: rgba(0, 0, 0, 0.7);
  color: #fff;
  padding: 10px;
  border-radius: 10px;
  text-align: left;
  position: absolute;
  z-index: 10;
  top: 30px;
  left: 0;
  right: 0;
}

/deep/ .markdown-body span a:hover span {
  visibility: visible;
}

/deep/ .v-note-wrapper .v-note-panel .v-note-show {
  overflow: visible;
}

/deep/ .v-note-wrapper .v-note-panel {
  overflow: visible;
}

/deep/ .v-note-wrapper .v-note-panel .v-note-show .v-show-content {
  overflow: visible !important;
}

/deep/ .v-note-wrapper .v-note-panel .v-note-show .v-show-content-html {
  overflow: visible !important;
}

/deep/ .fade-in {
  animation: fadeInAnimation ease 2s;
  animation-iteration-count: 1; /*设置动画播放次数*/
  animation-fill-mode: forwards; /*设置样式以在动画不播放时应用元素。forward是设置动画结束后，使用元素的结束属性值*/
}

@keyframes fadeInAnimation {

  0% {
    opacity: 0; /*设置不透明度*/
  }
  100% {
    opacity: 1;
  }
}

.content {
  position: absolute;
  width: 15px;
  height: 15px;
  border-radius: 50%;
  background: #333;
  animation: blink 0.8s infinite;
  opacity: 0;
  z-index: 1600;
}

@keyframes blink {
  from {
    opacity: 0;
  }

  to {
    opacity: 1;
  }
}
</style>