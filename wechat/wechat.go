/*
一个基于企业微信应用的消息推送服务
应用文档：https://work.weixin.qq.com/api/doc/90000/90135/90236
*/

package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Wechat struct {
	AgentId    string
	CorpSecret string
	CorpID     string
}

type accessTokenType struct {
	AccessToken string `json:"access_token"`
}

type sendDataTextType struct {
	Content string `json:"content"`
}

type sendDataType struct {
	Touser                 string           `json:"touser"`
	MsgType                string           `json:"msgtype"`
	AgentId                string           `json:"agentid"`
	Text                   sendDataTextType `json:"text"`
	Safe                   int              `json:"safe"`
	EnableIdTrans          int              `json:"enable_id_trans"`
	EnableDuplicateCheck   int              `json:"enable_duplicate_check"`
	DuplicateCheckInterval int              `json:"duplicate_check_interval"`
}

type returnDataType struct {
	ErrCode int `json:"errcode"`
}

// get access token
func (w *Wechat) getToken() (error, string) {
	url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + w.CorpID + "&corpsecret=" + w.CorpSecret
	resp, err := http.Get(url)
	if err != nil {
		return err, ""
	}
	var token accessTokenType
	respData, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(respData, &token)
	if err != nil {
		return err, ""
	}
	return nil, token.AccessToken
}

func (w *Wechat) SendText(content, users string) error {
	err, token := w.getToken()
	if err != nil {
		return err
	}
	url := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + token
	data := sendDataType{
		Touser:  users,
		MsgType: "text",
		AgentId: w.AgentId,
		Text: sendDataTextType{
			Content: content,
		},
		Safe:                   0,
		EnableDuplicateCheck:   0,
		EnableIdTrans:          0,
		DuplicateCheckInterval: 1800,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	var respCode returnDataType
	respData, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(respData, &respCode)
	if err != nil {
		return err
	}
	if respCode.ErrCode != 0 {
		return fmt.Errorf("return error code is %d\n", respCode.ErrCode)
	}
	return nil
}
