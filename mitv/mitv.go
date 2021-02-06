package mitv

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//MiInfo 电视信息
type MiInfo struct {
	DeviceName string `json:"devicename"`
	DeviceID   string `json:"deviceid"`
	Ptf        int    `json:"ptf"`
	Codever    int    `json:"codever"`
	WifiMac    string `json:"wifimac"`
	EthMac     string `json:"ethmac"`
}

//Mi 小米电视连接
type Mi struct {
	host string
}

//New 创建小米电视连接
func New(hostName string) *Mi {
	return &Mi{hostName + ":6095"}
}

//Ping ping
func (mi *Mi) Ping() bool {
	_, err := net.DialTimeout("tcp", mi.host, time.Second*2)
	return err == nil
}

//Power 关机
func (mi *Mi) Power() error {
	return mi.SendKey("power")
}

//Up 上
func (mi *Mi) Up() error {
	return mi.SendKey("up")
}

//Down 下
func (mi *Mi) Down() error {
	return mi.SendKey("down")
}

//Left 左
func (mi *Mi) Left() error {
	return mi.SendKey("left")
}

//Right 右
func (mi *Mi) Right() error {
	return mi.SendKey("right")
}

//Home 返回home
func (mi *Mi) Home() error {
	return mi.SendKey("home")
}

//Enter 确认
func (mi *Mi) Enter() error {
	return mi.SendKey("enter")
}

//Back 返回
func (mi *Mi) Back() error {
	return mi.SendKey("back")
}

//Menu 菜单
func (mi *Mi) Menu() error {
	return mi.SendKey("menu")
}

//VolumeUp 音量增加
func (mi *Mi) VolumeUp() error {
	return mi.SendKey("volumeup")
}

//VolumeDown 音量减少
func (mi *Mi) VolumeDown() error {
	return mi.SendKey("volumedown")
}

//VolumeSet 设置音量
func (mi *Mi) VolumeSet(volume int) error {
	info, err := mi.GetInfo()
	if err != nil {
		return err
	}
	volStr := strconv.Itoa(volume)
	ts := "11111" // 可以直接使用固定值
	signStr := "mitvsignsalt" + volStr + info.EthMac + ts
	sign := md5.Sum([]byte(signStr))
	v := url.Values{}
	v.Add("action", "setVolum")
	v.Add("volum", volStr)
	v.Add("ts", ts)
	v.Add("sign", hex.EncodeToString(sign[:]))
	data, err := mi.RequestAPI("/general", v)
	if err != nil {
		return err
	}
	if data.Ret != 200 {
		return errors.New("error")
	}
	return nil
}

type volStruct struct {
	Volume int `json:"volum"`
}

//GetVolume 获取当前音量
func (mi *Mi) GetVolume() (int, error) {
	v := url.Values{}
	v.Add("action", "getVolum")
	data, err := mi.RequestAPI("/general", v)
	if err != nil {
		return -1, err
	}
	if data.Ret != 200 {
		return -1, errors.New("error")
	}
	volume := &volStruct{}
	var s string
	str, err := json.Marshal(data.Data)
	if err != nil {
		return -1, err
	}
	err = json.Unmarshal(str, &s)
	if err != nil {
		return -1, err
	}
	err = json.Unmarshal([]byte(s), volume)
	return volume.Volume, err
}

//DataResp 接口返回数据
type DataResp struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Ret    int         `json:"request_result"`
}

//SendKey 发送按键信息
func (mi *Mi) SendKey(key string) error {
	v := url.Values{}
	v.Add("action", "keyevent")
	v.Add("keycode", key)
	data, err := mi.RequestAPI("/controller", v)
	if err != nil {
		return err
	} else if data.Status != 0 || data.Msg != "success" {
		return errors.New("error")
	}
	return nil
}

//GetInfo 获取电视信息
func (mi *Mi) GetInfo() (*MiInfo, error) {
	v := url.Values{}
	v.Add("action", "getsysteminfo")
	info := &MiInfo{}
	data, err := mi.RequestAPI("/controller", v)
	if err != nil {
		return info, err
	} else if data.Status != 0 || data.Msg != "success" {
		return info, errors.New("error")
	}
	infoData, err := json.Marshal(data.Data)
	if err == nil {
		err = json.Unmarshal(infoData, info)
	}
	return info, err
}

//RequestAPI 请求电视接口
func (mi *Mi) RequestAPI(path string, v url.Values) (*DataResp, error) {
	apiURL := "http://" + mi.host + path + "?" + v.Encode()
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	body := resp.Body
	defer body.Close()
	bodyData, err := ioutil.ReadAll(body)
	miResp := &DataResp{}
	err = json.Unmarshal(bodyData, miResp)
	return miResp, err
}
