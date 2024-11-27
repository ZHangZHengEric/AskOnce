/**
* create by zhangxiang on 2023-06-20 13:07
* 类注释：
* 备注：
*/
<template>
  <div class="content">
    <svg-icon class="loading" icon-class="loading"></svg-icon>
    <div style="margin-top: 10px">加载中...</div>
  </div>
</template>

<script setup>

</script>

<style scoped lang="less">
.content {
  text-align: center;

  .loading {
    animation: rotate 1s linear infinite;
  }
}


@keyframes rotate {
  0% {
    -webkit-transform: rotate(0deg);
  }
  100% {
    -webkit-transform: rotate(360deg);
  }
}

</style>