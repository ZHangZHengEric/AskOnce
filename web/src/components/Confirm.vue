/**
* create by wanghu on 2024/2/29 13:48
* 类注释：
* 备注：
*/
<template>
    <el-dialog v-model="data.visible"
        class="confirm bg-white w-60 md:w-80 lg:w-96 px-6 pb-4 rounded-md text-color333">
         <slot></slot>
        <div class="flex items-center justify-end">
            <div class="px-3 py-0.5 rounded-full bg-themeRed mr-3 font-base lg:text-lg text-theme cursor-pointer font-MiSans"
                @click="confirm">
                {{ $t('public.Confirm') }}
            </div>
            <div class="px-3 py-0.5 border border-solid border-themeRed rounded-full font-base lg:text-lg text-themeRed cursor-pointer font-MiSans"
                @click="cancel">
                {{ $t('public.Cancel') }}
            </div>
        </div>
    </el-dialog>
</template>

<script setup>
import { onMounted, reactive, defineExpose, defineProps, defineEmits } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
const { locale } = useI18n()

const data = reactive({
    visible: false,
})

const props = defineProps({
    text: String
})

const emit = defineEmits(["confirm"])

const confirm = () => {
    data.visible = false
    emit("confirm")
}

const cancel = () => {
    data.visible = false
}

defineExpose({
    data
})
</script>

<style lang="less" rel="stylesheet/less">
.confirm {

    .el-dialog__header,
    .el-dialog__body {
        padding: 0;
    }

    .el-dialog__headerbtn {
        @apply top-2 right-2 w-fit h-fit;
    }
}
</style>