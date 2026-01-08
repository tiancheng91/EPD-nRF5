#include "todolist.h"

#include <string.h>

#include "app_scheduler.h"
#include "fds.h"
#include "nordic_common.h"
#include "nrf_log.h"

#define TODOLIST_FILE_ID 0x0001
#define TODOLIST_REC_KEY 0x0001

static void run_fds_gc(void* p_event_data, uint16_t event_size) {
    NRF_LOG_DEBUG("run garbage collection (fds_gc)\n");
    fds_gc();
}

void todolist_init(void) {
    // FDS已经在epd_config_init中初始化，这里不需要重复初始化
    // 如果FDS未初始化，后续的读写操作会失败，这是可以接受的
}

uint32_t todolist_read(todolist_t* todolist) {
    if (todolist == NULL) {
        return NRF_ERROR_NULL;
    }

    fds_flash_record_t flash_record;
    fds_record_desc_t record_desc;
    fds_find_token_t ftok;

    memset(todolist, 0, sizeof(todolist_t));
    memset(&ftok, 0x00, sizeof(fds_find_token_t));

    ret_code_t ret = fds_record_find(TODOLIST_FILE_ID, TODOLIST_REC_KEY, &record_desc, &ftok);
    if (ret != NRF_SUCCESS) {
        NRF_LOG_DEBUG("todolist_read: record not found\n");
        return ret;
    }

    ret = fds_record_open(&record_desc, &flash_record);
    if (ret != NRF_SUCCESS) {
        NRF_LOG_ERROR("todolist_read: record open failed!\n");
        return ret;
    }

#ifdef S112
    uint32_t record_len = flash_record.p_header->length_words * sizeof(uint32_t);
#else
    uint32_t record_len = flash_record.p_header->tl.length_words * sizeof(uint32_t);
#endif
    memcpy(todolist, flash_record.p_data, MIN(sizeof(todolist_t), record_len));
    fds_record_close(&record_desc);

    return NRF_SUCCESS;
}

uint32_t todolist_write(todolist_t* todolist) {
    if (todolist == NULL) {
        return NRF_ERROR_NULL;
    }

    ret_code_t ret;
    fds_record_t record;
    fds_record_desc_t record_desc;
    fds_find_token_t ftok;

    record.file_id = TODOLIST_FILE_ID;
    record.key = TODOLIST_REC_KEY;
#ifdef S112
    record.data.p_data = (void*)todolist;
    record.data.length_words = BYTES_TO_WORDS(sizeof(todolist_t));
#else
    fds_record_chunk_t record_chunk;
    record_chunk.p_data = todolist;
    record_chunk.length_words = BYTES_TO_WORDS(sizeof(todolist_t));
    record.data.p_chunks = &record_chunk;
    record.data.num_chunks = 1;
#endif

    memset(&ftok, 0x00, sizeof(fds_find_token_t));
    ret = fds_record_find(TODOLIST_FILE_ID, TODOLIST_REC_KEY, &record_desc, &ftok);
    if (ret == NRF_SUCCESS) {
        ret = fds_record_update(&record_desc, &record);
    } else {
        ret = fds_record_write(&record_desc, &record);
    }

    if (ret != NRF_SUCCESS) {
        NRF_LOG_ERROR("todolist_write: record write/update failed, code=%d\n", ret);
        if (ret == FDS_ERR_NO_SPACE_IN_FLASH) {
            app_sched_event_put(NULL, 0, run_fds_gc);
        }
    }

    return ret;
}

uint32_t todolist_clear(void) {
    ret_code_t ret;
    fds_record_desc_t record_desc;
    fds_find_token_t ftok;

    memset(&ftok, 0x00, sizeof(fds_find_token_t));
    ret = fds_record_find(TODOLIST_FILE_ID, TODOLIST_REC_KEY, &record_desc, &ftok);
    if (ret != NRF_SUCCESS) {
        NRF_LOG_DEBUG("todolist_clear: record not found\n");
        return NRF_SUCCESS;  // 记录不存在，视为已清空
    }

    ret = fds_record_delete(&record_desc);
    if (ret != NRF_SUCCESS) {
        NRF_LOG_ERROR("todolist_clear: fds_record_delete failed, code=%d\n", ret);
    }

    return ret;
}

uint8_t todolist_get_count(void) {
    todolist_t todolist;
    if (todolist_read(&todolist) == NRF_SUCCESS) {
        return todolist.count;
    }
    return 0;
}

