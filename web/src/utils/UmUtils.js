/**
 * create by Administrator on 2023/11/17 10:36
 * 类注释：友盟统计工具
 * 备注：
 */

/**
 * 用于发送页面上按钮等交互元素被触发时的事件统计请求。如视频的“播放、暂停、调整音量”，页面上的“返回顶部”、“赞”、“收藏”等。也可用于发送Flash事件统计请求。
 * @param category string 表示事件发生在谁身上，如“视频”、“小说”、“轮显层”等等
 * @param action string 表示访客跟元素交互的行为动作，如“播放”、“收藏”、“翻层”等等
 * @param label string 选填 用于更详细的描述事件，如具体是哪个视频，哪部小说
 * @param value int 选填 用于填写打分型事件的分值，加载时间型事件的时长，订单型事件的价格 请填写整数数值，如果填写为其他形式(类型)，系统将按0处理，若填写为浮点小数，系统会自动取整，去掉小数点
 */
export const trackEvent = (category, action, label = '', value = 0) => {
    window._czc && window._czc.push(["_trackEvent", category, action, label, value]);
}
/**
 * 用于发送某个URL的PV统计请求，适用于统计AJAX、异步加载页面，友情链接，下载链接的流量。
 * @param content_url string 自定义页面url 路径
 * @param referer_url string 选填 自定义该受访页面路径
 */
export const trackPageview = (content_url, referer_url = '') => {
    window._czc && window._czc.push(["_trackPageview", content_url, referer_url]);
}