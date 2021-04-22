package weixin

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"go_gin_wechat/pkg/e"
	"go_gin_wechat/wechat"
	"io/ioutil"
	"net/http"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id" pg:",pk, notnull, type:uuid"`
	CreatedAt time.Time `json:"created_at" pg:"default:now()"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"-" pg:",soft_delete"`
	Phone     string    `json:"phone" pg:"type:varchar(50),unique"`
	Password  string    `json:"-" pg:"type:varchar(255)"`
	Email     string    `json:"email" pg:"type:varchar(255),unique"`
	Nickname  string    `json:"nickname" pg:"type:varchar(50)"`
	CompanyID string    `json:"company_id" pg:"type:varchar(100)"`
	Username  string    `json:"username" pg:"type:varchar(50),unique"`
	Role      string    `json:"role" pg:"type:varchar(100)"`
	AvatarUrl string    `json:"avatar_url" pg:"type:varchar(255)"`
	Openid    string    `json:"openid" pg:"type:varchar(255)"`
}

//微信公众号自定义菜单结构体
type Menu struct {
	Buttons   []Button   `json:"button,omitempty"`
	MatchRule *MatchRule `json:"matchrule,omitempty"`
	MenuId    int64      `json:"menuid,omitempty"` // 有个性化菜单时查询接口返回值包含这个字段
}

//微信公众号自定义菜单结构体
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

//微信公众号自定义菜单结构体
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

type MenuResp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

//微信自定义菜单生成方法
func CreateMenu(c *gin.Context) {

	menu := Menu{
		Buttons: []Button{
			{
				Type: "view",
				Name: "万息手机端",
				URL:  "https://beta-ones.fagougou.com/base/wx/redirect",
			},
		},
	}

	accessToken, err := wechat.FetchAccessToken("aaa", "bbb", "https://api.weixin.qq.com/cgi-bin/token")

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
}

//微信公众号服务器配置 需启用和停用
func ServeHTTP(c *gin.Context) {
	const token = "aaa"
	//获取参数 signature
	signature := c.Query("signature")
	//获取参数 timestamp
	timestamp := c.Query("timestamp")
	//获取参数 nonce
	nonce := c.Query("nonce")
	//获取参数 echostr
	echostr := c.Query("echostr")

	//微信公众号签名检查
	ok := wechat.CheckSignature(signature, timestamp, nonce, token)
	if !ok {
		log.Println("微信公众号接入校验失败!")
		return
	}

	log.Println("微信公众号接入校验成功!")
	_, _ = c.Writer.WriteString(echostr)

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

// WXMsgReceive 微信消息接收
func WXMsgReceive(c *gin.Context) {
	currentTime := time.Now().Local()
	newFormat := currentTime.Format("2006-01-02 15:04")
	errcode := e.SUCCESS

	//解析微信服务器发过来的 xml 格式数据
	var textMsg WXTextMsg
	err := c.ShouldBindXML(&textMsg)

	if err != nil {
		log.Printf("[消息接收] - XML数据包解析失败: %v\n", err)
		return
	}

	openid := textMsg.FromUserName
	//声明律师变量，来进行律师属性的获取
	var wxuser User
	var repo *database.Repo

	if len(textMsg.EventKey) > 0 && textMsg.Event == "subscribe" {
		//扫码关注公众号事件
		user := &User{
			Openid: openid,
		}
		const prefix = "qrscene_"
		scene := textMsg.EventKey[len(prefix):]
		user.ID, _ = uuid.FromString(scene)

		err := repo.SelectById(scene, &wxuser)

		var userRepo *UserRepo

		_, err = userRepo.PUpdateNotZero(user)

		if err != nil {
			errcode = e.INVALID_PARAMS
			c.JSON(http.StatusOK, gin.H{
				"code": errcode,
				"msg":  e.GetMsg(errcode),
			})
			return
		}

		err = wechat.Bind(openid, newFormat, wxuser.Phone)
		if err != nil {
			log.Error(err.Error())
		}

	} else if len(textMsg.EventKey) > 0 && textMsg.Event == "SCAN" {
		//已关注公众号再扫码事件
		user := &User{
			Openid: openid,
		}
		scene := textMsg.EventKey
		user.ID, _ = uuid.FromString(scene)

		err = repo.SelectById(scene, &wxuser)

		var userRepo *UserRepo
		_, err = userRepo.PUpdateNotZero(user)
		//对接收的消息进行被动回复
		WXMsgReply(c, textMsg.ToUserName, textMsg.FromUserName, newFormat)

		if err != nil {
			errcode = e.INVALID_PARAMS
			c.JSON(http.StatusOK, gin.H{
				"code": errcode,
				"msg":  e.GetMsg(errcode),
			})
			return
		}

		err = wechat.Bind(user.Openid, newFormat, wxuser.Phone)
		if err != nil {
			log.Error(err.Error())
		}
	}

}

// WXRepTextMsg 微信回复文本消息结构体
type WXRepTextMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   string
	MsgType      string
	Content      string
	// 若不标记XMLName, 则解析后的xml名为该结构体的名称
	XMLName xml.Name `xml:"xml"`
}

// WXMsgReply 微信消息回复
func WXMsgReply(c *gin.Context, fromUser, toUser string, time string) {
	repTextMsg := WXRepTextMsg{
		ToUserName:   toUser,
		FromUserName: fromUser,
		CreateTime:   time,
		MsgType:      "text",
		Content:      "你好，你的微信号已经绑定万息账号。",
	}

	msg, err := xml.Marshal(&repTextMsg)
	if err != nil {
		log.Printf("[消息回复] - 将对象进行XML编码出错: %v\n", err)
		return
	}
	_, err = c.Writer.Write(msg)
	return
}

//微信返回的 ticket + url
type PermQrcode struct {
	Ticket string `json:"ticket"`
	URL    string `json:"url"`
}

//生成带参数二维码（永久+字符串）
func GetCode(c *gin.Context) {
	code := e.SUCCESS
	//向微信服务器获取权限code
	accessToken, err := wechat.FetchAccessToken("aaa", "aaa", "https://api.weixin.qq.com/cgi-bin/token")

	if err != nil {
		fmt.Println("向微信服务器发送获取accessToken的get请求失败", err)
		return
	}

	session := sessions.Default(c)
	userID := session.Get("userID").(string)

	url := "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=" + accessToken
	//定义请求结构消息体
	var request struct {
		ActionName string `json:"action_name"`
		ActionInfo struct {
			Scene struct {
				SceneStr string `json:"scene_str"`
			} `json:"scene"`
		} `json:"action_info"`
	}

	request.ActionName = "QR_LIMIT_STR_SCENE"
	request.ActionInfo.Scene.SceneStr = userID

	var bodyJson []byte
	bodyJson, err = json.Marshal(request)
	fmt.Println(bodyJson)
	if err != nil {
		log.Error(err)
		return
	}
	//发送参数二维码所需格式数据
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

	//解析传过来的 ticket
	var result PermQrcode
	if err := json.Unmarshal(respbody, &result); err != nil {
		return
	}

	ticket := result.Ticket
	//进行二维码链接url的拼接
	url = "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + ticket
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": map[string]interface{}{
			"link": url,
		},
	})
}
