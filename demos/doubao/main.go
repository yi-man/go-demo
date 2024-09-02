package main

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

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

	api_key := viper.GetString("ARK_API_KEY")
	model_id := viper.GetString("MODEL_ID")

	fmt.Print("api_key: ", api_key)
	fmt.Print("model_id: ", model_id)

	client := arkruntime.NewClientWithApiKey(
		api_key,
	)

	ctx := context.Background()

	fmt.Println("----- standard request -----")
	req := model.ChatCompletionRequest{
		Model: model_id,
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("你是豆包，是由字节跳动开发的 AI 人工智能助手"),
				},
			},
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("常见的十字花科植物有哪些？"),
				},
			},
		},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("standard chat error: %v\n", err)
		return
	}
	fmt.Println(*resp.Choices[0].Message.Content.StringValue)
}
