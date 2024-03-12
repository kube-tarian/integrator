// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.19.6
// source: server.proto

package serverpb

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
	Server_NewClusterRegistration_FullMethodName     = "/serverpb.Server/NewClusterRegistration"
	Server_UpdateClusterRegistration_FullMethodName  = "/serverpb.Server/UpdateClusterRegistration"
	Server_DeleteClusterRegistration_FullMethodName  = "/serverpb.Server/DeleteClusterRegistration"
	Server_GetClusters_FullMethodName                = "/serverpb.Server/GetClusters"
	Server_GetCluster_FullMethodName                 = "/serverpb.Server/GetCluster"
	Server_GetClusterApps_FullMethodName             = "/serverpb.Server/GetClusterApps"
	Server_GetClusterAppLaunchConfigs_FullMethodName = "/serverpb.Server/GetClusterAppLaunchConfigs"
	Server_GetClusterApp_FullMethodName              = "/serverpb.Server/GetClusterApp"
	Server_GetClusterDetails_FullMethodName          = "/serverpb.Server/GetClusterDetails"
	Server_AddStoreApp_FullMethodName                = "/serverpb.Server/AddStoreApp"
	Server_UpdateStoreApp_FullMethodName             = "/serverpb.Server/UpdateStoreApp"
	Server_DeleteStoreApp_FullMethodName             = "/serverpb.Server/DeleteStoreApp"
	Server_GetStoreApp_FullMethodName                = "/serverpb.Server/GetStoreApp"
	Server_GetStoreApps_FullMethodName               = "/serverpb.Server/GetStoreApps"
	Server_GetStoreAppValues_FullMethodName          = "/serverpb.Server/GetStoreAppValues"
	Server_DeployStoreApp_FullMethodName             = "/serverpb.Server/DeployStoreApp"
	Server_UnDeployStoreApp_FullMethodName           = "/serverpb.Server/UnDeployStoreApp"
	Server_UpgradeStoreApp_FullMethodName            = "/serverpb.Server/UpgradeStoreApp"
	Server_StoreCredential_FullMethodName            = "/serverpb.Server/StoreCredential"
	Server_SetupDatabase_FullMethodName              = "/serverpb.Server/SetupDatabase"
	Server_RunMigrations_FullMethodName              = "/serverpb.Server/RunMigrations"
)

// ServerClient is the client API for Server service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServerClient interface {
	NewClusterRegistration(ctx context.Context, in *NewClusterRegistrationRequest, opts ...grpc.CallOption) (*NewClusterRegistrationResponse, error)
	UpdateClusterRegistration(ctx context.Context, in *UpdateClusterRegistrationRequest, opts ...grpc.CallOption) (*UpdateClusterRegistrationResponse, error)
	DeleteClusterRegistration(ctx context.Context, in *DeleteClusterRegistrationRequest, opts ...grpc.CallOption) (*DeleteClusterRegistrationResponse, error)
	GetClusters(ctx context.Context, in *GetClustersRequest, opts ...grpc.CallOption) (*GetClustersResponse, error)
	GetCluster(ctx context.Context, in *GetClusterRequest, opts ...grpc.CallOption) (*GetClusterResponse, error)
	GetClusterApps(ctx context.Context, in *GetClusterAppsRequest, opts ...grpc.CallOption) (*GetClusterAppsResponse, error)
	GetClusterAppLaunchConfigs(ctx context.Context, in *GetClusterAppLaunchConfigsRequest, opts ...grpc.CallOption) (*GetClusterAppLaunchConfigsResponse, error)
	GetClusterApp(ctx context.Context, in *GetClusterAppRequest, opts ...grpc.CallOption) (*GetClusterAppResponse, error)
	GetClusterDetails(ctx context.Context, in *GetClusterDetailsRequest, opts ...grpc.CallOption) (*GetClusterDetailsResponse, error)
	AddStoreApp(ctx context.Context, in *AddStoreAppRequest, opts ...grpc.CallOption) (*AddStoreAppResponse, error)
	UpdateStoreApp(ctx context.Context, in *UpdateStoreAppRequest, opts ...grpc.CallOption) (*UpdateStoreAppRsponse, error)
	DeleteStoreApp(ctx context.Context, in *DeleteStoreAppRequest, opts ...grpc.CallOption) (*DeleteStoreAppResponse, error)
	GetStoreApp(ctx context.Context, in *GetStoreAppRequest, opts ...grpc.CallOption) (*GetStoreAppResponse, error)
	GetStoreApps(ctx context.Context, in *GetStoreAppsRequest, opts ...grpc.CallOption) (*GetStoreAppsResponse, error)
	GetStoreAppValues(ctx context.Context, in *GetStoreAppValuesRequest, opts ...grpc.CallOption) (*GetStoreAppValuesResponse, error)
	DeployStoreApp(ctx context.Context, in *DeployStoreAppRequest, opts ...grpc.CallOption) (*DeployStoreAppResponse, error)
	UnDeployStoreApp(ctx context.Context, in *UnDeployStoreAppRequest, opts ...grpc.CallOption) (*UnDeployStoreAppResponse, error)
	UpgradeStoreApp(ctx context.Context, in *UpgradeStoreAppRequest, opts ...grpc.CallOption) (*UpgradeStoreAppResponse, error)
	StoreCredential(ctx context.Context, in *StoreCredentialRequest, opts ...grpc.CallOption) (*StoreCredentialResponse, error)
	SetupDatabase(ctx context.Context, in *DBSetupRequest, opts ...grpc.CallOption) (*DBSetupResponse, error)
	RunMigrations(ctx context.Context, in *DBMigrationRequest, opts ...grpc.CallOption) (*DBMigrationResponse, error)
}

