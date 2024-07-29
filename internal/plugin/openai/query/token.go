package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const TokenNonceLength = 32
const TokenApiAddress = "https://openaiapi.weixin.qq.com/v2/token"

type TokenResponse struct {
	Code int `json:"code"`
	Data struct {
		AccessToken string `json:"access_token"`
	} `json:"data"`
	Msg       string `json:"msg"`
	RequestID string `json:"request_id"`
}

func GetToken(account, appId, token string) string {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := GenRandomLetterString(TokenNonceLength)
	bodyStr := fmt.Sprintf(`{"account":"%s"}`, account)
	body := []byte(bodyStr)
	sign := CalcSign(token, timestamp, nonce, body)

	client := http.Client{}
	request, err := http.NewRequest("POST", TokenApiAddress, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	request.Header = http.Header{
		"X-APPID":    {appId},
		"request_id": {GetUUid()},
		"timestamp":  {timestamp},
		"nonce":      {nonce},
		"sign":       {sign},
	}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	var tokenResp TokenResponse
	err = json.Unmarshal(respBody, &tokenResp)
	if err != nil {
		fmt.Printf("request: %+v, respBody: %v, err: %v", request, string(respBody), err)
		panic(err)
	}
	return tokenResp.Data.AccessToken
}
