#ifndef __ICALENDAR_H
#define __ICALENDAR_H

#include <stdint.h>
#include "sdk_errors.h"

// 从 iCalendar 文本解析并存储待办事项
// data: iCalendar 格式的文本数据
// length: 数据长度
// 返回: NRF_SUCCESS 成功，其他值表示错误
ret_code_t todolist_parse_icalendar(uint8_t* data, uint16_t length);

#endif  // __ICALENDAR_H

