## 1. 数据结构和定义
- [x] 1.1 在`EPD/EPD_service.h`中定义待办事项数据结构（`todo_data_t`，使用简单字符串格式，长度60字符）
- [x] 1.2 在`EPD/EPD_service.h`中新增BLE命令`EPD_CMD_SET_TODO` (0x40)
- [x] 1.3 在`EPD/EPD_service.h`中新增BLE命令`EPD_CMD_SET_LOCATION` (0x41)
- [x] 1.4 在`EPD/EPD_service.h`中新增BLE命令`EPD_CMD_SET_WEATHER` (0x42)
- [x] 1.5 在`EPD/EPD_service.h`中新增日历显示模式枚举（完整日历、简洁日历+待办）
- [x] 1.6 在`GUI/GUI.h`中扩展`gui_data_t`结构体，添加待办事项字符串和显示模式字段
- [x] 1.7 在`GUI/GUI.h`中扩展`gui_data_t`结构体，添加地址位置字符串字段（`location_string[64]`）
- [x] 1.8 在`GUI/GUI.h`中扩展`gui_data_t`结构体，添加天气信息字符串字段（`weather_string[64]`）

## 2. BLE服务实现
- [x] 2.1 在`EPD/EPD_service.c`的`epd_service_on_write`函数中处理`EPD_CMD_SET_TODO`命令
- [x] 2.2 实现待办事项字符串接收逻辑（格式：`事项1; 事项2; ...`）
- [x] 2.3 实现待办事项数据验证（检查字符串长度，不超过60字符）
- [x] 2.4 将待办事项数据保存到Flash（使用FDS/FStorage）
- [x] 2.5 实现收到推送后直接全量显示逻辑
- [x] 2.6 实现日历显示模式切换功能（在`EPD_CMD_SET_TIME`命令中支持`mode=3`，设置`calendar_mode`为`CALENDAR_MODE_SIMPLE_TODO`）
- [x] 2.7 在`EPD/EPD_service.c`的`epd_service_on_write`函数中处理`EPD_CMD_SET_LOCATION`命令
- [x] 2.8 实现地址位置字符串接收逻辑（纯文本，直接原样显示）
- [x] 2.9 实现地址位置数据验证（检查字符串长度，不超过60字符）
- [x] 2.10 将地址位置数据保存到Flash，仅在MODE_CALENDAR_TODO模式下触发局部刷新
- [x] 2.11 在`EPD/EPD_service.c`的`epd_service_on_write`函数中处理`EPD_CMD_SET_WEATHER`命令
- [x] 2.12 实现天气信息字符串接收逻辑（格式：`温度,天气,湿度,风向,风力描述`）
- [x] 2.13 实现天气信息数据验证（检查字符串长度，不超过60字符）
- [x] 2.14 将天气信息数据保存到Flash，仅在MODE_CALENDAR_TODO模式下触发局部刷新

## 3. 数据存储
- [x] 3.1 扩展`EPD/EPD_config.c`以支持待办事项数据的读写
- [x] 3.2 在设备初始化时从Flash加载待办事项数据
- [x] 3.3 实现待办事项数据的持久化存储
- [x] 3.4 扩展`EPD/EPD_config.h`中的`epd_config_t`结构体，添加`location_string[64]`字段
- [x] 3.5 扩展`EPD/EPD_config.h`中的`epd_config_t`结构体，添加`weather_string[64]`字段
- [x] 3.6 扩展`EPD/EPD_config.c`以支持地址位置数据的读写（通过epd_config_read/write自动支持）
- [x] 3.7 扩展`EPD/EPD_config.c`以支持天气信息数据的读写（通过epd_config_read/write自动支持）
- [x] 3.8 在设备初始化时从Flash加载地址位置和天气信息数据（通过epd_config_read自动支持）
- [x] 3.9 实现地址位置和天气信息数据的持久化存储（通过epd_config_write自动支持）

## 4. GUI显示实现
- [x] 4.1 在`GUI/GUI.c`中创建`DrawSimpleCalendar`函数用于绘制简洁日历（左侧65%）
- [x] 4.2 在`GUI/GUI.c`中创建`DrawTodayInfo`函数用于绘制今日信息（右侧上方1/3：周几+农历）
- [x] 4.3 在`GUI/GUI.c`中创建`DrawTodoList`函数用于绘制待办事项列表（右侧下方2/3）
- [x] 4.4 更新`DrawCalendarTodo`函数以使用新的GUI样式和布局（参考React组件设计）
- [x] 4.5 在`GUI/GUI.c`中创建`DrawStatusBar`函数用于绘制顶部状态栏（电池电压左侧、设备ID右侧）
- [x] 4.6 更新`DrawTodayInfo`函数，实现新的今日信息区域布局
  - 左侧：年月（红色）、农历年份、日期（大号字体）、星期、农历日期（红色）
  - 右侧：天气信息（天气、温度、湿度、风向）
