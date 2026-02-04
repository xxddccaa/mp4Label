package annotation

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidateTimestamp 验证时间戳格式 (mm:ss.SSS 或 mm:ss)
func ValidateTimestamp(timestamp string) error {
	// 支持两种格式：mm:ss.SSS 或 mm:ss
	pattern := regexp.MustCompile(`^\d{2}:\d{2}(\.\d{3})?$`)
	if !pattern.MatchString(timestamp) {
		return fmt.Errorf("invalid timestamp format, should be mm:ss.SSS (e.g., 12:32.766)")
	}

	// 验证分钟和秒数范围
	parts := strings.Split(timestamp, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid timestamp format")
	}
	
	var minutes int
	fmt.Sscanf(parts[0], "%d", &minutes)
	
	// 解析秒和毫秒
	secondsParts := strings.Split(parts[1], ".")
	var seconds int
	fmt.Sscanf(secondsParts[0], "%d", &seconds)
	
	if seconds >= 60 {
		return fmt.Errorf("seconds cannot exceed 59")
	}
	
	// 如果有毫秒，验证毫秒范围
	if len(secondsParts) > 1 {
		var millis int
		fmt.Sscanf(secondsParts[1], "%d", &millis)
		if millis >= 1000 {
			return fmt.Errorf("milliseconds cannot exceed 999")
		}
	}

	return nil
}

// ValidateStep 验证步骤格式
func ValidateStep(step Step) error {
	if step.Number <= 0 {
		return fmt.Errorf("step number must be greater than 0")
	}

	if err := ValidateTimestamp(step.Timestamp); err != nil {
		return err
	}

	if strings.TrimSpace(step.Description) == "" {
		return fmt.Errorf("step description cannot be empty")
	}

	return nil
}

// ValidateAnnotation 验证标注内容
func ValidateAnnotation(ann *Annotation) error {
	if ann == nil {
		return fmt.Errorf("annotation cannot be nil")
	}

	if !ann.IsTutorial {
		// 非教学视频，无需验证
		return nil
	}

	// 教学视频验证
	if strings.TrimSpace(ann.Title) == "" {
		return fmt.Errorf("tutorial title cannot be empty")
	}

	if len(ann.Title) > 100 {
		return fmt.Errorf("tutorial title cannot exceed 100 characters")
	}

	if len(ann.Steps) == 0 {
		return fmt.Errorf("at least one step is required")
	}

	// 验证每个步骤
	for i, step := range ann.Steps {
		if err := ValidateStep(step); err != nil {
			return fmt.Errorf("step %d validation failed: %w", i+1, err)
		}

		// 验证步骤编号是否连续
		if step.Number != i+1 {
			return fmt.Errorf("step numbers not consecutive, expected %d, got %d", i+1, step.Number)
		}
	}

	return nil
}
