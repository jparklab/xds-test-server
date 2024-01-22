// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.6.1
// source: proto/config.proto

package config

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CommandType int32

const (
	CommandType_COMMAND_ADD_SERVICE    CommandType = 0
	CommandType_COMMAND_REMOVE_SERVICE CommandType = 1
	CommandType_COMMAND_SET_OPTION     CommandType = 2
)

// Enum value maps for CommandType.
var (
	CommandType_name = map[int32]string{
		0: "COMMAND_ADD_SERVICE",
		1: "COMMAND_REMOVE_SERVICE",
		2: "COMMAND_SET_OPTION",
	}
	CommandType_value = map[string]int32{
		"COMMAND_ADD_SERVICE":    0,
		"COMMAND_REMOVE_SERVICE": 1,
		"COMMAND_SET_OPTION":     2,
	}
)

func (x CommandType) Enum() *CommandType {
	p := new(CommandType)
	*p = x
	return p
}

func (x CommandType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CommandType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_config_proto_enumTypes[0].Descriptor()
}

func (CommandType) Type() protoreflect.EnumType {
	return &file_proto_config_proto_enumTypes[0]
}

func (x CommandType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CommandType.Descriptor instead.
func (CommandType) EnumDescriptor() ([]byte, []int) {
	return file_proto_config_proto_rawDescGZIP(), []int{0}
}

type AddServiceCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServiceName string `protobuf:"bytes,1,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
}

func (x *AddServiceCommand) Reset() {
	*x = AddServiceCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddServiceCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddServiceCommand) ProtoMessage() {}

func (x *AddServiceCommand) ProtoReflect() protoreflect.Message {
	mi := &file_proto_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddServiceCommand.ProtoReflect.Descriptor instead.
func (*AddServiceCommand) Descriptor() ([]byte, []int) {
	return file_proto_config_proto_rawDescGZIP(), []int{0}
}

func (x *AddServiceCommand) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

type RemoveServiceCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServiceName string `protobuf:"bytes,1,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
}

func (x *RemoveServiceCommand) Reset() {
	*x = RemoveServiceCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveServiceCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveServiceCommand) ProtoMessage() {}

func (x *RemoveServiceCommand) ProtoReflect() protoreflect.Message {
	mi := &file_proto_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveServiceCommand.ProtoReflect.Descriptor instead.
func (*RemoveServiceCommand) Descriptor() ([]byte, []int) {
	return file_proto_config_proto_rawDescGZIP(), []int{1}
}

func (x *RemoveServiceCommand) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

type SetOptionCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// Types that are assignable to Value:
	//
	//	*SetOptionCommand_StringValue
	//	*SetOptionCommand_Int64Value
	//	*SetOptionCommand_BoolValue
	Value isSetOptionCommand_Value `protobuf_oneof:"value"`
}

func (x *SetOptionCommand) Reset() {
	*x = SetOptionCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_config_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetOptionCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetOptionCommand) ProtoMessage() {}

func (x *SetOptionCommand) ProtoReflect() protoreflect.Message {
	mi := &file_proto_config_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetOptionCommand.ProtoReflect.Descriptor instead.
func (*SetOptionCommand) Descriptor() ([]byte, []int) {
	return file_proto_config_proto_rawDescGZIP(), []int{2}
}

func (x *SetOptionCommand) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (m *SetOptionCommand) GetValue() isSetOptionCommand_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *SetOptionCommand) GetStringValue() string {
	if x, ok := x.GetValue().(*SetOptionCommand_StringValue); ok {
		return x.StringValue
	}
	return ""
}

func (x *SetOptionCommand) GetInt64Value() int64 {
	if x, ok := x.GetValue().(*SetOptionCommand_Int64Value); ok {
		return x.Int64Value
	}
	return 0
}

func (x *SetOptionCommand) GetBoolValue() bool {
	if x, ok := x.GetValue().(*SetOptionCommand_BoolValue); ok {
		return x.BoolValue
	}
	return false
}

type isSetOptionCommand_Value interface {
	isSetOptionCommand_Value()
}

type SetOptionCommand_StringValue struct {
	StringValue string `protobuf:"bytes,2,opt,name=string_value,json=stringValue,proto3,oneof"`
}

