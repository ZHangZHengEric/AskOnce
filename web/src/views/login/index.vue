/**
* create by wanghu on 2024-2-22 9:58
* 类注释：
* 备注：
*/

<template>

  <el-container class="login h-screen">
    <el-header class="flex mt-5 ml-4">
      <div class="text-default text-size32 font-[400] logo flex-1">AskOnce</div>
      <div class="text-size16 h-10 mt-2.5 leading-10 font-normal text-default">{{ $t('login.noAccount') }}</div>
      <div class="px-12 ml-4 h-10 mt-2.5 leading-10 cursor-pointer border border-default border-solid rounded-full"
           @click="router.push({path:'/register'})">{{ $t('login.register') }}
      </div>
    </el-header>
    <el-main class="h-full w-full flex items-center justify-center text-default">
      <div class="rounded-xl  box bg-[#FFFFFF80] px-12 py-10">
        <el-form :model="formData" label-position="top" v-if="!data.showMobile">
          <div class="text-size22 font-[600] text-default">{{ $t('login.signTip') }}</div>
          <el-form-item class="mt-10" :label="$t('login.Account')">
            <el-input class="w-[30rem]" v-model="formData.account"></el-input>
          </el-form-item>
          <el-form-item :label="$t('login.Password')">
            <el-input class="w-[30rem]" show-password type="password" v-model="formData.password"></el-input>
          </el-form-item>
          <div class="flex mt-10" v-if="false">
            <div class="flex-1"></div>
            <div class="text-normal text-size16 border-b border-solid cursor-pointer" @click="data.showMobile=true">
              {{ $t('login.usePhone') }}
            </div>
          </div>
          <div
              class="w-full mt-10 cursor-pointer bg-default text-white text-size18
              font-semibold text-center rounded-full py-1.5"
              @click="toLogin">
            {{ $t('login.login') }}
          </div>
        </el-form>

        <el-form
            ref="refMobile"
            :model="mobileForm"
            v-if="data.showMobile"
            :rules="rules"
            hide-required-asterisk
            label-position="top">
          <div class="text-size22 font-[600] text-default text-center whitespace-pre-wrap"> {{
              $t('login.phoneTip')
            }}
          </div>
          <el-form-item label=" " class="mt-10" prop="phone">
            <el-input :placeholder="$t('login.phonep')" class="w-[30rem]" v-model="mobileForm.phone"></el-input>
          </el-form-item>
          <el-form-item class="mt-10" prop="smsCode">
            <el-input :placeholder="$t('login.smsp')" class="w-[20rem]" v-model="mobileForm.smsCode"></el-input>
            <button type="button" class="w-[8rem] ml-[2rem] text-default" @click="sendSms"> {{
                getCodeBtn
              }}
            </button>
          </el-form-item>
          <div
              class="w-full mt-10 cursor-pointer bg-default text-white text-size18
              font-semibold text-center rounded-full py-1.5"
              @click="mobileLogin">
            {{ $t('login.login') }}
          </div>
        </el-form>
      </div>
    </el-main>
  </el-container>

</template>

<script setup>
import {reactive, ref, toRefs} from 'vue'
import {login, loginByPhone, loginSendSms} from "@/http/api/user";
import {setSession} from "@/utils/tools";
import {useUserStore} from '@/store'
import {useRouter, useRoute} from 'vue-router'
import {ElMessage} from "element-plus";
import {useI18n} from 'vue-i18n'

const {t} = useI18n()

const userStore = useUserStore()
const router = useRouter()
const route = useRoute()
const refMobile = ref()

const formData = reactive({
  account: '',
  password: ""
})

const mobileForm = reactive({
  phone: '',
  smsCode: ''
})

const data = reactive({
  loading: false,
  showMobile: false,
  getCodeBtn: t('login.getsms'),
  codeSecond: 60,
})
let timer = null
const sendSms = () => {
  refMobile.value.validateField(["phone"], isOk => {
    if (isOk) {
      data.loading = true
      timer = setInterval(() => {
        if (data.codeSecond > 0) {
          data.getCodeBtn = data.codeSecond--
        } else {
          data.getCodeBtn = t('login.getsms')
          data.codeSecond = 60
          data.loading = false
          window.clearInterval(timer)
        }
      }, 1000)
      loginSendSms({
        phone: mobileForm.phone
      }).then(() => {
        ElMessage({
          message: '验证码发送成功',
          type: 'success',
        })
      })
    }
  })
}

const mobileLogin = () => {
  refMobile.value.validate(isOk => {
    if (isOk) {
      loginByPhone({
        phone: mobileForm.phone,
        smsCode: mobileForm.smsCode,
      }).then(res => {
        setSession(res.data.atomechoSession)
        userStore.setIsLogin(true)
        userStore.setInfo(res.data)
        if (route.query.redirect) {
          let url = []
          for (let key in route.query) {
            if (key !== 'redirect') {
              url.push(key = route.query[key])
            }
          }
          if (url.length) {
            router.push(route.query.redirect + '&' + url.join('&'))
          } else {
            router.push(route.query.redirect)
          }
        } else {
          router.push({path: '/'})
        }
      }).catch(() => {
        data.loginLoading = false
      })
    }
  })
}

const checkPhone = (rule, value, callback) => {
  if (!value) {
    return callback(new Error('手机号不能为空'));
  } else {
    const reg = /^1[3|4|5|6|7|8|9][0-9]\d{8}$/
    if (reg.test(value)) {
      callback();
    } else {
      return callback(new Error('请输入正确的手机号'));
    }
  }
};

const rules = reactive({
  phone: [
    {required: true, message: '请输入手机号', trigger: 'blur'},
    {validator: checkPhone, trigger: 'blur'}
  ],
  smsCode: [
    {required: true, message: '请输入验证码', trigger: 'blur'},
  ]
})

const toLogin = () => {
  login(formData).then(res => {
    setSession(res.data.atomechoSession)
    userStore.setIsLogin(true)
    userStore.setInfo(res.data)
    if (route.query.redirect) {
      let url = []
      for (let key in route.query) {
        if (key !== 'redirect') {
          const a = `${key}=${route.query[key]}`
          url.push(a)
        }
      }
      if (url.length) {
        router.push(route.query.redirect + '&' + url.join('&'))
      } else {
        router.push(route.query.redirect)
      }
    } else {
      router.push({path: '/'})
    }
  })
}

const {getCodeBtn} = toRefs(data)


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