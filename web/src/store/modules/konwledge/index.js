/**
 * create by zhangxiang on 2023-06-10 12:50
 * 类注释：
 * 备注：
 */
import {defineStore} from "pinia";

export const useKnowledgeStore = defineStore("knowledge-store", {
    state: () => {
        return {
            authType: -1,
        };
    },
    getters: {},
    actions: {
        setAuthType(val) {
            this.authType = val;
        }
    },
});
