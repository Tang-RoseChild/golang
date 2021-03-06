// Code generated by protoc-gen-go.
// source: msg.proto
// DO NOT EDIT!

/*
Package msgproto is a generated protocol buffer package.

It is generated from these files:
	msg.proto

It has these top-level messages:
	Msg
*/
package msgproto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Msg struct {
	Id               *int32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Topic            *string `protobuf:"bytes,2,opt,name=topic" json:"topic,omitempty"`
	Content          *string `protobuf:"bytes,3,opt,name=content" json:"content,omitempty"`
	Type             *int32  `protobuf:"varint,4,req,name=type" json:"type,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Msg) Reset()         { *m = Msg{} }
func (m *Msg) String() string { return proto.CompactTextString(m) }
func (*Msg) ProtoMessage()    {}

func (m *Msg) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Msg) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

func (m *Msg) GetContent() string {
	if m != nil && m.Content != nil {
		return *m.Content
	}
	return ""
}

func (m *Msg) GetType() int32 {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return 0
}
