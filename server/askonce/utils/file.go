// Package utils -----------------------------
// @file      : file.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/16 10:46
// -------------------------------------------
package utils

import (
	"askonce/components"
	"fmt"
	"github.com/xiangtao94/golib/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"os"
)

func IsNetWorkUrlOrLocal(input string) (network bool, err error) {
	if isNetworkURL(input) {
		return true, nil
	}
	local, err := pathExists(input)
	if !local || err != nil {
		return false, components.ErrorFileNoExist
	}
	return false, nil
}

func isNetworkURL(input string) bool {
	parsedURL, err := url.Parse(input)
	if err != nil {
		return false
	}
	// 判断是否有 scheme 并且 scheme 是 http 或 https
	if parsedURL.Scheme == "http" || parsedURL.Scheme == "https" {
		return true
	}
	return false
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 下载 ZIP 文件到内存
func DownloadZip(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewError(500, fmt.Sprintf("failed to download zip: status %s", resp.Status))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
