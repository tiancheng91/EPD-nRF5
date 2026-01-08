package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	deviceAddr  = flag.String("device", "", "蓝牙设备地址（必需）")
	url         = flag.String("url", "", "iCalendar URL（必需）")
	interval    = flag.String("interval", "3h", "更新间隔（默认: 3h，格式如: 3h, 30m）")
	daemon      = flag.Bool("daemon", false, "后台运行模式（也可使用 -D）")
	daemonShort = flag.Bool("D", false, "后台运行模式（短选项）")
	once        = flag.Bool("once", false, "执行一次后退出（用于测试）")
)

func main() {
	flag.Parse()

	// 验证必需参数
	if *deviceAddr == "" {
		log.Fatal("错误: 必须指定蓝牙设备地址 (--device)")
	}
	if *url == "" {
		log.Fatal("错误: 必须指定 iCalendar URL (--url)")
	}

	// 解析更新间隔
	updateInterval, err := time.ParseDuration(*interval)
	if err != nil {
		log.Fatalf("错误: 无效的更新间隔格式: %v", err)
	}

	// 后台运行模式
	if *daemon || *daemonShort {
		if err := daemonize(); err != nil {
			log.Fatalf("错误: 后台运行失败: %v", err)
		}
	}

	// 创建蓝牙客户端
	btClient := NewBluetoothClient(*deviceAddr)

	// 创建推送器
	pusher := NewPusher(*url, btClient)

	// 创建调度器
	scheduler := NewScheduler(pusher, updateInterval, *once)

	// 设置信号处理，优雅退出
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// 在 goroutine 中启动调度器
	go func() {
		scheduler.Start()
	}()

	// 等待信号
	<-sigCh
	log.Println("收到退出信号，正在关闭...")
	scheduler.Stop()
	log.Println("程序已退出")
}

// daemonize 将进程转为后台运行
func daemonize() error {
	// 简单的后台运行实现
	// 在实际应用中，可以使用更完善的守护进程库
	// 这里使用基本的 fork 方式

	// 检查是否已经是后台进程（通过环境变量标记）
	if os.Getenv("PUSHER_DAEMON") == "1" {
		// 已经是后台进程，继续执行
		return nil
	}

	// 准备新进程的参数（移除 daemon 相关参数）
	newArgs := []string{}
	for _, arg := range os.Args {
		if arg == "-D" || arg == "--daemon" {
			continue
		}
		newArgs = append(newArgs, arg)
	}

	// 设置环境变量标记
	env := os.Environ()
	env = append(env, "PUSHER_DAEMON=1")

	// 创建新的进程会话
	_, err := os.StartProcess(os.Args[0], newArgs, &os.ProcAttr{
		Files: []*os.File{nil, nil, nil}, // 关闭标准输入输出
		Env:   env,
	})
	if err != nil {
		return fmt.Errorf("创建后台进程失败: %w", err)
	}
	// 父进程退出
	os.Exit(0)

	return nil
}