- [x] 4.7 在`GUI/GUI.c`中创建`DrawWeatherInfo`函数用于解析和显示天气信息
  - 解析天气字符串格式：`温度,天气,湿度,风向,风力描述`
  - 显示天气（大号字体）、温度（大号字体）、湿度、风向（小号字体，竖排）
- [x] 4.8 修改`DrawSimpleCalendar`函数，调整布局为左侧65%（原50%）
- [x] 4.9 修改`DrawTodoList`函数，调整布局为右侧35%（原50%）
- [x] 4.10 在`GUI/GUI.c`中创建`DrawFooter`函数用于绘制底部信息（地址位置左侧、空气质量右侧）
- [x] 4.11 实现地址位置信息显示（纯文本，直接原样显示）
- [x] 4.12 实现待办事项字符串分割逻辑（按分号`;`分割）
- [x] 4.13 实现待办事项显示（每个事项一行，带复选框，已完成显示删除线）
- [x] 4.14 处理文本的UTF-8编码和字体渲染（待办事项、地址位置、天气信息）
- [x] 4.15 针对400x300最小分辨率优化布局和字体大小
- [x] 4.16 实现更大尺寸屏幕的按比例调整（日历文字间距、字体大小）

## 5. 局部刷新功能
- [x] 5.1 在`EPD/EPD_driver.h`中定义局部刷新函数指针
- [x] 5.2 在`EPD/UC81xx.c`中实现`UC81xx_PartialRefresh`函数
- [x] 5.3 在`EPD/SSD16xx.c`中实现`SSD16xx_PartialRefresh`函数（如果支持）
- [ ] 5.4 在`EPD/EPD_driver.c`中实现通用的局部刷新接口
- [ ] 5.5 在`EPD/EPD_service.c`中实现待办事项更新时的局部刷新逻辑（右侧35%区域）
- [ ] 5.6 在`EPD/EPD_service.c`中实现天气信息更新时的局部刷新逻辑（今日信息区域右侧）
- [ ] 5.7 在`EPD/EPD_service.c`中实现地址位置更新时的局部刷新逻辑（底部区域）
- [ ] 5.8 实现局部刷新区域坐标计算函数（根据屏幕尺寸和布局计算）
- [ ] 5.9 确保局部刷新仅在MODE_CALENDAR_TODO模式下触发

## 6. 集成和测试
- [ ] 6.1 测试BLE命令接收和字符串解析（待办事项、地址位置、天气信息）
- [ ] 6.2 测试待办事项数据的Flash存储和加载
- [ ] 6.3 测试地址位置和天气信息数据的Flash存储和加载
- [ ] 6.4 测试完整日历模式显示
- [ ] 6.5 测试简洁日历+待办模式显示（400x300屏幕，新布局）
- [ ] 6.6 测试待办事项字符串格式解析（`事项1; 事项2; ...`，按分号分割）
- [ ] 6.7 测试天气信息字符串格式解析（`温度,天气,湿度,风向,风力描述`，按逗号分割）
- [ ] 6.8 测试地址位置信息显示（纯文本，直接原样显示）
- [ ] 6.9 测试收到推送后直接全量显示功能
- [ ] 6.10 测试局部刷新功能（待办事项区域、天气信息区域、地址位置区域）
- [ ] 6.11 测试局部刷新仅在MODE_CALENDAR_TODO模式下触发
- [ ] 6.12 测试内存使用情况（确保不超过RAM限制，字符串长度60字符）
- [ ] 6.13 测试中文文本显示（待办事项、地址位置、天气信息）
- [ ] 6.14 测试不同屏幕尺寸的布局适配（400x300最小，更大尺寸按比例调整）
- [ ] 6.15 在Windows模拟器上测试GUI代码（400x300尺寸）

## 7. Web界面实现
- [x] 7.1 在`html/index.html`中添加日历显示模式选择下拉菜单
- [x] 7.2 在`html/js/main.js`中实现`setCalendarDisplayMode`函数，调用`syncTime(1)`或`syncTime(3)`设置日历显示模式
- [x] 7.3 在`html/js/main.js`的`handleNotify`函数中自动更新下拉菜单显示当前模式（从配置通知中读取）
- [x] 7.4 修改`syncTime`函数，支持`mode=3`（日历&日程模式）并显示相应日志
- [x] 7.5 在`html/index.html`中添加地址位置输入框
- [x] 7.6 在`html/index.html`中添加天气信息输入框
- [x] 7.7 在`html/js/main.js`中实现`setLocation`函数，发送`EPD_CMD_SET_LOCATION`命令
- [x] 7.8 在`html/js/main.js`中实现`setWeather`函数，发送`EPD_CMD_SET_WEATHER`命令
- [x] 7.9 在`html/js/main.js`的`handleNotify`函数中自动更新地址位置和天气信息显示（从配置通知中读取）

## 8. 文档和清理
- [ ] 8.1 更新代码注释
- [ ] 8.2 添加RTT日志输出用于调试
- [ ] 8.3 验证代码符合项目代码风格规范
