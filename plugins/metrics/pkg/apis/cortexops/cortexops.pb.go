// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1-devel
// 	protoc        v1.0.0
// source: github.com/rancher/opni/plugins/metrics/pkg/apis/cortexops/cortexops.proto

package cortexops

import (
	v1 "github.com/rancher/opni/pkg/apis/storage/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type InstallState int32

const (
	InstallState_Unknown      InstallState = 0
	InstallState_NotInstalled InstallState = 1
	InstallState_Updating     InstallState = 2
	InstallState_Installed    InstallState = 3
	InstallState_Uninstalling InstallState = 4
)

// Enum value maps for InstallState.
var (
	InstallState_name = map[int32]string{
		0: "Unknown",
		1: "NotInstalled",
		2: "Updating",
		3: "Installed",
		4: "Uninstalling",
	}
	InstallState_value = map[string]int32{
		"Unknown":      0,
		"NotInstalled": 1,
		"Updating":     2,
		"Installed":    3,
		"Uninstalling": 4,
	}
)

func (x InstallState) Enum() *InstallState {
	p := new(InstallState)
	*p = x
	return p
}

func (x InstallState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (InstallState) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_enumTypes[0].Descriptor()
}

func (InstallState) Type() protoreflect.EnumType {
	return &file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_enumTypes[0]
}

func (x InstallState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use InstallState.Descriptor instead.
func (InstallState) EnumDescriptor() ([]byte, []int) {
	return file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_rawDescGZIP(), []int{0}
}

type DeploymentMode int32

const (
	DeploymentMode_AllInOne        DeploymentMode = 0
	DeploymentMode_HighlyAvailable DeploymentMode = 1
)

// Enum value maps for DeploymentMode.
var (
	DeploymentMode_name = map[int32]string{
		0: "AllInOne",
		1: "HighlyAvailable",
	}
	DeploymentMode_value = map[string]int32{
		"AllInOne":        0,
		"HighlyAvailable": 1,
	}
)

func (x DeploymentMode) Enum() *DeploymentMode {
	p := new(DeploymentMode)
	*p = x
	return p
}

func (x DeploymentMode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DeploymentMode) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_enumTypes[1].Descriptor()
}

func (DeploymentMode) Type() protoreflect.EnumType {
	return &file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_enumTypes[1]
}

func (x DeploymentMode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DeploymentMode.Descriptor instead.
func (DeploymentMode) EnumDescriptor() ([]byte, []int) {
	return file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_rawDescGZIP(), []int{1}
}

type InstallStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State    InstallState      `protobuf:"varint,1,opt,name=state,proto3,enum=cortexops.InstallState" json:"state,omitempty"`
	Version  string            `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Metadata map[string]string `protobuf:"bytes,3,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *InstallStatus) Reset() {
	*x = InstallStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstallStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstallStatus) ProtoMessage() {}

func (x *InstallStatus) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstallStatus.ProtoReflect.Descriptor instead.
func (*InstallStatus) Descriptor() ([]byte, []int) {
	return file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_rawDescGZIP(), []int{0}
}

func (x *InstallStatus) GetState() InstallState {
	if x != nil {
		return x.State
	}
	return InstallState_Unknown
}

func (x *InstallStatus) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *InstallStatus) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

type ClusterConfiguration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Mode    DeploymentMode  `protobuf:"varint,1,opt,name=mode,proto3,enum=cortexops.DeploymentMode" json:"mode,omitempty"`
	Storage *v1.StorageSpec `protobuf:"bytes,2,opt,name=storage,proto3" json:"storage,omitempty"`
}

func (x *ClusterConfiguration) Reset() {
	*x = ClusterConfiguration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClusterConfiguration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterConfiguration) ProtoMessage() {}

func (x *ClusterConfiguration) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterConfiguration.ProtoReflect.Descriptor instead.
func (*ClusterConfiguration) Descriptor() ([]byte, []int) {
	return file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_rawDescGZIP(), []int{1}
}

func (x *ClusterConfiguration) GetMode() DeploymentMode {
	if x != nil {
		return x.Mode
	}
	return DeploymentMode_AllInOne
}

func (x *ClusterConfiguration) GetStorage() *v1.StorageSpec {
	if x != nil {
		return x.Storage
	}
	return nil
}

type AgentConfiguration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PrometheusImage string `protobuf:"bytes,1,opt,name=prometheusImage,proto3" json:"prometheusImage,omitempty"`
}

