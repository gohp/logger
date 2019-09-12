package logde

import (
	"github.com/wzyonggege/logde"
	"testing"
)


func BenchmarkLogger(b *testing.B)  {
	b.Logf("Logging at a disabled level with some accumulated context.")
	b.Run("logde logger without fields", func(b *testing.B) {
		c := logde.New()
		c.CloseConsoleDisplay()
		logger := c.InitLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info("1234")
			}
		})
	})
	b.Run("logde logger with fields", func(b *testing.B) {
		c := logde.New()
		c.CloseConsoleDisplay()
		logger := c.InitLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info("1234", logger.With("Trace", "1234455"))
			}
		})
	})
	b.Run("logde logger without fields write into file", func(b *testing.B) {
		c := logde.New()
		c.CloseConsoleDisplay()
		c.SetInfoFile("../logs/test_stdout.log")
		logger := c.InitLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info("1234")
			}
		})
	})
	b.Run("logde logger with fields write into file", func(b *testing.B) {
		c := logde.New()
		c.CloseConsoleDisplay()
		c.SetInfoFile("../logs/test_stdout.log")
		logger := c.InitLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info("1234", logger.With("Trace", "1234455"))
			}
		})
	})
}

