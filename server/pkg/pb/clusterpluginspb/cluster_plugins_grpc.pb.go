// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.19.6
// source: cluster_plugins.proto

package clusterpluginspb

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
	ClusterPlugins_DeployClusterPlugin_FullMethodName   = "/clusterpluginspb.ClusterPlugins/DeployClusterPlugin"
	ClusterPlugins_UnDeployClusterPlugin_FullMethodName = "/clusterpluginspb.ClusterPlugins/UnDeployClusterPlugin"
	ClusterPlugins_GetClusterPlugins_FullMethodName     = "/clusterpluginspb.ClusterPlugins/GetClusterPlugins"
)

// ClusterPluginsClient is the client API for ClusterPlugins service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClusterPluginsClient interface {
	DeployClusterPlugin(ctx context.Context, in *DeployClusterPluginRequest, opts ...grpc.CallOption) (*DeployClusterPluginResponse, error)
	UnDeployClusterPlugin(ctx context.Context, in *UnDeployClusterPluginRequest, opts ...grpc.CallOption) (*UnDeployClusterPluginResponse, error)
	GetClusterPlugins(ctx context.Context, in *GetClusterPluginsRequest, opts ...grpc.CallOption) (*GetClusterPluginsResponse, error)
}

type clusterPluginsClient struct {
	cc grpc.ClientConnInterface
}

func NewClusterPluginsClient(cc grpc.ClientConnInterface) ClusterPluginsClient {
	return &clusterPluginsClient{cc}
}

func (c *clusterPluginsClient) DeployClusterPlugin(ctx context.Context, in *DeployClusterPluginRequest, opts ...grpc.CallOption) (*DeployClusterPluginResponse, error) {
	out := new(DeployClusterPluginResponse)
	err := c.cc.Invoke(ctx, ClusterPlugins_DeployClusterPlugin_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterPluginsClient) UnDeployClusterPlugin(ctx context.Context, in *UnDeployClusterPluginRequest, opts ...grpc.CallOption) (*UnDeployClusterPluginResponse, error) {
	out := new(UnDeployClusterPluginResponse)
	err := c.cc.Invoke(ctx, ClusterPlugins_UnDeployClusterPlugin_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterPluginsClient) GetClusterPlugins(ctx context.Context, in *GetClusterPluginsRequest, opts ...grpc.CallOption) (*GetClusterPluginsResponse, error) {
	out := new(GetClusterPluginsResponse)
	err := c.cc.Invoke(ctx, ClusterPlugins_GetClusterPlugins_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClusterPluginsServer is the server API for ClusterPlugins service.
// All implementations must embed UnimplementedClusterPluginsServer
// for forward compatibility
type ClusterPluginsServer interface {
	DeployClusterPlugin(context.Context, *DeployClusterPluginRequest) (*DeployClusterPluginResponse, error)
	UnDeployClusterPlugin(context.Context, *UnDeployClusterPluginRequest) (*UnDeployClusterPluginResponse, error)
	GetClusterPlugins(context.Context, *GetClusterPluginsRequest) (*GetClusterPluginsResponse, error)
	mustEmbedUnimplementedClusterPluginsServer()
}

// UnimplementedClusterPluginsServer must be embedded to have forward compatible implementations.
type UnimplementedClusterPluginsServer struct {
}

func (UnimplementedClusterPluginsServer) DeployClusterPlugin(context.Context, *DeployClusterPluginRequest) (*DeployClusterPluginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeployClusterPlugin not implemented")
}
func (UnimplementedClusterPluginsServer) UnDeployClusterPlugin(context.Context, *UnDeployClusterPluginRequest) (*UnDeployClusterPluginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnDeployClusterPlugin not implemented")
}
func (UnimplementedClusterPluginsServer) GetClusterPlugins(context.Context, *GetClusterPluginsRequest) (*GetClusterPluginsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClusterPlugins not implemented")
}
func (UnimplementedClusterPluginsServer) mustEmbedUnimplementedClusterPluginsServer() {}

// UnsafeClusterPluginsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClusterPluginsServer will
// result in compilation errors.
type UnsafeClusterPluginsServer interface {
	mustEmbedUnimplementedClusterPluginsServer()
}

func RegisterClusterPluginsServer(s grpc.ServiceRegistrar, srv ClusterPluginsServer) {
	s.RegisterService(&ClusterPlugins_ServiceDesc, srv)
}

func _ClusterPlugins_DeployClusterPlugin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeployClusterPluginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterPluginsServer).DeployClusterPlugin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClusterPlugins_DeployClusterPlugin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterPluginsServer).DeployClusterPlugin(ctx, req.(*DeployClusterPluginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClusterPlugins_UnDeployClusterPlugin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnDeployClusterPluginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterPluginsServer).UnDeployClusterPlugin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClusterPlugins_UnDeployClusterPlugin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterPluginsServer).UnDeployClusterPlugin(ctx, req.(*UnDeployClusterPluginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClusterPlugins_GetClusterPlugins_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClusterPluginsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterPluginsServer).GetClusterPlugins(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClusterPlugins_GetClusterPlugins_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterPluginsServer).GetClusterPlugins(ctx, req.(*GetClusterPluginsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ClusterPlugins_ServiceDesc is the grpc.ServiceDesc for ClusterPlugins service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClusterPlugins_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "clusterpluginspb.ClusterPlugins",
	HandlerType: (*ClusterPluginsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeployClusterPlugin",
			Handler:    _ClusterPlugins_DeployClusterPlugin_Handler,
		},
		{
			MethodName: "UnDeployClusterPlugin",
			Handler:    _ClusterPlugins_UnDeployClusterPlugin_Handler,
		},
		{
			MethodName: "GetClusterPlugins",
			Handler:    _ClusterPlugins_GetClusterPlugins_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cluster_plugins.proto",
}