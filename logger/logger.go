package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var L *zap.Logger

func init() {
	atom := zap.NewAtomicLevel()

	// To keep the example deterministic, disable timestamps in the output.
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""

	L = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
}
