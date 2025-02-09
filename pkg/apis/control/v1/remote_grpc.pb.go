// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - ragu               v1.0.0
// source: github.com/rancher/opni/pkg/apis/control/v1/remote.proto

package v1

import (
	context "context"
	v1 "github.com/rancher/opni/pkg/apis/core/v1"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// HealthClient is the client API for Health service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HealthClient interface {
	GetHealth(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*v1.Health, error)
}

type healthClient struct {
	cc grpc.ClientConnInterface
}

func NewHealthClient(cc grpc.ClientConnInterface) HealthClient {
	return &healthClient{cc}
}

func (c *healthClient) GetHealth(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*v1.Health, error) {
	out := new(v1.Health)
	err := c.cc.Invoke(ctx, "/control.Health/GetHealth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HealthServer is the server API for Health service.
// All implementations must embed UnimplementedHealthServer
// for forward compatibility
type HealthServer interface {
	GetHealth(context.Context, *emptypb.Empty) (*v1.Health, error)
	mustEmbedUnimplementedHealthServer()
}

// UnimplementedHealthServer must be embedded to have forward compatible implementations.
type UnimplementedHealthServer struct {
}

func (UnimplementedHealthServer) GetHealth(context.Context, *emptypb.Empty) (*v1.Health, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHealth not implemented")
}
func (UnimplementedHealthServer) mustEmbedUnimplementedHealthServer() {}

// UnsafeHealthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HealthServer will
// result in compilation errors.
type UnsafeHealthServer interface {
	mustEmbedUnimplementedHealthServer()
}

func RegisterHealthServer(s grpc.ServiceRegistrar, srv HealthServer) {
	s.RegisterService(&Health_ServiceDesc, srv)
}

func _Health_GetHealth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthServer).GetHealth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/control.Health/GetHealth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthServer).GetHealth(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Health_ServiceDesc is the grpc.ServiceDesc for Health service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Health_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "control.Health",
	HandlerType: (*HealthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetHealth",
			Handler:    _Health_GetHealth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/rancher/opni/pkg/apis/control/v1/remote.proto",
}

// HealthListenerClient is the client API for HealthListener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HealthListenerClient interface {
	UpdateHealth(ctx context.Context, in *v1.Health, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type healthListenerClient struct {
	cc grpc.ClientConnInterface
}

func NewHealthListenerClient(cc grpc.ClientConnInterface) HealthListenerClient {
	return &healthListenerClient{cc}
}

func (c *healthListenerClient) UpdateHealth(ctx context.Context, in *v1.Health, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/control.HealthListener/UpdateHealth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HealthListenerServer is the server API for HealthListener service.
// All implementations must embed UnimplementedHealthListenerServer
// for forward compatibility
type HealthListenerServer interface {
	UpdateHealth(context.Context, *v1.Health) (*emptypb.Empty, error)
	mustEmbedUnimplementedHealthListenerServer()
}

// UnimplementedHealthListenerServer must be embedded to have forward compatible implementations.
type UnimplementedHealthListenerServer struct {
}

func (UnimplementedHealthListenerServer) UpdateHealth(context.Context, *v1.Health) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateHealth not implemented")
}
func (UnimplementedHealthListenerServer) mustEmbedUnimplementedHealthListenerServer() {}

// UnsafeHealthListenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HealthListenerServer will
// result in compilation errors.
type UnsafeHealthListenerServer interface {
	mustEmbedUnimplementedHealthListenerServer()
}

func RegisterHealthListenerServer(s grpc.ServiceRegistrar, srv HealthListenerServer) {
	s.RegisterService(&HealthListener_ServiceDesc, srv)
}

func _HealthListener_UpdateHealth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.Health)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthListenerServer).UpdateHealth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/control.HealthListener/UpdateHealth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthListenerServer).UpdateHealth(ctx, req.(*v1.Health))
	}
	return interceptor(ctx, in, info, handler)
}

// HealthListener_ServiceDesc is the grpc.ServiceDesc for HealthListener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var HealthListener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "control.HealthListener",
	HandlerType: (*HealthListenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateHealth",
			Handler:    _HealthListener_UpdateHealth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/rancher/opni/pkg/apis/control/v1/remote.proto",
}

