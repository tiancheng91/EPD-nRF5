## ADDED Requirements

### Requirement: 通过BLE接收待办事项数据
系统 SHALL 通过BLE服务接收待办事项数据。待办事项数据应通过`EPD_CMD_SET_TODO`命令传输，使用简单字符串格式，不需要和日期关联。

#### Scenario: 接收待办事项命令
- **WHEN** 客户端通过BLE发送`EPD_CMD_SET_TODO`命令
- **AND** 命令格式正确（包含待办事项字符串）
- **THEN** 系统应解析命令数据
- **AND** 验证待办事项字符串的有效性（字符串长度等）
- **AND** 将待办事项数据保存到Flash存储
- **AND** 墨水屏收到推送后直接全量显示

#### Scenario: 待办事项命令格式
- **WHEN** 接收`EPD_CMD_SET_TODO`命令
- **THEN** 命令格式应为：`[0x40] [string_len] [todo_string...]`
- **AND** `string_len`为待办事项字符串长度（1字节，0-60）
- **AND** `todo_string`为待办事项字符串（UTF-8编码）

#### Scenario: 待办事项字符串格式
- **WHEN** 接收待办事项字符串
- **THEN** 字符串格式应为：`事项1; 事项2; ...`
- **AND** 多个待办事项用分号`;`分隔
- **AND** 字符串长度不超过60个字符（避免占用内存）
- **AND** 系统应以推送内容为准，不需要二次处理
- **AND** 不需要和日期关联

#### Scenario: 数据验证失败
- **WHEN** 接收的待办事项数据无效（字符串长度超限等）
- **THEN** 系统应拒绝保存数据
- **AND** 记录错误日志（RTT输出）

#### Scenario: 待办事项数据持久化
- **WHEN** 成功接收并验证待办事项数据
- **THEN** 系统应将数据保存到Flash存储
- **AND** 设备重启后应能从Flash加载待办事项数据

#### Scenario: 字符串长度限制
- **WHEN** 接收的待办事项字符串长度超过60个字符
- **THEN** 系统应截断字符串至60个字符
- **AND** 记录警告日志
