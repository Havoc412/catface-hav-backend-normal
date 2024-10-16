package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Config 包含所有配置部分
type Config struct {
	MySQL MySQLConfig `json:"mysql"`
}

// MySQLConfig 用于存储 MySQL 数据库的配置
type MySQLConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Database string `json:"database"`
}

// LoadConfig 从文件中加载所有配置信息
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("could not decode config file: %v", err)
	}

	return &config, nil
}

func main() {
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 使用 MySQL 配置信息
	fmt.Printf("Connecting to MySQL database at %s\n", config.MySQL.Host)
	// 使用 config.MySQL.Username, config.MySQL.Password, config.MySQL.Database 来连接数据库
}