type SetOptionCommand_Int64Value struct {
	Int64Value int64 `protobuf:"varint,3,opt,name=int64_value,json=int64Value,proto3,oneof"`
}

type SetOptionCommand_BoolValue struct {
	BoolValue bool `protobuf:"varint,4,opt,name=bool_value,json=boolValue,proto3,oneof"`
}

func (*SetOptionCommand_StringValue) isSetOptionCommand_Value() {}

func (*SetOptionCommand_Int64Value) isSetOptionCommand_Value() {}

func (*SetOptionCommand_BoolValue) isSetOptionCommand_Value() {}

type Command struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type CommandType `protobuf:"varint,1,opt,name=type,proto3,enum=xdsserver.CommandType" json:"type,omitempty"`
	// Types that are assignable to Command:
	//
	//	*Command_AddService
	//	*Command_RemoveService
	//	*Command_SetOption
	Command isCommand_Command `protobuf_oneof:"command"`
}

func (x *Command) Reset() {
	*x = Command{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_config_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Command) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Command) ProtoMessage() {}

func (x *Command) ProtoReflect() protoreflect.Message {
	mi := &file_proto_config_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Command.ProtoReflect.Descriptor instead.
func (*Command) Descriptor() ([]byte, []int) {
	return file_proto_config_proto_rawDescGZIP(), []int{3}
}

func (x *Command) GetType() CommandType {
	if x != nil {
		return x.Type
	}
	return CommandType_COMMAND_ADD_SERVICE
}

func (m *Command) GetCommand() isCommand_Command {
	if m != nil {
		return m.Command
	}
	return nil
}

func (x *Command) GetAddService() *AddServiceCommand {
	if x, ok := x.GetCommand().(*Command_AddService); ok {
		return x.AddService
	}
	return nil
}

func (x *Command) GetRemoveService() *RemoveServiceCommand {
	if x, ok := x.GetCommand().(*Command_RemoveService); ok {
		return x.RemoveService
	}
	return nil
}

func (x *Command) GetSetOption() *SetOptionCommand {
	if x, ok := x.GetCommand().(*Command_SetOption); ok {
		return x.SetOption
	}
	return nil
}

type isCommand_Command interface {
	isCommand_Command()
}

type Command_AddService struct {
	AddService *AddServiceCommand `protobuf:"bytes,2,opt,name=add_service,json=addService,proto3,oneof"`
}

type Command_RemoveService struct {
	RemoveService *RemoveServiceCommand `protobuf:"bytes,3,opt,name=remove_service,json=removeService,proto3,oneof"`
}

type Command_SetOption struct {
	SetOption *SetOptionCommand `protobuf:"bytes,4,opt,name=set_option,json=setOption,proto3,oneof"`
}

func (*Command_AddService) isCommand_Command() {}

func (*Command_RemoveService) isCommand_Command() {}

func (*Command_SetOption) isCommand_Command() {}

var File_proto_config_proto protoreflect.FileDescriptor

var file_proto_config_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x78, 0x64, 0x73, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x36, 0x0a, 0x11,
	0x41, 0x64, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x4e, 0x61, 0x6d, 0x65, 0x22, 0x39, 0x0a, 0x14, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x21, 0x0a, 0x0c,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22,
	0x96, 0x01, 0x0a, 0x10, 0x53, 0x65, 0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x23, 0x0a, 0x0c, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0b,
	0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x21, 0x0a, 0x0b, 0x69,
	0x6e, 0x74, 0x36, 0x34, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x48, 0x00, 0x52, 0x0a, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1f,
	0x0a, 0x0a, 0x62, 0x6f, 0x6f, 0x6c, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x08, 0x48, 0x00, 0x52, 0x09, 0x62, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x42,
	0x07, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x89, 0x02, 0x0a, 0x07, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x12, 0x2a, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x16, 0x2e, 0x78, 0x64, 0x73, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x43,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x3f, 0x0a, 0x0b, 0x61, 0x64, 0x64, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x78, 0x64, 0x73, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2e, 0x41, 0x64, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x48, 0x00, 0x52, 0x0a, 0x61, 0x64, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x48, 0x0a, 0x0e, 0x72, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x78, 0x64, 0x73, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x48, 0x00, 0x52, 0x0d, 0x72, 0x65,
	0x6d, 0x6f, 0x76, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3c, 0x0a, 0x0a, 0x73,
	0x65, 0x74, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1b, 0x2e, 0x78, 0x64, 0x73, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x53, 0x65, 0x74, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x48, 0x00, 0x52, 0x09,
	0x73, 0x65, 0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x09, 0x0a, 0x07, 0x63, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x2a, 0x5a, 0x0a, 0x0b, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x17, 0x0a, 0x13, 0x43, 0x4f, 0x4d, 0x4d, 0x41, 0x4e, 0x44, 0x5f, 0x41,
	0x44, 0x44, 0x5f, 0x53, 0x45, 0x52, 0x56, 0x49, 0x43, 0x45, 0x10, 0x00, 0x12, 0x1a, 0x0a, 0x16,
	0x43, 0x4f, 0x4d, 0x4d, 0x41, 0x4e, 0x44, 0x5f, 0x52, 0x45, 0x4d, 0x4f, 0x56, 0x45, 0x5f, 0x53,
	0x45, 0x52, 0x56, 0x49, 0x43, 0x45, 0x10, 0x01, 0x12, 0x16, 0x0a, 0x12, 0x43, 0x4f, 0x4d, 0x4d,
	0x41, 0x4e, 0x44, 0x5f, 0x53, 0x45, 0x54, 0x5f, 0x4f, 0x50, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x02,
	0x32, 0x52, 0x0a, 0x0d, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x41, 0x0a, 0x11, 0x53, 0x65, 0x6e, 0x64, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x43,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x12, 0x2e, 0x78, 0x64, 0x73, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x22, 0x00, 0x42, 0x3a, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x6a, 0x70, 0x61, 0x72, 0x6b, 0x6c, 0x61, 0x62, 0x2f, 0x78, 0x64, 0x73, 0x2d,
	0x74, 0x65, 0x73, 0x74, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_config_proto_rawDescOnce sync.Once
	file_proto_config_proto_rawDescData = file_proto_config_proto_rawDesc
)

