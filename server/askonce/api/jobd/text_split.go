package jobd

type TextChunkInput struct {
	Text      string `json:"text"`
	Id        int64  `json:"id"`
	ChunkType string `json:"chunk_type"`
}

type MoveWindowTextChunkItem struct {
	PassageId      string `json:"passage_id"`
	PassageContent string `json:"passage_content"`
	Start          int    `json:"start"`
	End            int    `json:"end"`
}

func (entity *JobdApi) TextSplit(text string) (res []MoveWindowTextChunkItem, err error) {
	input := &TextChunkInput{
		ChunkType: "sign_mv_chunker",
		Text:      text,
	}

	res, err = doTaskProcess[*TextChunkInput, []MoveWindowTextChunkItem](entity, "textchunk", input, 100000)
	if res == nil {
		res = make([]MoveWindowTextChunkItem, 0)
	}
	return res, err
}
