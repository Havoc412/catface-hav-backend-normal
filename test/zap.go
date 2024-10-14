package test

import "go.uber.org/zap"

func main() {
	log := zap.NewExample()
	log.Debug("this is debug message")
	log.Info("this is info message")
	log.Info("this is info message with fileds",
		zap.Int("age", 24), zap.String("agender", "man"))
	log.Warn("this is warn message")
	log.Error("error message")
	log.Panic("panic message")
}
