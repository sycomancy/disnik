package main

import (
	"time"

	"github.com/sycomancy/glasnik/types"
)

func main() {

	requestProcessor := NewRequestProcessor(":3333", "/result-webhook", "http://localhost:3000/api/request-njuska")
	go sendRequest(requestProcessor)
	requestProcessor.Run()
}

func sendRequest(rp *RequestProcessor) {
	time.Sleep(time.Second * 5)
	err := rp.SendRequest(&types.Request{
		Filter:      "https://www.njuskalo.hr/prodaja-stanova?geo%5BlocationIds%5D=2698",
		Token:       "121345",
		CallbackURL: "http://localhost:3333/result-webhook",
	})
	if err != nil {
		panic(err)
	}
}
