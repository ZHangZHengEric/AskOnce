import {createApp} from "vue";
import App from "./App.vue";
import router from "./router";
import ElementPlus from "element-plus";
import "element-plus/dist/index.css";
import zhCn from "element-plus/dist/locale/zh-cn.mjs";
import "./assets/css/index.css";

const app = createApp(App);
import lazyPlugin from 'vue3-lazy'

app.use(ElementPlus, {
    locale: zhCn,
});
app.use(lazyPlugin, {
    // loading: require('@/assets/images/default.png'), // 图片加载时默认图片
    // error: require('@/assets/images/error.png')// 图片加载失败时默认图片
})

import {debounce} from "lodash";


//vue + element plus：ResizeObserver loop completed with undelivered notifications
const resizeObserver = window.ResizeObserver;
window.ResizeObserver = class ResizeObserver extends resizeObserver {
    constructor(callback) {
        callback = debounce(callback, 100);
        super(callback);
    }
};

import * as ElementPlusIconsVue from "@element-plus/icons-vue";

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component);
}

import "@/icons/index"; // 将src/icons/index.js引入
import SvgIcon from "./components/SvgIcon"; // 引入SvgIcon.vue
import {setupStore} from "./store";

import i18n from "./language";

app.use(i18n);

import directive from "./directive"; // directive
directive(app);

setupStore(app);
app.use(router);
app.component("svg-icon", SvgIcon);
app.mount("#app");
