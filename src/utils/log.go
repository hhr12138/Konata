package utils

import (
	"fmt"
	"github.com/hhr12138/Konata/src/consts"
	"log"
	"time"
)

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if consts.SHOW_LOG {
		log.Printf(format, a...)
	}
	return
}

func Printf(level consts.LogLevel, serviceId string, term int, offset int, msg string, a ...interface{}) {
	var (
		now     = time.Now()
		message = fmt.Sprintf(msg, a...)
	)
	if consts.SHOW_LOG {
		log.Printf("[%v] %v.%v.%v %v %v", level, serviceId, term, offset, message, now.Format(consts.TIME_FORMAT_MS))
	}
}
