package Events

type DemoQueue struct {
	EventQueue
}

// 事件监听
func (d DemoQueue) Handle(data string) {
	d.Init(data)
}
