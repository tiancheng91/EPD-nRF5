#ifndef __TODOLIST_H
#define __TODOLIST_H

#include <stdint.h>

// NRF51 内存较小，使用更小的配置
#ifdef NRF51
#define TODO_MAX_SUMMARY_LEN 32    // 摘要最大长度（NRF51 使用较小值）
#define TODO_MAX_ITEMS 10           // 最大待办事项数量（NRF51 使用较小值）
#else
#define TODO_MAX_SUMMARY_LEN 64    // 摘要最大长度
#define TODO_MAX_ITEMS 20           // 最大待办事项数量
#endif

typedef enum {
    TODO_STATUS_NEEDS_ACTION = 0,  // 待处理
    TODO_STATUS_COMPLETED = 1,      // 已完成
    TODO_STATUS_CANCELLED = 2       // 已取消
} todo_status_t;

typedef struct {
    char summary[TODO_MAX_SUMMARY_LEN];  // 待办事项摘要
    uint32_t start_timestamp;              // 开始时间 (Unix timestamp, DTSTART, 必需字段)
    todo_status_t status;                 // 状态 (NEEDS-ACTION, COMPLETED, CANCELLED)
    uint8_t priority;                     // 优先级 (0-9, 0为未设置)
} todo_item_t;

typedef struct {
    uint8_t count;                         // 待办事项数量
    todo_item_t items[TODO_MAX_ITEMS];     // 待办事项数组
} todolist_t;

// 初始化待办事项存储
void todolist_init(void);

// 从 Flash 读取待办事项列表（读取所有存储的事项）
uint32_t todolist_read(todolist_t* todolist);

// 写入待办事项列表到 Flash
uint32_t todolist_write(todolist_t* todolist);

// 清空待办事项列表
uint32_t todolist_clear(void);

// 获取待办事项数量
uint8_t todolist_get_count(void);

#endif  // __TODOLIST_H

