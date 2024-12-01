package jobd

type AtomResultReferenceReq struct {
	Id            string   `json:"id"`
	Result        string   `json:"result"`
	ReferenceList []string `json:"reference_list"`
}

type AtomResultReferenceRes struct {
	Id                       string           `json:"id"`
	ReferenceMap             []ReferenceMap   `json:"reference_map"`
	ReferenceListSelectIndex map[string][]int `json:"reference_list_select_index"`
}

type ReferenceMap struct {
	IndexRange    []int `json:"index_range"`
	ReferenceList []int `json:"reference_list"`
}

func (entity *JobdApi) AtomResultReference(result string, referenceList []string) (res *AtomResultReferenceRes, err error) {
	input := &AtomResultReferenceReq{
		Id:            "",
		Result:        result,
		ReferenceList: referenceList,
	}
	return doTaskProcess[*AtomResultReferenceReq, *AtomResultReferenceRes](entity, "atom_result_reference", input, 100000)
}
