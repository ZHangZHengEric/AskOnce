<template>
  <el-container class="overscroll-none">
    <el-aside class="aside-menu drawer-transition" :width="data.menuWidth+'px'" ref="refAside">
      <div class="menu relative">
        <div class="relative">
          <div class="text-center text-default mt-8 text-size32 font-[400] logo" v-if="data.showMenu">AskOnce</div>
          <div class="text-center text-default mt-8 text-size32 font-[400] logo" v-else>A</div>
          <div class="absolute right-3 top-3 cursor-pointer" @click="hideMenu" v-show="data.showMenu">
            <svg-icon icon-class="menu_hide" class="w-6 h-6"></svg-icon>
          </div>
        </div>
        <div class="mt-10">
          <div :class="[{select:data.selectIndex===index},{'px-8':data.showMenu}]"
               @click="toPage(item,index)"
               class="menu-item flex text-size18 cursor-pointer leading-10 h-12 mb-3 py-1"
               v-for="(item,index) in menuList" :key="item.name">
            <svg-icon :icon-class="item.icon" class="w-4 h-4 mt-3"
                      :class="[{'ml-6':!data.showMenu}]"></svg-icon>
            <span class="ml-4" v-if="data.showText">{{ item.name }}</span>
          </div>
        </div>
      </div>
      <div class="absolute inset-y-2/4 left-4  cursor-pointer" @click="showMenu" v-show="!data.showMenu">
        <svg-icon icon-class="menu_hide" class="w-6 h-6 rotate-180"></svg-icon>
      </div>
    </el-aside>
    <el-container class="h-screen container-main relative">
      <el-header>
        <div class="flex text-size16 mt-4 cursor-pointer">
          <div class="flex-1"></div>
          <div class="relative">
            <div @click="data.showLanguage = !data.showLanguage" v-clickOutside="closeAll">
              {{ locale === 'zh-cn' ? $t('app.Chinese') : $t('app.English') }}
            </div>
            <div
                class="absolute left-1/2 top-10 -translate-x-1/2 w-max bg-white lg:rounded-lg p-1 lg:p-3 z-[9999]"
                v-show="data.showLanguage">
              <div @click="changeLanguage('zh-cn')">
                简体中文
              </div>
              <div class="mt-2" @click="changeLanguage('en-us')">
                English
              </div>
            </div>
          </div>
          <div class="w-fit ml-4 relative">
            <div class="flex" @click="data.showUser=!data.showUser" v-clickOutside="closeAll">
              <svg-icon icon-class="user_info" class="w-4 h-4 mt-1"></svg-icon>
              <div class="ml-2">{{ account }}</div>
              <el-icon class="mt-1 ml-2">
                <ArrowDown/>
              </el-icon>
            </div>
            <div
                class="absolute left-1/2 top-10 -translate-x-1/2 w-max bg-white lg:rounded-lg p-1 lg:p-3 z-[9999]"
                v-show="data.showUser">
              <div @click="loginOut" class="px-4">退出</div>
            </div>
          </div>
        </div>
      </el-header>
      <el-main class="px-12">
        <router-view/>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import {reactive, watch, computed, onMounted, ref} from 'vue'
import {useRouter, useRoute} from 'vue-router'
import {useI18n} from 'vue-i18n'
import {useUserStore} from '@/store'
import {ClickOutside as vClickOutside} from 'element-plus'

import {clearSession, getLanguage, setLanguage} from "@/utils/tools";
import SvgIcon from "@/components/SvgIcon/index.vue";

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const {t, locale} = useI18n();
onMounted(() => {
  if (getLanguage()) {
    locale.value = getLanguage()
  }
})

const refAside = ref()
const data = reactive({
  selectIndex: 0,
  showLanguage: false,
  showUser: false,
  menuWidth: 280,
  showMenu: true,
  showText:true
})

const account = computed(() => {
  return userStore.getInfo.account
})

const hideMenu = () => {
  data.menuWidth = 64
  data.showMenu = false
  data.showText=false
}
const showMenu = () => {
  data.menuWidth = 280
  data.showMenu = true
  setTimeout(()=>{
    data.showText=true
  },200)
}

const menuList = computed(() => {
  return [
    {
      name: t('app.home'),
      icon: 'home',
      url: '/'
    },
    {
      name: t('app.knowledge'),
      icon: 'know',
      url: '/knowledge-manage'
    },
    // {
    //   name: t('app.searchConfig'),
    //   icon: 'search',
    //   url: '/search-config'
    // },
    {
      name: t('app.history'),
      icon: 'history',
      url: '/history'
    },
  ]
})
const closeAll = () => {
  data.showLanguage = false;
  data.showUser = false
}
const changeLanguage = (lang) => {
  setLanguage(lang)
  locale.value = lang
  data.showLanguage = false
}

const loginOut = () => {
  clearSession()
  router.push({path: "/login"})
}

watch(() => route.path, (newPath, oldPath) => {
  if (newPath.indexOf('knowledge') > -1) {
    data.selectIndex = 1
  } else if (newPath.indexOf('history') > -1) {
    data.selectIndex = 3
  } else if (newPath.indexOf('search') > -1) {
    data.selectIndex = 2
  } else {
    data.selectIndex = 0
  }

}, {immediate: true});
const toPage = (item, index) => {
  data.selectIndex = index
  router.push({path: item.url})
}
</script>

<style scoped lang="less">
.aside-menu {
  background-image: url("@/assets/img/layout/menu.png");
  background-size: 100% 100%;
  background-repeat: no-repeat;
}

.container-main {
  background-image: url("@/assets/img/layout/main.png");
  background-size: 100% 100%;
  background-repeat: no-repeat;
}

.select {
  background: white;
}

.logo {
  font-family: 'AlfaSlabOne'
}

:deep(.el-main ){
  padding-top: 0;
}

.drawer-transition {
  transition: 0.3s width ease-in-out, 0.3s padding-left ease-in-out, 0.3s padding-right ease-in-out;
}


</style>