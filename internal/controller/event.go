package controller

type Event struct {
	EventName string
	Data      interface{}
}

func (event *Event) GetEventName() string {
	return event.EventName
}

func (event *Event) GetData() interface{} {
	return event.Data
}

const (
	EXIT            string = "exit"
	ADD_BACKEND     string = "add"
	CHANGE_STRATEGY string = "change"
)
