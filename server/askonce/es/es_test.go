// Package es -----------------------------
// @file      : es_test.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/13 19:21
// -------------------------------------------
package es

import (
	"askonce/test"
	"testing"
)

func TestCommonIndexCreate(t *testing.T) {
	test.Init()
	err := DocIndexCreate(test.Ctx, "demo_test")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCommonIndexDelete(t *testing.T) {
	test.Init()
	err := DocIndexDelete(test.Ctx, "demo_test")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCommonDocumentDelete(t *testing.T) {
	test.Init()
	err := CommonDocumentDelete(test.Ctx, "demo_test", []int64{1})
	if err != nil {
		t.Fatal(err)
	}
}
