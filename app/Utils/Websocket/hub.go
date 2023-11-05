package Websocket

import (
	"log"
	"sync"
)

var h = hub{
	c: sync.Map{},
	u: make(chan *connection),
	b: make(chan Data),
	r: make(chan *connection),
}

type hub struct {
	c sync.Map
	b chan Data
	r chan *connection
	u chan *connection
}

func (h *hub) run() {
	defer func() {
		if err := recover(); err != nil {
			log.Print("hub广播出现异常", err)
		}
	}()
	for {
		select {
		case c := <-h.r:
			c.data.Message = "hi " + c.data.ClientId
			c.data.Type = "init"
			if _, ok := h.c.Load(c.data.ClientId); ok {
				c.sc <- c.data
			}
		case c := <-h.u:
			if _, ok := h.c.Load(c.data.ClientId); ok {
				close(c.sc)
				h.c.Delete(c.data.ClientId)
			}
		case data := <-h.b:
			if val, ok := h.c.Load(data.ClientId); ok {
				select {
				case val.(*connection).sc <- data:
				default:
					close(val.(*connection).sc)
					h.c.Delete(data.ClientId)
				}
			}
		}
	}
}
