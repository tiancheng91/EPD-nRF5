#include "icalendar.h"

#include <ctype.h>
#include <string.h>

#include "nrf_log.h"
#include "todolist.h"

// 将iCalendar日期时间字符串转换为Unix timestamp
// 支持格式: YYYYMMDDTHHMMSSZ 或 YYYYMMDD
static uint32_t parse_icalendar_datetime(const char* str, uint16_t len) {
    if (str == NULL || len < 8) {
        return 0;
    }

    // 解析日期部分 YYYYMMDD
    char year_str[5] = {0};
    char month_str[3] = {0};
    char day_str[3] = {0};
    char hour_str[3] = {0};
    char min_str[3] = {0};
    char sec_str[3] = {0};

    if (len < 8) {
        return 0;
    }

    memcpy(year_str, str, 4);
    memcpy(month_str, str + 4, 2);
    memcpy(day_str, str + 6, 2);

    uint16_t year = (uint16_t)atoi(year_str);
    uint8_t month = (uint8_t)atoi(month_str);
    uint8_t day = (uint8_t)atoi(day_str);
    uint8_t hour = 0;
    uint8_t min = 0;
    uint8_t sec = 0;

    // 如果是日期时间格式 YYYYMMDDTHHMMSSZ
    if (len >= 16 && str[8] == 'T') {
        memcpy(hour_str, str + 9, 2);
        memcpy(min_str, str + 11, 2);
        memcpy(sec_str, str + 13, 2);
        hour = (uint8_t)atoi(hour_str);
        min = (uint8_t)atoi(min_str);
        sec = (uint8_t)atoi(sec_str);
    }

    // 转换为Unix timestamp
    // 使用简单的计算方式（从1970-01-01开始）
    // 注意：这里使用简化算法，实际应该考虑时区等
    uint32_t timestamp = 0;
    uint16_t y;
    uint8_t m;

    // 计算从1970年到目标年份的天数
    for (y = 1970; y < year; y++) {
        if ((y % 4 == 0 && y % 100 != 0) || (y % 400 == 0)) {
            timestamp += 366 * 86400;  // 闰年
        } else {
            timestamp += 365 * 86400;  // 平年
        }
    }

    // 计算从1月到目标月份的天数
    uint8_t days_in_month[] = {31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31};
    bool is_leap = (year % 4 == 0 && year % 100 != 0) || (year % 400 == 0);
    if (is_leap) {
        days_in_month[1] = 29;
    }

    for (m = 1; m < month; m++) {
        timestamp += days_in_month[m - 1] * 86400;
    }

    // 计算天数
    timestamp += (day - 1) * 86400;

    // 计算小时、分钟、秒
    timestamp += hour * 3600;
    timestamp += min * 60;
    timestamp += sec;

    return timestamp;
}

// 不区分大小写的字符串比较
static bool strcasecmp_prefix(const char* str, const char* prefix, uint16_t len) {
    uint16_t i;
    for (i = 0; prefix[i] != '\0' && i < len; i++) {
        if (tolower(str[i]) != tolower(prefix[i])) {
            return false;
        }
    }
    return prefix[i] == '\0';
}

// 跳过空白字符
static const char* skip_whitespace(const char* str) {
    while (*str == ' ' || *str == '\t') {
        str++;
    }
    return str;
}

// 解析iCalendar文本中的一行（处理行折叠）
static const char* parse_line(const char* text, uint16_t text_len, uint16_t* pos, char* line_buf, uint16_t line_buf_size) {
    uint16_t line_pos = 0;
    bool continuation = false;

    while (*pos < text_len && line_pos < line_buf_size - 1) {
        char c = text[*pos];

        // 处理行折叠：以空格或制表符开头的行是上一行的续行
        if (c == '\r' || c == '\n') {
            (*pos)++;
            if (c == '\r' && *pos < text_len && text[*pos] == '\n') {
                (*pos)++;
            }

            // 检查下一行是否是续行
            if (*pos < text_len && (text[*pos] == ' ' || text[*pos] == '\t')) {
                continuation = true;
                // 跳过续行的前导空白
                while (*pos < text_len && (text[*pos] == ' ' || text[*pos] == '\t')) {
                    (*pos)++;
                }
                continue;
            } else {
                // 行结束
                break;
            }
        }

        line_buf[line_pos++] = c;
        (*pos)++;
    }

    line_buf[line_pos] = '\0';
    return line_buf;
}

