package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// TodoItem 表示一个待办事项
type TodoItem struct {
	UID         string
	Summary     string
	DTStart     string
	Status      string
	Priority    string
	DTStartFull string // 完整的 DTSTART 行（包括参数）
}

// ParseICalendar 解析 iCalendar 文本，返回过滤后的当日及之后的待办事项
func ParseICalendar(text string) (string, error) {
	// 处理行折叠，合并续行
	unfoldedLines := []string{}
	currentLine := ""
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		// 移除回车符
		line = strings.TrimRight(line, "\r")

		// 如果行以空格或制表符开头，是上一行的续行
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			currentLine += line[1:]
		} else {
			if currentLine != "" {
				unfoldedLines = append(unfoldedLines, currentLine)
			}
			currentLine = line
		}
	}
	if currentLine != "" {
		unfoldedLines = append(unfoldedLines, currentLine)
	}

	// 解析 iCalendar 文件
	todos := []TodoItem{}
	var currentTodo *TodoItem
	inVtodo := false

	for _, line := range unfoldedLines {
		colonIndex := strings.Index(line, ":")
		if colonIndex == -1 {
			continue
		}

		// 处理字段名中的参数（如 DTSTART;VALUE=DATE-TIME:...）
		fieldPart := line[:colonIndex]
		semicolonIndex := strings.Index(fieldPart, ";")
		fieldName := fieldPart
		if semicolonIndex >= 0 {
			fieldName = fieldPart[:semicolonIndex]
		}
		fieldName = strings.ToUpper(fieldName)
		fieldValue := line[colonIndex+1:]

		if fieldName == "BEGIN" && fieldValue == "VTODO" {
			inVtodo = true
			currentTodo = &TodoItem{}
		} else if fieldName == "END" && fieldValue == "VTODO" {
			if currentTodo != nil && currentTodo.Summary != "" && currentTodo.DTStart != "" {
				// 检查日期是否为今天或之后
				dtstart := parseIcalDateTime(currentTodo.DTStart)
				today := time.Now()
				today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
				todayTimestamp := today.Unix()

				if dtstart >= todayTimestamp {
					todos = append(todos, *currentTodo)
				}
			}
			inVtodo = false
			currentTodo = nil
		} else if inVtodo && currentTodo != nil {
			switch fieldName {
			case "SUMMARY":
				currentTodo.Summary = unescapeIcalValue(fieldValue)
			case "DTSTART":
				currentTodo.DTStart = fieldValue
				currentTodo.DTStartFull = line // 保存完整行以便后续使用
			case "STATUS":
				currentTodo.Status = fieldValue
			case "PRIORITY":
				currentTodo.Priority = fieldValue
			case "UID":
				currentTodo.UID = fieldValue
			}
		}
	}

	// 按开始时间排序
	sort.Slice(todos, func(i, j int) bool {
		timeA := parseIcalDateTime(todos[i].DTStart)
		timeB := parseIcalDateTime(todos[j].DTStart)
		return timeA < timeB
	})

	// 生成过滤后的 iCalendar 数据
	var icalContent strings.Builder
	icalContent.WriteString("BEGIN:VCALENDAR\r\n")
	icalContent.WriteString("VERSION:2.0\r\n")

	for _, todo := range todos {
		icalContent.WriteString("BEGIN:VTODO\r\n")
		if todo.UID != "" {
			icalContent.WriteString(fmt.Sprintf("UID:%s\r\n", todo.UID))
		}
		if todo.Summary != "" {
			icalContent.WriteString(fmt.Sprintf("SUMMARY:%s\r\n", escapeIcalValue(todo.Summary)))
		}
		if todo.DTStartFull != "" {
			// 使用完整的 DTSTART 行（包括参数）
			icalContent.WriteString(todo.DTStartFull + "\r\n")
		} else if todo.DTStart != "" {
			icalContent.WriteString(fmt.Sprintf("DTSTART:%s\r\n", todo.DTStart))
		}
		if todo.Status != "" {
			icalContent.WriteString(fmt.Sprintf("STATUS:%s\r\n", todo.Status))
		}
		if todo.Priority != "" {
			icalContent.WriteString(fmt.Sprintf("PRIORITY:%s\r\n", todo.Priority))
		}
		icalContent.WriteString("END:VTODO\r\n")
	}
	icalContent.WriteString("END:VCALENDAR\r\n")

	return icalContent.String(), nil
}

// parseIcalDateTime 解析 iCalendar 日期时间格式 (YYYYMMDDTHHMMSSZ 或 YYYYMMDD)
func parseIcalDateTime(dtstart string) int64 {
	if dtstart == "" {
		return 0
	}

	// 移除时区标识和T分隔符前的参数
	clean := dtstart
	if strings.Contains(clean, ":") {
		parts := strings.Split(clean, ":")
		if len(parts) > 1 {
			clean = parts[1] // 取冒号后的值
		}
	}
	// 移除时区标识
	clean = strings.TrimSuffix(clean, "Z")
	if len(clean) > 10 {
		// 移除时区偏移如 +0800 或 -0500
		if clean[len(clean)-5] == '+' || clean[len(clean)-5] == '-' {
			clean = clean[:len(clean)-5]
		}
	}

	var year, month, day, hour, minute, second int

	if len(clean) == 8 {
		// 日期格式 YYYYMMDD
		fmt.Sscanf(clean, "%4d%2d%2d", &year, &month, &day)
	} else if len(clean) >= 15 {
		// 日期时间格式 YYYYMMDDTHHMMSS
		fmt.Sscanf(clean, "%4d%2d%2dT%2d%2d%2d", &year, &month, &day, &hour, &minute, &second)
	} else {
		return 0
	}

	date := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
	return date.Unix()
}

// escapeIcalValue 转义 iCalendar 值
func escapeIcalValue(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, ";", "\\;")
	value = strings.ReplaceAll(value, ",", "\\,")
	value = strings.ReplaceAll(value, "\n", "\\n")
	return value
}

// unescapeIcalValue 反转义 iCalendar 值
func unescapeIcalValue(value string) string {
	value = strings.ReplaceAll(value, "\\n", "\n")
	value = strings.ReplaceAll(value, "\\,", ",")
	value = strings.ReplaceAll(value, "\\;", ";")
	value = strings.ReplaceAll(value, "\\\\", "\\")
	return value
}
