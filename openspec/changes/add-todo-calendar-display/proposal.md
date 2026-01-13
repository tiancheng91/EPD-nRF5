# Change: 日历模式下新增待办事项显示

## Why
用户希望在日历模式下能够显示待办事项，并通过蓝牙推送待办事项列表到墨水屏设备。当收到新的待办事项时，设备应使用局部刷新功能更新显示，避免全屏刷新带来的延迟和闪烁。同时支持两种日历显示模式：完整日历模式和简洁日历+待办模式，以适应不同使用场景。

## What Changes
- **ADDED**: 日历显示模式选择功能（完整日历模式、简洁日历+待办模式）
- **ADDED**: 简洁日历+待办模式的左右布局设计（左侧50%简洁日历，右侧50%当日明细和待办）
- **ADDED**: 日历模式下显示待办事项列表功能
- **ADDED**: 通过BLE接收和存储待办事项数据（简化字符串格式：`事项1; 事项2; ...`，长度不超过60字符）
- **ADDED**: 待办事项不需要和日期关联，收到推送后直接全量显示
- **ADDED**: 支持待办事项的局部刷新显示
- **ADDED**: 待办事项数据结构定义和存储机制（简化设计）
- **ADDED**: BLE命令用于推送待办事项数据
- **ADDED**: HTML网页界面支持配置日历显示模式（通过下拉菜单，复用`SET_TIME`命令，`mode=3`表示日历&日程模式）

## Impact
- 影响的规范：
  - `calendar-display`: 日历显示功能需要扩展以支持待办事项显示
  - `ble-service`: BLE服务需要新增命令处理待办事项数据
  - `epd-refresh`: EPD驱动需要支持局部刷新功能
- 影响的代码：
  - `GUI/GUI.c`: 修改`DrawCalendar`函数以支持两种显示模式，新增简洁日历绘制函数
  - `GUI/GUI.h`: 扩展`gui_data_t`结构体包含待办事项数据和显示模式
  - `EPD/EPD_service.c`: 新增BLE命令处理待办事项推送
  - `EPD/EPD_service.h`: 新增命令ID和数据结构定义
  - `EPD/EPD_driver.c`或`EPD/UC81xx.c`/`EPD/SSD16xx.c`: 实现局部刷新功能
  - `EPD/EPD_config.c`: 扩展配置存储以保存待办事项和显示模式
  - `EPD/EPD_config.h`: 扩展配置结构体包含`calendar_mode`和`todo_string`字段
  - `EPD/EPD_service.c`: 修改`EPD_CMD_SET_TIME`命令处理，支持`mode=3`设置日历&日程模式
  - `html/index.html`: 添加日历显示模式选择下拉菜单
  - `html/js/main.js`: 修改`syncTime`和`setCalendarDisplayMode`函数，使用`mode=3`设置日历&日程模式
