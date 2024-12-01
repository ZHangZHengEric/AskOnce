package web_search

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
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

func BingSearch(ctx *gin.Context, question string) (out []BingSearchResp, err error) {
	const endpoint = "https://api.bing.microsoft.com/v7.0/search"
	token := "29613074456942c090da3dc819c52ef4"
	// Declare a new GET request.
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	// Add the payload to the request.
	param := req.URL.Query()
	param.Add("q", question)
	req.URL.RawQuery = param.Encode()
	// Insert the request header.
	req.Header.Add("Ocp-Apim-Subscription-Key", token)
	// Create a new client.
	client := http.Client{Timeout: 15 * time.Second}
	// Send the request to Bing.
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// Close the response.
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// Create a new answer.
	ans := new(BingAnswer)
	err = json.Unmarshal(body, &ans)
	if err != nil {
		return nil, err
	}
	for _, result := range ans.WebPages.Value {
		out = append(out, BingSearchResp{
			Title:   result.Name,
			Url:     result.URL,
			Content: result.Snippet,
		})
	}
	return
}
