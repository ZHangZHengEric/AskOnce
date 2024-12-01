package jobd

import "askonce/components/dto/dto_search"

type GenerateRelateInfoReq struct {
	Id           string                          `json:"id"`
	Question     string                          `json:"question"`
	AnswerStyle  string                          `json:"answer_style"`
	SearchResult []dto_search.CommonSearchOutput `json:"search_result"`
	ModelName    string                          `json:"model_name"`
}

type GenerateRelateInfoRes struct {
	Id         string                     `json:"id"`
	EventsInfo []GenerateRelateEventsInfo `json:"events_info"`
	PeopleInfo []GenerateRelatePeopleInfo `json:"people_info"`
	OrgsInfo   []GenerateRelateOrgInfo    `json:"orgs_info"`
}

type GenerateRelateEventsInfo struct {
	EventName     string `json:"event_name"`
	EventDate     string `json:"event_date"`
	EventDescribe string `json:"event_describe"`
}

type GenerateRelatePeopleInfo struct {
	PersonName     string `json:"person_name"`
	PersonDescribe string `json:"person_describe"`
}

type GenerateRelateOrgInfo struct {
	OrgName     string `json:"org_name"`
	OrgDescribe string `json:"org_describe"`
}

func (entity *JobdApi) GenerateRelateInfo(question string, searchResult []dto_search.CommonSearchOutput) (res *GenerateRelateInfoRes, err error) {
	input := &GenerateRelateInfoReq{
		Id:           "",
		Question:     question,
		SearchResult: searchResult,
	}
	return doTaskProcess[*GenerateRelateInfoReq, *GenerateRelateInfoRes](entity, "atom_app_askonce_generate_relate_info", input, 100000)
}
