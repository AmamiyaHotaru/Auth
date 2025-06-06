// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v4.25.7
// source: google_auth.proto

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

type MigrationPayload_Algorithm int32

const (
	MigrationPayload_ALGO_INVALID MigrationPayload_Algorithm = 0
	MigrationPayload_ALGO_SHA1    MigrationPayload_Algorithm = 1
)

// Enum value maps for MigrationPayload_Algorithm.
var (
	MigrationPayload_Algorithm_name = map[int32]string{
		0: "ALGO_INVALID",
		1: "ALGO_SHA1",
	}
	MigrationPayload_Algorithm_value = map[string]int32{
		"ALGO_INVALID": 0,
		"ALGO_SHA1":    1,
	}
)

func (x MigrationPayload_Algorithm) Enum() *MigrationPayload_Algorithm {
	p := new(MigrationPayload_Algorithm)
	*p = x
	return p
}

func (x MigrationPayload_Algorithm) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MigrationPayload_Algorithm) Descriptor() protoreflect.EnumDescriptor {
	return file_google_auth_proto_enumTypes[0].Descriptor()
}

func (MigrationPayload_Algorithm) Type() protoreflect.EnumType {
	return &file_google_auth_proto_enumTypes[0]
}

func (x MigrationPayload_Algorithm) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MigrationPayload_Algorithm.Descriptor instead.
func (MigrationPayload_Algorithm) EnumDescriptor() ([]byte, []int) {
	return file_google_auth_proto_rawDescGZIP(), []int{0, 0}
}

type MigrationPayload_OtpType int32

const (
	MigrationPayload_OTP_INVALID MigrationPayload_OtpType = 0
	MigrationPayload_OTP_HOTP    MigrationPayload_OtpType = 1
	MigrationPayload_OTP_TOTP    MigrationPayload_OtpType = 2
)

// Enum value maps for MigrationPayload_OtpType.
var (
	MigrationPayload_OtpType_name = map[int32]string{
		0: "OTP_INVALID",
		1: "OTP_HOTP",
		2: "OTP_TOTP",
	}
	MigrationPayload_OtpType_value = map[string]int32{
		"OTP_INVALID": 0,
		"OTP_HOTP":    1,
		"OTP_TOTP":    2,
	}
)

func (x MigrationPayload_OtpType) Enum() *MigrationPayload_OtpType {
	p := new(MigrationPayload_OtpType)
	*p = x
	return p
}

func (x MigrationPayload_OtpType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MigrationPayload_OtpType) Descriptor() protoreflect.EnumDescriptor {
	return file_google_auth_proto_enumTypes[1].Descriptor()
}

func (MigrationPayload_OtpType) Type() protoreflect.EnumType {
	return &file_google_auth_proto_enumTypes[1]
}

func (x MigrationPayload_OtpType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MigrationPayload_OtpType.Descriptor instead.
func (MigrationPayload_OtpType) EnumDescriptor() ([]byte, []int) {
	return file_google_auth_proto_rawDescGZIP(), []int{0, 1}
}

type MigrationPayload struct {
	state         protoimpl.MessageState            `protogen:"open.v1"`
	OtpParameters []*MigrationPayload_OtpParameters `protobuf:"bytes,1,rep,name=otp_parameters,json=otpParameters,proto3" json:"otp_parameters,omitempty"`
	Version       int32                             `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
	BatchSize     int32                             `protobuf:"varint,3,opt,name=batch_size,json=batchSize,proto3" json:"batch_size,omitempty"`
	BatchIndex    int32                             `protobuf:"varint,4,opt,name=batch_index,json=batchIndex,proto3" json:"batch_index,omitempty"`
	BatchId       int32                             `protobuf:"varint,5,opt,name=batch_id,json=batchId,proto3" json:"batch_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MigrationPayload) Reset() {
	*x = MigrationPayload{}
	mi := &file_google_auth_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MigrationPayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MigrationPayload) ProtoMessage() {}

func (x *MigrationPayload) ProtoReflect() protoreflect.Message {
	mi := &file_google_auth_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MigrationPayload.ProtoReflect.Descriptor instead.
func (*MigrationPayload) Descriptor() ([]byte, []int) {
	return file_google_auth_proto_rawDescGZIP(), []int{0}
}

func (x *MigrationPayload) GetOtpParameters() []*MigrationPayload_OtpParameters {
	if x != nil {
		return x.OtpParameters
	}
	return nil
}

func (x *MigrationPayload) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *MigrationPayload) GetBatchSize() int32 {
	if x != nil {
		return x.BatchSize
	}
	return 0
}

func (x *MigrationPayload) GetBatchIndex() int32 {
	if x != nil {
		return x.BatchIndex
	}
	return 0
}

func (x *MigrationPayload) GetBatchId() int32 {
	if x != nil {
		return x.BatchId
	}
	return 0
}

type MigrationPayload_OtpParameters struct {
	state         protoimpl.MessageState     `protogen:"open.v1"`
	Secret        []byte                     `protobuf:"bytes,1,opt,name=secret,proto3" json:"secret,omitempty"`
	Name          string                     `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Issuer        string                     `protobuf:"bytes,3,opt,name=issuer,proto3" json:"issuer,omitempty"`
	Algorithm     MigrationPayload_Algorithm `protobuf:"varint,4,opt,name=algorithm,proto3,enum=proto.MigrationPayload_Algorithm" json:"algorithm,omitempty"`
	Digits        int32                      `protobuf:"varint,5,opt,name=digits,proto3" json:"digits,omitempty"`
	Type          MigrationPayload_OtpType   `protobuf:"varint,6,opt,name=type,proto3,enum=proto.MigrationPayload_OtpType" json:"type,omitempty"`
	Counter       int64                      `protobuf:"varint,7,opt,name=counter,proto3" json:"counter,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MigrationPayload_OtpParameters) Reset() {
	*x = MigrationPayload_OtpParameters{}
	mi := &file_google_auth_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MigrationPayload_OtpParameters) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MigrationPayload_OtpParameters) ProtoMessage() {}

