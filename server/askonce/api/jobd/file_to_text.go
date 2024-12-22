package jobd

import "time"

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

type EmlTextDetail struct {
	Sender            string   `json:"sender"`             //发件人
	Receiver          string   `json:"receiver"`           //收件人
	SenderName        string   `json:"sender_name"`        //发件人
	SenderAddress     string   `json:"sender_address"`     //发件人
	ReceiverNames     []string `json:"receiver_names"`     //收件人
	ReceiverAddresses []string `json:"receiver_addresses"` //收件人

	Subject    string        `json:"subject"`    //主题
	Date       any           `json:"date"`       //发送日期 时间戳秒
	Cc         string        `json:"cc"`         //抄送
	Body       EmlContent    `json:"body"`       //正文
	Attachment EmlAttachment `json:"attachment"` //附件
}

type EmlContent struct {
	TextContent       string `json:"text_content"`
	TextOriginContent string `json:"text_origin_content"`
	HtmlContent       string `json:"html_content"`
}

type EmlAttachment struct {
	Files  []EmlAttachmentFile `json:"files"`
	Images []string            `json:"images"`
}

type EmlAttachmentFile struct {
	AttachName        string        `json:"attach_name"`
	FileName          string        `json:"file_name"`
	FileContent       string        `json:"file_content"`
	FileOriginContent string        `json:"file_origin_content"`
	SourceCode        string        `json:"source_code"`
	EmlCode           string        `json:"eml_code"`
	TextDetail        EmlTextDetail `json:"text_detail"`
	EmlSendTime       time.Time     `json:"eml_send_time"`
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
	ress, err := doTaskProcess[[]FileToTextItem, []*FileToTextRes](entity, "convert_file_to_text", []FileToTextItem{inputReq}, 100000)
	if err != nil {
		return nil, err
	}
	res = ress[0]
	return
}
