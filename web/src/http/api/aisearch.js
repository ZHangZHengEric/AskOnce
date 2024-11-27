import Http from "@/http/http";

export const aiSearchSession = (params) => {
    return Http.get("/askOnce/search/genSession", {params});
}
export const aiSearchCase = (params) => {
    return Http.get("/askOnce/search/case", {params});
}
export const aiSearchHis = (data) => {
    return Http.post("/askOnce/search/his", data);
}
export const aiSearchRefer = (data) => {
    return Http.post("/askOnce/search/refer", data);
}
export const aiSearchOutLine = (data) => {
    return Http.post("/askOnce/search/outline", data);
}
export const historyAsk = (data) => {
    return Http.post("/askOnce/history/ask", data);
}
export const unlike = (data) => {
    return Http.post("/askOnce/search/unlike", data);
}
export const relation = (data) => {
    return Http.post("/askOnce/search/relation", data);
}
export const searchProcess = (data) => {
    return Http.post("/askOnce/search/process", data);
}
