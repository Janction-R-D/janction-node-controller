package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
)

type RecoveredErr struct {
	Content string
}

func (e RecoveredErr) Error() string {
	return e.Content
}

func RecoverAndRePanic(err *error) {
	if panicE := recover(); panicE != nil {
		var buf [10240]byte
		var n = runtime.Stack(buf[:], false)
		var errContent = fmt.Sprintf("%+v =====> [stack-info-when-panic]\n %s", panicE, string(buf[:n]))
		if err == nil {
			logrus.Error(errContent)
		} else {
			if (*err) != nil && (*err).Error() != "" {
				logrus.Errorf("=====> [err-info-before-panic] %s", (*err).Error())
			}
			*err = RecoveredErr{errContent}
		}
		panic(panicE)
	}
}

func PrintPanicAndRePanic() {
	if panicE := recover(); panicE != nil {
		var buf [10240]byte
		var n = runtime.Stack(buf[:], false)
		var errContent = fmt.Sprintf("%+v =====> [stack-info-when-panic]\n %s", panicE, string(buf[:n]))
		logrus.Error(errContent)
		panic(panicE)
	}
}

// RecoverUnexpectedErr 捕捉某个函数发生的未知panic，由于是使用recover()实现的，因此RecoverUnexpectedErr必须由defer调用才可生效
func RecoverUnexpectedErr(funcName string) {
	panicE := recover()
	if panicE == nil {
		return
	}

	var buf [10240]byte
	var n = runtime.Stack(buf[:], false)
	var errContent = fmt.Sprintf("%+v =====> [stack-info-when-panic]\n %s", panicE, string(buf[:n]))

	switch panicE.(type) {
	case runtime.Error:
		logrus.Errorf("runtime panicE in func %s(): %+v \n %s", funcName, panicE, errContent)
	default:
		logrus.Errorf("application error in func %s(): %+v \n %s", funcName, panicE, errContent)
	}
}

func GetStack(length int) string {
	var buf []byte
	if length <= 0 {
		buf = make([]byte, 0, 10240)
	} else {
		buf = make([]byte, 0, length)
	}
	var n = runtime.Stack(buf, false)
	return BytesToStr(buf[:n])
}
