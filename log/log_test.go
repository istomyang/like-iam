package log

import "testing"

func TestLogger(t *testing.T) {
	logger := Default()
	logger.Info("test")
}