func (x *AgentConfiguration) Reset() {
	*x = AgentConfiguration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AgentConfiguration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgentConfiguration) ProtoMessage() {}

func (x *AgentConfiguration) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AgentConfiguration.ProtoReflect.Descriptor instead.
func (*AgentConfiguration) Descriptor() ([]byte, []int) {
	return file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_rawDescGZIP(), []int{2}
}

func (x *AgentConfiguration) GetPrometheusImage() string {
	if x != nil {
		return x.PrometheusImage
	}
	return ""
}

var File_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto protoreflect.FileDescriptor

var file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_rawDesc = []byte{
	0x0a, 0x4a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x61, 0x6e,
	0x63, 0x68, 0x65, 0x72, 0x2f, 0x6f, 0x70, 0x6e, 0x69, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x73, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70,
	0x69, 0x73, 0x2f, 0x63, 0x6f, 0x72, 0x74, 0x65, 0x78, 0x6f, 0x70, 0x73, 0x2f, 0x63, 0x6f, 0x72,
	0x74, 0x65, 0x78, 0x6f, 0x70, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x63, 0x6f,
	0x72, 0x74, 0x65, 0x78, 0x6f, 0x70, 0x73, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x39, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x65, 0x72, 0x2f, 0x6f, 0x70, 0x6e, 0x69, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2f, 0x76,
	0x31, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xd9, 0x01, 0x0a, 0x0d, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x2d, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x17, 0x2e, 0x63, 0x6f, 0x72, 0x74, 0x65, 0x78, 0x6f, 0x70, 0x73, 0x2e, 0x49, 0x6e, 0x73,
	0x74, 0x61, 0x6c, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x42, 0x0a, 0x08, 0x6d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x63,
	0x6f, 0x72, 0x74, 0x65, 0x78, 0x6f, 0x70, 0x73, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x3b,
	0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x75, 0x0a, 0x14, 0x43,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x2d, 0x0a, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x19, 0x2e, 0x63, 0x6f, 0x72, 0x74, 0x65, 0x78, 0x6f, 0x70, 0x73, 0x2e, 0x44, 0x65,
	0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x6d, 0x6f,
	0x64, 0x65, 0x12, 0x2e, 0x0a, 0x07, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x53, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x53, 0x70, 0x65, 0x63, 0x52, 0x07, 0x73, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x22, 0x3e, 0x0a, 0x12, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x0a, 0x0f, 0x70, 0x72, 0x6f, 0x6d,
	0x65, 0x74, 0x68, 0x65, 0x75, 0x73, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0f, 0x70, 0x72, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x65, 0x75, 0x73, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x2a, 0x5c, 0x0a, 0x0c, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12,
	0x10, 0x0a, 0x0c, 0x4e, 0x6f, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x65, 0x64, 0x10,
	0x01, 0x12, 0x0c, 0x0a, 0x08, 0x55, 0x70, 0x64, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x10, 0x02, 0x12,
	0x0d, 0x0a, 0x09, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x65, 0x64, 0x10, 0x03, 0x12, 0x10,
	0x0a, 0x0c, 0x55, 0x6e, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x10, 0x04,
	0x2a, 0x33, 0x0a, 0x0e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x6f,
	0x64, 0x65, 0x12, 0x0c, 0x0a, 0x08, 0x41, 0x6c, 0x6c, 0x49, 0x6e, 0x4f, 0x6e, 0x65, 0x10, 0x00,
	0x12, 0x13, 0x0a, 0x0f, 0x48, 0x69, 0x67, 0x68, 0x6c, 0x79, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61,
	0x62, 0x6c, 0x65, 0x10, 0x01, 0x32, 0xe2, 0x01, 0x0a, 0x09, 0x43, 0x6f, 0x72, 0x74, 0x65, 0x78,
	0x4f, 0x70, 0x73, 0x12, 0x4b, 0x0a, 0x10, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65,
	0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x1f, 0x2e, 0x63, 0x6f, 0x72, 0x74, 0x65, 0x78,
	0x6f, 0x70, 0x73, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x12, 0x44, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x18, 0x2e, 0x63,
	0x6f, 0x72, 0x74, 0x65, 0x78, 0x6f, 0x70, 0x73, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x42, 0x0a, 0x10, 0x55, 0x6e, 0x69, 0x6e, 0x73, 0x74,
	0x61, 0x6c, 0x6c, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0xd9, 0x01, 0x0a, 0x08, 0x41,
	0x67, 0x65, 0x6e, 0x74, 0x4f, 0x70, 0x73, 0x12, 0x47, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x75, 0x72, 0x65, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x12, 0x1d, 0x2e, 0x63, 0x6f, 0x72, 0x74,
	0x65, 0x78, 0x6f, 0x70, 0x73, 0x2e, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x12, 0x42, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x18, 0x2e, 0x63, 0x6f, 0x72,
	0x74, 0x65, 0x78, 0x6f, 0x70, 0x73, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x40, 0x0a, 0x0e, 0x55, 0x6e, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6c,
	0x6c, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x3c, 0x5a, 0x3a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x65, 0x72, 0x2f, 0x6f, 0x70, 0x6e,
	0x69, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x73, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x63, 0x6f, 0x72, 0x74, 0x65,
	0x78, 0x6f, 0x70, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_rawDescOnce sync.Once
	file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_rawDescData = file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_rawDesc
)

