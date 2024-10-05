package log

import "os"

// CustomWriteSyncer implements zapcore.WriteSyncer
type CustomWriteSyncer struct {
	writer *os.File
}

func (cws *CustomWriteSyncer) Write(p []byte) (n int, err error) {
	return cws.writer.Write(p)
}

func (cws *CustomWriteSyncer) Sync() error {
	return cws.writer.Sync()
}
