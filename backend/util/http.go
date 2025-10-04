package util

import (
	"JetGet/backend/pget"
	"github.com/pkg/errors"
	"log"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

func DoHeadRequest(url, proxy string) (*http.Response, error) {
	// 查询文件大小
	client := pget.NewClientByProxy(16, proxy)
	r, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		log.Println("new request failed:", err)
	}
	res, err := client.Do(r)
	if err != nil {
		log.Println("failed to head request:", err)
		return res, err
	}

	if res.Header.Get("Accept-Ranges") != "bytes" {
		return res, errors.New("does not support range request")
	}
	if res.ContentLength <= 0 {
		return res, errors.New("invalid content length")
	}
	return res, nil
}

func GetFileName(resp *http.Response) string {
	// 优先检查 Content-Disposition header
	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		// 使用 mime 包来正确解析可能包含特殊字符的 header
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err == nil {
			// filename* 优先，因为它支持 UTF-8
			if filename, ok := params["filename*"]; ok {
				// filename* 的格式是 charset''encoded-value
				if parts := strings.SplitN(filename, "''", 2); len(parts) == 2 {
					decodedFilename, err := url.QueryUnescape(parts[1])
					if err == nil {
						return decodedFilename
					}
				}
			}
			// 其次是普通的 filename
			if filename, ok := params["filename"]; ok {
				return filename
			}
		}
	}
	// 如果 header 不存在，则从最终的 URL 路径中提取 (使用 resp.Request.URL 可以获取重定向之后的最终 URL)
	finalURL := resp.Request.URL.Path
	filename := filepath.Base(finalURL)
	if filename != "." && filename != "/" {
		return filename
	}
	return "unrecognized"
}
