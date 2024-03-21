// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: kvstore.proto

package protobuf

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

const (
	KVStoreService_Set_FullMethodName    = "/kvstore.KVStoreService/Set"
	KVStoreService_Get_FullMethodName    = "/kvstore.KVStoreService/Get"
	KVStoreService_Delete_FullMethodName = "/kvstore.KVStoreService/Delete"
)

// KVStoreServiceClient is the client API for KVStoreService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KVStoreServiceClient interface {
	// Set方法用于设置键值对
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error)
	// Get方法用于获取键值对的值
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	// Delete方法用于删除键值对
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
}

type kVStoreServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKVStoreServiceClient(cc grpc.ClientConnInterface) KVStoreServiceClient {
	return &kVStoreServiceClient{cc}
}

func (c *kVStoreServiceClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error) {
	out := new(SetResponse)
	err := c.cc.Invoke(ctx, KVStoreService_Set_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVStoreServiceClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, KVStoreService_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVStoreServiceClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, KVStoreService_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KVStoreServiceServer is the server API for KVStoreService service.
// All implementations must embed UnimplementedKVStoreServiceServer
// for forward compatibility
type KVStoreServiceServer interface {
	// Set方法用于设置键值对
	Set(context.Context, *SetRequest) (*SetResponse, error)
	// Get方法用于获取键值对的值
	Get(context.Context, *GetRequest) (*GetResponse, error)
	// Delete方法用于删除键值对
	Delete(context.Context, *DeleteRequest) (*DeleteResponse, error)
	mustEmbedUnimplementedKVStoreServiceServer()
}

// UnimplementedKVStoreServiceServer must be embedded to have forward compatible implementations.
type UnimplementedKVStoreServiceServer struct {
}

func (UnimplementedKVStoreServiceServer) Set(context.Context, *SetRequest) (*SetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedKVStoreServiceServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedKVStoreServiceServer) Delete(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedKVStoreServiceServer) mustEmbedUnimplementedKVStoreServiceServer() {}

// UnsafeKVStoreServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KVStoreServiceServer will
// result in compilation errors.
type UnsafeKVStoreServiceServer interface {
	mustEmbedUnimplementedKVStoreServiceServer()
}

func RegisterKVStoreServiceServer(s grpc.ServiceRegistrar, srv KVStoreServiceServer) {
	s.RegisterService(&KVStoreService_ServiceDesc, srv)
}

func _KVStoreService_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVStoreServiceServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KVStoreService_Set_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVStoreServiceServer).Set(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVStoreService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVStoreServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KVStoreService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVStoreServiceServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVStoreService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVStoreServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KVStoreService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVStoreServiceServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KVStoreService_ServiceDesc is the grpc.ServiceDesc for KVStoreService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KVStoreService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kvstore.KVStoreService",
	HandlerType: (*KVStoreServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Set",
			Handler:    _KVStoreService_Set_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _KVStoreService_Get_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _KVStoreService_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kvstore.proto",
}
