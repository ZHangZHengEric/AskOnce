import axios from "axios";
import {ElMessage} from "element-plus";
import {clearSession, getLanguage} from "@/utils/tools";
import {store} from "@/store";
import {useUserStore} from "@/store";
// import {useRouter, useRoute} from 'vue-router';
import router from '@/router'


const userStore = useUserStore(store);
// const router = useRouter();
// const route = useRoute()

const Http = axios.create({
    baseURL: "/serverApi",
    timeout: 400 * 1000,
    withCredentials: true,
});

Http.interceptors.request.use((config) => {
    config.headers["Accept-Language"] = getLanguage()
    return config;
});

Http.interceptors.response.use((response) => {
        const data = response.data;
        if (data.errNo === 0) {
            return data;
        } else if (data.errNo === 3) {
            ElMessage.error(data.errMsg);
            clearSession();
            router.push('/login?redirect=' + window.location.pathname + window.location.search)
        } else {
            ElMessage.error(data.errMsg);
            return Promise.reject(data);
        }
    },
    (error) => {
        return Promise.reject(error);
    }
);

export default Http;
