// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: authenticator.proto

package pb

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
	Authenticator_Register_FullMethodName     = "/authenticator.Authenticator/Register"
	Authenticator_NonceConfirm_FullMethodName = "/authenticator.Authenticator/NonceConfirm"
)

// AuthenticatorClient is the client API for Authenticator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthenticatorClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	NonceConfirm(ctx context.Context, in *NonceConfirmRequest, opts ...grpc.CallOption) (*NonceConfirmResponse, error)
}

type authenticatorClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthenticatorClient(cc grpc.ClientConnInterface) AuthenticatorClient {
	return &authenticatorClient{cc}
}

func (c *authenticatorClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, Authenticator_Register_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticatorClient) NonceConfirm(ctx context.Context, in *NonceConfirmRequest, opts ...grpc.CallOption) (*NonceConfirmResponse, error) {
	out := new(NonceConfirmResponse)
	err := c.cc.Invoke(ctx, Authenticator_NonceConfirm_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthenticatorServer is the server API for Authenticator service.
// All implementations must embed UnimplementedAuthenticatorServer
// for forward compatibility
type AuthenticatorServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	NonceConfirm(context.Context, *NonceConfirmRequest) (*NonceConfirmResponse, error)
	mustEmbedUnimplementedAuthenticatorServer()
}

// UnimplementedAuthenticatorServer must be embedded to have forward compatible implementations.
type UnimplementedAuthenticatorServer struct {
}

func (UnimplementedAuthenticatorServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedAuthenticatorServer) NonceConfirm(context.Context, *NonceConfirmRequest) (*NonceConfirmResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NonceConfirm not implemented")
}
func (UnimplementedAuthenticatorServer) mustEmbedUnimplementedAuthenticatorServer() {}

// UnsafeAuthenticatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthenticatorServer will
// result in compilation errors.
type UnsafeAuthenticatorServer interface {
	mustEmbedUnimplementedAuthenticatorServer()
}

func RegisterAuthenticatorServer(s grpc.ServiceRegistrar, srv AuthenticatorServer) {
	s.RegisterService(&Authenticator_ServiceDesc, srv)
}

func _Authenticator_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticatorServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Authenticator_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticatorServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authenticator_NonceConfirm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NonceConfirmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticatorServer).NonceConfirm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Authenticator_NonceConfirm_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticatorServer).NonceConfirm(ctx, req.(*NonceConfirmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Authenticator_ServiceDesc is the grpc.ServiceDesc for Authenticator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Authenticator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "authenticator.Authenticator",
	HandlerType: (*AuthenticatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Authenticator_Register_Handler,
		},
		{
			MethodName: "NonceConfirm",
			Handler:    _Authenticator_NonceConfirm_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "authenticator.proto",
}
