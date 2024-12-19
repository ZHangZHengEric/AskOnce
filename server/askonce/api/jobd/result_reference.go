package jobd

type AtomResultReferenceReq struct {
	Id            string   `json:"id"`
	Result        string   `json:"result"`
	Threshold     float32  `json:"threshold"`
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

func (entity *JobdApi) ResultAddReference(result string, referenceList []string, threshold float32) (res *AtomResultReferenceRes, err error) {
	defaultThreshold := float32(0.7)
	if threshold > 0 {
		defaultThreshold = threshold
	}
	input := &AtomResultReferenceReq{
		Id:            "",
		Threshold:     defaultThreshold,
		Result:        result,
		ReferenceList: referenceList,
	}
	return doTaskProcess[*AtomResultReferenceReq, *AtomResultReferenceRes](entity, "result_add_reference", input, 100000)
}
