package utils

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// GetDependentRepositories 获取依赖仓库数量
func GetDependentRepositories(url string) (int, error) {
	// 发起 GET 请求
	res, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	// 检查请求状态
	if res.StatusCode != 200 {
		return 0, fmt.Errorf("failed to fetch page: %d %s", res.StatusCode, res.Status)
	}

	// 使用 goquery 加载 HTML 文档
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return 0, err
	}

	// 正则表达式提取数字部分（去除逗号）
	repoCount := 0
	re := regexp.MustCompile(`(\d{1,3}(?:,\d{3})*)\s+Repositories`)

	// 查找包含依赖仓库数量的元素
	doc.Find("a.btn-link.selected").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			// 去除逗号后解析数字
			countStr := strings.ReplaceAll(matches[1], ",", "")
			fmt.Sscanf(countStr, "%d", &repoCount)
		}
	})

	return repoCount, nil
}