// PluginManifestClient is the client API for PluginManifest service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PluginManifestClient interface {
	SendManifestsOrKnownPatch(ctx context.Context, in *ManifestMetadataList, opts ...grpc.CallOption) (*ManifestList, error)
	GetPluginManifests(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ManifestMetadataList, error)
	UploadPatch(ctx context.Context, in *PatchSpec, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type pluginManifestClient struct {
	cc grpc.ClientConnInterface
}

func NewPluginManifestClient(cc grpc.ClientConnInterface) PluginManifestClient {
	return &pluginManifestClient{cc}
}

func (c *pluginManifestClient) SendManifestsOrKnownPatch(ctx context.Context, in *ManifestMetadataList, opts ...grpc.CallOption) (*ManifestList, error) {
	out := new(ManifestList)
	err := c.cc.Invoke(ctx, "/control.PluginManifest/SendManifestsOrKnownPatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginManifestClient) GetPluginManifests(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ManifestMetadataList, error) {
	out := new(ManifestMetadataList)
	err := c.cc.Invoke(ctx, "/control.PluginManifest/GetPluginManifests", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginManifestClient) UploadPatch(ctx context.Context, in *PatchSpec, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/control.PluginManifest/UploadPatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PluginManifestServer is the server API for PluginManifest service.
// All implementations must embed UnimplementedPluginManifestServer
// for forward compatibility
type PluginManifestServer interface {
	SendManifestsOrKnownPatch(context.Context, *ManifestMetadataList) (*ManifestList, error)
	GetPluginManifests(context.Context, *emptypb.Empty) (*ManifestMetadataList, error)
	UploadPatch(context.Context, *PatchSpec) (*emptypb.Empty, error)
	mustEmbedUnimplementedPluginManifestServer()
}

// UnimplementedPluginManifestServer must be embedded to have forward compatible implementations.
type UnimplementedPluginManifestServer struct {
}

func (UnimplementedPluginManifestServer) SendManifestsOrKnownPatch(context.Context, *ManifestMetadataList) (*ManifestList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendManifestsOrKnownPatch not implemented")
}
func (UnimplementedPluginManifestServer) GetPluginManifests(context.Context, *emptypb.Empty) (*ManifestMetadataList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPluginManifests not implemented")
}
func (UnimplementedPluginManifestServer) UploadPatch(context.Context, *PatchSpec) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadPatch not implemented")
}
func (UnimplementedPluginManifestServer) mustEmbedUnimplementedPluginManifestServer() {}

// UnsafePluginManifestServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PluginManifestServer will
// result in compilation errors.
type UnsafePluginManifestServer interface {
	mustEmbedUnimplementedPluginManifestServer()
}

func RegisterPluginManifestServer(s grpc.ServiceRegistrar, srv PluginManifestServer) {
	s.RegisterService(&PluginManifest_ServiceDesc, srv)
}

func _PluginManifest_SendManifestsOrKnownPatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ManifestMetadataList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginManifestServer).SendManifestsOrKnownPatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/control.PluginManifest/SendManifestsOrKnownPatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginManifestServer).SendManifestsOrKnownPatch(ctx, req.(*ManifestMetadataList))
	}
	return interceptor(ctx, in, info, handler)
}

func _PluginManifest_GetPluginManifests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginManifestServer).GetPluginManifests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/control.PluginManifest/GetPluginManifests",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginManifestServer).GetPluginManifests(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PluginManifest_UploadPatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PatchSpec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginManifestServer).UploadPatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/control.PluginManifest/UploadPatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginManifestServer).UploadPatch(ctx, req.(*PatchSpec))
	}
	return interceptor(ctx, in, info, handler)
}

// PluginManifest_ServiceDesc is the grpc.ServiceDesc for PluginManifest service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PluginManifest_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "control.PluginManifest",
	HandlerType: (*PluginManifestServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendManifestsOrKnownPatch",
			Handler:    _PluginManifest_SendManifestsOrKnownPatch_Handler,
		},
		{
			MethodName: "GetPluginManifests",
			Handler:    _PluginManifest_GetPluginManifests_Handler,
		},
		{
			MethodName: "UploadPatch",
			Handler:    _PluginManifest_UploadPatch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/rancher/opni/pkg/apis/control/v1/remote.proto",
}
