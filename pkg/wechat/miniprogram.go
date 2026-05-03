package wechat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type JsCode2SessionResp struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid,omitempty"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type MiniProgram struct {
	appID     string
	appSecret string
}

func NewMiniProgram(appID, appSecret string) *MiniProgram {
	return &MiniProgram{appID: appID, appSecret: appSecret}
}

func (mp *MiniProgram) JsCode2Session(jsCode string) (*JsCode2SessionResp, error) {
	api := "https://api.weixin.qq.com/sns/jscode2session"
	params := url.Values{}
	params.Set("appid", mp.appID)
	params.Set("secret", mp.appSecret)
	params.Set("js_code", jsCode)
	params.Set("grant_type", "authorization_code")

	resp, err := http.Get(api + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("wechat api request failed: %w", err)
	}
	defer resp.Body.Close()

	var result JsCode2SessionResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode wechat response: %w", err)
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("wechat api error: code=%d, msg=%s", result.ErrCode, result.ErrMsg)
	}

	return &result, nil
}
