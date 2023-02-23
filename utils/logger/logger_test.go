package logger

import "testing"

func TestLogger(t *testing.T) {
	Debug("Debug")
	Info("Info")
	Warn("Warn")
	Error("Error")
}
