package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

// Pusher 推送器，负责获取、解析和推送待办事项
type Pusher struct {
	url          string
	btClient     *BluetoothClient
	lastHash     string
	lastHashLock sync.Mutex
}

// NewPusher 创建新的推送器
func NewPusher(url string, btClient *BluetoothClient) *Pusher {
	return &Pusher{
		url:      url,
		btClient: btClient,
	}
}

// FetchAndPush 获取 iCalendar 数据并推送到设备
func (p *Pusher) FetchAndPush() error {
	log.Println("开始获取 iCalendar 数据...")

	// 从 URL 获取 iCalendar 数据
	resp, err := http.Get(p.url)
	if err != nil {
		return fmt.Errorf("获取 iCalendar 数据失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP 请求失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应数据失败: %w", err)
	}

	// 解析 iCalendar 数据
	icalText := string(body)
	filteredIcal, err := ParseICalendar(icalText)
	if err != nil {
		return fmt.Errorf("解析 iCalendar 数据失败: %w", err)
	}

	// 计算哈希值
	hash := p.calculateHash(filteredIcal)

	// 检查是否有变更
	p.lastHashLock.Lock()
	if p.lastHash == hash {
		p.lastHashLock.Unlock()
		log.Println("待办事项无变更，跳过推送")
		return nil
	}
	p.lastHashLock.Unlock()

	log.Printf("检测到待办事项变更，开始推送到设备...")

	// 推送到蓝牙设备
	err = p.btClient.SendTodolist(filteredIcal)
	if err != nil {
		return fmt.Errorf("推送失败: %w", err)
	}

	// 更新哈希值
	p.lastHashLock.Lock()
	p.lastHash = hash
	p.lastHashLock.Unlock()

	log.Println("推送成功！")
	return nil
}

// calculateHash 计算数据的 SHA256 哈希值
func (p *Pusher) calculateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
