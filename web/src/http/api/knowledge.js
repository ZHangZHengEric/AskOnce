import Http from "@/http/http";


/**
 * 新增知识库
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeAdd = (data) => {
    return Http.post("/askonce/kdb/add", data);
};
/**
 * 知识库列表
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeList = (data) => {
    return Http.post("/askonce/kdb/list", data);
};
/**
 * 知识库删除
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeDel = (data) => {
    return Http.post("/askonce/kdb/delete", data);
};
/**
 * 问答可选知识库列表
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const kdbList = (data) => {
    return Http.post("/askOnce/search/kdbList", data);
};
/**
 * 知识库详情
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeDetail = (params) => {
    return Http.get("/askonce/kdb/detail", {params});
};
/**
 * 知识库详情
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeUpdate = (data) => {
    return Http.post("/askonce/kdb/update", data);
};
/**
 * 知识库数据新增
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeDataAdd = (data, config) => {
    return Http.post("/askonce/kdb/doc/add", data, config);
};
/**
 * 知识库数据列表
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeDataList = (data) => {
    return Http.post("/askonce/kdb/doc/list", data);
};
/**
 * 知识库数据删除
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeDataDel = (data) => {
    return Http.post("/askonce/kdb/doc/delete", data);
};
/**
 * 知识库搜索
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeSearch = (data) => {
    return Http.post("/askonce/kdb/doc/recall", data);
}
/**
 * 知识库封面
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeCovers = () => {
    return Http.get("/askonce/kdb/covers");
}
/**
 * 知识库权限检查
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeAuth = (data) => {
    return Http.post("/askonce/kdb/auth", data);
}
/**
 * 已有用户列表
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeUserList = (data) => {
    return Http.post("/askonce/kdb/userList", data);
}
/**
 * 用户新增查询
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeUserQuery = (data) => {
    return Http.post("/askonce/kdb/userQuery", data);
}
/**
 * 用户新增
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeUserAdd = (data) => {
    return Http.post("/askonce/kdb/userAdd", data);
}
/**
 * 用户删除
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeUserDelete = (data) => {
    return Http.post("/askonce/kdb/userDelete", data);
}
/**
 * 用户删除自己
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeUserDeleteSelf = (data) => {
    return Http.post("/askonce/kdb/deleteSelf", data);
}
/**
 * 数据重做
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeDataRedo = (data) => {
    return Http.post("/askonce/kdb/doc/redo", data);
}
/**
 * 数据重做
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeShareCodeGen = (data) => {
    return Http.post("/askonce/kdb/shareCodeGen", data);
}
/**
 * 分享码验证
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeShareCodeVerify = (data) => {
    return Http.post("/askonce/kdb/shareCodeVerify", data);
}
/**
 * 分享码信息获取
 * @param data
 * @returns {Promise<AxiosResponse<any>>}
 */
export const knowledgeShareCodeInfo = (params) => {
    return Http.get("/askonce/kdb/shareCodeInfo", {params});
}
