package main

import (
	"testing"
	"time"

	"github.com/go-ping/ping"
	"github.com/stretchr/testify/assert"
)

func TestPingHost(t *testing.T) {
	// 要 ping 的 IP 地址或域名
	host := "www.baidu.com"

	// 创建一个新的 pinger
	pinger, err := ping.NewPinger(host)
	if err != nil {
		t.Fatalf("创建 pinger 失败: %v", err)
	}

	// 设置发送的 ICMP 包数
	pinger.Count = 4

	// 设置每个包发送的间隔时间
	pinger.Interval = time.Second

	// 设置超时时间
	pinger.Timeout = time.Second * 10

	// 运行 ping
	err = pinger.Run()
	if err != nil {
		t.Fatalf("Ping 失败: %v", err)
	}

	// 获取结果
	stats := pinger.Statistics()

	// 验证 ping 成功
	assert.Greater(t, stats.PacketsRecv, 0, "Ping 没有收到任何回应")

	// 输出统计信息
	t.Logf("\n--- %s ping statistics ---\n", stats.Addr)
	t.Logf("%d packets transmitted, %d packets received, %.2f%% packet loss\n",
		stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
	t.Logf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
		stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
}
