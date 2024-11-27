/**
* create by zhangxiang on 2023-06-19 15:12
* 类注释：
* 备注：
*/
<template>
  <section class="dots-container">
    <div class="dot"></div>
    <div class="dot"></div>
    <div class="dot"></div>
    <div class="dot"></div>
    <div class="dot"></div>
  </section>
</template>

<script setup>

</script>

<style scoped>
.dots-container {
  padding-top: 10px;
  display: flex;
  padding-left: 10px;
  height: 30px;
}

.dot {
  height: 10px;
  width: 10px;
  margin-right: 5px;
  border-radius: 5px;
  background-color: #b3d4fc;
  animation: pulse 1.5s infinite ease-in-out;
}

.dot:last-child {
  margin-right: 0;
}

.dot:nth-child(1) {
  animation-delay: -0.3s;
}

.dot:nth-child(2) {
  animation-delay: -0.1s;
}

.dot:nth-child(3) {
  animation-delay: 0.1s;
}

@keyframes pulse {
  0% {
    transform: scale(0.8);
    background-color: #b3d4fc;
    box-shadow: 0 0 0 0 rgba(178, 212, 252, 0.7);
  }

  50% {
    transform: scale(1.2);
    background-color: #6793fb;
    box-shadow: 0 0 0 10px rgba(178, 212, 252, 0);
  }

  100% {
    transform: scale(0.8);
    background-color: #b3d4fc;
    box-shadow: 0 0 0 0 rgba(178, 212, 252, 0.7);
  }
}
</style>