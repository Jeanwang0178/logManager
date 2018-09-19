package models

import (
	"fmt"
	"github.com/beego/bee/logger"
	"strconv"
	"strings"
)

type GoRoutineManager struct {
	grchannelMap *GoroutineChannelMap
}

func NewGoRoutineManager() *GoRoutineManager {
	gm := &GoroutineChannelMap{}
	return &GoRoutineManager{grchannelMap: gm}
}

func (gm GoRoutineManager) StopLoopGoroutine(name string) error {
	stopChannel, ok := gm.grchannelMap.grchannels[name]
	if !ok {
		return fmt.Errorf("not found goroutine name :" + name)
	}
	gm.grchannelMap.grchannels[name].msg <- STOP + strconv.Itoa(int(stopChannel.gid))
	return nil
}

func (gm *GoRoutineManager) RegisterGoroutine(name string) {
	go func(n string) {
		//register channel
		err := gm.grchannelMap.register(n)
		if err != nil {
			beeLogger.Log.Errorf("grchannelMap register: %v", err)
			return
		}
	}(name)

}

func (gm *GoRoutineManager) NewLoopGoroutine(name string, fc interface{}, args ...interface{}) {

	go func(this *GoRoutineManager, n string, fc interface{}, args ...interface{}) {
		//register channel
		err := this.grchannelMap.register(n)
		if err != nil {
			return
		}
		for {
			select {
			case info := <-this.grchannelMap.grchannels[name].msg:
				taskInfo := strings.Split(info, ":")
				signal, gid := taskInfo[0], taskInfo[1]
				if gid == strconv.Itoa(int(this.grchannelMap.grchannels[name].gid)) {
					if signal == "_STOP" {

						beeLogger.Log.Info(name + "：gid[" + gid + "] quit")
						this.grchannelMap.unregister(name)
						return
					} else {
						beeLogger.Log.Info("unknow signal")
					}
				}
			default:
				beeLogger.Log.Info("no signal")
			}

			if len(args) > 1 {
				fc.(func(...interface{}))(args)
			} else if len(args) == 1 {
				fc.(func(interface{}))(args[0])
			} else {
				fc.(func())()
			}
		}
	}(gm, name, fc, args...)

}

func (gm *GoRoutineManager) NewGoroutine(name string, fc interface{}, args ...interface{}) {
	go func(n string, fc interface{}, args ...interface{}) {
		//register channel
		err := gm.grchannelMap.register(n)
		if err != nil {
			beeLogger.Log.Errorf("grchannelMap register: %v", err)
			return
		}
		if len(args) > 1 {
			fc.(func(...interface{}))(args)
		} else if len(args) == 1 {
			fc.(func(interface{}))(args[0])
		} else {
			fc.(func())()
		}
		gm.grchannelMap.unregister(name)
	}(name, fc, args...)

}
