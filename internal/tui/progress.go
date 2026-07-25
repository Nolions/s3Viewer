package tui

import (
	"fmt"
	"strings"
)

// FormatBytes 將 byte 數轉換為可讀性佳的容量單位 (B, KB, MB, GB, TB)
func FormatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

// FormatProgressBar 生成終端機進度條字串
// 例如: Uploading test.zip [████████████░░░░░░░░] 60.0% (1.2 MB / 2.0 MB)
func FormatProgressBar(label string, currentBytes, totalBytes int64) string {
	if totalBytes <= 0 {
		return fmt.Sprintf("%s %s", label, FormatBytes(currentBytes))
	}

	percent := float64(currentBytes) / float64(totalBytes) * 100.0
	if percent > 100 {
		percent = 100
	}

	const barWidth = 20
	completed := int(percent / 100.0 * float64(barWidth))
	if completed > barWidth {
		completed = barWidth
	}

	bar := strings.Repeat("█", completed) + strings.Repeat("░", barWidth-completed)
	return fmt.Sprintf("%s [[yellow]%s[-]] [cyan]%.1f%%[-] (%s / %s)", label, bar, percent, FormatBytes(currentBytes), FormatBytes(totalBytes))
}
