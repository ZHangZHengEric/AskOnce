package jobd

type DocumentSplitInput struct {
	Text               string `json:"text"`
	Id                 int64  `json:"id"`
	WindowSize         int    `json:"window_size"`
	Stride             int    `json:"stride"`
	FixLengthList      []int  `json:"fix_length_list"`
	TextCuttingVersion string `json:"text_cutting_version"`
}
type DocumentSplitResp struct {
	SentencesList []TextChunkItem `json:"sentences_list"`
}

type TextChunkItem struct {
	PassageId      string `json:"passage_id"`
	PassageContent string `json:"passage_content"`
	Start          int    `json:"start"`
	End            int    `json:"end"`
}

func (entity *JobdApi) DocumentSplit(content string) (res DocumentSplitResp, err error) {
	render := DocumentSplitResp{}
	inputReq := DocumentSplitInput{
		Text:               content,
		Id:                 1,
		TextCuttingVersion: "move_window_cutting",
	}
	inputReq.WindowSize = 256
	inputReq.Stride = 170
	inputReq.FixLengthList = []int{256}
	render, err = doTaskProcess[DocumentSplitInput, DocumentSplitResp](entity, "document_split", inputReq, 10000)
	if err != nil {
		return
	}
	return render, nil
}
