## ADDED Requirements

### Requirement: 待办事项显示区域的局部刷新
系统 SHALL 在更新待办事项显示时使用局部刷新功能，仅刷新待办事项显示区域，而不刷新整个屏幕。

#### Scenario: 局部刷新待办事项区域
- **WHEN** 收到新的待办事项数据
- **AND** 待办事项显示需要更新
- **THEN** 系统应仅刷新待办事项显示区域
- **AND** 日历部分不应被刷新
- **AND** 刷新时间应明显短于全屏刷新（目标：< 3秒）

#### Scenario: 局部刷新区域定义
- **WHEN** 执行局部刷新
- **THEN** 刷新区域应为待办事项显示区域
- **AND** 刷新区域坐标应根据屏幕尺寸和布局计算
- **AND** 刷新区域应包含所有待办事项文本

#### Scenario: 驱动芯片支持检查
- **WHEN** EPD驱动芯片支持局部刷新（如UC81xx系列）
- **THEN** 系统应使用局部刷新功能
- **WHEN** EPD驱动芯片不支持局部刷新
- **THEN** 系统应回退到全屏刷新
- **AND** 记录警告日志

#### Scenario: 局部刷新实现（UC81xx）
- **WHEN** 使用UC81xx系列驱动芯片
- **THEN** 系统应使用`UC81xx_PTIN`命令进入局部刷新模式
- **AND** 使用`_setPartialRamArea`设置刷新区域
- **AND** 写入待办事项显示数据
- **AND** 使用`UC81xx_PTOUT`命令退出局部刷新模式
- **AND** 触发刷新命令（`UC81xx_DRF`）

#### Scenario: 局部刷新实现（SSD16xx）
- **WHEN** 使用SSD16xx系列驱动芯片
- **AND** 驱动支持局部刷新
- **THEN** 系统应使用`_setPartialRamArea`设置刷新区域
- **AND** 写入待办事项显示数据
- **AND** 触发局部刷新命令

#### Scenario: 刷新完成等待
- **WHEN** 触发局部刷新
- **THEN** 系统应等待刷新完成（通过`EPD_WaitBusy`）
- **AND** 刷新完成后应恢复正常状态
