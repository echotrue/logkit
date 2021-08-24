package logkit

import (
	"log"
	"testing"
	"time"
)

func TestNewConsoleTarget(t *testing.T) {
	logKit := NewLogKitByOptions(WithCallStackDepth(0), WithBuffer(1024), WithTargets(NewConsoleTarget()))
	err := logKit.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer logKit.Close()
	logKit.Error("Message, current time: %d",time.Now().Nanosecond())
	// for i := 0; i <= 10000; i++ {
	// 	logKit.Error("Message %d, current time: %d", i, time.Now().Nanosecond())
	// }
}

func BenchmarkNewConsoleTarget(b *testing.B) {
	logKit := NewLogKitByOptions(WithCallStackDepth(0), WithBuffer(1024), WithTargets(NewConsoleTarget()))
	err := logKit.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer logKit.Close()

	logKit.Error("Message, current time: %d",time.Now().Nanosecond())
	// for i := 0; i <= 10000; i++ {
	// 	logKit.Error("Message %d, current time: %d", i, time.Now().Nanosecond())
	// }
}
