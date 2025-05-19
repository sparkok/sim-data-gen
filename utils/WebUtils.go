package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

// PostRequest 发送POST请求
// urlStr: 请求地址
// data: 请求参数（map类型）
// contentType: 请求体类型（支持 form 和 json）
func PostRequest(urlStr string, data map[string]interface{}, contentType string) (*http.Response, error) {
	var (
		req *http.Request
		err error
	)

	switch contentType {
	case "json":
		// JSON格式处理
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest("POST", urlStr, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

	case "form": // 默认表单格式处理
		formData := url.Values{}
		for k, v := range data {
			formData.Add(k, v.(string))
		}
		req, err = http.NewRequest("POST", urlStr, strings.NewReader(formData.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	default:
		// 自定义contentType处理
		req, err = http.NewRequest("POST", urlStr, strings.NewReader(url.Values{}.Encode()))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
