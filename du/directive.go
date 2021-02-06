package du

import "reflect"

//Directive 指令
type Directive struct {
	Header  Header      `json:"header"`
	Payload interface{} `json:"payload"`
}

//Header 指令header
type Header struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Version   int    `json:"payloadVersion"`
}

//CtPayload 控制类 payload
type CtPayload struct {
	App payloadApp `json:"appliance"`
}
type payloadApp struct {
	ID []string `json:"applianceId"`
}

func genCtPayload(devID string, payload interface{}) interface{} {
	ctPayload := &CtPayload{}
	pba := payloadApp{[]string{devID}}

	if payload == nil {
		ctPayload.App = pba
		payload = ctPayload
	} else {
		v := reflect.ValueOf(payload)
		switch v.Kind() {
		case reflect.Ptr:
			ps := reflect.ValueOf(payload).Elem()
			pn := ps.FieldByName("App")
			if pn.IsValid() && pn.IsZero() {
				ppb := pn.Addr().Interface().(*payloadApp)
				*ppb = pba
			}
		}
	}
	return payload
}
