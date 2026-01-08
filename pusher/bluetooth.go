package main

import (
	"context"
	"fmt"
	"time"

	"tinygo.org/x/bluetooth"
)

var (
	// EPD Service UUID
	epdServiceUUID bluetooth.UUID
	// EPD Characteristic UUID
	epdCharacteristicUUID bluetooth.UUID
)

const (
	// SET_TODOLIST 命令
	epdCmdSetTodolist = 0x22
)

func init() {
	// 初始化 UUID
	var err error
	epdServiceUUID, err = bluetooth.ParseUUID("62750001-d828-918d-fb46-b6c11c675aec")
	if err != nil {
		panic(fmt.Sprintf("解析服务 UUID 失败: %v", err))
	}
	epdCharacteristicUUID, err = bluetooth.ParseUUID("62750002-d828-918d-fb46-b6c11c675aec")
	if err != nil {
		panic(fmt.Sprintf("解析特征 UUID 失败: %v", err))
	}
}

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
	// 启用 BLE 适配器
	adapter := bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		return fmt.Errorf("启用蓝牙适配器失败: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 解析设备地址
	mac, err := bluetooth.ParseMAC(c.deviceAddr)
	if err != nil {
		return fmt.Errorf("解析设备地址失败: %w", err)
	}

	// 将 MAC 地址转换为 Address
	var addr bluetooth.Address
	addr.Set(mac.String())

	// 扫描并连接设备
	var foundAddr bluetooth.Address
	ch := make(chan bluetooth.Address, 1)
	errCh := make(chan error, 1)

	// 启动扫描
	go func() {
		err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
			if result.Address.String() == addr.String() {
				foundAddr = result.Address
				ch <- foundAddr
				adapter.StopScan()
			}
		})
		if err != nil {
			errCh <- fmt.Errorf("扫描设备失败: %w", err)
		}
	}()

	// 等待找到设备或超时
	select {
	case foundAddr = <-ch:
		// 设备已找到，扫描已停止
	case err := <-errCh:
		adapter.StopScan()
		return err
	case <-ctx.Done():
		adapter.StopScan()
		return fmt.Errorf("扫描设备超时")
	}

	// 连接设备
	device, err := adapter.Connect(foundAddr, bluetooth.ConnectionParams{})
	if err != nil {
		return fmt.Errorf("连接设备失败: %w", err)
	}
	defer func() {
		_ = device.Disconnect()
	}()

	// 发现服务
	services, err := device.DiscoverServices([]bluetooth.UUID{epdServiceUUID})
	if err != nil {
		return fmt.Errorf("发现服务失败: %w", err)
	}

	if len(services) == 0 {
		return fmt.Errorf("未找到 EPD 服务")
	}
	service := services[0]

	// 发现特征
	chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{epdCharacteristicUUID})
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

		// 使用 WriteWithoutResponse（兼容所有平台）
		// 注意：linux上不支持 Write（带响应），统一使用 WriteWithoutResponse
		_, err = char.WriteWithoutResponse(chunk)

		if err != nil {
			return fmt.Errorf("写入数据失败: %w", err)
		}

		// 添加小延迟避免发送过快
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}
