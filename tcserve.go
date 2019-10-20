package main

import (
	"encoding/base64"
	"fmt"
	asr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/asr/v20190614"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

func getTcClient() *asr.Client {
	credential := common.NewCredential(
		"",
		"",
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "asr.tencentcloudapi.com"
	client, _ := asr.NewClient(credential, "", cpf)
	return client
}
func TcTrans(seg []byte, auKey string, client *asr.Client) {

	request := asr.NewSentenceRecognitionRequest()
	data := base64.StdEncoding.EncodeToString(seg)
	params := fmt.Sprintf(`{"ProjectId":1258311667,"SubServiceType":2,"EngSerViceType":"8k","SourceType":1,"VoiceFormat":"wav","UsrAudioKey":"%s","Data":"%s","DataLen":%d}`, auKey, data, len(seg))
	err := request.FromJsonString(params)
	if err != nil {
		fmt.Println("params parse error")
		panic(err)
	}
	response, err := client.SentenceRecognition(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", response.ToJsonString())
}
