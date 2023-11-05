package Events

import (
	"log"
	"wangxin2.0/app/Listeners"
	"wangxin2.0/databases"
)

type EventQueue struct {
}

// 事件监听
func (e EventQueue) Init(str string) {
	qm := databases.QueueInit{}.InitEvent()
	// 注册队列任务处理器，TOPIC::GROUP 方式命名，和入栈队列payload一致
	err := qm.RegisterQueue("EVENT", "EVENT", Listeners.ListenseQueue{})
	check_err(err)
	// 注册队列任务处理器，TOPIC::GROUP 方式命名，和入栈队列payload一致
	qm.QueuePublish(&databases.QueuePayload{
		IsFast: true,
		Topic:  "EVENT",
		Group:  "EVENT",
		Body:   str,
	})
}

func check_err(err error) {
	if err != nil {
		log.Println("队列异常, err = ", err)
	}
}
