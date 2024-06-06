package main

import (
	"log"

	"github.com/hibiken/asynq"
)

const redisAddr = "127.0.0.1:6379"

const (
	EmailQueue = "notification_email"
	ImageQueue = "notification_image_resize"
)

func main() {
	queues := make(map[string]int)
	queues[EmailQueue] = 4
	queues[ImageQueue] = 2

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: queues,
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmailDelivery, HandleEmailDeliveryTask)
	mux.Handle(TypeImageResize, NewImageProcessor())
	// ...register other handlers...

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
