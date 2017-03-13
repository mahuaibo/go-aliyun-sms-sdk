package client

import (
	"fmt"
	"test/util"
	"net/http"
	"net/url"
	"time"
)

//常量
const (
	ApiSMSVersion = "2016-09-27"
	ApiSMSHost    = "https://sms.aliyuncs.com/"

	DefaultSignatureVersion = "1.0"
	DefaultSignatureMethod  = "HMAC-SHA1"
	JSONResponseFormat      = "JSON"
	XMLResponseFormat       = "XML"
	ECSRequestMethod        = "POST"
)

//请求接口
//所有请求对象继承的接口，也是Client接受处理的请求接口
//签名方式和必要参数信息。
type Request interface {
	//签名
	Sign(*Credentials)
	//返回*http.Request
	HttpRequestInstance() (*http.Request, error)
	//返回值的类型，支持JSON与XML.
	ResponseFormat() string
	//
	String() string
	//
	Clone() interface{}
	////返回请求处理超时限制时长
	//DeadLine() time.Duration
}

//  SMS 请求对象。实现 Request 接口
type SMSRequest struct {
	Format          	string
	Version         	string
	AccessKeyId     	string
	Signature       	string
	SignatureMethod 	string
	Timestamp       	util.ISO6801Time
	SignatureVersion	string
	SignatureNonce   	string
	RegionId 		string

	Action 			string
	SignName 		string
	TemplateCode 		string
	RecNum 			string
	ParamString 		string

	// http
	Host   string
	Method string
	Url    string
	Args   url.Values
}

// SMSRequest的必要字段转成参数
func (Sms *SMSRequest) StructToArgs() {
	Sms.SignatureNonce = util.CreateRandomString()
	Sms.Args.Set("Format", Sms.Format)
	Sms.Args.Set("Version", Sms.Version)
	Sms.Args.Set("AccessKeyId", Sms.AccessKeyId)
	Sms.Args.Set("SignatureMethod", Sms.SignatureMethod)
	Sms.Args.Set("Timestamp", Sms.Timestamp.String())
	Sms.Args.Set("SignatureVersion", Sms.SignatureVersion)
	Sms.Args.Set("SignatureNonce", Sms.SignatureNonce)

	Sms.Args.Set("Action", Sms.Action)
	Sms.Args.Set("SignName",Sms.SignName)
	Sms.Args.Set("TemplateCode",Sms.TemplateCode)
	Sms.Args.Set("RecNum",Sms.RecNum)
	Sms.Args.Set("ParamString",Sms.ParamString)
}

// 签名
func (Sms *SMSRequest) Sign(cert *Credentials) {
	Sms.AccessKeyId = cert.AccessKeyId
	Sms.StructToArgs()
	// 生成签名
	Sms.Signature = util.CreateSignatureForRequest(Sms.Method, &Sms.Args, cert.AccessKeySecret+"&")

}

func (Sms *SMSRequest) HttpRequestInstance() (httpReq *http.Request, err error) {
	// 生成请求url
	Sms.Url = Sms.Host + "?" + Sms.Args.Encode() + "&Signature=" + url.QueryEscape(Sms.Signature)
	httpReq, err = http.NewRequest(Sms.Method, Sms.Url, nil)
	return
}

func (Sms *SMSRequest) ResponseFormat() string {
	return Sms.Format
}

// A Timeout of zero means no timeout.
func (Sms *SMSRequest) DeadLine() time.Duration {
	return 0
}

func (Sms *SMSRequest) String() string {
	return fmt.Sprintf("Method:%s,Url:%s", Sms.Method, Sms.Url)
}

// 克隆
func (l *SMSRequest) Clone() interface{} {
	new_obj := (*l)
	//清空数据
	new_obj.Args = url.Values{}
	return &new_obj
}

func (Sms *SMSRequest) SetArgs(key, value string) {
	Sms.Args.Set(key, value)
}

func (Sms *SMSRequest) DelArgs(key string) {
	Sms.Args.Del(key)
}

// 生成SMSRequest
func NewSMSRequest(action, signName, templateCode, recNum, paramString string) *SMSRequest {
	return &SMSRequest{
		Format:			JSONResponseFormat,
		Version:		ApiSMSVersion,
		SignatureNonce:		util.CreateRandomString(),
		SignatureMethod:	DefaultSignatureMethod,
		SignatureVersion:	DefaultSignatureVersion,
		Timestamp:		util.NewISO6801Time(time.Now().UTC()),

		Action:			action,
		SignName:		signName,
		TemplateCode:		templateCode,
		RecNum:			recNum,
		ParamString:		paramString,

		Host:			ApiSMSHost,
		Method:			ECSRequestMethod,
		Args:			url.Values{},
	}
}