/**
 * create by zhangxiang on 2023-06-14 10:41
 * 类注释：
 * 备注：
 */

import Http from "@/http/http";

/**
 * 登录接口
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const login = (data) => {
    return Http.post("/askonce/user/loginByAccount", data);
};


/**
 * 注册
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const register = (data) => {
    return Http.post("/askonce/user/registerByAccount", data);
};

/**
 * 获取用户信息
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const loginInfo = () => {
    return Http.get("/askonce/user/loginInfo");
};
/**
 * 手机号验证码
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const loginSendSms = (data) => {
    return Http.post("/askonce/user/loginSendSms", data);
};
/**
 * 手机号登录
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const loginByPhone = (data) => {
    return Http.post("/askonce/user/loginByPhone", data);
};

