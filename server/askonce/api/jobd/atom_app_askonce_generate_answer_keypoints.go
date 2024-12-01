package jobd

import (
	"askonce/components/dto/dto_search"
	"fmt"
)

type generateAnswerKeyPointsReq struct {
	Id           string                          `json:"id"`
	Question     string                          `json:"question"`
	SearchResult []dto_search.CommonSearchOutput `json:"search_result"`
	ModelName    string                          `json:"model_name"`
}

type AnswerKeyPointsRes struct {
	Id              string                `json:"id"`
	AnswerKeyPoints []AnswerKeyPointsItem `json:"answer_keypoints"`
}

type AnswerKeyPointsItem struct {
	Level   string `json:"level"`
	Content string `json:"content"`
}

type KeyPointNode struct {
	AnswerKeyPointsItem
	FullPath string // 完整路径的内容
	Children []*KeyPointNode
}

func ReturnTree(items []AnswerKeyPointsItem) (result []*KeyPointNode, allPath []string) {
	// 构建树形结构
	root := &KeyPointNode{
		AnswerKeyPointsItem: AnswerKeyPointsItem{
			Level:   "",
			Content: "root",
		},
	}
	currentH1 := root
	for _, item := range items {
		level := item.Level
		content := item.Content

		switch level {
		case "h1":
			allPath = append(allPath, content)
			newNode := &KeyPointNode{AnswerKeyPointsItem: AnswerKeyPointsItem{
				Level:   level,
				Content: content,
			}, FullPath: content}
			root.Children = append(root.Children, newNode)
			currentH1 = newNode
		case "h2":

			fullPath := fmt.Sprintf("%s-%s", currentH1.Content, content)
			allPath = append(allPath, fullPath)

			newNode := &KeyPointNode{AnswerKeyPointsItem: AnswerKeyPointsItem{
				Level:   level,
				Content: content,
			}, FullPath: fullPath}
			currentH1.Children = append(currentH1.Children, newNode)
		case "h3":
			fullPath := fmt.Sprintf("%s-%s-%s", currentH1.Content, currentH1.Children[len(currentH1.Children)-1].Content, content)
			allPath = append(allPath, fullPath)

			newNode := &KeyPointNode{AnswerKeyPointsItem: AnswerKeyPointsItem{
				Level:   level,
				Content: content,
			}, FullPath: fullPath}
			currentH1.Children[len(currentH1.Children)-1].Children = append(currentH1.Children[len(currentH1.Children)-1].Children, newNode)
		}
	}
	return root.Children, allPath
}

func (entity *JobdApi) GenerateAnswerKeyPoints(question string, searchResult []dto_search.CommonSearchOutput) (res *AnswerKeyPointsRes, err error) {
	input := &generateAnswerKeyPointsReq{
		Id:           "",
		Question:     question,
		SearchResult: searchResult,
	}
	return doTaskProcess[*generateAnswerKeyPointsReq, *AnswerKeyPointsRes](entity, "atom_app_askonce_generate_answer_keypoints", input, 100000)
}
