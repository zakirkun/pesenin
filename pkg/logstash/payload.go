package logstash

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"runtime"
	"strings"

	"github.com/zakirkun/pesnin/pkg/config"
)

type LogPayload struct {
	Appname string `json:"appname"`
	Level   string `json:"level"`
	Index   any    `json:"index"`
	Message string `json:"message"`
	Meta    any    `json:"meta"`
	Data    any    `json:"data"`
}

type metaPayload struct {
	FuncName  string `json:"func"`
	ShortFunc string `json:"short_func_name"`
	File      string `json:"file"`
	Line      int    `json:"line"`
}

func (l *LogPayload) ToJson() []byte {
	body, err := json.Marshal(&l)
	if err != nil {
		fmt.Println("Error marshal payload Logstash:", err)
		return nil
	}

	return body
}

func (l *LogPayload) WriteCaller(level string, message string) error {

	// Set App Name
	if l.Appname == "" {
		l.Appname = config.GetString("logstash.app_name")
	}

	l.Level = level
	l.Message = message

	// Set Meta File
	l.setMeta()

	conn, err := LOGSTASH.Open()
	if err != nil {
		log.Printf("Error connect Log Stash: %v", err)
		return err
	}
	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			log.Printf("Error closing Logstash connection: %v", err.Error())
		}
	}(conn)

	_, err = conn.Write(l.ToJson())
	if err != nil {
		log.Printf("Error write Log Stash after reconnect: %v", err)
		return err
	}

	return nil
}

func (l *LogPayload) SetAppName(app string) *LogPayload {
	l.Appname = app
	return l
}

func (l *LogPayload) SetIndex(k any) *LogPayload {
	l.Index = k
	return l
}

func (l *LogPayload) SetData(data any) *LogPayload {
	l.Data = data
	return l
}

func (l *LogPayload) setMeta() *LogPayload {

	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		log.Print("Could not get caller information")
		return l
	}

	// Get function name
	funcName := runtime.FuncForPC(pc).Name()
	funcNameParts := strings.Split(funcName, "/")
	shortFuncName := funcNameParts[len(funcNameParts)-1]

	meta := metaPayload{
		FuncName:  funcName,
		File:      file,
		Line:      line,
		ShortFunc: shortFuncName,
	}
	l.Meta = &meta

	return l
}
