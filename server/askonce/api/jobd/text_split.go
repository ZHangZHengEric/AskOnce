package jobd

type TextSplitReq struct {
	Text      string `json:"text"`
	Id        int64  `json:"id"`
	ChunkType string `json:"chunk_type"`
}

type TextSplitRes struct {
	PassageId      string `json:"passage_id"`
	PassageContent string `json:"passage_content"`
	Start          int    `json:"start"`
	End            int    `json:"end"`
}

func (entity *JobdApi) TextSplit(text string) (res []TextSplitRes, err error) {
	input := &TextSplitReq{
		ChunkType: "sign_mv_chunker",
		Text:      text,
	}

	res, err = doTaskProcess[*TextSplitReq, []TextSplitRes](entity, "textchunk", input, 100000)
	if res == nil {
		res = make([]TextSplitRes, 0)
	}
	return res, err
}