func file_proto_config_proto_rawDescGZIP() []byte {
	file_proto_config_proto_rawDescOnce.Do(func() {
		file_proto_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_config_proto_rawDescData)
	})
	return file_proto_config_proto_rawDescData
}

var file_proto_config_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_config_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_config_proto_goTypes = []interface{}{
	(CommandType)(0),             // 0: xdsserver.CommandType
	(*AddServiceCommand)(nil),    // 1: xdsserver.AddServiceCommand
	(*RemoveServiceCommand)(nil), // 2: xdsserver.RemoveServiceCommand
	(*SetOptionCommand)(nil),     // 3: xdsserver.SetOptionCommand
	(*Command)(nil),              // 4: xdsserver.Command
	(*empty.Empty)(nil),          // 5: google.protobuf.Empty
}
var file_proto_config_proto_depIdxs = []int32{
	0, // 0: xdsserver.Command.type:type_name -> xdsserver.CommandType
	1, // 1: xdsserver.Command.add_service:type_name -> xdsserver.AddServiceCommand
	2, // 2: xdsserver.Command.remove_service:type_name -> xdsserver.RemoveServiceCommand
	3, // 3: xdsserver.Command.set_option:type_name -> xdsserver.SetOptionCommand
	4, // 4: xdsserver.ConfigService.SendConfigCommand:input_type -> xdsserver.Command
	5, // 5: xdsserver.ConfigService.SendConfigCommand:output_type -> google.protobuf.Empty
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proto_config_proto_init() }
func file_proto_config_proto_init() {
	if File_proto_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddServiceCommand); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveServiceCommand); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_config_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetOptionCommand); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_config_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Command); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_proto_config_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*SetOptionCommand_StringValue)(nil),
		(*SetOptionCommand_Int64Value)(nil),
		(*SetOptionCommand_BoolValue)(nil),
	}
	file_proto_config_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*Command_AddService)(nil),
		(*Command_RemoveService)(nil),
		(*Command_SetOption)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_config_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_config_proto_goTypes,
		DependencyIndexes: file_proto_config_proto_depIdxs,
		EnumInfos:         file_proto_config_proto_enumTypes,
		MessageInfos:      file_proto_config_proto_msgTypes,
	}.Build()
	File_proto_config_proto = out.File
	file_proto_config_proto_rawDesc = nil
	file_proto_config_proto_goTypes = nil
	file_proto_config_proto_depIdxs = nil
}
