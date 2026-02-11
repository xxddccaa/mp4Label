package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/xd/mp4label/pkg/annotation"
	"github.com/xd/mp4label/pkg/config"
	"github.com/xd/mp4label/pkg/video"
)

// Server represents the web server
type Server struct {
	config   *config.Config
	webFS    embed.FS
	readOnly bool // When true, all write operations are blocked (view mode)
}

// NewServer creates a new server instance in edit mode (loads config from file)
func NewServer(webFS embed.FS) (*Server, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &Server{
		config:   cfg,
		webFS:    webFS,
		readOnly: false,
	}, nil
}

// NewServerWithConfig creates a new server instance with a provided config.
// If readOnly is true, all write operations (save/delete annotation, save config, dialogs) are blocked.
func NewServerWithConfig(webFS embed.FS, cfg *config.Config, readOnly bool) (*Server, error) {
	return &Server{
		config:   cfg,
		webFS:    webFS,
		readOnly: readOnly,
	}, nil
}

// Start 启动服务器
func (s *Server) Start(port string) error {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/api/videos", s.handleVideos)
	http.HandleFunc("/api/annotation/", s.handleAnnotation)
	http.HandleFunc("/api/model-annotation/", s.handleModelAnnotation)
	http.HandleFunc("/api/video/", s.handleVideo)
	http.HandleFunc("/api/config", s.handleConfig)
	http.HandleFunc("/api/dialog", s.handleDialog)
	http.HandleFunc("/api/mode", s.handleMode)
	http.HandleFunc("/api/task-groups", s.handleTaskGroups)

	// 静态文件服务 - 使用嵌入的文件系统
	staticFS, err := fs.Sub(s.webFS, "web/static")
	if err != nil {
		return fmt.Errorf("failed to load static files: %w", err)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	addr := ":" + port
	log.Printf("服务器启动在 http://localhost%s", addr)
	log.Printf("请在浏览器中访问: http://localhost%s", addr)
	return http.ListenAndServe(addr, nil)
}

// handleIndex 处理主页
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// 从嵌入的文件系统读取 index.html
	data, err := s.webFS.ReadFile("web/index.html")
	if err != nil {
		http.Error(w, "Failed to load page", http.StatusInternalServerError)
		log.Printf("Failed to read index.html: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

// handleVideos handles video list requests
// Supports ?task=filename.txt query parameter to filter by a task file in TaskDir.
// When TaskDir is configured and no ?task= parameter is given, shows ALL videos (ignores TaskFile).
func (s *Server) handleVideos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Determine which task file to use for filtering
	taskFilePath := ""
	taskParam := r.URL.Query().Get("task")

	if taskParam != "" && s.config.TaskDir != "" {
		// Validate task parameter: no path traversal
		if strings.Contains(taskParam, "..") || strings.Contains(taskParam, "/") || strings.Contains(taskParam, "\\") {
			http.Error(w, "Invalid task parameter", http.StatusBadRequest)
			return
		}
		if !strings.HasSuffix(taskParam, ".txt") {
			http.Error(w, "Task file must be a .txt file", http.StatusBadRequest)
			return
		}
		// Build full path and verify it's within TaskDir
		fullPath := filepath.Join(s.config.TaskDir, taskParam)
		absTaskDir, _ := filepath.Abs(s.config.TaskDir)
		absFullPath, _ := filepath.Abs(fullPath)
		if !strings.HasPrefix(absFullPath, absTaskDir+string(filepath.Separator)) {
			http.Error(w, "Invalid task parameter", http.StatusBadRequest)
			return
		}
		// Verify file exists
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			http.Error(w, fmt.Sprintf("Task file not found: %s", taskParam), http.StatusNotFound)
			return
		}
		taskFilePath = fullPath
	} else if taskParam == "" && s.config.TaskDir != "" {
		// TaskDir is configured but no task selected => show ALL videos (ignore TaskFile)
		taskFilePath = ""
	} else {
		// No TaskDir configured, use single TaskFile if set
		taskFilePath = s.config.TaskFile
	}

	videos, err := video.ScanVideos(s.config.VideoDir, taskFilePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to scan videos: %v", err), http.StatusInternalServerError)
		return
	}

	// 匹配预标注和已有标注
	video.MatchAnnotations(videos, s.config.PreAnnotationDir, s.config.OutputDir)

	// 计算统计信息
	totalCount := len(videos)
	annotatedCount := 0
	preAnnotatedCount := 0

	for _, v := range videos {
		if v.HasAnnotation {
			annotatedCount++
		} else if v.HasPreAnnotation {
			preAnnotatedCount++
		}
	}

	// 返回视频列表和统计信息
	response := map[string]interface{}{
		"videos": videos,
		"stats": map[string]int{
			"total":         totalCount,
			"annotated":     annotatedCount,
			"pre_annotated": preAnnotatedCount,
			"unannotated":   totalCount - annotatedCount - preAnnotatedCount,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleAnnotation 处理标注请求
func (s *Server) handleAnnotation(w http.ResponseWriter, r *http.Request) {
	// 提取文件名
	filename := strings.TrimPrefix(r.URL.Path, "/api/annotation/")
	if filename == "" {
		http.Error(w, "Filename cannot be empty", http.StatusBadRequest)
		return
	}

	stem := strings.TrimSuffix(filename, filepath.Ext(filename))

	switch r.Method {
	case http.MethodGet:
		s.getAnnotation(w, r, stem)
	case http.MethodPost:
		s.saveAnnotation(w, r, stem)
	case http.MethodDelete:
		s.deleteAnnotation(w, r, stem)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getAnnotation 获取标注
func (s *Server) getAnnotation(w http.ResponseWriter, r *http.Request, stem string) {
	var ann *annotation.Annotation
	var err error

	// 优先从输出目录读取
	if s.config.OutputDir != "" {
		outputPath := video.GetAnnotationPath(stem, s.config.OutputDir)
		if outputPath != "" {
			if ann, err = annotation.ParseFile(outputPath); err == nil {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(ann)
				return
			}
		}
	}

	// 从预标注目录读取
	if s.config.PreAnnotationDir != "" {
		prePath := video.GetPreAnnotationPath(stem, s.config.PreAnnotationDir)
		if prePath != "" {
			if ann, err = annotation.ParseFile(prePath); err == nil {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(ann)
				return
			}
		}
	}

	// 都没有，返回空标注
	ann = &annotation.Annotation{
		Title:      "",
		IsTutorial: true,
		Steps:      []annotation.Step{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ann)
}

// saveAnnotation saves an annotation (blocked in read-only mode)
func (s *Server) saveAnnotation(w http.ResponseWriter, r *http.Request, stem string) {
	if s.readOnly {
		http.Error(w, "Cannot save annotation: server is in read-only view mode", http.StatusForbidden)
		return
	}
	if s.config.OutputDir == "" {
		http.Error(w, "Output directory not set", http.StatusBadRequest)
		return
	}

	var ann annotation.Annotation
	if err := json.NewDecoder(r.Body).Decode(&ann); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse request: %v", err), http.StatusBadRequest)
		return
	}

	// 验证标注
	if err := annotation.ValidateAnnotation(&ann); err != nil {
		http.Error(w, fmt.Sprintf("Annotation validation failed: %v", err), http.StatusBadRequest)
		return
	}

	// 保存文件
	outputPath := video.GetAnnotationPath(stem, s.config.OutputDir)
	if err := ann.Save(outputPath); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// deleteAnnotation deletes an annotation (blocked in read-only mode)
func (s *Server) deleteAnnotation(w http.ResponseWriter, r *http.Request, stem string) {
	if s.readOnly {
		http.Error(w, "Cannot delete annotation: server is in read-only view mode", http.StatusForbidden)
		return
	}
	if s.config.OutputDir == "" {
		http.Error(w, "Output directory not set", http.StatusBadRequest)
		return
	}

	outputPath := video.GetAnnotationPath(stem, s.config.OutputDir)
	if outputPath == "" {
		http.Error(w, "Annotation file does not exist", http.StatusNotFound)
		return
	}

	// 删除文件
	if err := os.Remove(outputPath); err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Annotation file does not exist", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to delete: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// handleVideo 处理视频文件请求
func (s *Server) handleVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := strings.TrimPrefix(r.URL.Path, "/api/video/")
	if filename == "" {
		http.Error(w, "Filename cannot be empty", http.StatusBadRequest)
		return
	}

	videoPath := filepath.Join(s.config.VideoDir, filename)

	// 安全检查：确保文件在视频目录内（使用绝对路径比较）
	absVideoDir, err := filepath.Abs(s.config.VideoDir)
	if err != nil {
		http.Error(w, "Failed to get video directory path", http.StatusInternalServerError)
		return
	}
	absVideoPath, err := filepath.Abs(videoPath)
	if err != nil {
		http.Error(w, "Failed to get video path", http.StatusInternalServerError)
		return
	}
	if !strings.HasPrefix(absVideoPath, absVideoDir+string(filepath.Separator)) && absVideoPath != absVideoDir {
		http.Error(w, "Illegal path", http.StatusForbidden)
		return
	}

	http.ServeFile(w, r, videoPath)
}

// handleModelAnnotation 处理模型标注请求（只读）
func (s *Server) handleModelAnnotation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 提取文件名
	filename := strings.TrimPrefix(r.URL.Path, "/api/model-annotation/")
	if filename == "" {
		http.Error(w, "Filename cannot be empty", http.StatusBadRequest)
		return
	}

	stem := strings.TrimSuffix(filename, filepath.Ext(filename))
	s.getModelAnnotation(w, r, stem)
}

// getModelAnnotation 获取模型标注
func (s *Server) getModelAnnotation(w http.ResponseWriter, r *http.Request, stem string) {
	// 如果没有配置模型标注目录，返回空结果
	if s.config.ModelAnnotationDir == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"available": false,
			"message":   "Model annotation directory not configured",
		})
		return
	}

	// 从模型标注目录读取
	modelPath := video.GetAnnotationPath(stem, s.config.ModelAnnotationDir)
	if modelPath == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"available": false,
			"message":   "Model annotation not found",
		})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"available": false,
			"message":   "Model annotation not found",
		})
		return
	}

	// 解析标注文件
	ann, err := annotation.ParseFile(modelPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse model annotation: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"available":  true,
		"annotation": ann,
	})
}

