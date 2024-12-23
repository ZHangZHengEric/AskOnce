import {createRouter, createWebHashHistory, createWebHistory} from "vue-router";
import {store} from "@/store";
import {useUserStore} from "@/store";
import {getUserInfo, loginInfo} from "@/http/api/user";
import {getSession} from "@/utils/tools";
import NProgress from 'nprogress' // progress bar
import 'nprogress/nprogress.css' // progress bar style
NProgress.configure({showSpinner: false}) // NProgress Configuration


import layout from '@/layout'


const userStore = useUserStore(store);

const constantRoutes = [
    {
        path: "/login",
        name: "login",
        meta: {
            white: true,
        },
        component: () => import("@/views/login"),
    },
    {
        path: "/register",
        name: "register",
        meta: {
            white: true,
        },
        component: () => import("@/views/register"),
    },
    {
        path: "/online-file",
        name: "onlineFile",
        meta: {
            white: true,
        },
        component: () => import("@/views/onLine"),
    },
    {
        path: "/knowledge-share",
        name: "knowledgeShare",
        component: () => import("@/views/knowledge/share"),
    },
    {
        path: '/',
        name: 'home',
        component: layout,
        children: [
            {
                path: '/',
                name: 'home',
                component: () => import('@/views/home')
            },
            {
                path: "/detail/:id",
                name: "detail",
                component: () => import("@/views/detail"),
            },
            {
                path: "knowledge-manage",
                name: "knowledgeManage",
                component: () => import("@/views/knowledge"),
            },
            {
                path: "knowledge-create",
                name: "knowledgeCreate",
                component: () => import("@/views/knowledge/create.vue"),
            },
            {
                path: "knowledge-add",
                name: "knowledgeAdd",
                component: () => import("@/views/knowledge/config/add-file.vue"),
            },
            {
                path: "database-add",
                name: "databaseAdd",
                component: () => import("@/views/knowledge/config/add-database.vue"),
            },
            {
                path: "database-detail",
                name: "databaseDetail",
                component: () => import("@/views/knowledge/config/database-detail.vue"),
            },
            {
                path: "knowledge-config",
                name: "knowledgeConfig",
                component: () => import("@/views/knowledge/config/index.vue"),
                children: [
                    {
                        path: "detail",
                        name: "knowledgeDetail",
                        component: () => import("@/views/knowledge/config/detail.vue"),
                    },
                    {
                        path: 'base',
                        name: "knowledgeConfigBase",
                        component: () => import("@/views/knowledge/config/base.vue"),

                    },
                    {
                        path: 'member',
                        name: "knowledgeConfigMember",
                        component: () => import("@/views/knowledge/config/member.vue"),

                    },
                    {
                        path: "search",
                        name: "knowledgeSearch",
                        component: () => import("@/views/knowledge/config/search.vue"),
                    },
                    {
                        path: "setting",
                        name: "knowledgeSetting",
                        component: () => import("@/views/knowledge/config/setting.vue"),
                    },
                ]
            },
            {
                path: "history",
                name: "history",
                component: () => import("@/views/history/index.vue"),
            },
            {
                path: "search-config",
                name: "searchConfig",
                component: () => import("@/views/config/index.vue"),
            },
        ]
    }
];

export const asyRouter = []

export const routes = [...constantRoutes, ...asyRouter];

const router = createRouter({
    history: createWebHistory(),
    scrollBehavior: () => ({y: 0}),
    routes,
});

router.beforeEach(async (to, from, next) => {
    NProgress.start()
    if (to.meta.white) {
        await next()
        NProgress.done()
    } else {
        const session = getSession()
        if (session) {
            if (userStore.isLogin) {
                await next()
            } else {
                const res = await loginInfo()
                if (res.data.userId) {
                    await next()
                    userStore.setIsLogin(true)
                    userStore.setInfo(res.data)
                } else {
                    await next(`/login?redirect=${to.fullPath}`)
                }
            }
            NProgress.done()
        } else {
            await next(`/login?redirect=${to.fullPath}`)
            NProgress.done()
        }
    }
});


export default router;
