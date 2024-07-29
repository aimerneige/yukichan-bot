package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const QueryNonceLength = 32
const QueryApiAddress = "https://openaiapi.weixin.qq.com/v2/bot/query"

type ApiReq struct {
	Query    string `json:"query"`
	Env      string `json:"env"`
	UserName string `json:"user_name"`
	Avatar   string `json:"avatar"`
	Userid   string `json:"userid"`
}

type ApiRsp struct {
	Code      int32      `json:"code"`
	Msg       string     `json:"msg"`
	Data      ApiRspData `json:"data"`
	RequestId string     `json:"request_id"`
}

type ApiRspData struct {
	Answer     string       `json:"answer"`
	AnswerType string       `json:"answer_type"`
	IntentName string       `json:"intent_name"`
	MsgId      string       `json:"msg_id"`
	Options    []Option     `json:"options"`
	SkillName  string       `json:"skill_name"`
	Slots      []SlotDetail `json:"slots"`
	Status     string       `json:"status"`
}

type Option struct {
	AnsNodeName string  `json:"ans_node_name"`
	Title       string  `json:"title"`
	Answer      string  `json:"answer"`
	Confidence  float64 `json:"confidence"`
}

type SlotDetail struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Norm  string `json:"norm"`
}

func SendQueryRequest(apiReq ApiReq, accessToken, signToken, aeskey string) string {
	data, err := json.Marshal(apiReq)
	if err != nil {
		panic(err)
	}
	reqJson := string(data)
	cipher, err := Encrypt(aeskey, reqJson)
	if err != nil {
		panic(err)
	}
	response := apiCall(accessToken, signToken, aeskey, []byte(cipher))
	return response
}

func apiCall(accessToken, signToken, aeskey string, body []byte) string {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := GenRandomLetterString(QueryNonceLength)
	sign := CalcSign(signToken, timestamp, nonce, body)

	client := http.Client{}
	request, err := http.NewRequest("POST", QueryApiAddress, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	request.Header = http.Header{
		"X-OPENAI-TOKEN": {accessToken},
		"request_id":     {GetUUid()},
		"timestamp":      {timestamp},
		"nonce":          {nonce},
		"sign":           {sign},
		"Content-Type":   {"text/plain"},
	}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	cipher := string(respBody)
	plain, err := Decrypt(aeskey, cipher)
	if err != nil {
		panic(err)
	}
	apiResp := ApiRsp{}
	if err := json.Unmarshal([]byte(plain), &apiResp); err != nil {
		panic(err)
	}
	return apiResp.Data.Answer
}
