package core

import (
	"fmt"
	"reflect"
)

// Reply interface
type Reply interface {
	ReplyName() string
}

// SerializeReply serializes replies with a registered marshaller
func SerializeReply(v Reply) ([]byte, error) {
	return marshal(v.ReplyName(), v)
}

// DeserializeReply deserializes the reply data using a registered marshaller returning a *Reply
func DeserializeReply(replyName string, data []byte) (Reply, error) {
	reply, err := unmarshal(replyName, data)
	if err != nil {
		return nil, err
	}

	if reply != nil {
		if _, ok := reply.(Reply); !ok {
			return nil, fmt.Errorf("`%s` was registered but not registered as a reply", replyName)
		}
	}

	return reply.(Reply), nil
}

// RegisterReplies registers one or more replies with a registered marshaller
//
// Register replies using any form desired "&MyReply{}", "MyReply{}", "(*MyReply)(nil)"
//
// Replies must be registered after first registering a marshaller you wish to use
func RegisterReplies(replies ...Reply) {
	for _, reply := range replies {
		if v := reflect.ValueOf(reply); v.Kind() == reflect.Ptr && v.Pointer() == 0 {
			replyName := reflect.Zero(reflect.TypeOf(reply).Elem()).Interface().(Reply).ReplyName()
			registerType(replyName, reply)
		} else {
			registerType(reply.ReplyName(), reply)
		}
	}
}