// handleConfig 处理配置请求
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getConfig(w, r)
	case http.MethodPost:
		s.saveConfig(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getConfig 获取配置
func (s *Server) getConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.config)
}

// saveConfig saves configuration (blocked in read-only mode)
func (s *Server) saveConfig(w http.ResponseWriter, r *http.Request) {
	if s.readOnly {
		http.Error(w, "Cannot save config: server is in read-only view mode", http.StatusForbidden)
		return
	}
	var cfg config.Config
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse request: %v", err), http.StatusBadRequest)
		return
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("Config validation failed: %v", err), http.StatusBadRequest)
		return
	}

	// 保存配置
	if err := config.Save(&cfg); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save config: %v", err), http.StatusInternalServerError)
		return
	}

	s.config = &cfg

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// handleDialog opens a native file/folder picker dialog (blocked in read-only mode)
func (s *Server) handleDialog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.readOnly {
		http.Error(w, "Cannot open dialog: server is in read-only view mode", http.StatusForbidden)
		return
	}

	mode := r.URL.Query().Get("mode")
	if mode != "file" && mode != "directory" {
		http.Error(w, "Invalid mode, must be 'file' or 'directory'", http.StatusBadRequest)
		return
	}

	path, err := openNativeDialog(mode)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to open dialog: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"path": path})
}

