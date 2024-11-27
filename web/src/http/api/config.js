/**
 * create by zhangxiang on 2023-06-14 10:41
 * 类注释：
 * 备注：
 */

import Http from "@/http/http";

/**
 * 搜索配置详情
 * @param
 * @returns {Promise<AxiosResponse<any>>}
 */
export const configDetail = (params) => {
    return Http.get("/askOnce/config/detail", {params});
};
/**
 * 模型配置字典
 * @param
 * @returns {Promise<AxiosResponse<any>>}
 */
export const configDict = (params) => {
    return Http.get("/askOnce/config/dict", {params});
}


/**
 * 搜索配置保存
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const configSave = (data) => {
    return Http.post("/askOnce/config/save", data);
}