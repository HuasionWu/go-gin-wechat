package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	e "go-gin-weixin/pkg/error"
	"io/ioutil"
	"net/http"
	_ "time"
)

//微信返回的 ticket + url
type PermQrcode struct {
	Ticket string `json:"ticket"`
	URL    string `json:"url"`
}

// WXTextMsg 微信文本消息结构体
type WXTextMsg struct {
	ToUserName   string `xml:"ToUserName"   json:"ToUserName"`
	FromUserName string `xml:"FromUserName" json:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"   json:"CreateTime"`
	MsgType      string
	Event        string  `xml:"Event" 		json:"Event"`
	MsgId        int64   `xml:"MsgId"        json:"MsgId"`        // request
	Content      string  `xml:"Content"      json:"Content"`      // request
	MediaId      string  `xml:"MediaId"      json:"MediaId"`      // request
	PicURL       string  `xml:"PicUrl"       json:"PicUrl"`       // request
	Format       string  `xml:"Format"       json:"Format"`       // request
	Recognition  string  `xml:"Recognition"  json:"Recognition"`  // request
	ThumbMediaId string  `xml:"ThumbMediaId" json:"ThumbMediaId"` // request
	LocationX    float64 `xml:"Location_X"   json:"Location_X"`   // request
	LocationY    float64 `xml:"Location_Y"   json:"Location_Y"`   // request
	Scale        int     `xml:"Scale"        json:"Scale"`        // request
	Label        string  `xml:"Label"        json:"Label"`        // request
	Title        string  `xml:"Title"        json:"Title"`        // request
	Description  string  `xml:"Description"  json:"Description"`  // request
	URL          string  `xml:"Url"          json:"Url"`          // request
	EventKey     string  `xml:"EventKey"     json:"EventKey"`     // request, menu
	Ticket       string  `xml:"Ticket"       json:"Ticket"`       // request
	Latitude     float64 `xml:"Latitude"     json:"Latitude"`     // request
	Longitude    float64 `xml:"Longitude"    json:"Longitude"`    // request
	Precision    float64 `xml:"Precision"    json:"Precision"`    // request

	// menu
	MenuId int64 `xml:"MenuId" json:"MenuId"`
}

//自定义微信公众号菜单栏接口
type Menu struct {
	Buttons   []Button   `json:"button,omitempty"`
	MatchRule *MatchRule `json:"matchrule,omitempty"`
	MenuId    int64      `json:"menuid,omitempty"` // 有个性化菜单时查询接口返回值包含这个字段
}

type MatchRule struct {
	GroupId            string `json:"group_id,omitempty"`
	Sex                string `json:"sex,omitempty"`
	Country            string `json:"country,omitempty"`
	Province           string `json:"province,omitempty"`
	City               string `json:"city,omitempty"`
	ClientPlatformType string `json:"client_platform_type,omitempty"`
	Language           string `json:"language,omitempty"`
	TagId              string `json:"tag_id,omitempty"`
}

type Button struct {
	Type       string   `json:"type,omitempty"`       // 非必须; 菜单的响应动作类型
	Name       string   `json:"name,omitempty"`       // 必须;  菜单标题
	Key        string   `json:"key,omitempty"`        // 非必须; 菜单KEY值, 用于消息接口推送
	URL        string   `json:"url,omitempty"`        // 非必须; 网页链接, 用户点击菜单可打开链接
	MediaId    string   `json:"media_id,omitempty"`   // 非必须; 调用新增永久素材接口返回的合法media_id
	AppId      string   `json:"appid,omitempty"`      // 非必须; 跳转到小程序的appid
	PagePath   string   `json:"pagepath,omitempty"`   // 非必须; 跳转到小程序的path
	SubButtons []Button `json:"sub_button,omitempty"` // 非必须; 二级菜单数组
}

//自定义微信公众号菜单栏的回传信息
type MenuResp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

//微信公众号自定义菜单接口
func CreateMenu(c *gin.Context) {
	code := e.SUCCESS

	menu := Menu{
		Buttons: []Button{
			{
				Type: "view",
				Name: "万息手机端",
				URL:  "https://ones.fagougou.com/base/wx/redirect",
			},
		},
	}

	accessToken, err := FetchAccessToken("wx21e84ec720ccf278", "ffb8358ca7cc576e351141860be2e185", "https://api.weixin.qq.com/cgi-bin/token")

	url := "https://api.weixin.qq.com/cgi-bin/menu/create?access_token=" + accessToken

	var bodyJson []byte
	bodyJson, err = json.Marshal(menu)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println(bodyJson)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))

	if err != nil {
		log.Error(err)
		return
	}
	req.Header.Set("Content-type", "application/json")

	//http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}

	var result MenuResp
	if err := json.Unmarshal(respbody, &result); err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"errcode": result.Errcode,
		"errmsg":  result.Errmsg,
	})
}

func RedirectUrl(c *gin.Context) {
	//微信公众号网页授权链接，用来做万息后台手机端的跳转
	url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx21e84ec720ccf278&redirect_uri=https%3A%2F%2Fones.fagougou.com%2Fmobilelogin&response_type=code&scope=snsapi_userinfo&state=123&connect_redirect=1#wechat_redirect"
	c.Redirect(http.StatusMovedPermanently, url)
}

//微信公众号服务器配置 需启用和停用
func ServeHTTP(c *gin.Context) {
	const token = "huashengtoken"
	signature := c.Query("signature")

	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	ok := CheckSignature(signature, timestamp, nonce, token)
	if !ok {
		log.Println("微信公众号接入校验失败!")
		return
	}

	log.Println("微信公众号接入校验成功!")
	_, _ = c.Writer.WriteString(echostr)

}

//生成带参数二维码
func GetCode(c *gin.Context) {
	code := e.SUCCESS
	//向微信服务器获取权限code
	accessToken, err := FetchAccessToken("wx266f49a71875db29", "9f0a1a3de5efa1f38610130ed1421ab3", "https://api.weixin.qq.com/cgi-bin/token")

	if err != nil {
		fmt.Println("向微信服务器发送获取accessToken的get请求失败", err)
		return
	}

	//自定义二维码参数
	a := "aaa"

	url := "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=" + accessToken

	var request struct {
		ActionName string `json:"action_name"`
		ActionInfo struct {
			Scene struct {
				SceneStr string `json:"scene_str"`
			} `json:"scene"`
		} `json:"action_info"`
	}

	request.ActionName = "QR_LIMIT_STR_SCENE"
	request.ActionInfo.Scene.SceneStr = a

	var bodyJson []byte
	//将数据编码成json字符串
	bodyJson, err = json.Marshal(request)
	fmt.Println(bodyJson)
	if err != nil {
		code = e.ERROR_TOJSON_FAIL
		log.Error(err)
		return
	}

	//发送请求参数二维码所需格式数据
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Error(err)
		return
	}
	req.Header.Set("Content-type", "application/json")

	//http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		code = e.ERROR_AUTH_TOKEN
		log.Error(err)
		return
	}
	defer resp.Body.Close()
	//读取传过来的值
	respbody, _ := ioutil.ReadAll(resp.Body)

	//解析传过来的值
	var result PermQrcode
	if err := json.Unmarshal(respbody, &result); err != nil {
		code = e.ERROR_JSON_FAIL
		return
	}

	//拿到ticket
	ticket := result.Ticket
	//进行二维码链接的拼接
	url = "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + ticket
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": map[string]interface{}{
			"link": url,
		},
	})
}