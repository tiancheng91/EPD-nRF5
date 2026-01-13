# Change: 日历模式下新增待办事项显示

## Why
用户希望在日历模式下能够显示待办事项，并通过蓝牙推送待办事项列表到墨水屏设备。当收到新的待办事项时，设备应使用局部刷新功能更新显示，避免全屏刷新带来的延迟和闪烁。同时支持两种日历显示模式：完整日历模式和简洁日历+待办模式，以适应不同使用场景。

## What Changes
- **ADDED**: 日历显示模式选择功能（完整日历模式、简洁日历+待办模式）
- **ADDED**: 简洁日历+待办模式的左右布局设计（左侧50%简洁日历，右侧50%待办事项）
- **ADDED**: 日历模式下显示待办事项列表功能
- **ADDED**: 通过BLE接收和存储待办事项数据（简化字符串格式：`事项1; 事项2; ...`，长度不超过60字符）
- **ADDED**: 待办事项不需要和日期关联，收到推送后直接全量显示
- **ADDED**: 支持待办事项的局部刷新显示
- **ADDED**: 待办事项数据结构定义和存储机制（简化设计）
- **ADDED**: BLE命令用于推送待办事项数据（`EPD_CMD_SET_TODO`）
- **ADDED**: HTML网页界面支持配置日历显示模式（通过下拉菜单，复用`SET_TIME`命令，`mode=3`表示日历&日程模式）
- **ADDED**: 更新MODE_CALENDAR_TODO模式的GUI样式和布局（参考React组件设计）
  - 顶部状态栏：电池电压（左侧）和设备ID（右侧）
  - 今日信息区域：左侧显示年月、日期、星期、农历；右侧显示天气信息（天气、温度、湿度、风向）
  - 中间区域：左侧50%显示简洁日历，右侧50%显示待办事项列表
  - 底部：地址位置信息和空气质量
- **ADDED**: BLE命令`EPD_CMD_SET_LOCATION` (0x41)用于更新左下角地址位置信息（纯文本，直接原样显示，最大长度60字符）
- **ADDED**: BLE命令`EPD_CMD_SET_WEATHER` (0x42)用于更新天气信息（格式：`温度,天气,湿度,风向,风力描述`，示例：`7.6,晴,56,东南,东南风 二级`，最大长度60字符）
- **ADDED**: 地址和天气信息仅在MODE_CALENDAR_TODO模式下生效，更新后使用局部刷新相应区域
- **ADDED**: HTML网页界面支持设置地址位置和天气信息
- **ADDED**: 支持400x300最小分辨率，更大尺寸按比例调整日历文字间距

## Impact
- 影响的规范：
  - `calendar-display`: 日历显示功能需要扩展以支持待办事项显示、地址位置和天气信息
  - `ble-service`: BLE服务需要新增命令处理待办事项数据、地址位置和天气信息
  - `epd-refresh`: EPD驱动需要支持局部刷新功能（待办事项、地址、天气区域）
- 影响的代码：
  - `GUI/GUI.c`: 
    - 修改`DrawCalendarTodo`函数以使用新的GUI样式和布局
    - 新增状态栏绘制函数（电池电压、设备ID）
    - 新增今日信息区域绘制函数（左侧：年月、日期、星期、农历；右侧：天气信息）
    - 修改日历和待办事项布局（左侧50%日历，右侧50%待办）
    - 新增底部地址位置信息绘制函数
    - 支持400x300最小分辨率，更大尺寸按比例调整
  - `GUI/GUI.h`: 扩展`gui_data_t`结构体包含待办事项数据、地址位置字符串、天气信息字符串
  - `EPD/EPD_service.c`: 
    - 新增BLE命令处理待办事项推送（`EPD_CMD_SET_TODO`）
    - 新增BLE命令处理地址位置推送（`EPD_CMD_SET_LOCATION`）
    - 新增BLE命令处理天气信息推送（`EPD_CMD_SET_WEATHER`）
    - 命令处理时检查当前显示模式，仅在MODE_CALENDAR_TODO模式下触发局部刷新
  - `EPD/EPD_service.h`: 新增命令ID定义（`EPD_CMD_SET_LOCATION = 0x41`, `EPD_CMD_SET_WEATHER = 0x42`）
  - `EPD/EPD_driver.c`或`EPD/UC81xx.c`/`EPD/SSD16xx.c`: 实现局部刷新功能（支持指定区域刷新）
  - `EPD/EPD_config.c`: 扩展配置存储以保存待办事项、地址位置、天气信息和显示模式
  - `EPD/EPD_config.h`: 扩展配置结构体包含`todo_string`、`location_string`、`weather_string`字段
  - `EPD/EPD_service.c`: 修改`EPD_CMD_SET_TIME`命令处理，支持`mode=3`设置日历&日程模式
  - `html/index.html`: 添加日历显示模式选择下拉菜单、地址位置输入框、天气信息输入框
  - `html/js/main.js`: 
    - 修改`syncTime`和`setCalendarDisplayMode`函数，使用`mode=3`设置日历&日程模式
    - 新增函数发送地址位置和天气信息命令
