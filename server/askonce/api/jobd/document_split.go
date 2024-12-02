package jobd

type DocumentSplitInput struct {
	Text               string `json:"text"`
	Id                 int64  `json:"id"`
	DocName            string `json:"doc_name"`
	WindowSize         int    `json:"window_size"`
	Stride             int    `json:"stride"`
	FixLengthList      string `json:"fix_length_list"`
	TextCuttingVersion string `json:"text_cutting_version"`
}
type DocumentSplitRender struct {
	S DocumentSplitV2Res `json:"sentences_list"`
}

type DocumentSplitV2Res struct {
	DocId               string          `json:"doc_id"`
	DocName             string          `json:"doc_name"`
	DocTitle            string          `json:"doc_title"`
	DocContent          string          `json:"doc_content"`
	DocSummary          string          `json:"doc_summary"`
	MoveWindowTextChunk []TextChunkItem `json:"move_window_text_chunk"`
}

type TextChunkItem struct {
	PassageId      string `json:"passage_id"`
	PassageContent string `json:"passage_content"`
	Start          int    `json:"start"`
	End            int    `json:"end"`
}

func (entity *JobdApi) DocumentSplit(docName, content string) (res DocumentSplitV2Res, err error) {
	render := DocumentSplitRender{}
	inputReq := DocumentSplitInput{
		Text:               content,
		Id:                 1,
		DocName:            docName,
		TextCuttingVersion: "text_cutting_v2",
	}
	inputReq.WindowSize = 128
	inputReq.Stride = 64
	inputReq.FixLengthList = "128 256 512"
	render, err = doTaskProcess[DocumentSplitInput, DocumentSplitRender](entity, "document_split", inputReq, 10000)
	if err != nil {
		return
	}
	return render.S, nil
}
