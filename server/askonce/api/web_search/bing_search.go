package web_search

import (
	"askonce/conf"
	"encoding/json"
	"github.com/minio/pkg/env"
	"github.com/xiangtao94/golib/flow"
	http2 "github.com/xiangtao94/golib/pkg/http"
	"time"
)

// This struct formats the answers provided by the Bing Web Search API.
type BingAnswer struct {
	Type         string `json:"_type"`
	QueryContext struct {
		OriginalQuery string `json:"originalQuery"`
	} `json:"queryContext"`
	WebPages struct {
		WebSearchURL          string `json:"webSearchUrl"`
		TotalEstimatedMatches int    `json:"totalEstimatedMatches"`
		Value                 []struct {
			ID               string    `json:"id"`
			Name             string    `json:"name"`
			URL              string    `json:"url"`
			IsFamilyFriendly bool      `json:"isFamilyFriendly"`
			DisplayURL       string    `json:"displayUrl"`
			Snippet          string    `json:"snippet"`
			DateLastCrawled  time.Time `json:"dateLastCrawled"`
			SearchTags       []struct {
				Name    string `json:"name"`
				Content string `json:"content"`
			} `json:"searchTags,omitempty"`
			About []struct {
				Name string `json:"name"`
			} `json:"about,omitempty"`
			Summary string `json:"summary"`
		} `json:"value"`
	} `json:"webPages"`
	RelatedSearches struct {
		ID    string `json:"id"`
		Value []struct {
			Text         string `json:"text"`
			DisplayText  string `json:"displayText"`
			WebSearchURL string `json:"webSearchUrl"`
		} `json:"value"`
	} `json:"relatedSearches"`
	RankingResponse struct {
		Mainline struct {
			Items []struct {
				AnswerType  string `json:"answerType"`
				ResultIndex int    `json:"resultIndex"`
				Value       struct {
					ID string `json:"id"`
				} `json:"value"`
			} `json:"items"`
		} `json:"mainline"`
		Sidebar struct {
			Items []struct {
				AnswerType string `json:"answerType"`
				Value      struct {
					ID string `json:"id"`
				} `json:"value"`
			} `json:"items"`
		} `json:"sidebar"`
	} `json:"rankingResponse"`
}

type BingSearchResp struct {
	Title   string `json:"title"`
	Url     string `json:"url"`
	Content string `json:"content"`
}
type WebSearchReq struct {
	Query   string `json:"query"`
	Summary bool   `json:"summary"`
}
type WebSearchApi struct {
	flow.Api
}

func (a *WebSearchApi) OnCreate() {
	a.Client = conf.WebConf.Api["bochaai"]
}

func (a *WebSearchApi) Search(question string) (out []BingSearchResp, err error) {
	req := &WebSearchReq{
		Query:   question,
		Summary: true,
	}
	requestOpt := http2.HttpRequestOptions{
		RequestBody: req,
		Encode:      http2.EncodeJson,
		Headers: map[string]string{"Content-Type": "application/json",
			"Authorization": "Bearer " + env.Get("BACKEND_WEB_SEARCH_KEY", "")},
	}
	httpResult, err := a.ApiPostWithOpts("/v1/web-search", requestOpt)
	if err != nil {
		return nil, err
	}
	ans := new(BingAnswer)
	err = json.Unmarshal(httpResult.Data, &ans)
	if err != nil {
		return nil, err
	}
	for _, result := range ans.WebPages.Value {
		out = append(out, BingSearchResp{
			Title:   result.Name,
			Url:     result.URL,
			Content: result.Summary,
		})
	}
	return
}
