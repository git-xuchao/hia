package lab

import (
	"fmt"
)

type BasicEvent struct {
	EventId int
}

func (ev *BasicEvent) updateEventID(id int) {
	ev.EventId = id
}

func Lab3Command() {
	ev1 := &BasicEvent{EventId: 1}
	fmt.Printf("before update id = %d\n", ev1.EventId)
	ev1.updateEventID(2)
	fmt.Printf("After update id = %d\n", ev1.EventId)
	ev2 := BasicEvent{EventId: 1}
	fmt.Printf("before update id = %d\n", ev2.EventId)
	ev2.updateEventID(2)
	fmt.Printf("After update id = %d\n", ev2.EventId)
}
