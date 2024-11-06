package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type OpenAIRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func callOpenAI(prompt string) (string, error) {

	apiKey := os.Getenv("GPT_API") // 确保将API密钥存入环境变量
	if apiKey == "" {
		log.Fatal("API key not found")
	}
	url := "https://api.openai.com/v1/chat/completions"

	// 使用GPT-4模型构造请求数据，并添加system消息
	requestData := OpenAIRequest{
		Model: "gpt-4o", // 使用GPT-4模型
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant to guess the nation and confidence_level for a github user by multiple data from this user."}, // 添加system消息
			{Role: "user", Content: prompt},
		},
		MaxTokens: 100, // 设置返回内容的最大字符数
	}

	// 将请求数据转换为JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求并获取响应
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析响应
	var response OpenAIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	// 返回模型生成的文本
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("no response from OpenAI")
}

func main() {
	response, err := callOpenAI(prompt + promptJson)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Response:", response)
	}
}

const promptJson = `[ { "id": "ad1ba224-3b41-4129-992a-7606fb3c2a37", "username": "Livinglist", "name": "Jojo Feng", "company": "", "blog": "", "location": "", "email": "", "bio": "A software engineer passionate about and specializing in mobile application development. Check out the pinned repos for apps I built. ", "followers": 112, "organizations": "null", "readme": "# 出入\n\n![iOS](https://img.shields.io/badge/iOS-13%20-blue)\n[![App Store](https://img.shields.io/itunes/v/1522619720?label=App%20Store)](https://apps.apple.com/us/app/出入-简易记账/id152261972...![image](https://github.com/Livinglist/CoronaCounter/blob/master/screenshot.png?raw=true)\n\n## Welcome to CoronaCounter\n\n### installation\n\n- Clone this repo and build it from Xcode.\n\nOR\n\n- Download the...# Dumbbell\n\n![iOS](https://img.shields.io/badge/iOS-12%20-blue)\n[![App Store](https://img.shields.io/itunes/v/1462586545?label=App%20Store)](https://apps.apple.com/us/app/dumbbell-workout-planner/id14...# feature_discovery\n\nThis Flutter package implements Feature Discovery following the [Material Design guidelines](https://material.io/archive/guidelines/growth-communications/feature-discovery.html). ...# flutter_siri_suggestions\n\nFlutter plugin for exposure on Siri Suggestions.\n\n<img src=\"https://img.shields.io/pub/v/flutter_siri_suggestions.svg\" />\n<img src=\"https://img.shields.io/github/license/my...\n# <img width=\"64\" src=\"https://user-images.githubusercontent.com/7277662/167775086-0b234f28-dee4-44f6-aae4-14a28ed4bbb6.png\"> Hacki for Hacker News\n\nA [Hacker News](https://news.ycombinator.com/) cli...# Kaby\n\n![iOS](https://img.shields.io/badge/iOS-12%20-blue)\n[![App Store](https://img.shields.io/itunes/v/1504762505?label=App%20Store)](https://apps.apple.com/us/app/kaby-easy-task-management/id15047...## Stat\n<p align=\"center\">\n <a href=\"https://github.com/livinglist?tab=repositories\">\n<img src=\"https://github-readme-stats.vercel.app/api?username=livinglist&&show_icons=true&title_color=ffffff&icon...# Manji\n\n\n![iOS](https://img.shields.io/badge/iOS-11%20-blue)\n[![App Store](https://img.shields.io/itunes/v/1464774967?label=App%20Store)](https://apps.apple.com/us/app/manji-learn-kanji/id1464774967#...", "commits": "Repo: Churu\nChanged to auto size text\nChanged formatter.\nChanged padding.\nChanged font size.\nFixed button color.\nFixed potential scrollview exception,\nMerge branch 'master' of https://github.com/Livinglist/Churu\n\nRepo: CoronaCounter\nMerge pull request #2 from gaibb/master\n\ncode optimization\nFix bug\nFetch data from wikipedia instead of google's corona site since the later one is javascript rendered.\n\nRepo: CoronaCounterDMG\n1.2\n1.1\n1.0.6\n1.0.5+\n1.0.5\n1.0.4\n1.0.3\n1.0.1\n1.0.0\n\nRepo: Dumbbell\nChanged fontsize.\nFixed the bug where empty part_card still shows up after the user quits the edit page.\nIncreased font size. Minor UI changes.\nUpdate README.md\nUpdate README.md\nChanged setsLeft.\n\nRepo: feature_discovery\nupdate.\n\nRepo: flutter_siri_suggestions\nbump kotlin version.\nbump kotlin version.\nremove android support.\n\nRepo: Hacki\nchore: bump fastlane version. (#486)\nfix comment parser. (#485)\nfeat: show msg if no favorites. (#481)\nfix: favorites screen. (#480)\n\nRepo: Kaby\nFixed the issue where removing project does not refresh the main_page\nMerge branch 'master' of https://github.com/Livinglist/kanban\nMigrated to flutter v2.\nUpdate README.md\nUI changes\nUpdate README.md...\n\nRepo: Livinglist\nupdate README.md\nupdate README.md\nupdate README.md.\nupdate README.md.\nupdate README.md\nupdate README.md\nupdate README.md\nupdated README.md\nupdated README.md\nupdated README.md\n\nRepo: Manji\nUpdate README.md\nMerge branch 'master' of github.com:Livinglist/Manji\nFix format.\nUpdate README.md\nUpdate README.md\nAdd effective dart package.\nEnforce effective dart rules.\n\n", "score": 44285.61129253463, "possible_nation": "", "confidence_level": "" } ]`

const prompt = `Based on the provided JSON information of a GitHub user (id, username, name, company, blog, location, email, bio, followers, organizations, readme, commits, and score), if the location field is empty, analyze and infer the user’s possible nation (location) based on the other available fields. Only respond with a JSON object in the following format: 

{
    "possible_nation": "<country or 'N/A'>",
    "confidence_level": <percentage as a number>
}

If the information is insufficient to determine the user's nation, set "possible_nation" to "N/A" and "confidence_level" to 0.`
