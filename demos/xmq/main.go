package main

import (
	"context"
	"fmt"
	"git.basebit.me/xobj/xmq-client-go/v2"
)

var host string
var port string

type MessageBody struct {
	Text string `json:"text"`
}

func init() {
	host = "172.18.18.82"
	port = "32701"
}

func main() {

	fmt.Println("basic info", host, port)

	client, err := xmq.NewClient(host, port)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	writer, err := client.NewWriter(ctx, "topic-wade-test")
	if err != nil {
		panic(err)
	}
	body := &MessageBody{Text: "test body"}
	body1 := &MessageBody{Text: "test body1"}

	err = writer.Write(
		ctx,
		xmq.NewMessageBuilder().MustSetBodyWithMarshalJSON(body).Build(),
	)
	err = writer.Write(
		ctx,
		xmq.NewMessageBuilder().SetType("type-1").MustSetBodyWithMarshalJSON(body1).Build(),
	)
	if err != nil {
		panic(err)
	}
}
