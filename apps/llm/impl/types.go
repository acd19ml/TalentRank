package impl

// InputData 定义输入 JSON 的结构
type InputData struct {
	Model    string              `json:"model"`
	Messages []map[string]string `json:"messages"`
}
