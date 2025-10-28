package logx

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func Init(logPath string) error {
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	multi := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multi)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return nil
}
