package main

import (
	"context"
	"fmt"
	"time"

	horizonclient "github.com/stellar/go/exp/clients/horizon"
)

func main() {
	exampleClientStream()
}

func exampleClientStream() {
	// stream effects

	client := horizonclient.DefaultPublicNetClient
	effectRequest := horizonclient.EffectRequest{Cursor: "now"}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		// Stop streaming after 60 seconds.
		time.Sleep(60 * time.Second)
		cancel()
	}()

	// Valid handler 1
	ericHandler := func(e interface{}) {
		fmt.Println("Hello from handler func")
	}

	// Valid handler 2
	// usefulHandler := func(e interface{}) {
	// 	resp, ok := e.(effects.Base)
	// 	if ok {
	// 		fmt.Println(resp.Type)
	// 	}
	// }

	err := client.Stream(ctx, effectRequest, ericHandler)

	// Options:
	// 0) Generic stream (handler has interface{} signature)
	//		  err := client.Stream(ctx, effectRequest, ericHandler)

	// 1) Stream by type (handler has typed signature)
	// 			err := client.StreamEffects(ctx, effectRequest, ericHandler)

	// 2) Set handler on type (handler has typed signature)
	// 			effectRequest.Handler = ericHandler
	// 			err := client.Stream(ctx, effectRequest)

	// 3) Set handler on client (handler has typed signature)
	// 			client.EffectHandler = ericHandler
	// 			err := client.Stream(ctx, effectRequest)

	if err != nil {
		fmt.Println(err)
	}
}
