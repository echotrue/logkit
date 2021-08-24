package logkit

import (
	log2 "log"
	"testing"
	"time"
)

func TestNewFileTarget(t *testing.T) {
	ft := NewFileTarget()
	ft.FilePath = "./temp"
	log := NewLogKitByOptions(WithBuffer(1024), WithTargets(ft), WithCallStackDepth(0))
	err := log.Open()
	if err != nil {
		log2.Fatal(err)
	}
	defer log.Close()
	log.Error("Message, current time: %d",time.Now().Nanosecond())
	// for i := 0; i <= 10000; i++ {
	// 	log.Error("Message %d, current time: %d",i,time.Now().Nanosecond())
	// }
}

func BenchmarkNewFileTarget(b *testing.B) {
	ft := NewFileTarget()
	ft.FilePath = "./temp"
	log := NewLogKitByOptions(WithBuffer(1024), WithTargets(ft), WithCallStackDepth(0))
	err := log.Open()
	if err != nil {
		log2.Fatal(err)
	}
	defer log.Close()

	log.Error("Message, current time: %d",time.Now().Nanosecond())
	// for i := 0; i <= 10000; i++ {
	// 	log.Error("Message %d, current time: %d",i,time.Now().Nanosecond())
	// }
}
