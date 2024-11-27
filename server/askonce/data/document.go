package data

import (
	"askonce/components"
	"askonce/components/defines"
	"askonce/components/dto/dto_kdb"
	"askonce/models"
	"github.com/xiangtao94/golib/flow"
)

/*
*
文档处理
*/
type DocumentData struct {
	flow.Layer
	fileDao *models.FileDao
	jobdApi *jobd.JobdApi
	ragApi  *rag.RagApi
}

func (d *DocumentData) OnCreate() {
	d.fileDao = flow.Create(d.GetCtx(), new(models.FileDao))
	d.jobdApi = flow.Create(d.GetCtx(), new(jobd.JobdApi))
	d.ragApi = flow.Create(d.GetCtx(), new(rag.RagApi))
}

// 默认的知识库文档处理规则
var DefaultProcessRuleDetail = dto_kdb.DocProcessRuleDetail{
	PreProcessRules: []dto_kdb.PreProcessRule{
		{
			Id:      "remove_extra_spaces",
			Enabled: false,
		},
		{
			Id:      "remove_urls_emails",
			Enabled: false,
		},
	},
	Segmentation: dto_kdb.SegmentationRule{
		Separator:    "\n",
		ChunkSize:    500,
		ChunkOverlap: 50,
	},
}

// 文件转文本
func (d *DocumentData) ConvertFileToText(fileId string) (fileName string, output string, err error) {
	// 获取文件
	file, err := d.fileDao.GetById(fileId)
	if err != nil {
		return
	}
	if file == nil { // 文件不存在
		return "", "", components.ErrorFileNoExist
	}
	fileToText, err := d.jobdApi.FileToText(file.Path)
	if err != nil {
		return
	}
	return file.Name, fileToText.Text, nil
}

// 文本切分
func (d *DocumentData) TextSplit(processRule dto_kdb.ProcessRule, text string, inFields []string) (segments []map[defines.StructuredKey]any, isStructured bool, fields []string, err error) {
	// 文本处理规则
	processRuleDetail := DefaultProcessRuleDetail
	if processRule.Mode == dto_kdb.DocProcessModeCustom {
		processRuleDetail = processRule.Rules
	}
	isStructured = false
	fields = []string{}
	segments = make([]map[defines.StructuredKey]any, 0)
	if len(processRuleDetail.Segmentation.SegmentMode) == 0 || processRuleDetail.Segmentation.SegmentMode == dto_kdb.DocSegmentModeSimple { //不考虑结构化
		ragRes, err := d.ragApi.TextSplit(text, processRuleDetail)
		if err != nil {
			return segments, isStructured, fields, err
		}
		for _, t := range ragRes.Text {
			segments = append(segments, map[defines.StructuredKey]any{
				defines.StructuredContent: t,
			})
		}
		return segments, isStructured, fields, nil
	}
	ragRes, err := d.ragApi.TextSplitWithIndex(text, inFields, processRuleDetail)
	if err != nil {
		return
	}
	isStructured = ragRes.IsStructured
	for _, value := range ragRes.Values {
		segments = append(segments, value)
	}
	fields = ragRes.Fields
	return
}
