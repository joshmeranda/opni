// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - ragu               v1.0.0
// source: github.com/rancher/opni/plugins/logging/pkg/apis/loggingadmin/loggingadmin.proto

package loggingadmin

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// LoggingAdminClient is the client API for LoggingAdmin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LoggingAdminClient interface {
	GetOpensearchCluster(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*OpensearchCluster, error)
	DeleteOpensearchCluster(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	CreateOrUpdateOpensearchCluster(ctx context.Context, in *OpensearchCluster, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UpgradeAvailable(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*UpgradeAvailableResponse, error)
	DoUpgrade(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetStorageClasses(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StorageClassResponse, error)
	GetOpensearchStatus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StatusResponse, error)
}

type loggingAdminClient struct {
	cc grpc.ClientConnInterface
}

func NewLoggingAdminClient(cc grpc.ClientConnInterface) LoggingAdminClient {
	return &loggingAdminClient{cc}
}

func (c *loggingAdminClient) GetOpensearchCluster(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*OpensearchCluster, error) {
	out := new(OpensearchCluster)
	err := c.cc.Invoke(ctx, "/loggingadmin.LoggingAdmin/GetOpensearchCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loggingAdminClient) DeleteOpensearchCluster(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/loggingadmin.LoggingAdmin/DeleteOpensearchCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loggingAdminClient) CreateOrUpdateOpensearchCluster(ctx context.Context, in *OpensearchCluster, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/loggingadmin.LoggingAdmin/CreateOrUpdateOpensearchCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loggingAdminClient) UpgradeAvailable(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*UpgradeAvailableResponse, error) {
	out := new(UpgradeAvailableResponse)
	err := c.cc.Invoke(ctx, "/loggingadmin.LoggingAdmin/UpgradeAvailable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loggingAdminClient) DoUpgrade(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/loggingadmin.LoggingAdmin/DoUpgrade", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loggingAdminClient) GetStorageClasses(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StorageClassResponse, error) {
	out := new(StorageClassResponse)
	err := c.cc.Invoke(ctx, "/loggingadmin.LoggingAdmin/GetStorageClasses", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loggingAdminClient) GetOpensearchStatus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/loggingadmin.LoggingAdmin/GetOpensearchStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoggingAdminServer is the server API for LoggingAdmin service.
// All implementations must embed UnimplementedLoggingAdminServer
// for forward compatibility
type LoggingAdminServer interface {
	GetOpensearchCluster(context.Context, *emptypb.Empty) (*OpensearchCluster, error)
	DeleteOpensearchCluster(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	CreateOrUpdateOpensearchCluster(context.Context, *OpensearchCluster) (*emptypb.Empty, error)
	UpgradeAvailable(context.Context, *emptypb.Empty) (*UpgradeAvailableResponse, error)
	DoUpgrade(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	GetStorageClasses(context.Context, *emptypb.Empty) (*StorageClassResponse, error)
	GetOpensearchStatus(context.Context, *emptypb.Empty) (*StatusResponse, error)
	mustEmbedUnimplementedLoggingAdminServer()
}

// UnimplementedLoggingAdminServer must be embedded to have forward compatible implementations.
type UnimplementedLoggingAdminServer struct {
}

func (UnimplementedLoggingAdminServer) GetOpensearchCluster(context.Context, *emptypb.Empty) (*OpensearchCluster, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOpensearchCluster not implemented")
}
func (UnimplementedLoggingAdminServer) DeleteOpensearchCluster(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteOpensearchCluster not implemented")
}
func (UnimplementedLoggingAdminServer) CreateOrUpdateOpensearchCluster(context.Context, *OpensearchCluster) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrUpdateOpensearchCluster not implemented")
}
func (UnimplementedLoggingAdminServer) UpgradeAvailable(context.Context, *emptypb.Empty) (*UpgradeAvailableResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpgradeAvailable not implemented")
}
func (UnimplementedLoggingAdminServer) DoUpgrade(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DoUpgrade not implemented")
}
func (UnimplementedLoggingAdminServer) GetStorageClasses(context.Context, *emptypb.Empty) (*StorageClassResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStorageClasses not implemented")
}
func (UnimplementedLoggingAdminServer) GetOpensearchStatus(context.Context, *emptypb.Empty) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOpensearchStatus not implemented")
}
func (UnimplementedLoggingAdminServer) mustEmbedUnimplementedLoggingAdminServer() {}

// UnsafeLoggingAdminServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LoggingAdminServer will
// result in compilation errors.
type UnsafeLoggingAdminServer interface {
	mustEmbedUnimplementedLoggingAdminServer()
}

func RegisterLoggingAdminServer(s grpc.ServiceRegistrar, srv LoggingAdminServer) {
	s.RegisterService(&LoggingAdmin_ServiceDesc, srv)
}

func _LoggingAdmin_GetOpensearchCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggingAdminServer).GetOpensearchCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loggingadmin.LoggingAdmin/GetOpensearchCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggingAdminServer).GetOpensearchCluster(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoggingAdmin_DeleteOpensearchCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggingAdminServer).DeleteOpensearchCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loggingadmin.LoggingAdmin/DeleteOpensearchCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggingAdminServer).DeleteOpensearchCluster(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoggingAdmin_CreateOrUpdateOpensearchCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpensearchCluster)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggingAdminServer).CreateOrUpdateOpensearchCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loggingadmin.LoggingAdmin/CreateOrUpdateOpensearchCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggingAdminServer).CreateOrUpdateOpensearchCluster(ctx, req.(*OpensearchCluster))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoggingAdmin_UpgradeAvailable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggingAdminServer).UpgradeAvailable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loggingadmin.LoggingAdmin/UpgradeAvailable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggingAdminServer).UpgradeAvailable(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoggingAdmin_DoUpgrade_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggingAdminServer).DoUpgrade(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loggingadmin.LoggingAdmin/DoUpgrade",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggingAdminServer).DoUpgrade(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoggingAdmin_GetStorageClasses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggingAdminServer).GetStorageClasses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loggingadmin.LoggingAdmin/GetStorageClasses",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggingAdminServer).GetStorageClasses(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoggingAdmin_GetOpensearchStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggingAdminServer).GetOpensearchStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loggingadmin.LoggingAdmin/GetOpensearchStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggingAdminServer).GetOpensearchStatus(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// LoggingAdmin_ServiceDesc is the grpc.ServiceDesc for LoggingAdmin service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LoggingAdmin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "loggingadmin.LoggingAdmin",
	HandlerType: (*LoggingAdminServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOpensearchCluster",
			Handler:    _LoggingAdmin_GetOpensearchCluster_Handler,
		},
		{
			MethodName: "DeleteOpensearchCluster",
			Handler:    _LoggingAdmin_DeleteOpensearchCluster_Handler,
		},
		{
			MethodName: "CreateOrUpdateOpensearchCluster",
			Handler:    _LoggingAdmin_CreateOrUpdateOpensearchCluster_Handler,
		},
		{
			MethodName: "UpgradeAvailable",
			Handler:    _LoggingAdmin_UpgradeAvailable_Handler,
		},
		{
			MethodName: "DoUpgrade",
			Handler:    _LoggingAdmin_DoUpgrade_Handler,
		},
		{
			MethodName: "GetStorageClasses",
			Handler:    _LoggingAdmin_GetStorageClasses_Handler,
		},
		{
			MethodName: "GetOpensearchStatus",
			Handler:    _LoggingAdmin_GetOpensearchStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/rancher/opni/plugins/logging/pkg/apis/loggingadmin/loggingadmin.proto",
}

// LoggingAdminV2Client is the client API for LoggingAdminV2 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LoggingAdminV2Client interface {
	GetOpensearchCluster(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*OpensearchClusterV2, error)
	DeleteOpensearchCluster(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	CreateOrUpdateOpensearchCluster(ctx context.Context, in *OpensearchClusterV2, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UpgradeAvailable(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*UpgradeAvailableResponse, error)
	DoUpgrade(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetStorageClasses(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StorageClassResponse, error)
	GetOpensearchStatus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StatusResponse, error)
}

type loggingAdminV2Client struct {
	cc grpc.ClientConnInterface
}

func NewLoggingAdminV2Client(cc grpc.ClientConnInterface) LoggingAdminV2Client {
	return &loggingAdminV2Client{cc}
}

func (c *loggingAdminV2Client) GetOpensearchCluster(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*OpensearchClusterV2, error) {
	out := new(OpensearchClusterV2)
	err := c.cc.Invoke(ctx, "/loggingadmin.LoggingAdminV2/GetOpensearchCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loggingAdminV2Client) DeleteOpensearchCluster(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/loggingadmin.LoggingAdminV2/DeleteOpensearchCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loggingAdminV2Client) CreateOrUpdateOpensearchCluster(ctx context.Context, in *OpensearchClusterV2, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/loggingadmin.LoggingAdminV2/CreateOrUpdateOpensearchCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loggingAdminV2Client) UpgradeAvailable(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*UpgradeAvailableResponse, error) {
	out := new(UpgradeAvailableResponse)
	err := c.cc.Invoke(ctx, "/loggingadmin.LoggingAdminV2/UpgradeAvailable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loggingAdminV2Client) DoUpgrade(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/loggingadmin.LoggingAdminV2/DoUpgrade", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loggingAdminV2Client) GetStorageClasses(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StorageClassResponse, error) {
	out := new(StorageClassResponse)
	err := c.cc.Invoke(ctx, "/loggingadmin.LoggingAdminV2/GetStorageClasses", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loggingAdminV2Client) GetOpensearchStatus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/loggingadmin.LoggingAdminV2/GetOpensearchStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoggingAdminV2Server is the server API for LoggingAdminV2 service.
// All implementations must embed UnimplementedLoggingAdminV2Server
// for forward compatibility
type LoggingAdminV2Server interface {
	GetOpensearchCluster(context.Context, *emptypb.Empty) (*OpensearchClusterV2, error)
	DeleteOpensearchCluster(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	CreateOrUpdateOpensearchCluster(context.Context, *OpensearchClusterV2) (*emptypb.Empty, error)
	UpgradeAvailable(context.Context, *emptypb.Empty) (*UpgradeAvailableResponse, error)
	DoUpgrade(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	GetStorageClasses(context.Context, *emptypb.Empty) (*StorageClassResponse, error)
	GetOpensearchStatus(context.Context, *emptypb.Empty) (*StatusResponse, error)
	mustEmbedUnimplementedLoggingAdminV2Server()
}

// UnimplementedLoggingAdminV2Server must be embedded to have forward compatible implementations.
type UnimplementedLoggingAdminV2Server struct {
}

func (UnimplementedLoggingAdminV2Server) GetOpensearchCluster(context.Context, *emptypb.Empty) (*OpensearchClusterV2, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOpensearchCluster not implemented")
}
func (UnimplementedLoggingAdminV2Server) DeleteOpensearchCluster(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteOpensearchCluster not implemented")
}
func (UnimplementedLoggingAdminV2Server) CreateOrUpdateOpensearchCluster(context.Context, *OpensearchClusterV2) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrUpdateOpensearchCluster not implemented")
}
func (UnimplementedLoggingAdminV2Server) UpgradeAvailable(context.Context, *emptypb.Empty) (*UpgradeAvailableResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpgradeAvailable not implemented")
}
func (UnimplementedLoggingAdminV2Server) DoUpgrade(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DoUpgrade not implemented")
}
func (UnimplementedLoggingAdminV2Server) GetStorageClasses(context.Context, *emptypb.Empty) (*StorageClassResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStorageClasses not implemented")
}
func (UnimplementedLoggingAdminV2Server) GetOpensearchStatus(context.Context, *emptypb.Empty) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOpensearchStatus not implemented")
}
func (UnimplementedLoggingAdminV2Server) mustEmbedUnimplementedLoggingAdminV2Server() {}

// UnsafeLoggingAdminV2Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LoggingAdminV2Server will
// result in compilation errors.
type UnsafeLoggingAdminV2Server interface {
	mustEmbedUnimplementedLoggingAdminV2Server()
}

func RegisterLoggingAdminV2Server(s grpc.ServiceRegistrar, srv LoggingAdminV2Server) {
	s.RegisterService(&LoggingAdminV2_ServiceDesc, srv)
}

func _LoggingAdminV2_GetOpensearchCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggingAdminV2Server).GetOpensearchCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loggingadmin.LoggingAdminV2/GetOpensearchCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggingAdminV2Server).GetOpensearchCluster(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoggingAdminV2_DeleteOpensearchCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggingAdminV2Server).DeleteOpensearchCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loggingadmin.LoggingAdminV2/DeleteOpensearchCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggingAdminV2Server).DeleteOpensearchCluster(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoggingAdminV2_CreateOrUpdateOpensearchCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpensearchClusterV2)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggingAdminV2Server).CreateOrUpdateOpensearchCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loggingadmin.LoggingAdminV2/CreateOrUpdateOpensearchCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggingAdminV2Server).CreateOrUpdateOpensearchCluster(ctx, req.(*OpensearchClusterV2))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoggingAdminV2_UpgradeAvailable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggingAdminV2Server).UpgradeAvailable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loggingadmin.LoggingAdminV2/UpgradeAvailable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggingAdminV2Server).UpgradeAvailable(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoggingAdminV2_DoUpgrade_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggingAdminV2Server).DoUpgrade(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loggingadmin.LoggingAdminV2/DoUpgrade",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggingAdminV2Server).DoUpgrade(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoggingAdminV2_GetStorageClasses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggingAdminV2Server).GetStorageClasses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loggingadmin.LoggingAdminV2/GetStorageClasses",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggingAdminV2Server).GetStorageClasses(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoggingAdminV2_GetOpensearchStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggingAdminV2Server).GetOpensearchStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/loggingadmin.LoggingAdminV2/GetOpensearchStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggingAdminV2Server).GetOpensearchStatus(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// LoggingAdminV2_ServiceDesc is the grpc.ServiceDesc for LoggingAdminV2 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LoggingAdminV2_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "loggingadmin.LoggingAdminV2",
	HandlerType: (*LoggingAdminV2Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOpensearchCluster",
			Handler:    _LoggingAdminV2_GetOpensearchCluster_Handler,
		},
		{
			MethodName: "DeleteOpensearchCluster",
			Handler:    _LoggingAdminV2_DeleteOpensearchCluster_Handler,
		},
		{
			MethodName: "CreateOrUpdateOpensearchCluster",
			Handler:    _LoggingAdminV2_CreateOrUpdateOpensearchCluster_Handler,
		},
		{
			MethodName: "UpgradeAvailable",
			Handler:    _LoggingAdminV2_UpgradeAvailable_Handler,
		},
		{
			MethodName: "DoUpgrade",
			Handler:    _LoggingAdminV2_DoUpgrade_Handler,
		},
		{
			MethodName: "GetStorageClasses",
			Handler:    _LoggingAdminV2_GetStorageClasses_Handler,
		},
		{
			MethodName: "GetOpensearchStatus",
			Handler:    _LoggingAdminV2_GetOpensearchStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/rancher/opni/plugins/logging/pkg/apis/loggingadmin/loggingadmin.proto",
}
