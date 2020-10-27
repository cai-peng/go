package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

const (
	KB = 1 << 10
	MB = KB << 10
	GB = MB << 10
)

type EXIT bool

type level int

const (
	DetailLog level = iota
	DebugLog
	InfoLog
	WarnLog
	ErrorLog
	PanicLog
)

var levelName = []string{
	DetailLog: "DETAULT",
	DebugLog:  "DEBUG",
	InfoLog:   "INFO",
	WarnLog:   "WARN",
	ErrorLog:  "ERROR",
	PanicLog:  "PANIC",
}

type loggingT struct {
	toStderr    bool
	outputLevel level
	mutex       sync.Mutex
	filename    string
	writer      *os.File
	backCount   int
	maxBytes    int64
}

var LogPath = "./error.log"
var logging = &loggingT{toStderr: true, outputLevel: DetailLog, writer: os.Stderr}

func setLevel(l level) {
	setRotatingFile()
	logging.outputLevel = l
}

func setRotatingFile() bool {
	dir := filepath.Dir(LogPath)
	if !strings.EqualFold(dir, ".") {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("ERROR - MkdirAll(%s): %s\n", dir, err.Error())
			return false
		}
	}
	f, err := os.OpenFile(LogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("ERROR - OpenFile(%s): %s\n", LogPath, err.Error())
		return false
	}

	logging.toStderr = false
	logging.filename = LogPath
	logging.writer = f
	logging.maxBytes = 10 * GB
	logging.backCount = 10
	return true
}

func (self *loggingT) rotatingFile() {
	self.writer.Close()
	self.writer = nil

	for i := self.backCount; i >= 0; i-- {
		backFilename := self.filename + "." + strconv.Itoa(i)
		backFilename_obj := self.filename + "." + strconv.Itoa(i+1)
		if i > 0 && !FileExist(backFilename) {
			continue
		}
		if i == self.backCount {
			os.Remove(backFilename)
		} else if i == 0 {
			os.Rename(self.filename, backFilename_obj)
		} else {
			os.Rename(backFilename, backFilename_obj)
		}
	}

	f, err := os.OpenFile(self.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	self.writer = f
}

func (self *loggingT) printf(l level, format string, v ...interface{}) {
	if self.outputLevel > l {
		return
	}
	msg := fmt.Sprintf(format, v...)
	self.mutex.Lock()
	defer self.mutex.Unlock()
	if self.toStderr == false {
		info, err := self.writer.Stat()
		if err == nil && info.Size() >= self.maxBytes {
			self.rotatingFile()
		}
	}
	_, file, _, ok := runtime.Caller(2)
	if !ok {
		file = "Runtime.Caller occurred an unknown exception???"
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	if l == 0 {
		fmt.Fprintf(self.writer, "%s %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
	} else {
		fmt.Fprintf(self.writer, "%s %s - %s [%s]\n",
			time.Now().Format("2006-01-02 15:04:05"),
			levelName[l], file, msg)
	}
}

func Default(format string, v ...interface{}) {
	setLevel(0)
	logging.printf(DetailLog, format, v...)
}

func Debug(format string, v ...interface{}) {
	setLevel(1)
	logging.printf(DebugLog, format, v...)
}

func Info(format string, v ...interface{}) {
	setLevel(2)
	logging.printf(InfoLog, format, v...)
}
func Warn(format string, v ...interface{}) {
	setLevel(3)
	logging.printf(WarnLog, format, v...)
}

func Error(format string, v ...interface{}) {
	setLevel(4)
	logging.printf(ErrorLog, format, v...)
}

func Panic(format string, v ...interface{}) {
	setLevel(5)
	logging.printf(PanicLog, format, v...)
	os.Exit(1)
}