func (x *MigrationPayload_OtpParameters) ProtoReflect() protoreflect.Message {
	mi := &file_google_auth_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MigrationPayload_OtpParameters.ProtoReflect.Descriptor instead.
func (*MigrationPayload_OtpParameters) Descriptor() ([]byte, []int) {
	return file_google_auth_proto_rawDescGZIP(), []int{0, 0}
}

func (x *MigrationPayload_OtpParameters) GetSecret() []byte {
	if x != nil {
		return x.Secret
	}
	return nil
}

func (x *MigrationPayload_OtpParameters) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *MigrationPayload_OtpParameters) GetIssuer() string {
	if x != nil {
		return x.Issuer
	}
	return ""
}

func (x *MigrationPayload_OtpParameters) GetAlgorithm() MigrationPayload_Algorithm {
	if x != nil {
		return x.Algorithm
	}
	return MigrationPayload_ALGO_INVALID
}

func (x *MigrationPayload_OtpParameters) GetDigits() int32 {
	if x != nil {
		return x.Digits
	}
	return 0
}

func (x *MigrationPayload_OtpParameters) GetType() MigrationPayload_OtpType {
	if x != nil {
		return x.Type
	}
	return MigrationPayload_OTP_INVALID
}

func (x *MigrationPayload_OtpParameters) GetCounter() int64 {
	if x != nil {
		return x.Counter
	}
	return 0
}

var File_google_auth_proto protoreflect.FileDescriptor

const file_google_auth_proto_rawDesc = "" +
	"\n" +
	"\x11google_auth.proto\x12\x05proto\"\xb9\x04\n" +
	"\x10MigrationPayload\x12L\n" +
	"\x0eotp_parameters\x18\x01 \x03(\v2%.proto.MigrationPayload.OtpParametersR\rotpParameters\x12\x18\n" +
	"\aversion\x18\x02 \x01(\x05R\aversion\x12\x1d\n" +
	"\n" +
	"batch_size\x18\x03 \x01(\x05R\tbatchSize\x12\x1f\n" +
	"\vbatch_index\x18\x04 \x01(\x05R\n" +
	"batchIndex\x12\x19\n" +
	"\bbatch_id\x18\x05 \x01(\x05R\abatchId\x1a\xfb\x01\n" +
	"\rOtpParameters\x12\x16\n" +
	"\x06secret\x18\x01 \x01(\fR\x06secret\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x16\n" +
	"\x06issuer\x18\x03 \x01(\tR\x06issuer\x12?\n" +
	"\talgorithm\x18\x04 \x01(\x0e2!.proto.MigrationPayload.AlgorithmR\talgorithm\x12\x16\n" +
	"\x06digits\x18\x05 \x01(\x05R\x06digits\x123\n" +
	"\x04type\x18\x06 \x01(\x0e2\x1f.proto.MigrationPayload.OtpTypeR\x04type\x12\x18\n" +
	"\acounter\x18\a \x01(\x03R\acounter\",\n" +
	"\tAlgorithm\x12\x10\n" +
	"\fALGO_INVALID\x10\x00\x12\r\n" +
	"\tALGO_SHA1\x10\x01\"6\n" +
	"\aOtpType\x12\x0f\n" +
	"\vOTP_INVALID\x10\x00\x12\f\n" +
	"\bOTP_HOTP\x10\x01\x12\f\n" +
	"\bOTP_TOTP\x10\x02B\tZ\a./protob\x06proto3"

var (
	file_google_auth_proto_rawDescOnce sync.Once
	file_google_auth_proto_rawDescData []byte
)

func file_google_auth_proto_rawDescGZIP() []byte {
	file_google_auth_proto_rawDescOnce.Do(func() {
		file_google_auth_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_google_auth_proto_rawDesc), len(file_google_auth_proto_rawDesc)))
	})
	return file_google_auth_proto_rawDescData
}

var file_google_auth_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_google_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_google_auth_proto_goTypes = []any{
	(MigrationPayload_Algorithm)(0),        // 0: proto.MigrationPayload.Algorithm
	(MigrationPayload_OtpType)(0),          // 1: proto.MigrationPayload.OtpType
	(*MigrationPayload)(nil),               // 2: proto.MigrationPayload
	(*MigrationPayload_OtpParameters)(nil), // 3: proto.MigrationPayload.OtpParameters
}
var file_google_auth_proto_depIdxs = []int32{
	3, // 0: proto.MigrationPayload.otp_parameters:type_name -> proto.MigrationPayload.OtpParameters
	0, // 1: proto.MigrationPayload.OtpParameters.algorithm:type_name -> proto.MigrationPayload.Algorithm
	1, // 2: proto.MigrationPayload.OtpParameters.type:type_name -> proto.MigrationPayload.OtpType
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_google_auth_proto_init() }
func file_google_auth_proto_init() {
	if File_google_auth_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_google_auth_proto_rawDesc), len(file_google_auth_proto_rawDesc)),
			NumEnums:      2,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_google_auth_proto_goTypes,
		DependencyIndexes: file_google_auth_proto_depIdxs,
		EnumInfos:         file_google_auth_proto_enumTypes,
		MessageInfos:      file_google_auth_proto_msgTypes,
	}.Build()
	File_google_auth_proto = out.File
	file_google_auth_proto_goTypes = nil
	file_google_auth_proto_depIdxs = nil
}
