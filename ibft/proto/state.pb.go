// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/bloxapp/ssv/ibft/proto/state.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type State struct {
	Stage RoundState `protobuf:"varint,1,opt,name=stage,proto3,enum=proto.RoundState" json:"stage,omitempty"`
	// lambda is an instance unique identifier, much like a block hash in a blockchain
	Lambda []byte `protobuf:"bytes,2,opt,name=lambda,proto3" json:"lambda,omitempty"`
	// sequence number is an incremental number for each instance, much like a block number would be in a blockchain
	SeqNumber            uint64   `protobuf:"varint,3,opt,name=seq_number,json=seqNumber,proto3" json:"seq_number,omitempty"`
	PreviousLambda       []byte   `protobuf:"bytes,4,opt,name=previous_lambda,json=previousLambda,proto3" json:"previous_lambda,omitempty"`
	InputValue           []byte   `protobuf:"bytes,5,opt,name=input_value,json=inputValue,proto3" json:"input_value,omitempty"`
	Round                uint64   `protobuf:"varint,6,opt,name=round,proto3" json:"round,omitempty"`
	PreparedRound        uint64   `protobuf:"varint,7,opt,name=prepared_round,json=preparedRound,proto3" json:"prepared_round,omitempty"`
	PreparedValue        []byte   `protobuf:"bytes,8,opt,name=prepared_value,json=preparedValue,proto3" json:"prepared_value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *State) Reset()         { *m = State{} }
func (m *State) String() string { return proto.CompactTextString(m) }
func (*State) ProtoMessage()    {}
func (*State) Descriptor() ([]byte, []int) {
	return fileDescriptor_f98bd33b792b1786, []int{0}
}

func (m *State) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_State.Unmarshal(m, b)
}
func (m *State) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_State.Marshal(b, m, deterministic)
}
func (m *State) XXX_Merge(src proto.Message) {
	xxx_messageInfo_State.Merge(m, src)
}
func (m *State) XXX_Size() int {
	return xxx_messageInfo_State.Size(m)
}
func (m *State) XXX_DiscardUnknown() {
	xxx_messageInfo_State.DiscardUnknown(m)
}

var xxx_messageInfo_State proto.InternalMessageInfo

func (m *State) GetStage() RoundState {
	if m != nil {
		return m.Stage
	}
	return RoundState_NotStarted
}

func (m *State) GetLambda() []byte {
	if m != nil {
		return m.Lambda
	}
	return nil
}

func (m *State) GetSeqNumber() uint64 {
	if m != nil {
		return m.SeqNumber
	}
	return 0
}

func (m *State) GetPreviousLambda() []byte {
	if m != nil {
		return m.PreviousLambda
	}
	return nil
}

func (m *State) GetInputValue() []byte {
	if m != nil {
		return m.InputValue
	}
	return nil
}

func (m *State) GetRound() uint64 {
	if m != nil {
		return m.Round
	}
	return 0
}

func (m *State) GetPreparedRound() uint64 {
	if m != nil {
		return m.PreparedRound
	}
	return 0
}

func (m *State) GetPreparedValue() []byte {
	if m != nil {
		return m.PreparedValue
	}
	return nil
}

func init() {
	proto.RegisterType((*State)(nil), "proto.State")
}

func init() {
	proto.RegisterFile("github.com/bloxapp/ssv/ibft/proto/state.proto", fileDescriptor_f98bd33b792b1786)
}

var fileDescriptor_f98bd33b792b1786 = []byte{
	// 260 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x50, 0xcd, 0x4a, 0xc4, 0x30,
	0x10, 0xa6, 0xeb, 0xb6, 0xea, 0xa8, 0x2b, 0x06, 0x91, 0x20, 0x88, 0x55, 0x91, 0xed, 0x41, 0x1b,
	0xd0, 0x37, 0xf0, 0x2c, 0x1e, 0x22, 0x78, 0xf0, 0x52, 0x12, 0x1b, 0x6b, 0xa1, 0x6d, 0xb2, 0xf9,
	0x29, 0xbe, 0x88, 0xef, 0x2b, 0x9d, 0xec, 0x7a, 0xf1, 0xb0, 0xa7, 0xc9, 0x7c, 0xbf, 0x61, 0xe0,
	0xbe, 0x69, 0xfd, 0x57, 0x90, 0xe5, 0x87, 0xee, 0x99, 0xec, 0xf4, 0xb7, 0x30, 0x86, 0x39, 0x37,
	0xb2, 0x56, 0x7e, 0x7a, 0x66, 0xac, 0xf6, 0x9a, 0x39, 0x2f, 0xbc, 0x2a, 0xf1, 0x4d, 0x52, 0x1c,
	0xe7, 0x77, 0xdb, 0x5d, 0xbd, 0x6b, 0x5c, 0x34, 0x5d, 0xff, 0xcc, 0x20, 0x7d, 0x9d, 0x42, 0xc8,
	0x12, 0x52, 0xe7, 0x45, 0xa3, 0x68, 0x92, 0x27, 0xc5, 0xe2, 0xe1, 0x24, 0x0a, 0x4a, 0xae, 0xc3,
	0x50, 0xa3, 0x82, 0x47, 0x9e, 0x9c, 0x41, 0xd6, 0x89, 0x5e, 0xd6, 0x82, 0xce, 0xf2, 0xa4, 0x38,
	0xe4, 0xeb, 0x8d, 0x5c, 0x00, 0x38, 0xb5, 0xaa, 0x86, 0xd0, 0x4b, 0x65, 0xe9, 0x4e, 0x9e, 0x14,
	0x73, 0xbe, 0xef, 0xd4, 0xea, 0x05, 0x01, 0xb2, 0x84, 0x63, 0x63, 0xd5, 0xd8, 0xea, 0xe0, 0xaa,
	0xb5, 0x7f, 0x8e, 0xfe, 0xc5, 0x06, 0x7e, 0x8e, 0x39, 0x97, 0x70, 0xd0, 0x0e, 0x26, 0xf8, 0x6a,
	0x14, 0x5d, 0x50, 0x34, 0x45, 0x11, 0x20, 0xf4, 0x36, 0x21, 0xe4, 0x14, 0x52, 0x3b, 0xfd, 0x8a,
	0x66, 0xd8, 0x11, 0x17, 0x72, 0x0b, 0x53, 0x90, 0x11, 0x56, 0xd5, 0x55, 0xa4, 0x77, 0x91, 0x3e,
	0xda, 0xa0, 0xfc, 0x9f, 0x2c, 0x16, 0xec, 0x61, 0xc1, 0x9f, 0x0c, 0x3b, 0x9e, 0x6e, 0xde, 0xaf,
	0xb6, 0xde, 0x51, 0x66, 0x38, 0x1e, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0xf3, 0xa6, 0x14, 0x5c,
	0xa9, 0x01, 0x00, 0x00,
}
