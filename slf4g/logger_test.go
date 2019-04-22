package slf4g

import (
	"fmt"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	SetLogPath("/data/logs/test")
	//logger := api.GetLogger("test1")
	//SetConsole()
	for i := 0; i < 20000; i++  {
		Infof("info-%d", i)
		time.Sleep(time.Millisecond * 1)
	}
	Infof("测试中文\n")
	Debug("debug")
	Error("error")
	Warn("warn")
	fmt.Println("ok")
	//slf4g.FlushLogger()
	//FlushLogger()
}
