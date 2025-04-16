// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: payment.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type FailureCode int32

const (
	FailureCode_GENERAL_ERROR        FailureCode = 0
	FailureCode_INSUFFICIENT_BALANCE FailureCode = 1
	FailureCode_MISSING_DATA         FailureCode = 2
)

// Enum value maps for FailureCode.
var (
	FailureCode_name = map[int32]string{
		0: "GENERAL_ERROR",
		1: "INSUFFICIENT_BALANCE",
		2: "MISSING_DATA",
	}
	FailureCode_value = map[string]int32{
		"GENERAL_ERROR":        0,
		"INSUFFICIENT_BALANCE": 1,
		"MISSING_DATA":         2,
	}
)

func (x FailureCode) Enum() *FailureCode {
	p := new(FailureCode)
	*p = x
	return p
}

func (x FailureCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FailureCode) Descriptor() protoreflect.EnumDescriptor {
	return file_payment_proto_enumTypes[0].Descriptor()
}

func (FailureCode) Type() protoreflect.EnumType {
	return &file_payment_proto_enumTypes[0]
}

func (x FailureCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FailureCode.Descriptor instead.
func (FailureCode) EnumDescriptor() ([]byte, []int) {
	return file_payment_proto_rawDescGZIP(), []int{0}
}

type Failure struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	FailureCode    FailureCode            `protobuf:"varint,1,opt,name=failureCode,proto3,enum=proto.FailureCode" json:"failureCode,omitempty"`
	FailureMessage string                 `protobuf:"bytes,2,opt,name=failureMessage,proto3" json:"failureMessage,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *Failure) Reset() {
	*x = Failure{}
	mi := &file_payment_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Failure) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Failure) ProtoMessage() {}

func (x *Failure) ProtoReflect() protoreflect.Message {
	mi := &file_payment_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Failure.ProtoReflect.Descriptor instead.
func (*Failure) Descriptor() ([]byte, []int) {
	return file_payment_proto_rawDescGZIP(), []int{0}
}

func (x *Failure) GetFailureCode() FailureCode {
	if x != nil {
		return x.FailureCode
	}
	return FailureCode_GENERAL_ERROR
}

func (x *Failure) GetFailureMessage() string {
	if x != nil {
		return x.FailureMessage
	}
	return ""
}

type BalanceReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BalanceReq) Reset() {
	*x = BalanceReq{}
	mi := &file_payment_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BalanceReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BalanceReq) ProtoMessage() {}

func (x *BalanceReq) ProtoReflect() protoreflect.Message {
	mi := &file_payment_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BalanceReq.ProtoReflect.Descriptor instead.
func (*BalanceReq) Descriptor() ([]byte, []int) {
	return file_payment_proto_rawDescGZIP(), []int{1}
}

func (x *BalanceReq) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type BalanceRes struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Result:
	//
	//	*BalanceRes_Success_
	//	*BalanceRes_Failure
	Result        isBalanceRes_Result `protobuf_oneof:"result"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BalanceRes) Reset() {
	*x = BalanceRes{}
	mi := &file_payment_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BalanceRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BalanceRes) ProtoMessage() {}

