package frame

import (
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync/atomic"
	"syscall"

	"github.com/wspGreen/skyee/slog"
)

func GetClassName(refType reflect.Type) string {
	sType := refType.String()
	index := strings.Index(sType, ".")
	if index != -1 {
		sType = sType[index+1:]
	}
	return sType
}

func IsStopSvr() bool {
	return stop == 1
}

func IsRuning() bool {
	return stop == 0
}

func Stop() {
	if !atomic.CompareAndSwapInt32(&stop, 0, 1) {
		return
	}

	// close(stopChanForGo)
	// for sc := 0; !waitAll.TryWait(); sc++ {
	// 	Sleep(1)
	// 	if sc >= 3000 {
	// 		slog.Error("Server Stop Timeout")
	// 		// stopCheckMap.Lock()
	// 		// for _, v := range stopCheckMap.M {
	// 		// 	LogError("Server Stop Timeout:%v", v)
	// 		// }
	// 		// stopCheckMap.Unlock()
	// 		sc = 0
	// 	}
	// }

	slog.Info("Server Stop")
	close(stopChanForSys)
}

func WaitForSystemExit() {

	signal.Notify(stopChanForSys, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-stopChanForSys:
		Stop()
	}

	stopLog()
}

func stopLog() {
	if !atomic.CompareAndSwapInt32(&stopForLog, 0, 1) {
		return
	}
	close(stopChanForLog)
	waitAllForLog.Wait()
}
