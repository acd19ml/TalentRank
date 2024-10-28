package main

import (
	"context"
	"fmt"
	"log"

	"github.com/acd19ml/TalentRank/utils/git"

	"github.com/acd19ml/TalentRank/utils"
)

func main() {
	service := git.NewGitClient()
	names, err := service.GetOrganizations(context.Background(), utils.Username)
	if err != nil {
		log.Fatalf("Error getting name: %v", err)
	}
	for _, name := range names {
		fmt.Printf("Organization for user %s: %s\n", utils.Username, name)
	}
}

// 测试进度条
//package main

//import (
//	"fmt"
//	"net/http"
//	"time"
//
//	"github.com/gin-gonic/gin"
//)
//
//var progressChan = make(chan int)
//
//func main() {
//	router := gin.Default()
//	router.LoadHTMLFiles("index.html")
//
//	router.GET("/", func(c *gin.Context) {
//		c.HTML(http.StatusOK, "index.html", nil)
//	})
//
//	router.POST("/start", func(c *gin.Context) {
//		go func() {
//			for i := 0; i <= 10; i++ {
//				time.Sleep(200 * time.Millisecond) // 增加延迟到200毫秒
//				progressChan <- i                  // 将进度发送到 channel
//			}
//		}()
//		c.Status(http.StatusOK)
//	})
//
//	router.GET("/progress", func(c *gin.Context) {
//		c.Writer.Header().Set("Content-Type", "text/event-stream")
//		c.Writer.Header().Set("Cache-Control", "no-cache")
//		c.Writer.Header().Set("Connection", "keep-alive")
//
//		for progress := range progressChan {
//			fmt.Fprintf(c.Writer, "data: %d%%\n\n", progress) // 发送进度
//			if flusher, ok := c.Writer.(http.Flusher); ok {
//				flusher.Flush() // 刷新数据到客户端
//			}
//		}
//	})
//
//	router.Run(":8080")
//}
