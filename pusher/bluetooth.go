package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/darwin"
)

const (
	// EPD Service UUID
	epdServiceUUID = "62750001-d828-918d-fb46-b6c11c675aec"
	// EPD Characteristic UUID
	epdCharacteristicUUID = "62750002-d828-918d-fb46-b6c11c675aec"
	// SET_TODOLIST 命令
	epdCmdSetTodolist = 0x22
)

// BluetoothClient 蓝牙客户端
type BluetoothClient struct {
	deviceAddr string
}

// NewBluetoothClient 创建新的蓝牙客户端
func NewBluetoothClient(deviceAddr string) *BluetoothClient {
	return &BluetoothClient{
		deviceAddr: deviceAddr,
	}
}

// SendTodolist 连接设备并发送待办事项数据
func (c *BluetoothClient) SendTodolist(icalData string) error {
	// 初始化 BLE 设备
	device, err := darwin.NewDevice()
	if err != nil {
		return fmt.Errorf("初始化蓝牙设备失败: %w", err)
	}
	ble.SetDefaultDevice(device)

	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), 30*time.Second))
	defer ctx.Done()

	// 扫描并连接设备
	addr, err := ble.Parse(c.deviceAddr)
	if err != nil {
		return fmt.Errorf("解析设备地址失败: %w", err)
	}

	cln, err := ble.Connect(ctx, func(a ble.Advertisement) bool {
		return a.Addr().String() == addr.String()
	})
	if err != nil {
		return fmt.Errorf("连接设备失败: %w", err)
	}
	defer func() {
		_ = cln.CancelConnection()
	}()

	// 发现服务
	serviceUUID, err := ble.Parse(epdServiceUUID)
	if err != nil {
		return fmt.Errorf("解析服务 UUID 失败: %w", err)
	}

	// 发现服务
	services, err := cln.DiscoverServices([]ble.UUID{serviceUUID})
	if err != nil {
		return fmt.Errorf("发现服务失败: %w", err)
	}

	if len(services) == 0 {
		return fmt.Errorf("未找到 EPD 服务")
	}
	service := services[0]

	// 发现特征
	charUUID, err := ble.Parse(epdCharacteristicUUID)
	if err != nil {
		return fmt.Errorf("解析特征 UUID 失败: %w", err)
	}

	chars, err := cln.DiscoverCharacteristics([]ble.UUID{charUUID}, service)
	if err != nil {
		return fmt.Errorf("发现特征失败: %w", err)
	}

	if len(chars) == 0 {
		return fmt.Errorf("未找到 EPD 特征")
	}
	char := chars[0]

	// 准备数据：命令字节 + iCalendar 文本（UTF-8 编码）
	data := append([]byte{epdCmdSetTodolist}, []byte(icalData)...)

	// 发送数据
	// BLE 有 MTU 限制，需要分片发送
	mtu := 20            // 默认 MTU，实际应该从设备获取
	chunkSize := mtu - 3 // 减去命令字节和开销

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]

		// 第一个块包含命令字节，后续块只包含数据
		// noRsp=false 表示需要响应，noRsp=true 表示不需要响应（更快但不可靠）
		err = cln.WriteCharacteristic(char, chunk, i > 0)

		if err != nil {
			return fmt.Errorf("写入数据失败: %w", err)
		}

		// 添加小延迟避免发送过快
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}
