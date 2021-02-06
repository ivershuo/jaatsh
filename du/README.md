# 小度控制

使用小度控制其他设备，只实现简单的获取设备信息和控制设备。

 ## EXAMPLES
 ```go
import "github.com/ivershuo/jaatsh/du"

// 以下信息自行通过抓包获取
const (
  cid   = "your cid"
  cuid  = "your cuid"
  bdid  = "your bdid"
  bduss = "your bduss"
)
xiaodu :=  NewDu(cid, cuid, bdid, bduss)

// 加载一个设备。appliance id 可提前抓包配置好或通过 GetDevs() 获取
device := xiaodu.NewDevice("a_appliance_id")

// 发送打开设备指令
err := device.Act("turnOn")

/* 将设备（空调）调至制热模式 */
type payloadMode struct {
  CtPayload //嵌套 CtPayload 后可自动添加 appliance 信息
  Mode struct {
    Value string `json:"value"`
  } `json:"mode"`
}

mode := &payloadMode{}
mode.Mode.Value = "HEAT"
err := device.Control("setMode", mode)
 ```