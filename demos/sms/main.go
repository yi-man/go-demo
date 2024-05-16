package main

import (
	"fmt"
	"github.com/baidubce/bce-sdk-go/services/sms"
	"github.com/baidubce/bce-sdk-go/services/sms/api"
	"github.com/spf13/viper"
	"log"
)

type Data_Sms struct {
	AccessId    string `protobuf:"bytes,1,opt,name=access_id,json=accessId,proto3" json:"access_id,omitempty"`
	AccessKey   string `protobuf:"bytes,2,opt,name=access_key,json=accessKey,proto3" json:"access_key,omitempty"`
	SignatureId string `protobuf:"bytes,3,opt,name=signature_id,json=signatureId,proto3" json:"signature_id,omitempty"`
	Endpoint    string `protobuf:"bytes,4,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Port        string `protobuf:"bytes,5,opt,name=port,proto3" json:"port,omitempty"`
}

type SmsClient struct {
	clientWithPort *sms.Client
	client         *sms.Client
	Conf           *Data_Sms
}

func NewSmsClient(confData *Data_Sms) *SmsClient {
	// send api 的 client 没有端口号。本地环境所有的都不需要端口号
	client, sendErr := sms.NewClient(confData.AccessId, confData.AccessKey, confData.Endpoint)

	if sendErr != nil {
		panic(fmt.Sprintf("init sms send client error, %s", sendErr))
	}

	// 除了 send api，百度的其他api在线上环境需要端口号
	endpoint := confData.Endpoint
	if confData.Port != "" {
		endpoint = confData.Endpoint + ":" + confData.Port
	}
	clientWithPort, err := sms.NewClient(confData.AccessId, confData.AccessKey, endpoint)

	if err != nil {
		panic(fmt.Sprintf("init sms client error, %s", err))
	}

	client.Config.ConnectionTimeoutInMillis = 30 * 1000

	return &SmsClient{
		client:         client,
		clientWithPort: clientWithPort,
		Conf:           confData,
	}
}

type SendSmsParam struct {
	ContentVar map[string]interface{}
	Template   string
	Mobile     string
}

// work
func (sms *SmsClient) SendSms(param *SendSmsParam) (*api.SendSmsResult, error) {
	log.Println("send sms: start, params are ", param)

	sendSmsArgs := &api.SendSmsArgs{
		Mobile:      param.Mobile,
		SignatureId: sms.Conf.SignatureId,
		Template:    param.Template,
		ContentVar:  param.ContentVar,
	}

	var sendSmsResult *api.SendSmsResult
	var err error

	sendSmsResult, err = sms.client.SendSms(sendSmsArgs)

	if err != nil {
		log.Fatalln("send sms error: ", err)

		return nil, err
	}

	log.Println("send sms success, params are ", param)

	return sendSmsResult, nil
}

func (sms *SmsClient) GetSmsTemplate(templateId string) (*api.GetTemplateResult, error) {
	getTemplateArgs := &api.GetTemplateArgs{
		TemplateId: templateId,
	}
	getTemplateResult, err := sms.clientWithPort.GetTemplate(getTemplateArgs)
	if err != nil {
		log.Fatalf("获取模版详情失败, %s", err)
		return nil, err
	}

	log.Printf("获取模版详情成功. %s\n", getTemplateResult)

	return getTemplateResult, nil
}

func (sms *SmsClient) QueryQuotaAndRateLimit() (*api.QueryQuotaRateResult, error) {
	result, err := sms.clientWithPort.QueryQuotaAndRateLimit()
	if err != nil {
		log.Printf("query quota or rate limit error, %s", err)

		return nil, err
	}

	log.Printf("query quota or rate limit success. %v", result)

	return result, nil
}

func (sms *SmsClient) GetSignature(signatureId string) (*api.GetSignatureResult, error) {
	result, err := sms.client.GetSignature(&api.GetSignatureArgs{
		SignatureId: signatureId,
	})
	if err != nil {
		fmt.Printf("get signature error, %s", err)
		return nil, err
	}
	fmt.Printf("get signature success. %s", result)

	return result, nil
}

func (sms *SmsClient) ListStatistics() (*api.ListStatisticsResponse, error) {
	res, err := sms.client.ListStatistics(&api.ListStatisticsArgs{
		SmsType:   "CommonSale",
		StartTime: "2024-05-14",
		EndTime:   "2024-05-14",
	})
	if err != nil {
		fmt.Printf("GetMobileBlack error, %s", err)
		return nil, err
	}
	// fmt.Printf("GetMobileBlack success. %s", res)

	return res, nil
}

func (sms *SmsClient) GetMobileBlack() (*api.GetMobileBlackResult, error) {
	res, err := sms.client.GetMobileBlack(&api.GetMobileBlackArgs{
		SmsType:        "CommonNotice",
		SignatureIdStr: "sddd",
		Phone:          "12345678901",
		StartTime:      "2024-05-1",
		EndTime:        "2024-05-14",
		PageNo:         "1",
		PageSize:       "10",
	})
	if err != nil {
		fmt.Printf("GetMobileBlack error, %s", err)
		return nil, err
	}
	fmt.Printf("GetMobileBlack success. %v", res)
	return res, nil
}

func main() {
	// 设置配置文件名和路径
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	conf := &Data_Sms{
		AccessId:    viper.GetString("sms.access_id"),
		AccessKey:   viper.GetString("sms.access_key"),
		SignatureId: viper.GetString("sms.signature_id"),
		Endpoint:    viper.GetString("sms.endpoint"),
		Port:        viper.GetString("sms.port"),
	}

	smsClient := NewSmsClient(conf)
	sendSmsParam := &SendSmsParam{
		ContentVar: map[string]interface{}{
			"jobType": "docker22",
		},
		Template: "sms-tmpl-CvsFUg07712",
		Mobile:   "18601762076",
	}

	fmt.Println(sendSmsParam)

	// 发送短信
	// result, err := smsClient.SendSms(sendSmsParam)
	// 获取模版
	// result, err := smsClient.GetSmsTemplate(sendSmsParam.Template)

	// 查询额度  not work
	result, err := smsClient.QueryQuotaAndRateLimit()

	// 查询签名
	// result, err := smsClient.GetSignature(smsClient.Conf.SignatureId)

	// 查询黑名单
	// result, err := smsClient.GetMobileBlack()

	// 业务统计
	// result, err := smsClient.ListStatistics()

	if err != nil {
		log.Fatalln("send sms error: ", err)
	}
	fmt.Println(result)
}
