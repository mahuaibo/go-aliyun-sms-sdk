package client

import (
	"testing"
	"fmt"
)

const (
	conf_accessKey_id = "AccessKeyId"
	conf_accessKey_secret = "AccessKeySecret"
)

func TestClient(t *testing.T) {
	//1
	client := NewClient(&Credentials{conf_accessKey_id, conf_accessKey_secret})
	client.SetDebug(true)
	//2
	paramString := ` { "companyName" : "XXXX" , "code" : "666666" } `;
	request := NewSMSRequest("SingleSendSms", "XXXX", "SMS_35020166", "189****5460", paramString)
	response := &ErrorResponse{}
	//3
	err := client.Query(request,response)
	if err == nil {
		fmt.Println("发送成功")
	}else{
		fmt.Println(err)
	}
}