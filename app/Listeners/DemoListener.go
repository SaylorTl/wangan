package Listeners

import (
	"fmt"
	"log"
	"wangxin2.0/databases"
)

type ListenseQueue struct {
	databases.QueueInit
}

func (b ListenseQueue) Execute(payload *databases.QueuePayload) *databases.QueueResult {

	fmt.Printf("liscense:%v\n", payload)
	log.Println("DemoDemoQueue:", payload)
	//测试panic
	if payload.Body == "" {
		panic("empty body")
	}
	fmt.Printf("liscense:%v\n", payload.Body)
	return databases.NewQueueResult(true, "DemoDemoQueue.ok", nil)
}
