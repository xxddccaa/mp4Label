package video

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// VideoInfo 表示视频文件信息
type VideoInfo struct {
	Filename      string `json:"filename"`       // 文件名（含扩展名）
	Stem         string `json:"stem"`          // 文件名（不含扩展名）
	Path         string `json:"path"`          // 完整路径
	HasPreAnnotation bool `json:"has_pre_annotation"` // 是否有预标注
	HasAnnotation    bool `json:"has_annotation"`     // 是否已有标注
}

// ScanVideos 扫描视频目录，返回视频列表
// 如果 taskFile 不为空，则只返回 taskFile 中列出的视频
func ScanVideos(videoDir string, taskFile string) ([]VideoInfo, error) {
	if videoDir == "" {
		return []VideoInfo{}, nil
	}

	if _, err := os.Stat(videoDir); os.IsNotExist(err) {
		return []VideoInfo{}, fmt.Errorf("video directory does not exist: %s", videoDir)
	}

	// 读取任务文件（如果提供）
	var taskVideos map[string]bool
	if taskFile != "" {
		var err error
		taskVideos, err = loadTaskFile(taskFile)
		if err != nil {
			return []VideoInfo{}, fmt.Errorf("failed to load task file: %w", err)
		}
	}

	var videos []VideoInfo
	err := filepath.Walk(videoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".mp4" {
			filename := filepath.Base(path)
			stem := strings.TrimSuffix(filename, filepath.Ext(filename))

			// 如果有任务文件，检查该视频是否在任务列表中
			if taskVideos != nil {
				if !taskVideos[stem] {
					return nil // 跳过不在任务列表中的视频
				}
			}

			videos = append(videos, VideoInfo{
				Filename: filename,
				Stem:     stem,
				Path:     path,
			})
		}

		return nil
	})

	return videos, err
}

// loadTaskFile 加载任务文件，返回视频名称（不含扩展名）的集合
func loadTaskFile(taskFile string) (map[string]bool, error) {
	file, err := os.Open(taskFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	taskVideos := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 跳过空行
		if line == "" {
			continue
		}
		// 移除可能的 .mp4 后缀
		stem := strings.TrimSuffix(line, ".mp4")
		taskVideos[stem] = true
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return taskVideos, nil
}

// MatchAnnotations 匹配预标注和已有标注
func MatchAnnotations(videos []VideoInfo, preAnnotationDir, outputDir string) {
	// 创建预标注文件映射
	preAnnotationMap := make(map[string]bool)
	if preAnnotationDir != "" {
		if entries, err := os.ReadDir(preAnnotationDir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.ToLower(filepath.Ext(entry.Name())) == ".txt" {
					stem := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
					preAnnotationMap[stem] = true
				}
			}
		}
	}

	// 创建已有标注文件映射
	annotationMap := make(map[string]bool)
	if outputDir != "" {
		if entries, err := os.ReadDir(outputDir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.ToLower(filepath.Ext(entry.Name())) == ".txt" {
					stem := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
					annotationMap[stem] = true
				}
			}
		}
	}

	// 更新视频信息
	for i := range videos {
		videos[i].HasPreAnnotation = preAnnotationMap[videos[i].Stem]
		videos[i].HasAnnotation = annotationMap[videos[i].Stem]
	}
}

// GetAnnotationPath 获取标注文件路径
func GetAnnotationPath(stem, dir string) string {
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, stem+".txt")
}

// GetPreAnnotationPath 获取预标注文件路径
func GetPreAnnotationPath(stem, dir string) string {
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, stem+".txt")
}
