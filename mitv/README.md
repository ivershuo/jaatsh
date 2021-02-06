# 小米电视控制

使用小米电视自带的接口服务控制电视，能控制除“开机”之外的其他所有正常操作。

 ## EXAMPLE
 ```go
import "github.com/ivershuo/jaatsh/mitv"

tv := mitv.New("192.168.0.201") // 电视IP
err := tv.VolumeUp()
 ```