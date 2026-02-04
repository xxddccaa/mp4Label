package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config 表示应用配置
type Config struct {
	VideoDir        string `json:"video_dir"`         // 视频目录
	PreAnnotationDir string `json:"pre_annotation_dir"` // 预标注目录（可选）
	OutputDir       string `json:"output_dir"`         // 输出目录
	TaskFile        string `json:"task_file"`          // 子任务文件（可选），用于指定要标注的视频列表
}

var defaultConfig = Config{
	VideoDir: "/Users/xd/Downloads/process_mp4",
}

// GetConfigPath 获取配置文件路径
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".mp4label")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.Join(configDir, "config.json"), nil
}

// Load 加载配置
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := defaultConfig
		return &config, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// Save 保存配置
func Save(config *Config) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to save config file: %w", err)
	}

	return nil
}

// CleanPath 清理路径中的引号
func CleanPath(path string) string {
	path = strings.TrimSpace(path)
	
	// 移除前后的单引号或双引号
	if len(path) >= 2 {
		if (path[0] == '\'' && path[len(path)-1] == '\'') ||
			(path[0] == '"' && path[len(path)-1] == '"') {
			path = path[1 : len(path)-1]
		}
	}
	
	return path
}

// Normalize 规范化配置路径
func (c *Config) Normalize() {
	c.VideoDir = CleanPath(c.VideoDir)
	c.PreAnnotationDir = CleanPath(c.PreAnnotationDir)
	c.OutputDir = CleanPath(c.OutputDir)
	c.TaskFile = CleanPath(c.TaskFile)
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 先规范化路径
	c.Normalize()
	
	if c.VideoDir == "" {
		return fmt.Errorf("video directory cannot be empty")
	}

	if c.OutputDir == "" {
		return fmt.Errorf("output directory cannot be empty")
	}

	// 验证目录是否存在（如果已设置）
	if c.VideoDir != "" {
		if _, err := os.Stat(c.VideoDir); os.IsNotExist(err) {
			return fmt.Errorf("video directory does not exist: %s", c.VideoDir)
		}
	}

	if c.PreAnnotationDir != "" {
		if _, err := os.Stat(c.PreAnnotationDir); os.IsNotExist(err) {
			return fmt.Errorf("pre-annotation directory does not exist: %s", c.PreAnnotationDir)
		}
	}

	if c.TaskFile != "" {
		if _, err := os.Stat(c.TaskFile); os.IsNotExist(err) {
			return fmt.Errorf("task file does not exist: %s", c.TaskFile)
		}
	}

	return nil
}
