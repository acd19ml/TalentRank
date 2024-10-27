package git

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ContributorInfo 包含贡献者信息
type ContributorInfo struct {
	Name         string
	Contribution int // 贡献次数
	LinesAdded   int // 新增行数
	LinesDeleted int // 删除行数
}

// GetContributors 获取指定仓库的贡献者及其代码贡献数量
func GetContributors(username, repo string) (map[string]ContributorInfo, error) {
	url := fmt.Sprintf("https://github.com/%s/%s/graphs/contributors", username, repo)

	// 发起 GET 请求
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 检查请求状态
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch page: %d %s", res.StatusCode, res.Status)
	}

	// 使用 goquery 加载 HTML 文档
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	// 存储贡献者及其贡献信息的映射
	contributors := make(map[string]ContributorInfo)

	// 查找每个贡献者的元素
	doc.Find(".js-issue-row").Each(func(i int, s *goquery.Selection) {
		name := s.Find(".author").Text()
		contribCountText := s.Find(".contrib-number").Text()
		linesAddedText := s.Find(".contrib-graph").AttrOr("data-lines-added", "0")
		linesDeletedText := s.Find(".contrib-graph").AttrOr("data-lines-deleted", "0")

		contribCount := parseContributionCount(contribCountText)
		linesAdded := parseContributionCount(linesAddedText)
		linesDeleted := parseContributionCount(linesDeletedText)

		if name != "" {
			info := ContributorInfo{
				Name:         strings.TrimSpace(name),
				Contribution: contribCount,
				LinesAdded:   linesAdded,
				LinesDeleted: linesDeleted,
			}
			contributors[info.Name] = info
		}
	})

	return contributors, nil
}

// parseContributionCount 解析贡献数量
func parseContributionCount(countStr string) int {
	countStr = strings.TrimSpace(countStr)
	countStr = strings.ReplaceAll(countStr, ",", "") // 去除逗号

	var count int
	fmt.Sscanf(countStr, "%d", &count)
	return count
}

// GetTotalLines 获取总行数
func GetTotalLines(contributors map[string]ContributorInfo) (int, int) {
	totalLinesAdded := 0
	totalLinesDeleted := 0

	for _, info := range contributors {
		totalLinesAdded += info.LinesAdded
		totalLinesDeleted += info.LinesDeleted
	}

	return totalLinesAdded, totalLinesDeleted
}
