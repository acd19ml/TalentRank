package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// DatabaseConfig 用于存储数据库配置
type DatabaseConfig struct {
	Database struct {
		Host     string `xml:"host"`
		Port     string `xml:"port"`
		Username string `xml:"username"`
		Password string `xml:"password"`
		Dbname   string `xml:"dbname"`
	} `xml:"database"`
}

type Ranking struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Score    int    `json:"score"`
	RankNo   int    `json:"rankno"`
}

type LocationCount struct {
	Location string `json:"location"`
	Count    int    `json:"count"`
}

// 从 XML 文件中读取数据库配置
func loadConfig(filename string) (DatabaseConfig, error) {
	var config DatabaseConfig
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = xml.Unmarshal(data, &config)
	return config, err
}

func main() {
	r := gin.Default()
	// CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 加载数据库配置
	config, err := loadConfig("./conf/db_config.xml")
	if err != nil {
		panic(err)
	}

	// 数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.Database.Username, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Dbname)
	fmt.Println("Connecting to database with DSN:", dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	r.GET("/api/rankings", func(c *gin.Context) {
		// 调用存储过程
		rows, err := db.Query("CALL GetRankings()")
		if err != nil {
			fmt.Println("Error calling stored procedure:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		defer rows.Close()

		var rankings []Ranking
		for rows.Next() {
			var ranking Ranking
			if err := rows.Scan(&ranking.Name, &ranking.Location, &ranking.Score, &ranking.RankNo); err != nil {
				fmt.Println("Error scanning row:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
			rankings = append(rankings, ranking)
		}

		c.JSON(http.StatusOK, rankings)
	})

	// 新增的路由来获取地区人数
	r.GET("/api/locations", func(c *gin.Context) {
		// 调用存储过程
		rows, err := db.Query("CALL GetLocation()")
		if err != nil {
			fmt.Println("Error calling stored procedure:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		defer rows.Close()

		var locationCounts []LocationCount
		for rows.Next() {
			var locationCount LocationCount
			if err := rows.Scan(&locationCount.Location, &locationCount.Count); err != nil {
				fmt.Println("Error scanning row:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
			locationCounts = append(locationCounts, locationCount)
		}

		c.JSON(http.StatusOK, locationCounts)
	})

	r.Run(":8080") // 在8080端口启动
}
