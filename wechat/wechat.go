package wechat

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-gin-weixin/config"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

type Token struct {
	AccessToken string `json:"access_token"`
	ExpressIn   string `json:"express_in"`
}

type BindMessage struct {
	Touser     string   `json:"touser"`
	TemplateId string   `json:"template_id"`
	Data       BindData `json:"data"`
}

type BindData struct {
	First    FormData `json:"first"`
	Keyword1 FormData `json:"keyword1"`
	Keyword2 FormData `json:"keyword2"`
	Keyword3 FormData `json:"keyword3"`
	Keyword4 FormData `json:"keyword4"`
	Remark   FormData `json:"remark"`
}

type FormData struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

// CheckSignature 微信公众号签名检查
func CheckSignature(signature, timestamp, nonce, token string) bool {
	arr := []string{timestamp, nonce, token}
	// 字典序排序
	sort.Strings(arr)

	n := len(timestamp) + len(nonce) + len(token)
	var b strings.Builder
	b.Grow(n)

	for i := 0; i < len(arr); i++ {
		b.WriteString(arr[i])
	}

	return Sha1(b.String()) == signature
}

// Sha1 进行Sha1编码
func Sha1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

type AccessTokenResponse struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float64 `json:"expires_in"`
}

type AccessTokenErrorResponse struct {
	Errcode float64
	Errmsg  string
}

// FetchAccessToken 获取 AccessToken 调用业务接口时需要
func FetchAccessToken(appID, appSecret, accessTokenFetchUrl string) (string, error) {

	requestLine := strings.Join([]string{accessTokenFetchUrl,
		"?grant_type=client_credential&appid=",
		appID,
		"&secret=",
		appSecret}, "")

	resp, err := http.Get(requestLine)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("发送get请求获取 atoken 错误", err)
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("发送get请求获取 atoken 读取返回body错误", err)
		return "", err
	}

	if bytes.Contains(body, []byte("access_token")) {
		atr := AccessTokenResponse{}
		err = json.Unmarshal(body, &atr)
		if err != nil {
			fmt.Println("发送get请求获取 atoken 返回数据json解析错误", err)
			return "", err
		}
		return atr.AccessToken, nil
	} else {
		fmt.Println("发送get请求获取 微信返回 err")
		ater := AccessTokenErrorResponse{}
		err = json.Unmarshal(body, &ater)
		fmt.Printf("发送get请求获取 微信返回 的错误信息 %+v\n", ater)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("%s", ater.Errmsg)
	}
}

// Model 模板消息推送
func Model(wxID string, time string, phone string) error {

	FirstData := "你好!欢迎使用微信公众号测试版"
	RemarkData := "时间：" + time

	var bodyJson []byte

	accessToken, _ := FetchAccessToken(config.APP_ID, config.APP_SECRECT, "https://api.weixin.qq.com/cgi-bin/token")
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + accessToken
	bm := BindMessage{
		Touser:     wxID,
		TemplateId: "PXccCZB8xgC_FVIIS6yHyooM6fXtX7kPySyNBSkjkOc",
		Data: BindData{
			First: FormData{
				Value: FirstData,
				Color: "#173177",
			},
			Keyword1: FormData{
				Value: phone,
				Color: "#173177",
			},
			Keyword2: FormData{
				Value: "aaa",
				Color: "#173177",
			},
			Remark: FormData{
				Value: RemarkData,
				Color: "#173177",
			},
		},
	}

	var err error
	bodyJson, err = json.Marshal(bm)
	if err != nil {
		log.Error(err)
		return errors.New("http post body to json failed")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Error(err)
		return errors.New("new request is fail: %v \n")
	}
	req.Header.Set("Content-type", "application/json")

	//http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return errors.New("response is fail: %v \n")
	}
	defer resp.Body.Close()

	respbody, err := ioutil.ReadAll(resp.Body)
	log.Println(string(respbody), err)

	return err
}