type serverClient struct {
	cc grpc.ClientConnInterface
}

func NewServerClient(cc grpc.ClientConnInterface) ServerClient {
	return &serverClient{cc}
}

func (c *serverClient) NewClusterRegistration(ctx context.Context, in *NewClusterRegistrationRequest, opts ...grpc.CallOption) (*NewClusterRegistrationResponse, error) {
	out := new(NewClusterRegistrationResponse)
	err := c.cc.Invoke(ctx, Server_NewClusterRegistration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) UpdateClusterRegistration(ctx context.Context, in *UpdateClusterRegistrationRequest, opts ...grpc.CallOption) (*UpdateClusterRegistrationResponse, error) {
	out := new(UpdateClusterRegistrationResponse)
	err := c.cc.Invoke(ctx, Server_UpdateClusterRegistration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) DeleteClusterRegistration(ctx context.Context, in *DeleteClusterRegistrationRequest, opts ...grpc.CallOption) (*DeleteClusterRegistrationResponse, error) {
	out := new(DeleteClusterRegistrationResponse)
	err := c.cc.Invoke(ctx, Server_DeleteClusterRegistration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) GetClusters(ctx context.Context, in *GetClustersRequest, opts ...grpc.CallOption) (*GetClustersResponse, error) {
	out := new(GetClustersResponse)
	err := c.cc.Invoke(ctx, Server_GetClusters_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) GetCluster(ctx context.Context, in *GetClusterRequest, opts ...grpc.CallOption) (*GetClusterResponse, error) {
	out := new(GetClusterResponse)
	err := c.cc.Invoke(ctx, Server_GetCluster_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) GetClusterApps(ctx context.Context, in *GetClusterAppsRequest, opts ...grpc.CallOption) (*GetClusterAppsResponse, error) {
	out := new(GetClusterAppsResponse)
	err := c.cc.Invoke(ctx, Server_GetClusterApps_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) GetClusterAppLaunchConfigs(ctx context.Context, in *GetClusterAppLaunchConfigsRequest, opts ...grpc.CallOption) (*GetClusterAppLaunchConfigsResponse, error) {
	out := new(GetClusterAppLaunchConfigsResponse)
	err := c.cc.Invoke(ctx, Server_GetClusterAppLaunchConfigs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) GetClusterApp(ctx context.Context, in *GetClusterAppRequest, opts ...grpc.CallOption) (*GetClusterAppResponse, error) {
	out := new(GetClusterAppResponse)
	err := c.cc.Invoke(ctx, Server_GetClusterApp_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) GetClusterDetails(ctx context.Context, in *GetClusterDetailsRequest, opts ...grpc.CallOption) (*GetClusterDetailsResponse, error) {
	out := new(GetClusterDetailsResponse)
	err := c.cc.Invoke(ctx, Server_GetClusterDetails_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) AddStoreApp(ctx context.Context, in *AddStoreAppRequest, opts ...grpc.CallOption) (*AddStoreAppResponse, error) {
	out := new(AddStoreAppResponse)
	err := c.cc.Invoke(ctx, Server_AddStoreApp_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) UpdateStoreApp(ctx context.Context, in *UpdateStoreAppRequest, opts ...grpc.CallOption) (*UpdateStoreAppRsponse, error) {
	out := new(UpdateStoreAppRsponse)
	err := c.cc.Invoke(ctx, Server_UpdateStoreApp_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) DeleteStoreApp(ctx context.Context, in *DeleteStoreAppRequest, opts ...grpc.CallOption) (*DeleteStoreAppResponse, error) {
	out := new(DeleteStoreAppResponse)
	err := c.cc.Invoke(ctx, Server_DeleteStoreApp_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) GetStoreApp(ctx context.Context, in *GetStoreAppRequest, opts ...grpc.CallOption) (*GetStoreAppResponse, error) {
	out := new(GetStoreAppResponse)
	err := c.cc.Invoke(ctx, Server_GetStoreApp_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) GetStoreApps(ctx context.Context, in *GetStoreAppsRequest, opts ...grpc.CallOption) (*GetStoreAppsResponse, error) {
	out := new(GetStoreAppsResponse)
	err := c.cc.Invoke(ctx, Server_GetStoreApps_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) GetStoreAppValues(ctx context.Context, in *GetStoreAppValuesRequest, opts ...grpc.CallOption) (*GetStoreAppValuesResponse, error) {
	out := new(GetStoreAppValuesResponse)
	err := c.cc.Invoke(ctx, Server_GetStoreAppValues_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) DeployStoreApp(ctx context.Context, in *DeployStoreAppRequest, opts ...grpc.CallOption) (*DeployStoreAppResponse, error) {
	out := new(DeployStoreAppResponse)
	err := c.cc.Invoke(ctx, Server_DeployStoreApp_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) UnDeployStoreApp(ctx context.Context, in *UnDeployStoreAppRequest, opts ...grpc.CallOption) (*UnDeployStoreAppResponse, error) {
	out := new(UnDeployStoreAppResponse)
	err := c.cc.Invoke(ctx, Server_UnDeployStoreApp_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) UpgradeStoreApp(ctx context.Context, in *UpgradeStoreAppRequest, opts ...grpc.CallOption) (*UpgradeStoreAppResponse, error) {
	out := new(UpgradeStoreAppResponse)
	err := c.cc.Invoke(ctx, Server_UpgradeStoreApp_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) StoreCredential(ctx context.Context, in *StoreCredentialRequest, opts ...grpc.CallOption) (*StoreCredentialResponse, error) {
	out := new(StoreCredentialResponse)
	err := c.cc.Invoke(ctx, Server_StoreCredential_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) SetupDatabase(ctx context.Context, in *DBSetupRequest, opts ...grpc.CallOption) (*DBSetupResponse, error) {
	out := new(DBSetupResponse)
	err := c.cc.Invoke(ctx, Server_SetupDatabase_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverClient) RunMigrations(ctx context.Context, in *DBMigrationRequest, opts ...grpc.CallOption) (*DBMigrationResponse, error) {
	out := new(DBMigrationResponse)
	err := c.cc.Invoke(ctx, Server_RunMigrations_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServerServer is the server API for Server service.
// All implementations must embed UnimplementedServerServer
// for forward compatibility
type ServerServer interface {
	NewClusterRegistration(context.Context, *NewClusterRegistrationRequest) (*NewClusterRegistrationResponse, error)
	UpdateClusterRegistration(context.Context, *UpdateClusterRegistrationRequest) (*UpdateClusterRegistrationResponse, error)
	DeleteClusterRegistration(context.Context, *DeleteClusterRegistrationRequest) (*DeleteClusterRegistrationResponse, error)
	GetClusters(context.Context, *GetClustersRequest) (*GetClustersResponse, error)
	GetCluster(context.Context, *GetClusterRequest) (*GetClusterResponse, error)
	GetClusterApps(context.Context, *GetClusterAppsRequest) (*GetClusterAppsResponse, error)
	GetClusterAppLaunchConfigs(context.Context, *GetClusterAppLaunchConfigsRequest) (*GetClusterAppLaunchConfigsResponse, error)
	GetClusterApp(context.Context, *GetClusterAppRequest) (*GetClusterAppResponse, error)
	GetClusterDetails(context.Context, *GetClusterDetailsRequest) (*GetClusterDetailsResponse, error)
	AddStoreApp(context.Context, *AddStoreAppRequest) (*AddStoreAppResponse, error)
	UpdateStoreApp(context.Context, *UpdateStoreAppRequest) (*UpdateStoreAppRsponse, error)
	DeleteStoreApp(context.Context, *DeleteStoreAppRequest) (*DeleteStoreAppResponse, error)
	GetStoreApp(context.Context, *GetStoreAppRequest) (*GetStoreAppResponse, error)
	GetStoreApps(context.Context, *GetStoreAppsRequest) (*GetStoreAppsResponse, error)
	GetStoreAppValues(context.Context, *GetStoreAppValuesRequest) (*GetStoreAppValuesResponse, error)
	DeployStoreApp(context.Context, *DeployStoreAppRequest) (*DeployStoreAppResponse, error)
	UnDeployStoreApp(context.Context, *UnDeployStoreAppRequest) (*UnDeployStoreAppResponse, error)
	UpgradeStoreApp(context.Context, *UpgradeStoreAppRequest) (*UpgradeStoreAppResponse, error)
	StoreCredential(context.Context, *StoreCredentialRequest) (*StoreCredentialResponse, error)
	SetupDatabase(context.Context, *DBSetupRequest) (*DBSetupResponse, error)
	RunMigrations(context.Context, *DBMigrationRequest) (*DBMigrationResponse, error)
	mustEmbedUnimplementedServerServer()
}

// UnimplementedServerServer must be embedded to have forward compatible implementations.
type UnimplementedServerServer struct {
}

func (UnimplementedServerServer) NewClusterRegistration(context.Context, *NewClusterRegistrationRequest) (*NewClusterRegistrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewClusterRegistration not implemented")
}
func (UnimplementedServerServer) UpdateClusterRegistration(context.Context, *UpdateClusterRegistrationRequest) (*UpdateClusterRegistrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateClusterRegistration not implemented")
}
func (UnimplementedServerServer) DeleteClusterRegistration(context.Context, *DeleteClusterRegistrationRequest) (*DeleteClusterRegistrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteClusterRegistration not implemented")
}
func (UnimplementedServerServer) GetClusters(context.Context, *GetClustersRequest) (*GetClustersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClusters not implemented")
}
func (UnimplementedServerServer) GetCluster(context.Context, *GetClusterRequest) (*GetClusterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCluster not implemented")
}
func (UnimplementedServerServer) GetClusterApps(context.Context, *GetClusterAppsRequest) (*GetClusterAppsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClusterApps not implemented")
}
func (UnimplementedServerServer) GetClusterAppLaunchConfigs(context.Context, *GetClusterAppLaunchConfigsRequest) (*GetClusterAppLaunchConfigsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClusterAppLaunchConfigs not implemented")
}
func (UnimplementedServerServer) GetClusterApp(context.Context, *GetClusterAppRequest) (*GetClusterAppResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClusterApp not implemented")
}
func (UnimplementedServerServer) GetClusterDetails(context.Context, *GetClusterDetailsRequest) (*GetClusterDetailsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClusterDetails not implemented")
}
func (UnimplementedServerServer) AddStoreApp(context.Context, *AddStoreAppRequest) (*AddStoreAppResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddStoreApp not implemented")
}
func (UnimplementedServerServer) UpdateStoreApp(context.Context, *UpdateStoreAppRequest) (*UpdateStoreAppRsponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStoreApp not implemented")
}
func (UnimplementedServerServer) DeleteStoreApp(context.Context, *DeleteStoreAppRequest) (*DeleteStoreAppResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteStoreApp not implemented")
}
func (UnimplementedServerServer) GetStoreApp(context.Context, *GetStoreAppRequest) (*GetStoreAppResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStoreApp not implemented")
}
func (UnimplementedServerServer) GetStoreApps(context.Context, *GetStoreAppsRequest) (*GetStoreAppsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStoreApps not implemented")
}
func (UnimplementedServerServer) GetStoreAppValues(context.Context, *GetStoreAppValuesRequest) (*GetStoreAppValuesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStoreAppValues not implemented")
}
func (UnimplementedServerServer) DeployStoreApp(context.Context, *DeployStoreAppRequest) (*DeployStoreAppResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeployStoreApp not implemented")
}
func (UnimplementedServerServer) UnDeployStoreApp(context.Context, *UnDeployStoreAppRequest) (*UnDeployStoreAppResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnDeployStoreApp not implemented")
}
func (UnimplementedServerServer) UpgradeStoreApp(context.Context, *UpgradeStoreAppRequest) (*UpgradeStoreAppResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpgradeStoreApp not implemented")
}
func (UnimplementedServerServer) StoreCredential(context.Context, *StoreCredentialRequest) (*StoreCredentialResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreCredential not implemented")
}
func (UnimplementedServerServer) SetupDatabase(context.Context, *DBSetupRequest) (*DBSetupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetupDatabase not implemented")
}
func (UnimplementedServerServer) RunMigrations(context.Context, *DBMigrationRequest) (*DBMigrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunMigrations not implemented")
}
func (UnimplementedServerServer) mustEmbedUnimplementedServerServer() {}

// UnsafeServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServerServer will
// result in compilation errors.
type UnsafeServerServer interface {
	mustEmbedUnimplementedServerServer()
}

func RegisterServerServer(s grpc.ServiceRegistrar, srv ServerServer) {
	s.RegisterService(&Server_ServiceDesc, srv)
}

func _Server_NewClusterRegistration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewClusterRegistrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).NewClusterRegistration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_NewClusterRegistration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).NewClusterRegistration(ctx, req.(*NewClusterRegistrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_UpdateClusterRegistration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateClusterRegistrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).UpdateClusterRegistration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_UpdateClusterRegistration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).UpdateClusterRegistration(ctx, req.(*UpdateClusterRegistrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_DeleteClusterRegistration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteClusterRegistrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).DeleteClusterRegistration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_DeleteClusterRegistration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).DeleteClusterRegistration(ctx, req.(*DeleteClusterRegistrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_GetClusters_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClustersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).GetClusters(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_GetClusters_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).GetClusters(ctx, req.(*GetClustersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_GetCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClusterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).GetCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_GetCluster_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).GetCluster(ctx, req.(*GetClusterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_GetClusterApps_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClusterAppsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).GetClusterApps(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_GetClusterApps_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).GetClusterApps(ctx, req.(*GetClusterAppsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_GetClusterAppLaunchConfigs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClusterAppLaunchConfigsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).GetClusterAppLaunchConfigs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_GetClusterAppLaunchConfigs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).GetClusterAppLaunchConfigs(ctx, req.(*GetClusterAppLaunchConfigsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_GetClusterApp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClusterAppRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).GetClusterApp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_GetClusterApp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).GetClusterApp(ctx, req.(*GetClusterAppRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_GetClusterDetails_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClusterDetailsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).GetClusterDetails(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_GetClusterDetails_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).GetClusterDetails(ctx, req.(*GetClusterDetailsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_AddStoreApp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddStoreAppRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).AddStoreApp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_AddStoreApp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).AddStoreApp(ctx, req.(*AddStoreAppRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_UpdateStoreApp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateStoreAppRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).UpdateStoreApp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_UpdateStoreApp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).UpdateStoreApp(ctx, req.(*UpdateStoreAppRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_DeleteStoreApp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteStoreAppRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).DeleteStoreApp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_DeleteStoreApp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).DeleteStoreApp(ctx, req.(*DeleteStoreAppRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_GetStoreApp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStoreAppRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).GetStoreApp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_GetStoreApp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).GetStoreApp(ctx, req.(*GetStoreAppRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_GetStoreApps_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStoreAppsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).GetStoreApps(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_GetStoreApps_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).GetStoreApps(ctx, req.(*GetStoreAppsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_GetStoreAppValues_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStoreAppValuesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).GetStoreAppValues(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_GetStoreAppValues_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).GetStoreAppValues(ctx, req.(*GetStoreAppValuesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_DeployStoreApp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeployStoreAppRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).DeployStoreApp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_DeployStoreApp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).DeployStoreApp(ctx, req.(*DeployStoreAppRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_UnDeployStoreApp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnDeployStoreAppRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).UnDeployStoreApp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_UnDeployStoreApp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).UnDeployStoreApp(ctx, req.(*UnDeployStoreAppRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_UpgradeStoreApp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpgradeStoreAppRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).UpgradeStoreApp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_UpgradeStoreApp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).UpgradeStoreApp(ctx, req.(*UpgradeStoreAppRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_StoreCredential_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreCredentialRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).StoreCredential(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_StoreCredential_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).StoreCredential(ctx, req.(*StoreCredentialRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_SetupDatabase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DBSetupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).SetupDatabase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_SetupDatabase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).SetupDatabase(ctx, req.(*DBSetupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Server_RunMigrations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DBMigrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).RunMigrations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Server_RunMigrations_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).RunMigrations(ctx, req.(*DBMigrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Server_ServiceDesc is the grpc.ServiceDesc for Server service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Server_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "serverpb.Server",
	HandlerType: (*ServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewClusterRegistration",
			Handler:    _Server_NewClusterRegistration_Handler,
		},
		{
			MethodName: "UpdateClusterRegistration",
			Handler:    _Server_UpdateClusterRegistration_Handler,
		},
		{
			MethodName: "DeleteClusterRegistration",
			Handler:    _Server_DeleteClusterRegistration_Handler,
		},
		{
			MethodName: "GetClusters",
			Handler:    _Server_GetClusters_Handler,
		},
		{
			MethodName: "GetCluster",
			Handler:    _Server_GetCluster_Handler,
		},
		{
			MethodName: "GetClusterApps",
			Handler:    _Server_GetClusterApps_Handler,
		},
		{
			MethodName: "GetClusterAppLaunchConfigs",
			Handler:    _Server_GetClusterAppLaunchConfigs_Handler,
		},
		{
			MethodName: "GetClusterApp",
			Handler:    _Server_GetClusterApp_Handler,
		},
		{
			MethodName: "GetClusterDetails",
			Handler:    _Server_GetClusterDetails_Handler,
		},
		{
			MethodName: "AddStoreApp",
			Handler:    _Server_AddStoreApp_Handler,
		},
		{
			MethodName: "UpdateStoreApp",
			Handler:    _Server_UpdateStoreApp_Handler,
		},
		{
			MethodName: "DeleteStoreApp",
			Handler:    _Server_DeleteStoreApp_Handler,
		},
		{
			MethodName: "GetStoreApp",
			Handler:    _Server_GetStoreApp_Handler,
		},
		{
			MethodName: "GetStoreApps",
			Handler:    _Server_GetStoreApps_Handler,
		},
		{
			MethodName: "GetStoreAppValues",
			Handler:    _Server_GetStoreAppValues_Handler,
		},
		{
			MethodName: "DeployStoreApp",
			Handler:    _Server_DeployStoreApp_Handler,
		},
		{
			MethodName: "UnDeployStoreApp",
			Handler:    _Server_UnDeployStoreApp_Handler,
		},
		{
			MethodName: "UpgradeStoreApp",
			Handler:    _Server_UpgradeStoreApp_Handler,
		},
		{
			MethodName: "StoreCredential",
			Handler:    _Server_StoreCredential_Handler,
		},
		{
			MethodName: "SetupDatabase",
			Handler:    _Server_SetupDatabase_Handler,
		},
		{
			MethodName: "RunMigrations",
			Handler:    _Server_RunMigrations_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server.proto",
}
