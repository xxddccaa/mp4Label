package annotation

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Annotation 表示一个标注文件的内容
type Annotation struct {
	Title      string   `json:"title"`       // 教程题目（教学视频）或空（非教学视频）
	IsTutorial bool     `json:"is_tutorial"` // 是否为教学视频
	Steps      []Step   `json:"steps"`       // 步骤列表（仅教学视频）
}

// Step 表示一个操作步骤
type Step struct {
	Number      int    `json:"number"`      // 步骤编号
	Timestamp   string `json:"timestamp"`   // 时间戳 (mm:ss.SSS)
	Description string `json:"description"` // 步骤描述
}

// ParseFile 解析标注文件
func ParseFile(filePath string) (*Annotation, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return ParseLines(lines)
}

// ParseLines 解析标注内容（行列表）
func ParseLines(lines []string) (*Annotation, error) {
	if len(lines) == 0 {
		return &Annotation{
			Title:      "",
			IsTutorial: true,
			Steps:      []Step{},
		}, nil
	}

	// 检查是否为非教学视频
	firstLine := strings.TrimSpace(lines[0])
	if strings.ToLower(firstLine) == "[not tutorial]" {
		return &Annotation{
			Title:      "",
			IsTutorial: false,
			Steps:      []Step{},
		}, nil
	}

	// 解析教学视频格式
	ann := &Annotation{
		IsTutorial: true,
		Steps:      []Step{},
	}

	// 第一行是题目
	if len(lines) > 0 {
		ann.Title = strings.TrimSpace(lines[0])
	}

	// 步骤格式：1) 00:11 描述 或 1) 00:11.123 描述
	// 支持两种格式：mm:ss 和 mm:ss.SSS
	stepPattern := regexp.MustCompile(`^(\d+)\)\s+(\d{2}:\d{2}(?:\.\d{3})?)\s+(.+)$`)

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		matches := stepPattern.FindStringSubmatch(line)
		if len(matches) == 4 {
			var number int
			fmt.Sscanf(matches[1], "%d", &number)
			timestamp := matches[2]
			
			// 如果时间戳不包含毫秒，自动添加 .000
			if !strings.Contains(timestamp, ".") {
				timestamp = timestamp + ".000"
			}
			
			description := matches[3]

			ann.Steps = append(ann.Steps, Step{
				Number:      number,
				Timestamp:   timestamp,
				Description: description,
			})
		}
	}

	return ann, nil
}

// Format 将标注格式化为文本
func (a *Annotation) Format() string {
	if !a.IsTutorial {
		return "[not tutorial]\n"
	}

	var sb strings.Builder
	if a.Title != "" {
		sb.WriteString(a.Title)
		sb.WriteString("\n\n")
	}

	for _, step := range a.Steps {
		sb.WriteString(fmt.Sprintf("%d) %s %s\n", step.Number, step.Timestamp, step.Description))
	}

	return sb.String()
}

// Save 保存标注到文件
func (a *Annotation) Save(filePath string) error {
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	content := a.Format()
	return os.WriteFile(filePath, []byte(content), 0644)
}
