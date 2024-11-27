<template>
  <div>
    <router-view v-if="isRouterAlive"></router-view>
  </div>
</template>

<script setup>
import {ref, provide, nextTick, onMounted} from 'vue';
import {useRouter} from 'vue-router'

const route = useRouter()
const isRouterAlive = ref(true);

provide("reload", () => {
  isRouterAlive.value = false;
  nextTick(() => {
    isRouterAlive.value = true;
  })
});
</script>

<style scoped lang="less" rel="stylesheet/less">
.app {
  height: 100vh;
}
</style>
