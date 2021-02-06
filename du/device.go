package du

import (
	"errors"

	"github.com/tidwall/gjson"
)

//Device 设备
type Device struct {
	du *Du
	id string
}

//Act 发送指令。适合除开 applianceId 不需要其他 payload 信息的情况
func (dev *Device) Act(action string) error {
	return dev.Control(action, nil)
}

//Control 控制
func (dev *Device) Control(action string, payload interface{}) error {
	if dev.du == nil || dev.id == "" {
		return errors.New("you need connect this device to du first")
	}
	// 处理指令名称：首字母大写，加上"Request"后缀
	action += "Request"
	name := []byte{action[0] - 32}
	name = append(name, action[1:]...)
	// 合成控制条件下的 payload
	payload = genCtPayload(dev.id, payload)
	directive := Directive{
		Header:  Header{Namespace: "DuerOS.ConnectedHome.Control", Name: string(name), Version: 1},
		Payload: payload,
	}
	return dev.du.Do(directive)
}

//GetInfo 获取设备信息
func (dev *Device) GetInfo() (info gjson.Result, err error) {
	if dev.du == nil || dev.id == "" {
		err = errors.New("you need connect this device to du first")
		return
	}
	devsInfo, err := dev.du.GetDevs()
	if err != nil {
		return
	}
	var ok bool
	if info, ok = devsInfo[dev.id]; !ok {
		err = errors.New("the device info not found")
	}
	return
}