func (x *BalanceRes) ProtoReflect() protoreflect.Message {
	mi := &file_payment_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BalanceRes.ProtoReflect.Descriptor instead.
func (*BalanceRes) Descriptor() ([]byte, []int) {
	return file_payment_proto_rawDescGZIP(), []int{2}
}

func (x *BalanceRes) GetResult() isBalanceRes_Result {
	if x != nil {
		return x.Result
	}
	return nil
}

func (x *BalanceRes) GetSuccess() *BalanceRes_Success {
	if x != nil {
		if x, ok := x.Result.(*BalanceRes_Success_); ok {
			return x.Success
		}
	}
	return nil
}

func (x *BalanceRes) GetFailure() *Failure {
	if x != nil {
		if x, ok := x.Result.(*BalanceRes_Failure); ok {
			return x.Failure
		}
	}
	return nil
}

type isBalanceRes_Result interface {
	isBalanceRes_Result()
}

type BalanceRes_Success_ struct {
	Success *BalanceRes_Success `protobuf:"bytes,1,opt,name=success,proto3,oneof"`
}

type BalanceRes_Failure struct {
	Failure *Failure `protobuf:"bytes,2,opt,name=failure,proto3,oneof"`
}

func (*BalanceRes_Success_) isBalanceRes_Result() {}

func (*BalanceRes_Failure) isBalanceRes_Result() {}

type BalanceRes_Success struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Balance       float32                `protobuf:"fixed32,1,opt,name=Balance,proto3" json:"Balance,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BalanceRes_Success) Reset() {
	*x = BalanceRes_Success{}
	mi := &file_payment_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BalanceRes_Success) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BalanceRes_Success) ProtoMessage() {}

func (x *BalanceRes_Success) ProtoReflect() protoreflect.Message {
	mi := &file_payment_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BalanceRes_Success.ProtoReflect.Descriptor instead.
func (*BalanceRes_Success) Descriptor() ([]byte, []int) {
	return file_payment_proto_rawDescGZIP(), []int{2, 0}
}

func (x *BalanceRes_Success) GetBalance() float32 {
	if x != nil {
		return x.Balance
	}
	return 0
}

var File_payment_proto protoreflect.FileDescriptor

const file_payment_proto_rawDesc = "" +
	"\n" +
	"\rpayment.proto\x12\x05proto\"g\n" +
	"\aFailure\x124\n" +
	"\vfailureCode\x18\x01 \x01(\x0e2\x12.proto.FailureCodeR\vfailureCode\x12&\n" +
	"\x0efailureMessage\x18\x02 \x01(\tR\x0efailureMessage\"$\n" +
	"\n" +
	"BalanceReq\x12\x16\n" +
	"\x06UserId\x18\x01 \x01(\tR\x06UserId\"\x9e\x01\n" +
	"\n" +
	"BalanceRes\x125\n" +
	"\asuccess\x18\x01 \x01(\v2\x19.proto.BalanceRes.SuccessH\x00R\asuccess\x12*\n" +
	"\afailure\x18\x02 \x01(\v2\x0e.proto.FailureH\x00R\afailure\x1a#\n" +
	"\aSuccess\x12\x18\n" +
	"\aBalance\x18\x01 \x01(\x02R\aBalanceB\b\n" +
	"\x06result*L\n" +
	"\vFailureCode\x12\x11\n" +
	"\rGENERAL_ERROR\x10\x00\x12\x18\n" +
	"\x14INSUFFICIENT_BALANCE\x10\x01\x12\x10\n" +
	"\fMISSING_DATA\x10\x022J\n" +
	"\x0ePaymentService\x128\n" +
	"\x0eBalanceInquiry\x12\x11.proto.BalanceReq\x1a\x11.proto.BalanceRes\"\x00B\tZ\a/;protob\x06proto3"

var (
	file_payment_proto_rawDescOnce sync.Once
	file_payment_proto_rawDescData []byte
)

func file_payment_proto_rawDescGZIP() []byte {
	file_payment_proto_rawDescOnce.Do(func() {
		file_payment_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_payment_proto_rawDesc), len(file_payment_proto_rawDesc)))
	})
	return file_payment_proto_rawDescData
}

var file_payment_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_payment_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_payment_proto_goTypes = []any{
	(FailureCode)(0),           // 0: proto.FailureCode
	(*Failure)(nil),            // 1: proto.Failure
	(*BalanceReq)(nil),         // 2: proto.BalanceReq
	(*BalanceRes)(nil),         // 3: proto.BalanceRes
	(*BalanceRes_Success)(nil), // 4: proto.BalanceRes.Success
}
var file_payment_proto_depIdxs = []int32{
	0, // 0: proto.Failure.failureCode:type_name -> proto.FailureCode
	4, // 1: proto.BalanceRes.success:type_name -> proto.BalanceRes.Success
	1, // 2: proto.BalanceRes.failure:type_name -> proto.Failure
	2, // 3: proto.PaymentService.BalanceInquiry:input_type -> proto.BalanceReq
	3, // 4: proto.PaymentService.BalanceInquiry:output_type -> proto.BalanceRes
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_payment_proto_init() }
func file_payment_proto_init() {
	if File_payment_proto != nil {
		return
	}
	file_payment_proto_msgTypes[2].OneofWrappers = []any{
		(*BalanceRes_Success_)(nil),
		(*BalanceRes_Failure)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_payment_proto_rawDesc), len(file_payment_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_payment_proto_goTypes,
		DependencyIndexes: file_payment_proto_depIdxs,
		EnumInfos:         file_payment_proto_enumTypes,
		MessageInfos:      file_payment_proto_msgTypes,
	}.Build()
	File_payment_proto = out.File
	file_payment_proto_goTypes = nil
	file_payment_proto_depIdxs = nil
}
