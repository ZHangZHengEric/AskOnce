package dto_knowledge

type ListResp struct {
	List  []ListItem `json:"list"`
	Total int64      `json:"total"`
}

type ListItem struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Cover        string `json:"cover"`
	DefaultColor bool   `json:"defaultColor"`
	CreateTime   string `json:"createTime"`
	Creator      string `json:"creator"`
	Type         int    `json:"type"` // 1 公共数据
}

type DetailRes struct {
	KdbId int64  `json:"kdbId"`
	Name  string `json:"name"`  // 知识库名称
	Intro string `json:"intro"` // 知识库介绍
	Cover string `json:"cover"` // 知识库封面
}

type AddRes struct {
	KdbId int64 `json:"kdbId"`
}

type DataListResp struct {
	List  []DataListItem `json:"list"`
	Total int64          `json:"total"`
}

type DataListItem struct {
	Id         int64  `json:"id"`
	Type       string `json:"type"`
	DataSuffix string `json:"dataSuffix"`
	DataName   string `json:"dataName"`
	DataPath   string `json:"dataPath"`
	Status     int    `json:"status"` // 0正在构建到知识库 1 成功 2 失败
	CreateTime string `json:"createTime"`
}

type SearchRes struct {
	List []SearchItem `json:"list"`
}

type SearchItem struct {
	DataName      string `json:"dataName"`
	DataPath      string `json:"dataPath"`
	SearchContent string `json:"searchContent"`
	DataContent   string `json:"dataContent"`
}

type CoversRes struct {
	List []CoverItem `json:"list"`
}

type CoverItem struct {
	Id   int64  `json:"id"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

type AuthRes struct {
	AuthType int `json:"authType"`
}

type UserListRes struct {
	List  []UserListItem `json:"list"`
	Total int64          `json:"total"`
}

type UserListItem struct {
	UserId   int64  `json:"userId"`
	UserName string `json:"userName"`
	JoinTime string `json:"joinTime"`
}

type UserQueryRes struct {
	List  []UserQueryItem `json:"list"`
	Total int64           `json:"total"`
}

type UserQueryItem struct {
	UserId   int64  `json:"userId"`
	UserName string `json:"userName"`
}

type GenShareCodeRes struct {
	ShareCode string `json:"shareCode" `
}

type ShareCodeInfoRes struct {
	Creator  string `json:"creator"`
	KdbName  string `json:"kdbName"`
	AuthType int    `json:"authType"`
}
