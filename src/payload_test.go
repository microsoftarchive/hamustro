package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/wunderlist/hamustro/src/payload"
)

// Returns a Collection
func GetTestPayloadCollection(userId uint32, numberOfPayloads int) *payload.Collection {
	var payloads []*payload.Payload
	env := payload.Environment_DEVELOPMENT
	parameter := payload.Parameter{
		Name:  proto.String("parameter"),
		Value: proto.String("test_parameter")}
	for i := 0; i < numberOfPayloads; i++ {
		payloads = append(payloads, &payload.Payload{
			At:         proto.Uint64(1454681104),
			Event:      proto.String("Client.CreateUser"),
			Nr:         proto.Uint32(1),
			UserId:     proto.String(fmt.Sprint(userId + uint32(i*100))),
			Ip:         proto.String("214.160.227.22"),
			Parameters: []*payload.Parameter{&parameter}})
	}
	c := payload.Collection{
		DeviceId:       proto.String("a73b1c37-2c24-4786-af7a-16de88fbe23a"),
		ClientId:       proto.String("bce44f67b2661fd445d469b525b04f68"),
		Session:        proto.String("244f056dee6d475ec673ea0d20b69bab"),
		SystemVersion:  proto.String("10.10"),
		ProductVersion: proto.String("1.1.2"),
		System:         proto.String("OSX"),
		ProductGitHash: proto.String("5416a5889392d509e3bafcf40f6388e83aab23e6"),
		Env:            &env,
		Payloads:       payloads}
	return &c
}
