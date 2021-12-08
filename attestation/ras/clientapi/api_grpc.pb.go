// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package clientapi

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RasClient is the client API for Ras service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RasClient interface {
	CreateIKCert(ctx context.Context, in *CreateIKCertRequest, opts ...grpc.CallOption) (*CreateIKCertReply, error)
	GenerateEKCert(ctx context.Context, in *GenerateEKCertRequest, opts ...grpc.CallOption) (*GenerateEKCertReply, error)
	RegisterClient(ctx context.Context, in *RegisterClientRequest, opts ...grpc.CallOption) (*RegisterClientReply, error)
	UnregisterClient(ctx context.Context, in *UnregisterClientRequest, opts ...grpc.CallOption) (*UnregisterClientReply, error)
	SendHeartbeat(ctx context.Context, in *SendHeartbeatRequest, opts ...grpc.CallOption) (*SendHeartbeatReply, error)
	SendReport(ctx context.Context, in *SendReportRequest, opts ...grpc.CallOption) (*SendReportReply, error)
}

type rasClient struct {
	cc grpc.ClientConnInterface
}

func NewRasClient(cc grpc.ClientConnInterface) RasClient {
	return &rasClient{cc}
}

func (c *rasClient) CreateIKCert(ctx context.Context, in *CreateIKCertRequest, opts ...grpc.CallOption) (*CreateIKCertReply, error) {
	out := new(CreateIKCertReply)
	err := c.cc.Invoke(ctx, "/Ras/CreateIKCert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rasClient) GenerateEKCert(ctx context.Context, in *GenerateEKCertRequest, opts ...grpc.CallOption) (*GenerateEKCertReply, error) {
	out := new(GenerateEKCertReply)
	err := c.cc.Invoke(ctx, "/Ras/GenerateEKCert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rasClient) RegisterClient(ctx context.Context, in *RegisterClientRequest, opts ...grpc.CallOption) (*RegisterClientReply, error) {
	out := new(RegisterClientReply)
	err := c.cc.Invoke(ctx, "/Ras/RegisterClient", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rasClient) UnregisterClient(ctx context.Context, in *UnregisterClientRequest, opts ...grpc.CallOption) (*UnregisterClientReply, error) {
	out := new(UnregisterClientReply)
	err := c.cc.Invoke(ctx, "/Ras/UnregisterClient", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rasClient) SendHeartbeat(ctx context.Context, in *SendHeartbeatRequest, opts ...grpc.CallOption) (*SendHeartbeatReply, error) {
	out := new(SendHeartbeatReply)
	err := c.cc.Invoke(ctx, "/Ras/SendHeartbeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rasClient) SendReport(ctx context.Context, in *SendReportRequest, opts ...grpc.CallOption) (*SendReportReply, error) {
	out := new(SendReportReply)
	err := c.cc.Invoke(ctx, "/Ras/SendReport", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RasServer is the server API for Ras service.
// All implementations must embed UnimplementedRasServer
// for forward compatibility
type RasServer interface {
	CreateIKCert(context.Context, *CreateIKCertRequest) (*CreateIKCertReply, error)
	GenerateEKCert(context.Context, *GenerateEKCertRequest) (*GenerateEKCertReply, error)
	RegisterClient(context.Context, *RegisterClientRequest) (*RegisterClientReply, error)
	UnregisterClient(context.Context, *UnregisterClientRequest) (*UnregisterClientReply, error)
	SendHeartbeat(context.Context, *SendHeartbeatRequest) (*SendHeartbeatReply, error)
	SendReport(context.Context, *SendReportRequest) (*SendReportReply, error)
	mustEmbedUnimplementedRasServer()
}

// UnimplementedRasServer must be embedded to have forward compatible implementations.
type UnimplementedRasServer struct {
}

func (UnimplementedRasServer) CreateIKCert(context.Context, *CreateIKCertRequest) (*CreateIKCertReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateIKCert not implemented")
}
func (UnimplementedRasServer) GenerateEKCert(context.Context, *GenerateEKCertRequest) (*GenerateEKCertReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateEKCert not implemented")
}
func (UnimplementedRasServer) RegisterClient(context.Context, *RegisterClientRequest) (*RegisterClientReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterClient not implemented")
}
func (UnimplementedRasServer) UnregisterClient(context.Context, *UnregisterClientRequest) (*UnregisterClientReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnregisterClient not implemented")
}
func (UnimplementedRasServer) SendHeartbeat(context.Context, *SendHeartbeatRequest) (*SendHeartbeatReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendHeartbeat not implemented")
}
func (UnimplementedRasServer) SendReport(context.Context, *SendReportRequest) (*SendReportReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendReport not implemented")
}
func (UnimplementedRasServer) mustEmbedUnimplementedRasServer() {}

// UnsafeRasServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RasServer will
// result in compilation errors.
type UnsafeRasServer interface {
	mustEmbedUnimplementedRasServer()
}

func RegisterRasServer(s grpc.ServiceRegistrar, srv RasServer) {
	s.RegisterService(&Ras_ServiceDesc, srv)
}

func _Ras_CreateIKCert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateIKCertRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RasServer).CreateIKCert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Ras/CreateIKCert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RasServer).CreateIKCert(ctx, req.(*CreateIKCertRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ras_GenerateEKCert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateEKCertRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RasServer).GenerateEKCert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Ras/GenerateEKCert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RasServer).GenerateEKCert(ctx, req.(*GenerateEKCertRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ras_RegisterClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterClientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RasServer).RegisterClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Ras/RegisterClient",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RasServer).RegisterClient(ctx, req.(*RegisterClientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ras_UnregisterClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnregisterClientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RasServer).UnregisterClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Ras/UnregisterClient",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RasServer).UnregisterClient(ctx, req.(*UnregisterClientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ras_SendHeartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendHeartbeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RasServer).SendHeartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Ras/SendHeartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RasServer).SendHeartbeat(ctx, req.(*SendHeartbeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ras_SendReport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RasServer).SendReport(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Ras/SendReport",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RasServer).SendReport(ctx, req.(*SendReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Ras_ServiceDesc is the grpc.ServiceDesc for Ras service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Ras_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Ras",
	HandlerType: (*RasServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateIKCert",
			Handler:    _Ras_CreateIKCert_Handler,
		},
		{
			MethodName: "GenerateEKCert",
			Handler:    _Ras_GenerateEKCert_Handler,
		},
		{
			MethodName: "RegisterClient",
			Handler:    _Ras_RegisterClient_Handler,
		},
		{
			MethodName: "UnregisterClient",
			Handler:    _Ras_UnregisterClient_Handler,
		},
		{
			MethodName: "SendHeartbeat",
			Handler:    _Ras_SendHeartbeat_Handler,
		},
		{
			MethodName: "SendReport",
			Handler:    _Ras_SendReport_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "clientapi/api.proto",
}
