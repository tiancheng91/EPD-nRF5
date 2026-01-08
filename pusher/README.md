# 蓝牙待办事项推送 CLI

这是一个用于定时从 iCalendar URL 获取待办事项并通过蓝牙推送到墨水屏设备的 CLI 程序。

## 功能特性

- 定时从 iCalendar URL 获取待办事项
- 自动过滤当日及之后的待办事项
- 变更检测：仅在待办事项有变更时推送
- 蓝牙连接管理：推送完成后自动断开以节省电量
- 支持后台运行模式
- 可配置的更新间隔

## 编译

```bash
cd pusher
go build -o pusher
```

## 使用方法

### 基本用法

```bash
# 前台运行，每3小时更新一次（默认）
./pusher --device AA:BB:CC:DD:EE:FF --url "https://ext.todoist.com/export/ical/..."

# 自定义更新间隔（30分钟）
./pusher --device AA:BB:CC:DD:EE:FF --url "..." --interval 30m

# 后台运行
./pusher --device AA:BB:CC:DD:EE:FF --url "..." -D

# 执行一次后退出（用于测试）
./pusher --device AA:BB:CC:DD:EE:FF --url "..." --once
```

### 参数说明

- `--device` / `-d`: 蓝牙设备地址（必需，格式：AA:BB:CC:DD:EE:FF）
- `--url`: iCalendar URL（必需）
- `--interval`: 更新间隔（默认: 3h，格式如: 3h, 30m, 1h30m）
- `-D` / `--daemon`: 后台运行模式
- `--once`: 执行一次后退出（用于测试）

### 示例

```bash
# 每3小时检查一次 Todoist 日历并推送到设备
./pusher \
  --device "AA:BB:CC:DD:EE:FF" \
  --url "https://ext.todoist.com/export/ical/project?user_id=469931&project_id=6CrfqqFQj6vXpHR3&ical_token=f067b2e7&r_factor=9954" \
  --interval 3h

# 后台运行，每30分钟检查一次
./pusher \
  --device "AA:BB:CC:DD:EE:FF" \
  --url "https://ext.todoist.com/export/ical/..." \
  --interval 30m \
  -D
```

## 工作原理

1. 程序启动后，根据配置的间隔定时执行任务
2. 从指定的 iCalendar URL 获取日历数据
3. 解析 iCalendar 格式，提取 VTODO 组件
4. 过滤出当日及之后的待办事项
5. 计算数据哈希值，与上一次的数据对比
6. 如果有变更，连接蓝牙设备并推送数据
7. 推送完成后立即断开连接
8. 更新内存中的哈希值

## 注意事项

1. **蓝牙权限**：macOS 需要授予蓝牙权限
2. **设备地址**：确保设备地址格式正确（AA:BB:CC:DD:EE:FF）
3. **网络连接**：需要能够访问 iCalendar URL
4. **后台运行**：使用 `-D` 或 `--daemon` 参数时，程序会 fork 到后台运行

## 错误处理

- HTTP 请求失败：记录日志，等待下次执行
- 蓝牙连接失败：记录日志，等待下次执行
- 数据解析失败：记录错误，跳过本次执行

## 日志输出

程序会输出以下类型的日志：
- 调度器启动信息
- 获取和解析 iCalendar 数据的状态
- 变更检测结果
- 蓝牙连接和推送状态
- 错误信息

## 依赖

- `tinygo.org/x/bluetooth` - 蓝牙通信库（跨平台支持）
- Go 1.16+

## 许可证

与主项目相同

