package helpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/xiangtao94/golib/pkg/zlog"
	"time"
)

func BaiduHotTest(ctx *gin.Context) (data []string) {
	collector := colly.NewCollector(
		func(collector *colly.Collector) {
			extensions.RandomUserAgent(collector)
		},
		func(c *colly.Collector) {
			c.OnRequest(func(request *colly.Request) {
				fmt.Println(request.URL, ", User-Agent:", request.Headers.Get("User-Agent"))
			})
		},
	)
	collector.SetRequestTimeout(time.Second * 60)

	data = []string{}

	collector.OnHTML(".container-bg_lQ801", func(element *colly.HTMLElement) {
		element.ForEach(".category-wrap_iQLoo", func(i int, element *colly.HTMLElement) {
			title := element.ChildText(".content_1YWBm .c-single-text-ellipsis")
			data = append(data, title)
		})
	})
	if err := collector.Visit("https://top.baidu.com/board?tab=realtime"); err != nil {
		zlog.Errorf(ctx, "查询热榜失败")
		return
	}
	return
}
