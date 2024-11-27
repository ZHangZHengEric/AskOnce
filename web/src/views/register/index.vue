/**
* create by wanghu on 2024-2-22 9:58
* 类注释：
* 备注：
*/

<template>
  <el-container class="login h-screen">
    <el-header class="flex mt-5 ml-4">
      <div class="text-default text-size32 font-[400] logo flex-1">AskOnce</div>
      <div class="text-size16 h-10 mt-2.5 leading-10 font-normal text-default">{{  $t('login.hasAccount') }}</div>
      <div class="px-12 ml-4 h-10 mt-2.5 leading-10 border border-default border-solid rounded-full cursor-pointer"
           @click="router.push({path:'/login'})">{{  $t('login.login') }}
      </div>
    </el-header>
    <el-main class="h-full w-full flex items-center justify-center text-default">
      <div class="rounded-xl  box bg-[#FFFFFF80] px-12 py-10">
        <div class="text-size22 font-black">{{  $t('login.signUpTip') }}</div>
        <el-form ref="refElForm" :model="formData" class="mt-10" label-position="top">
          <el-form-item :label="$t('login.Account')">
            <el-input class="w-[30rem]" v-model="formData.account"></el-input>
          </el-form-item>
          <el-form-item :label="$t('login.Password')">
            <el-input class="w-[30rem]" show-password type="password" v-model="formData.password"></el-input>
          </el-form-item>
          <el-form-item :label="$t('login.rePassword')">
            <el-input class="w-[30rem]" show-password type="password" v-model="formData.repassword"></el-input>
          </el-form-item>
          <el-form-item>
            <el-checkbox v-model="formData.check">{{$t('login.checked')}}</el-checkbox>
          </el-form-item>
          <div
              @click="toRegister"
              class="w-full mt-10 cursor-pointer bg-default text-white text-size18 font-semibold text-center rounded-full py-1.5">
           {{$t('login.register')}}
          </div>
        </el-form>
      </div>
    </el-main>
  </el-container>

</template>

<script setup>
import {reactive, ref} from 'vue'
import {useRouter, useRoute} from 'vue-router'
import {register} from "@/http/api/user";
import {ElMessage} from "element-plus";

const router = useRouter()
const route = useRoute()

const refElForm = ref()

const formData = reactive({
  check: false
})

const toRegister = () => {
  if (!formData.account) {
    ElMessage.error("请输入用户名")

    return;
  }
  if (formData.password !== formData.repassword) {
    ElMessage.error("两次密码不一致")
    return
  }
  if (!formData.check) {
    ElMessage.error("请勾选协议")
    return;
  }
  register(formData).then(res => {
    ElMessage.success('注册成功')
    router.push({path:'/login'})
  })
}


</script>

<style scoped lang="less">
.login {
  background-image: url("@/assets/img/layout/main.png");
  background-size: 100% 100%;
  background-repeat: no-repeat;

  .box {
    box-shadow: 0px 0px 8px 0px #0000001A;

  }
}
</style>