// 解析字段值（处理转义字符）
static void unescape_value(const char* src, char* dst, uint16_t dst_size) {
    uint16_t i = 0;
    uint16_t j = 0;

    while (src[i] != '\0' && j < dst_size - 1) {
        if (src[i] == '\\') {
            i++;
            if (src[i] == 'n') {
                dst[j++] = '\n';
                i++;
            } else if (src[i] == '\\') {
                dst[j++] = '\\';
                i++;
            } else if (src[i] == ';') {
                dst[j++] = ';';
                i++;
            } else if (src[i] == ',') {
                dst[j++] = ',';
                i++;
            } else {
                dst[j++] = src[i++];
            }
        } else {
            dst[j++] = src[i++];
        }
    }
    dst[j] = '\0';
}

ret_code_t todolist_parse_icalendar(uint8_t* data, uint16_t length) {
    if (data == NULL || length == 0) {
        return NRF_ERROR_INVALID_PARAM;
    }

    todolist_t todolist;
    memset(&todolist, 0, sizeof(todolist_t));

    char line_buf[256];
    uint16_t pos = 0;
    bool in_vtodo = false;
    todo_item_t current_item;
    memset(&current_item, 0, sizeof(todo_item_t));

    while (pos < length && todolist.count < TODO_MAX_ITEMS) {
        const char* line = parse_line((const char*)data, length, &pos, line_buf, sizeof(line_buf));
        if (line[0] == '\0') {
            continue;
        }

        // 查找字段名和值的分隔符 ':'
        const char* colon = strchr(line, ':');
        if (colon == NULL) {
            continue;
        }

        uint16_t field_name_len = colon - line;
        const char* field_value = colon + 1;

        // 处理字段名中的参数（如 DTSTART;VALUE=DATE-TIME:...）
        char field_name[64] = {0};
        uint16_t i;
        for (i = 0; i < field_name_len && i < sizeof(field_name) - 1; i++) {
            if (line[i] == ';') {
                break;  // 遇到分号，停止（忽略参数）
            }
            field_name[i] = line[i];
        }
        field_name[i] = '\0';

        // 处理 BEGIN:VTODO
        if (strcasecmp_prefix(field_name, "BEGIN", 5) && strcasecmp_prefix(field_value, "VTODO", 5)) {
            in_vtodo = true;
            memset(&current_item, 0, sizeof(todo_item_t));
            continue;
        }

        // 处理 END:VTODO
        if (strcasecmp_prefix(field_name, "END", 3) && strcasecmp_prefix(field_value, "VTODO", 5)) {
            if (in_vtodo) {
                // 保存当前项（如果有SUMMARY和DTSTART）
                if (current_item.summary[0] != '\0' && current_item.start_timestamp > 0) {
                    if (todolist.count < TODO_MAX_ITEMS) {
                        memcpy(&todolist.items[todolist.count], &current_item, sizeof(todo_item_t));
                        todolist.count++;
                    }
                }
            }
            in_vtodo = false;
            continue;
        }

        if (!in_vtodo) {
            continue;
        }

        // 解析 SUMMARY
        if (strcasecmp_prefix(field_name, "SUMMARY", 7)) {
            unescape_value(field_value, current_item.summary, TODO_MAX_SUMMARY_LEN);
            continue;
        }

        // 解析 DTSTART
        if (strcasecmp_prefix(field_name, "DTSTART", 7)) {
            // 检查是否有 VALUE=DATE 参数
            bool is_date_only = false;
            for (i = 0; i < field_name_len; i++) {
                if (line[i] == ';') {
                    const char* param = &line[i + 1];
                    if (strcasecmp_prefix(param, "VALUE=DATE", 10)) {
                        is_date_only = true;
                        break;
                    }
                }
            }

            uint16_t value_len = strlen(field_value);
            current_item.start_timestamp = parse_icalendar_datetime(field_value, value_len);
            continue;
        }

        // 解析 STATUS
        if (strcasecmp_prefix(field_name, "STATUS", 6)) {
            if (strcasecmp_prefix(field_value, "COMPLETED", 9)) {
                current_item.status = TODO_STATUS_COMPLETED;
            } else if (strcasecmp_prefix(field_value, "CANCELLED", 9)) {
                current_item.status = TODO_STATUS_CANCELLED;
            } else {
                current_item.status = TODO_STATUS_NEEDS_ACTION;
            }
            continue;
        }

        // 解析 PRIORITY
        if (strcasecmp_prefix(field_name, "PRIORITY", 8)) {
            uint8_t priority = (uint8_t)atoi(field_value);
            if (priority >= 1 && priority <= 9) {
                current_item.priority = priority;
            }
            continue;
        }
    }

    // 保存解析结果
    return todolist_write(&todolist);
}

