## 1. 数据结构和定义
- [x] 1.1 在`EPD/EPD_service.h`中定义待办事项数据结构（`todo_data_t`，使用简单字符串格式，长度60字符）
- [x] 1.2 在`EPD/EPD_service.h`中新增BLE命令`EPD_CMD_SET_TODO` (0x40)
- [x] 1.3 在`EPD/EPD_service.h`中新增日历显示模式枚举（完整日历、简洁日历+待办）
- [x] 1.4 在`GUI/GUI.h`中扩展`gui_data_t`结构体，添加待办事项字符串和显示模式字段

## 2. BLE服务实现
- [x] 2.1 在`EPD/EPD_service.c`的`epd_service_on_write`函数中处理`EPD_CMD_SET_TODO`命令
- [x] 2.2 实现待办事项字符串接收逻辑（格式：`事项1; 事项2; ...`）
- [x] 2.3 实现待办事项数据验证（检查字符串长度，不超过60字符）
- [x] 2.4 将待办事项数据保存到Flash（使用FDS/FStorage）
- [x] 2.5 实现收到推送后直接全量显示逻辑
- [x] 2.6 实现日历显示模式切换功能（在`EPD_CMD_SET_TIME`命令中支持`mode=3`，设置`calendar_mode`为`CALENDAR_MODE_SIMPLE_TODO`）

## 3. 数据存储
- [x] 3.1 扩展`EPD/EPD_config.c`以支持待办事项数据的读写
- [x] 3.2 在设备初始化时从Flash加载待办事项数据
- [x] 3.3 实现待办事项数据的持久化存储

## 4. GUI显示实现
- [x] 4.1 在`GUI/GUI.c`中创建`DrawSimpleCalendar`函数用于绘制简洁日历（左侧50%）
- [x] 4.2 在`GUI/GUI.c`中创建`DrawTodayInfo`函数用于绘制今日信息（右侧上方1/3：周几+农历）
- [x] 4.3 在`GUI/GUI.c`中创建`DrawTodoList`函数用于绘制待办事项列表（右侧下方2/3）
- [x] 4.4 修改`DrawCalendar`函数，根据显示模式选择完整日历或简洁日历+待办布局
- [x] 4.5 实现待办事项字符串分割逻辑（按分号`;`分割）
- [x] 4.6 实现待办事项显示（每个事项一行，不需要显示日期和时间）
- [x] 4.7 处理待办事项文本的UTF-8编码和字体渲染
- [x] 4.8 针对400x300屏幕优化布局和字体大小

## 5. 局部刷新功能
- [x] 5.1 在`EPD/EPD_driver.h`中定义局部刷新函数指针
- [x] 5.2 在`EPD/UC81xx.c`中实现`UC81xx_PartialRefresh`函数
- [x] 5.3 在`EPD/SSD16xx.c`中实现`SSD16xx_PartialRefresh`函数（如果支持）
- [ ] 5.4 在`EPD/EPD_driver.c`中实现通用的局部刷新接口
- [ ] 5.5 在`EPD/EPD_service.c`中实现待办事项更新时的局部刷新逻辑

## 6. 集成和测试
- [ ] 6.1 测试BLE命令接收和字符串解析
- [ ] 6.2 测试待办事项数据的Flash存储和加载
- [ ] 6.3 测试完整日历模式显示
- [ ] 6.4 测试简洁日历+待办模式显示（400x300屏幕）
- [ ] 6.5 测试待办事项字符串格式解析（`事项1; 事项2; ...`，按分号分割）
- [ ] 6.6 测试收到推送后直接全量显示功能
- [ ] 6.7 测试局部刷新功能（UC81xx和SSD16xx驱动）
- [ ] 6.8 测试内存使用情况（确保不超过RAM限制，字符串长度60字符）
- [ ] 6.9 测试中文文本显示
- [ ] 6.10 在Windows模拟器上测试GUI代码（400x300尺寸）

## 7. Web界面实现
- [x] 7.1 在`html/index.html`中添加日历显示模式选择下拉菜单
- [x] 7.2 在`html/js/main.js`中实现`setCalendarDisplayMode`函数，调用`syncTime(1)`或`syncTime(3)`设置日历显示模式
- [x] 7.3 在`html/js/main.js`的`handleNotify`函数中自动更新下拉菜单显示当前模式（从配置通知中读取）
- [x] 7.4 修改`syncTime`函数，支持`mode=3`（日历&日程模式）并显示相应日志

## 8. 文档和清理
- [ ] 8.1 更新代码注释
- [ ] 8.2 添加RTT日志输出用于调试
- [ ] 8.3 验证代码符合项目代码风格规范
