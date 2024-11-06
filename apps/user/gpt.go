package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Define a struct to represent the full request format
type Prompt struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type RequestBody struct {
	Prompt Prompt `json:"prompt"`
}

// PostAnalyze sends a POST request with a structured JSON request body
func PostAnalyze(input []byte) ([]byte, error) {
	url := "https://talent-rank-97292895702.asia-east2.run.app/analyze"

	// Print the JSON data for debugging
	// fmt.Println("Request Body:", string(input))

	// Send the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(input))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	// Read the response
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

// ExtractFields 从 JSON 字节数组中提取 possible_nation 和 confidence_level 的值
func ExtractFields(input []byte) (string, string) {
	// 定义一个结构来解析外层 JSON
	var outerResponse struct {
		Response string `json:"response"`
	}

	// 解析外层 JSON 以获取 response 字段的内容
	err := json.Unmarshal(input, &outerResponse)
	if err != nil {
		fmt.Printf("Error parsing outer JSON: %v\n", err)
		return "", ""
	}

	// 清理 response 内容，去除多余的字符如反引号和换行符
	innerJSON := outerResponse.Response
	innerJSON = innerJSON[8 : len(innerJSON)-4] // 去除 ```json\n 和 \n```

	// 定义一个结构来解析内部 JSON，将 confidence_level 设置为 json.Number
	var innerResponse struct {
		PossibleNation  string      `json:"possible_nation"`
		ConfidenceLevel json.Number `json:"confidence_level"`
	}

	// 解析内部 JSON
	err = json.Unmarshal([]byte(innerJSON), &innerResponse)
	if err != nil {
		fmt.Printf("Error parsing inner JSON: %v\n", err)
		return "", ""
	}

	// 将 confidence_level 转换为字符串
	confidenceLevel := innerResponse.ConfidenceLevel.String()

	return innerResponse.PossibleNation, confidenceLevel
}
