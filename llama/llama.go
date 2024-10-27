package llama

import (
	"encoding/json"
	// "fmt"

	"github.com/go-resty/resty/v2"
)

func analyzeDataWithLlama(data []string, apiKey string) ([]float64, error) {
    client := resty.New()	
    var scores []float64

    for _, text := range data {
        resp, err := client.R().
            SetHeader("Authorization", "Bearer "+apiKey).
            SetBody(map[string]interface{}{"text": text}).
            Post("https://api.llama.ai/analyze")

        if err != nil {
            return nil, err
        }

        // 假设API返回JSON并包含评分字段'score'
        var result map[string]float64
        if err := json.Unmarshal(resp.Body(), &result); err != nil {
            return nil, err
        }
        
        scores = append(scores, result["score"])
    }
    return scores, nil
}
