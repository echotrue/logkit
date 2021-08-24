package main

import (
	"fmt"
	log2 "log"
	"logkit"
	"runtime"
	"time"
)

func main() {
	// t1 := logkit.NewConsoleTarget()
	t2 := logkit.NewFileTarget()
	t2.FilePath = "./temp"
	//
	// t3 := logkit.NewNetworkTarget()
	// t3.Network = "tcp"
	// t3.Address = "192.168.1.58:9200"

	log := logkit.NewLogKitByOptions(logkit.WithBuffer(1024), logkit.WithTargets(t2), logkit.WithCallStackDepth(0))

	err := log.Open()
	if err != nil {
		log2.Fatal(err)
	}
	defer log.Close()

	// log.Log(logkit.SeverityDebug, "%d--%s", time.Now().Unix(), "当前时间戳")
	// log.Log(logkit.SeverityInfo, "%d--%s", time.Now().Unix(), "当前时间戳")
	// log.Log(logkit.SeverityNotice, "%d--%s", time.Now().Unix(), "当前时间戳")
	// log.Log(logkit.SeverityWarning, "%d--%s", time.Now().Unix(), "当前时间戳")
	// log.Log(logkit.SeverityError, "%d--%s", time.Now().Unix(), "当前时间戳")
	// log.Log(logkit.SeverityCritical, "%d--%s", time.Now().Unix(), "当前时间戳")
	// log.Log(logkit.SeverityAlert, "%d--%s", time.Now().Unix(), "当前时间戳")
	// log.Log(logkit.SeverityEmergency, "%d--%s", time.Now().Unix(), "当前时间戳")

	log.Debug("%d--%s", time.Now().Unix(), "当前时间戳")
	log.Info("%d--%s", time.Now().Unix(), "当前时间戳")
	log.Notice("%d--%s", time.Now().Unix(), "当前时间戳")
	log.Warning("%d--%s", time.Now().Unix(), "当前时间戳")
	log.Error("%d--%s", time.Now().Unix(), "当前时间戳")
	log.Critical("%d--%s", time.Now().Unix(), "当前时间戳")
	log.Alert("%d--%s", time.Now().Unix(), "当前时间戳")
	log.Emergency("%d--%s", time.Now().Unix(), "当前时间戳")

	l2 := log.Copy("axlrose", func(kit *logkit.LogKit, entry *logkit.Entry) string {
		return fmt.Sprintf("[%v] %s", kit.Category, entry.Message)
	})
	for i := 0; i <= 10000; i++ {
		l2.Error("new category message %d",i)
	}
	fmt.Println(runtime.NumGoroutine())
}
