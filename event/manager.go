package event

import "github.com/csimplestring/go-json-schema/schema"

type EventArg struct {
	Path    string
	ErrCode schema.ErrorCode
}

type EventHandler interface {
	Handle(arg EventArg)
}

type EventManager struct {
}

func (em *EventManager) Dispatch(eventName string, arg EventArg) {

}