func file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_rawDescGZIP() []byte {
	file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_rawDescOnce.Do(func() {
		file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_rawDescData)
	})
	return file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_rawDescData
}

var file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_goTypes = []interface{}{
	(InstallState)(0),            // 0: cortexops.InstallState
	(DeploymentMode)(0),          // 1: cortexops.DeploymentMode
	(*InstallStatus)(nil),        // 2: cortexops.InstallStatus
	(*ClusterConfiguration)(nil), // 3: cortexops.ClusterConfiguration
	(*AgentConfiguration)(nil),   // 4: cortexops.AgentConfiguration
	nil,                          // 5: cortexops.InstallStatus.MetadataEntry
	(*v1.StorageSpec)(nil),       // 6: storage.StorageSpec
	(*emptypb.Empty)(nil),        // 7: google.protobuf.Empty
}
var file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_depIdxs = []int32{
	0,  // 0: cortexops.InstallStatus.state:type_name -> cortexops.InstallState
	5,  // 1: cortexops.InstallStatus.metadata:type_name -> cortexops.InstallStatus.MetadataEntry
	1,  // 2: cortexops.ClusterConfiguration.mode:type_name -> cortexops.DeploymentMode
	6,  // 3: cortexops.ClusterConfiguration.storage:type_name -> storage.StorageSpec
	3,  // 4: cortexops.CortexOps.ConfigureCluster:input_type -> cortexops.ClusterConfiguration
	7,  // 5: cortexops.CortexOps.GetClusterStatus:input_type -> google.protobuf.Empty
	7,  // 6: cortexops.CortexOps.UninstallCluster:input_type -> google.protobuf.Empty
	4,  // 7: cortexops.AgentOps.ConfigureAgent:input_type -> cortexops.AgentConfiguration
	7,  // 8: cortexops.AgentOps.GetAgentStatus:input_type -> google.protobuf.Empty
	7,  // 9: cortexops.AgentOps.UninstallAgent:input_type -> google.protobuf.Empty
	7,  // 10: cortexops.CortexOps.ConfigureCluster:output_type -> google.protobuf.Empty
	2,  // 11: cortexops.CortexOps.GetClusterStatus:output_type -> cortexops.InstallStatus
	7,  // 12: cortexops.CortexOps.UninstallCluster:output_type -> google.protobuf.Empty
	7,  // 13: cortexops.AgentOps.ConfigureAgent:output_type -> google.protobuf.Empty
	2,  // 14: cortexops.AgentOps.GetAgentStatus:output_type -> cortexops.InstallStatus
	7,  // 15: cortexops.AgentOps.UninstallAgent:output_type -> google.protobuf.Empty
	10, // [10:16] is the sub-list for method output_type
	4,  // [4:10] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_init() }
func file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_init() {
	if File_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstallStatus); i {
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
		file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClusterConfiguration); i {
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
		file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AgentConfiguration); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_goTypes,
		DependencyIndexes: file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_depIdxs,
		EnumInfos:         file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_enumTypes,
		MessageInfos:      file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_msgTypes,
	}.Build()
	File_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto = out.File
	file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_rawDesc = nil
	file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_goTypes = nil
	file_github_com_rancher_opni_plugins_metrics_pkg_apis_cortexops_cortexops_proto_depIdxs = nil
}
