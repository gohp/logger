package logger

//import (
//	logger "github.com/wzyonggege/logger"
//	"testing"
//)
//
//
//func BenchmarkLogger(b *testing.B)  {
//	b.Logf("Logging at a disabled level with some accumulated context.")
//	b.Run("logde logger without fields", func(b *testing.B) {
//		c := logger.New()
//		c.CloseConsoleDisplay()
//		c.InitLogger()
//		b.ResetTimer()
//		b.RunParallel(func(pb *testing.PB) {
//			for pb.Next() {
//				logger.Info("1234")
//			}
//		})
//	})
//	b.Run("logde logger with fields", func(b *testing.B) {
//		c := logger.New()
//		c.CloseConsoleDisplay()
//		c.InitLogger()
//		b.ResetTimer()
//		b.RunParallel(func(pb *testing.PB) {
//			for pb.Next() {
//				logger.Info("1234", logger.With("Trace", "1234455"))
//			}
//		})
//	})
//	b.Run("logde logger without fields write into file", func(b *testing.B) {
//		c := logger.New()
//		c.CloseConsoleDisplay()
//		c.SetInfoFile("../logs/test_stdout.log")
//		c.InitLogger()
//		b.ResetTimer()
//		b.RunParallel(func(pb *testing.PB) {
//			for pb.Next() {
//				logger.Info("1234")
//			}
//		})
//	})
//	b.Run("logde logger with fields write into file", func(b *testing.B) {
//		c := logger.New()
//		c.CloseConsoleDisplay()
//		c.SetInfoFile("../logs/test_stdout.log")
//		c.InitLogger()
//		b.ResetTimer()
//		b.RunParallel(func(pb *testing.PB) {
//			for pb.Next() {
//				logger.Info("1234", logger.With("Trace", "1234455"))
//			}
//		})
//	})
//}

