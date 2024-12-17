// Package data -----------------------------
// @file      : parse.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/17 18:00
// -------------------------------------------
package content_parse

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type ParseResult struct {
}

type Parser interface {
	parse(ctx *gin.Context) (result *ParseResult, err error)
}

var allowExtension = []string{".pdf", ".doc", ".docx", ".txt", ".ppt", ".pptx", ".xlsx", ".xls", ".json"}

func NewTxtContentParser(ext string) (Parser, error) {
	switch ext {
	case "txt":
		return new(TxtParser), nil
	case "doc":
	case "docx":
		return new(DocParser), nil
	case "pdf":
		return new(PdfParser), nil
	default:
		return nil, fmt.Errorf("unknown content type: %s", ext)
	}
}
