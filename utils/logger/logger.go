/*
日志组件，如果需要统一，可以随时替换底层实现
*/

package logger

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type logWriter interface {
	WriteString(string)
	Flush()
}

const DEBUG = true

var writer logWriter = stdout{out: bufio.NewWriter(os.Stderr)}
var logLock = sync.Mutex{}

func Debug(messages ...interface{}) {
	if DEBUG {
		doLog(writer, 2, _debug, messages...)
	}
}
func Info(messages ...interface{})  { doLog(writer, 2, _info, messages...) }
func Warn(messages ...interface{})  { doLog(writer, 2, _warn, messages...) }
func Error(messages ...interface{}) { doLog(writer, 2, _error, messages...) }

const (
	cNORMAL   = "\033[m"
	cRED      = "\033[0;32;31m"
	cLightRed = "\033[1;31m"
	cGREEN    = "\033[0;32;32m"
	//cLIGHT_GREEN  = "\033[1;32m"
	//cBLUE         = "\033[0;32;34m"
	//cLIGHT_BLUE   = "\033[1;34m"
	//cDARY_GRAY    = "\033[1;30m"
	//cCYAN         = "\033[0;36m"
	//cLIGHT_CYAN   = "\033[1;36m"
	//cPURPLE       = "\033[0;35m"
	//cLIGHT_PURPLE = "\033[1;35m"
	cBROWN = "\033[0;33m"
	//cYELLOW       = "\033[1;33m"
	//cLIGHT_GRAY   = "\033[0;37m"
	//cWHITE        = "\033[1;37m"
)

const (
	_debug      = " [" + "DEBUG" + "] "
	_info       = " [" + cGREEN + "INFO" + cNORMAL + "] "
	_warn       = " [" + cRED + "WARN" + cNORMAL + "] "
	_error      = " [" + cLightRed + "ERROR" + cNORMAL + "] "
	serviceName = "[" + "storage" + "] "
)

var shortFileName = bytes.Compare([]byte(runtime.GOOS), []byte("linux")) == 0

func doLog(writer logWriter, callerBack int, level string, messages ...interface{}) {
	date := time.Now().Format("01-02 15:04:05")
	_, file, lineNo, ok := runtime.Caller(callerBack)
	if !ok {
		file = "unknown"
		lineNo = 0
	} else if shortFileName {
		_, file = path.Split(file)
		i := strings.LastIndex(file, ".go")
		if i > 0 {
			file = file[:i]
		}
	}

	sb := strings.Builder{}

	sb.WriteString(date)
	sb.WriteString(level)
	sb.WriteString(serviceName)
	sb.WriteString(cBROWN + file + cNORMAL)
	sb.WriteString(":")
	sb.WriteString(strconv.Itoa(lineNo))
	for _, msg := range messages {
		if msg == nil {
			msg = "nil"
		}
		msgVal := reflect.ValueOf(msg)
		if msgVal.Kind() == reflect.Ptr && msgVal.IsNil() {
			msg = "nil"
		}
		sb.WriteString(" ")
		switch msg.(type) {
		case fmt.Stringer:
			sb.WriteString(msg.(fmt.Stringer).String())
		default:
			sb.WriteString(fmt.Sprintf("%v", msg))
		}
	}
	sb.WriteString("\n")

	// 减少临界区
	logLock.Lock()
	defer logLock.Unlock()
	writer.WriteString(sb.String())
	writer.Flush()
}

type stdout struct {
	out *bufio.Writer
}

func (w stdout) WriteString(s string) {
	_, err := w.out.WriteString(s)
	if err != nil {
		_ = fmt.Errorf("Logger Error %s ", err.Error())
	}
}

func (w stdout) Flush() {
	err := w.out.Flush()
	if err != nil {
		_ = fmt.Errorf("Logger Error %s ", err.Error())
	}
}
