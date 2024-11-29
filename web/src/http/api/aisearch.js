import Http from "@/http/http";

export const aiSearchSession = (params) => {
    return Http.get("/askonce/search/genSession", {params});
}
export const aiSearchCase = (params) => {
    return Http.get("/askonce/search/case", {params});
}
export const aiSearchHis = (data) => {
    return Http.post("/askonce/search/his", data);
}
export const aiSearchRefer = (data) => {
    return Http.post("/askonce/search/refer", data);
}
export const aiSearchOutLine = (data) => {
    return Http.post("/askonce/search/outline", data);
}
export const historyAsk = (data) => {
    return Http.post("/askonce/history/ask", data);
}
export const unlike = (data) => {
    return Http.post("/askonce/search/unlike", data);
}
export const relation = (data) => {
    return Http.post("/askonce/search/relation", data);
}
export const searchProcess = (data) => {
    return Http.post("/askonce/search/process", data);
}
