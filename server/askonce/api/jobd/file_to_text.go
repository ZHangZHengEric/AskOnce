package jobd

type FileToTextItem struct {
	FilePath      string `json:"file_path"`
	NeedOcr       bool   `json:"need_ocr"`
	NeedTable     bool   `json:"need_table"`
	NeedParagraph bool   `json:"need_paragraph"`
	StartPage     int    `json:"start_page"`
	EndPage       int    `json:"end_page"`
	Id            int    `json:"id"`
	NeedDetail    bool   `json:"need_detail"`
	RemoveWrap    bool   `json:"remove_wrap"`
}

type FileToTextRes struct {
	Text string `json:"text"`
}

type TextDetail struct {
	Text            string `json:"text"`
	TextIndex       int    `json:"text_index"`
	PageIndex       int    `json:"page_index"`
	IndexInDocument []int  `json:"index_in_document"`
}

func (entity *JobdApi) FileToText(path string) (res *FileToTextRes, err error) {
	inputReq := FileToTextItem{
		FilePath:      path,
		NeedOcr:       true,
		NeedTable:     true,
		NeedParagraph: true,
		Id:            1,
		NeedDetail:    false,
		StartPage:     -1,
		EndPage:       -1,
		RemoveWrap:    true,
	}
	ress, err := doTaskProcess[[]FileToTextItem, []*FileToTextRes](entity, "convert_file_to_text", []FileToTextItem{inputReq}, 10000)
	if err != nil {
		return nil, err
	}
	res = ress[0]
	return
}
