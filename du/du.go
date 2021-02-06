package du

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

//APIHost 小度接口地址
const APIHost = "https://xiaodu.baidu.com/saiya/smarthome"
const referer = "https://xiaodu.baidu.com/saiya/smarthome/index.html"
const userAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 ios/14.4 sdk/0.9.0 FromApp/XiaoDuApp oneapp/3.37.0 bundleID/com.baidu.hec.XiaoduAIAssistant build/100"

//Du 小度
type Du struct {
	cid   string
	cuid  string
	bdid  string
	bduss string
	Debug bool
}

//NewDu 创建小度
func NewDu(cid, cuid, bdid, bduss string) *Du {
	return &Du{cid, cuid, bdid, bduss, false}
}

//GetDevs 获取小度已添加的设备信息
func (du *Du) GetDevs() (map[string]gjson.Result, error) {
	data, err := du.request("GET", "/devicelist?from=h5_control", nil)
	if err != nil {
		return nil, err
	}
	jsonData := gjson.ParseBytes(data)
	appliances := jsonData.Get("appliances").Array()
	devs := make(map[string]gjson.Result)
	for _, dev := range appliances {
		id := dev.Get("applianceId").String()
		devs[id] = dev
	}
	return devs, nil
}

//NewDevice 加载一个小度控制的设备
func (du *Du) NewDevice(id string) *Device {
	return &Device{
		id: id,
		du: du,
	}
}

//Do 发送指令
func (du *Du) Do(d Directive) error {
	body, _ := json.Marshal(d)
	_, err := du.request("POST", "/directivesend?from=h5_control", bytes.NewReader(body))
	return err
}

// 小度接口返回
type duResp struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func (du *Du) request(method, path string, body io.Reader) ([]byte, error) {
	if du.Debug {
		fmt.Printf("Req %s \"%s\"\n", method, path)
		if body != nil {
			fmt.Printf("Req body:\n%s\n", body)
		}
	}
	req, err := http.NewRequest(method, APIHost+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Referer", referer)
	req.Header.Add("User-Agent", userAgent)
	if body != nil {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}
	du.addCookie(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if du.Debug {
		fmt.Printf("Resp body:\n%s\n\n", respBody)
	}
	if err != nil {
		return nil, err
	}
	data := &duResp{}
	if err := json.Unmarshal(respBody, data); err != nil {
		return nil, err
	}
	if data.Status != 0 {
		return nil, errors.New(data.Msg)
	}
	return json.Marshal(data.Data)
}
func (du *Du) addCookie(req *http.Request) {
	req.AddCookie(&http.Cookie{
		Name:  "client-id",
		Value: du.cid,
	})
	req.AddCookie(&http.Cookie{
		Name:  "device-id",
		Value: du.cuid,
	})
	req.AddCookie(&http.Cookie{
		Name:  "BAIDUID",
		Value: du.bdid,
	})
	req.AddCookie(&http.Cookie{
		Name:  "BDUSS",
		Value: du.bduss,
	})
}
