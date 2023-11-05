package Websocket

type Data struct {
	ClientId string      `json:"client_id"`
	Data     interface{} `json:"data"`
	Message  string      `json:"message"`
	Success  bool        `json:"success"`
	Type     string      `json:"type"`
	From     string      `json:"from"`
	To       string      `json:"to"`
}
