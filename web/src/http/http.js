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
        if (data.code === 0) {
            return data;
        } else if (data.code === 3) {
            ElMessage.error(data.message);
            clearSession();
            router.push('/login?redirect=' + window.location.pathname + window.location.search)
        } else {
            ElMessage.error(data.message);
            return Promise.reject(data);
        }
    },
    (error) => {
        return Promise.reject(error);
    }
);

export default Http;
