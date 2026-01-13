#ifndef __GUI_H
#define __GUI_H

#include "Adafruit_GFX.h"

typedef enum {
    MODE_PICTURE = 0,
    MODE_CALENDAR = 1,
    MODE_CLOCK = 2,
    MODE_CALENDAR_TODO = 3, /**< Calendar + todo mode (simple calendar with todo list) */
} display_mode_t;

/**@brief Calendar display mode.
 */
typedef enum {
    CALENDAR_MODE_FULL = 0,      /**< Full calendar mode */
    CALENDAR_MODE_SIMPLE_TODO = 1, /**< Simple calendar + todo mode */
} calendar_display_mode_t;

typedef struct {
    display_mode_t mode;
    uint16_t color;
    uint16_t width;
    uint16_t height;
    uint32_t timestamp;
    uint8_t week_start;  // 0: Sunday, 1: Monday
    int8_t temperature;
    uint16_t voltage;
    char ssid[20];
    calendar_display_mode_t calendar_mode; /**< Calendar display mode */
    char todo_string[20];                  /**< Todo items string (format: "item1\nitem2\n...", max 19 chars including \n) */
    char location_string[20];              /**< Location string (plain text, max 19 chars) */
    char weather_string[20];               /**< Weather info string (format: "temp,weather,humidity,wind_dir,wind_desc", max 19 chars) */
} gui_data_t;

void DrawGUI(gui_data_t* data, buffer_callback callback, void* callback_data);

#endif