// handleTaskGroups returns the list of task group files in TaskDir
func (s *Server) handleTaskGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// If TaskDir is not configured, return unavailable
	if s.config.TaskDir == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"available": false,
			"groups":    []string{},
		})
		return
	}

	// Scan TaskDir for .txt files
	entries, err := os.ReadDir(s.config.TaskDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read task directory: %v", err), http.StatusInternalServerError)
		return
	}

	var groups []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".txt") {
			groups = append(groups, entry.Name())
		}
	}

	// groups are already sorted alphabetically by os.ReadDir

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"available": true,
		"groups":    groups,
	})
}

// handleMode returns the server mode (read-only or edit)
func (s *Server) handleMode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"readonly": s.readOnly})
}

func openNativeDialog(mode string) (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return openDialogMac(mode)
	case "windows":
		return openDialogWindows(mode)
	default:
		return openDialogLinux(mode)
	}
}

func openDialogMac(mode string) (string, error) {
	var script string
	if mode == "directory" {
		script = `POSIX path of (choose folder)`
	} else {
		script = `POSIX path of (choose file)`
	}
	out, err := runDialogCommand("osascript", "-e", script)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func openDialogWindows(mode string) (string, error) {
	var script string
	if mode == "directory" {
		script = `
Add-Type -AssemblyName System.Windows.Forms
$dialog = New-Object System.Windows.Forms.FolderBrowserDialog
if ($dialog.ShowDialog() -eq "OK") { $dialog.SelectedPath }
`
	} else {
		script = `
Add-Type -AssemblyName System.Windows.Forms
$dialog = New-Object System.Windows.Forms.OpenFileDialog
if ($dialog.ShowDialog() -eq "OK") { $dialog.FileName }
`
	}
	out, err := runDialogCommand("powershell", "-NoProfile", "-Command", script)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func openDialogLinux(mode string) (string, error) {
	var args []string
	if mode == "directory" {
		args = []string{"--file-selection", "--directory"}
	} else {
		args = []string{"--file-selection"}
	}
	out, err := runDialogCommand("zenity", args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func runDialogCommand(name string, args ...string) ([]byte, error) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return []byte(""), nil
		}
		return nil, err
	}
	return out, nil
}
