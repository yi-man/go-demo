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

// func OptionFunc(concurrentNum int32) xmq.ReaderOptions {
// 	return func(o *xmq.ReaderOptions) error {
// 		o.concurrentNum = 1
// 		return nil
// 	}
// }

// func WithGroupName(concurrentNum int32) xmq.ReaderOption {
// 	return func(o *xmq.ReaderOptions) error {
// 		o.concurrentNum = concurrentNum
// 		return nil
// 	}
// }

func main() {

	fmt.Println("basic info", host, port)

	client, err := xmq.NewClient(host, port)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	// func options(o *ReaderOptions) error{
	// 	return nil
	// }

	// optionFunc := func(opt *xmq.ReaderOptions) error {

	// 	// o.ack = true
	// 	// o.concurrentNum = 1
	// 	return nil
	// }

	subscriber, err := client.NewSubscriber(ctx, "topic-wade-test",
		xmq.WithGroupName("wade1"),
		xmq.WithConcurrenctNum(10),
		xmq.WithSequentialRead(false),
		xmq.WithAutoCommit(false),
		xmq.WithAck(true),
	)
	if err != nil {
		panic(err)
	}
	count := 1

	subscriber.Handle(ctx, func(ctx context.Context, m xmq.Message) error {
		body, err := xmq.MessageBodyUnmarshalJSON[MessageBody](m)
		if err != nil {
			return err
		}

		count += 1
		fmt.Println(count, body.Text)
		return nil
	}, xmq.WithMessageTypeFilter("xmq.type-1"))

	<-ctx.Done()

}
