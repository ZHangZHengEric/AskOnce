package dto_search

type CaseRes struct {
	Cases []string `json:"cases"`
}

type GenSessionRes struct {
	SessionId string `json:"sessionId"`
}

type KdbListRes struct {
	List []KdbListItem `json:"list"`
}

type KdbListItem struct {
	KdbId      int64  `json:"kdbId"`
	KdbName    string `json:"kdbName"`
	CreateTime string `json:"createTime"`
}

type AskRes struct {
	Stage string `json:"stage"` // analyze, search, generate, complete
	Text  string `json:"text"`
}

type HisRes struct {
	SessionId    string `json:"sessionId"`
	Question     string `json:"question"`
	Result       string `json:"result"`
	ResultRefers string `json:"resultRefers"`
	Unlike       bool   `json:"unlike"`
}

type ReferenceRes struct {
	List []ReferenceItem `json:"list"`
}

type ReferenceItem struct {
	Title   string `json:"title"`
	Url     string `json:"url"`
	Content string `json:"content"`
}

type OutlineRes struct {
	List []OutlineItem `json:"list"`
}

type OutlineItem struct {
	Level   string `json:"level"`
	Content string `json:"content"`
}

type RelationRes struct {
	EventsInfo []RelateEventsInfo `json:"eventsInfo"`
	PeopleInfo []RelatePeopleInfo `json:"peopleInfo"`
	OrgsInfo   []RelateOrgInfo    `json:"orgsInfo"`
}

type RelateEventsInfo struct {
	EventName     string `json:"eventName"`
	EventDate     string `json:"eventDate"`
	EventDescribe string `json:"eventDescribe"`
}

type RelatePeopleInfo struct {
	PersonName     string `json:"personName"`
	PersonDescribe string `json:"personDescribe"`
}

type RelateOrgInfo struct {
	OrgName     string `json:"orgName"`
	OrgDescribe string `json:"orgDescribe"`
}

type ProcessRes struct {
	List []ProcessItem `json:"list"`
}

type ProcessItem struct {
	StageName string `json:"stageName"`
	StageType string `json:"stageType"`
	Content   string `json:"content"`
	Time      int64  `json:"time"`
}