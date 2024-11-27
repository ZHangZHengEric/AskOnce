/**
 * create by zhangxiang on 2023-06-10 12:50
 * 类注释：
 * 备注：
 */
import { defineStore } from "pinia";

export const useUserStore = defineStore("user-store", {
	state: () => {
		return {
			info: {},
			isLogin: null,
			needLogin: false,
			hasAuth: false,
		};
	},
	getters: {
		getIsLogin: (state) => {
			return state.isLogin;
		},
		getNeedLogin: (state) => {
			return state.needLogin;
		},
		getAuth: (state) => {
			return state.hasAuth;
		},
		getInfo: (state) => {
			return state.info;
		},
	},
	actions: {
		setIsLogin(val) {
			this.isLogin = val;
		},
		setNeedLogin(val) {
			this.needLogin = val;
		},
		setAuth(val) {
			this.hasAuth = val;
		},
		setInfo(val) {
			this.info = val;
		},
	},
});
