// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: pb/git.proto

package git

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	GitService_GetDependentRepositoriesByRepo_FullMethodName  = "/git.GitService/GetDependentRepositoriesByRepo"
	GitService_GetStarsByRepo_FullMethodName                  = "/git.GitService/GetStarsByRepo"
	GitService_GetForksByRepo_FullMethodName                  = "/git.GitService/GetForksByRepo"
	GitService_GetTotalIssuesByRepo_FullMethodName            = "/git.GitService/GetTotalIssuesByRepo"
	GitService_GetUserSolvedIssuesByRepo_FullMethodName       = "/git.GitService/GetUserSolvedIssuesByRepo"
	GitService_GetTotalPullRequestsByRepo_FullMethodName      = "/git.GitService/GetTotalPullRequestsByRepo"
	GitService_GetUserMergedPullRequestsByRepo_FullMethodName = "/git.GitService/GetUserMergedPullRequestsByRepo"
	GitService_GetTotalCodeReviewsByRepo_FullMethodName       = "/git.GitService/GetTotalCodeReviewsByRepo"
	GitService_GetUserCodeReviewsByRepo_FullMethodName        = "/git.GitService/GetUserCodeReviewsByRepo"
	GitService_GetLineChangesByRepo_FullMethodName            = "/git.GitService/GetLineChangesByRepo"
	GitService_GetName_FullMethodName                         = "/git.GitService/GetName"
	GitService_GetCompany_FullMethodName                      = "/git.GitService/GetCompany"
	GitService_GetLocation_FullMethodName                     = "/git.GitService/GetLocation"
	GitService_GetEmail_FullMethodName                        = "/git.GitService/GetEmail"
	GitService_GetBio_FullMethodName                          = "/git.GitService/GetBio"
	GitService_GetOrganizations_FullMethodName                = "/git.GitService/GetOrganizations"
	GitService_GetFollowers_FullMethodName                    = "/git.GitService/GetFollowers"
	GitService_GetReadme_FullMethodName                       = "/git.GitService/GetReadme"
	GitService_GetCommits_FullMethodName                      = "/git.GitService/GetCommits"
)

// GitServiceClient is the client API for GitService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GitServiceClient interface {
	// Repo table methods
	GetDependentRepositoriesByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoIntMapResponse, error)
	GetStarsByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoIntMapResponse, error)
	GetForksByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoIntMapResponse, error)
	GetTotalIssuesByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoIntMapResponse, error)
	GetUserSolvedIssuesByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoIntMapResponse, error)
	GetTotalPullRequestsByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoIntMapResponse, error)
	GetUserMergedPullRequestsByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoIntMapResponse, error)
	GetTotalCodeReviewsByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoIntMapResponse, error)
	GetUserCodeReviewsByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoIntMapResponse, error)
	GetLineChangesByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoLineChangesResponse, error)
	// User table methods
	GetName(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*StringResponse, error)
	GetCompany(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*StringResponse, error)
	GetLocation(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*StringResponse, error)
	GetEmail(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*StringResponse, error)
	GetBio(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*StringResponse, error)
	GetOrganizations(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*OrgListResponse, error)
	GetFollowers(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*IntResponse, error)
	GetReadme(ctx context.Context, in *GetReadmeRequest, opts ...grpc.CallOption) (*StringResponse, error)
	GetCommits(ctx context.Context, in *GetCommitsRequest, opts ...grpc.CallOption) (*StringResponse, error)
}

type gitServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGitServiceClient(cc grpc.ClientConnInterface) GitServiceClient {
	return &gitServiceClient{cc}
}

func (c *gitServiceClient) GetDependentRepositoriesByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoIntMapResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RepoIntMapResponse)
	err := c.cc.Invoke(ctx, GitService_GetDependentRepositoriesByRepo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitServiceClient) GetStarsByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoIntMapResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RepoIntMapResponse)
	err := c.cc.Invoke(ctx, GitService_GetStarsByRepo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitServiceClient) GetForksByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoIntMapResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RepoIntMapResponse)
	err := c.cc.Invoke(ctx, GitService_GetForksByRepo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitServiceClient) GetTotalIssuesByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoIntMapResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RepoIntMapResponse)
	err := c.cc.Invoke(ctx, GitService_GetTotalIssuesByRepo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitServiceClient) GetUserSolvedIssuesByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoIntMapResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RepoIntMapResponse)
	err := c.cc.Invoke(ctx, GitService_GetUserSolvedIssuesByRepo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitServiceClient) GetTotalPullRequestsByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoIntMapResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RepoIntMapResponse)
	err := c.cc.Invoke(ctx, GitService_GetTotalPullRequestsByRepo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitServiceClient) GetUserMergedPullRequestsByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoIntMapResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RepoIntMapResponse)
	err := c.cc.Invoke(ctx, GitService_GetUserMergedPullRequestsByRepo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitServiceClient) GetTotalCodeReviewsByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoIntMapResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RepoIntMapResponse)
	err := c.cc.Invoke(ctx, GitService_GetTotalCodeReviewsByRepo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitServiceClient) GetUserCodeReviewsByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoIntMapResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RepoIntMapResponse)
	err := c.cc.Invoke(ctx, GitService_GetUserCodeReviewsByRepo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitServiceClient) GetLineChangesByRepo(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*RepoLineChangesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RepoLineChangesResponse)
	err := c.cc.Invoke(ctx, GitService_GetLineChangesByRepo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitServiceClient) GetName(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*StringResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StringResponse)
	err := c.cc.Invoke(ctx, GitService_GetName_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitServiceClient) GetCompany(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*StringResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StringResponse)
	err := c.cc.Invoke(ctx, GitService_GetCompany_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitServiceClient) GetLocation(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*StringResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StringResponse)
	err := c.cc.Invoke(ctx, GitService_GetLocation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitServiceClient) GetEmail(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*StringResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StringResponse)
	err := c.cc.Invoke(ctx, GitService_GetEmail_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitServiceClient) GetBio(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*StringResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StringResponse)
	err := c.cc.Invoke(ctx, GitService_GetBio_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitServiceClient) GetOrganizations(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*OrgListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrgListResponse)
	err := c.cc.Invoke(ctx, GitService_GetOrganizations_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitServiceClient) GetFollowers(ctx context.Context, in *GetUsernameRequest, opts ...grpc.CallOption) (*IntResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IntResponse)
	err := c.cc.Invoke(ctx, GitService_GetFollowers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitServiceClient) GetReadme(ctx context.Context, in *GetReadmeRequest, opts ...grpc.CallOption) (*StringResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StringResponse)
	err := c.cc.Invoke(ctx, GitService_GetReadme_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gitServiceClient) GetCommits(ctx context.Context, in *GetCommitsRequest, opts ...grpc.CallOption) (*StringResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StringResponse)
	err := c.cc.Invoke(ctx, GitService_GetCommits_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GitServiceServer is the server API for GitService service.
// All implementations must embed UnimplementedGitServiceServer
// for forward compatibility.
type GitServiceServer interface {
	// Repo table methods
	GetDependentRepositoriesByRepo(context.Context, *GetUsernameRequest) (*RepoIntMapResponse, error)
	GetStarsByRepo(context.Context, *GetUsernameRequest) (*RepoIntMapResponse, error)
	GetForksByRepo(context.Context, *GetUsernameRequest) (*RepoIntMapResponse, error)
	GetTotalIssuesByRepo(context.Context, *GetUsernameRequest) (*RepoIntMapResponse, error)
	GetUserSolvedIssuesByRepo(context.Context, *GetUsernameRequest) (*RepoIntMapResponse, error)
	GetTotalPullRequestsByRepo(context.Context, *GetUsernameRequest) (*RepoIntMapResponse, error)
	GetUserMergedPullRequestsByRepo(context.Context, *GetUsernameRequest) (*RepoIntMapResponse, error)
	GetTotalCodeReviewsByRepo(context.Context, *GetUsernameRequest) (*RepoIntMapResponse, error)
	GetUserCodeReviewsByRepo(context.Context, *GetUsernameRequest) (*RepoIntMapResponse, error)
	GetLineChangesByRepo(context.Context, *GetUsernameRequest) (*RepoLineChangesResponse, error)
	// User table methods
	GetName(context.Context, *GetUsernameRequest) (*StringResponse, error)
	GetCompany(context.Context, *GetUsernameRequest) (*StringResponse, error)
	GetLocation(context.Context, *GetUsernameRequest) (*StringResponse, error)
	GetEmail(context.Context, *GetUsernameRequest) (*StringResponse, error)
	GetBio(context.Context, *GetUsernameRequest) (*StringResponse, error)
	GetOrganizations(context.Context, *GetUsernameRequest) (*OrgListResponse, error)
	GetFollowers(context.Context, *GetUsernameRequest) (*IntResponse, error)
	GetReadme(context.Context, *GetReadmeRequest) (*StringResponse, error)
	GetCommits(context.Context, *GetCommitsRequest) (*StringResponse, error)
	mustEmbedUnimplementedGitServiceServer()
}

// UnimplementedGitServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGitServiceServer struct{}

func (UnimplementedGitServiceServer) GetDependentRepositoriesByRepo(context.Context, *GetUsernameRequest) (*RepoIntMapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDependentRepositoriesByRepo not implemented")
}
func (UnimplementedGitServiceServer) GetStarsByRepo(context.Context, *GetUsernameRequest) (*RepoIntMapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStarsByRepo not implemented")
}
func (UnimplementedGitServiceServer) GetForksByRepo(context.Context, *GetUsernameRequest) (*RepoIntMapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetForksByRepo not implemented")
}
func (UnimplementedGitServiceServer) GetTotalIssuesByRepo(context.Context, *GetUsernameRequest) (*RepoIntMapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTotalIssuesByRepo not implemented")
}
func (UnimplementedGitServiceServer) GetUserSolvedIssuesByRepo(context.Context, *GetUsernameRequest) (*RepoIntMapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserSolvedIssuesByRepo not implemented")
}
func (UnimplementedGitServiceServer) GetTotalPullRequestsByRepo(context.Context, *GetUsernameRequest) (*RepoIntMapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTotalPullRequestsByRepo not implemented")
}
func (UnimplementedGitServiceServer) GetUserMergedPullRequestsByRepo(context.Context, *GetUsernameRequest) (*RepoIntMapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserMergedPullRequestsByRepo not implemented")
}
func (UnimplementedGitServiceServer) GetTotalCodeReviewsByRepo(context.Context, *GetUsernameRequest) (*RepoIntMapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTotalCodeReviewsByRepo not implemented")
}
func (UnimplementedGitServiceServer) GetUserCodeReviewsByRepo(context.Context, *GetUsernameRequest) (*RepoIntMapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserCodeReviewsByRepo not implemented")
}
func (UnimplementedGitServiceServer) GetLineChangesByRepo(context.Context, *GetUsernameRequest) (*RepoLineChangesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLineChangesByRepo not implemented")
}
func (UnimplementedGitServiceServer) GetName(context.Context, *GetUsernameRequest) (*StringResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetName not implemented")
}
func (UnimplementedGitServiceServer) GetCompany(context.Context, *GetUsernameRequest) (*StringResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCompany not implemented")
}
func (UnimplementedGitServiceServer) GetLocation(context.Context, *GetUsernameRequest) (*StringResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLocation not implemented")
}
func (UnimplementedGitServiceServer) GetEmail(context.Context, *GetUsernameRequest) (*StringResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEmail not implemented")
}
func (UnimplementedGitServiceServer) GetBio(context.Context, *GetUsernameRequest) (*StringResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBio not implemented")
}
func (UnimplementedGitServiceServer) GetOrganizations(context.Context, *GetUsernameRequest) (*OrgListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrganizations not implemented")
}
func (UnimplementedGitServiceServer) GetFollowers(context.Context, *GetUsernameRequest) (*IntResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFollowers not implemented")
}
func (UnimplementedGitServiceServer) GetReadme(context.Context, *GetReadmeRequest) (*StringResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReadme not implemented")
}
func (UnimplementedGitServiceServer) GetCommits(context.Context, *GetCommitsRequest) (*StringResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCommits not implemented")
}
func (UnimplementedGitServiceServer) mustEmbedUnimplementedGitServiceServer() {}
func (UnimplementedGitServiceServer) testEmbeddedByValue()                    {}

// UnsafeGitServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GitServiceServer will
// result in compilation errors.
type UnsafeGitServiceServer interface {
	mustEmbedUnimplementedGitServiceServer()
}

func RegisterGitServiceServer(s grpc.ServiceRegistrar, srv GitServiceServer) {
	// If the following call pancis, it indicates UnimplementedGitServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&GitService_ServiceDesc, srv)
}

func _GitService_GetDependentRepositoriesByRepo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitServiceServer).GetDependentRepositoriesByRepo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitService_GetDependentRepositoriesByRepo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitServiceServer).GetDependentRepositoriesByRepo(ctx, req.(*GetUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitService_GetStarsByRepo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitServiceServer).GetStarsByRepo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitService_GetStarsByRepo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitServiceServer).GetStarsByRepo(ctx, req.(*GetUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitService_GetForksByRepo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitServiceServer).GetForksByRepo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitService_GetForksByRepo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitServiceServer).GetForksByRepo(ctx, req.(*GetUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitService_GetTotalIssuesByRepo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitServiceServer).GetTotalIssuesByRepo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitService_GetTotalIssuesByRepo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitServiceServer).GetTotalIssuesByRepo(ctx, req.(*GetUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitService_GetUserSolvedIssuesByRepo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitServiceServer).GetUserSolvedIssuesByRepo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitService_GetUserSolvedIssuesByRepo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitServiceServer).GetUserSolvedIssuesByRepo(ctx, req.(*GetUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitService_GetTotalPullRequestsByRepo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitServiceServer).GetTotalPullRequestsByRepo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitService_GetTotalPullRequestsByRepo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitServiceServer).GetTotalPullRequestsByRepo(ctx, req.(*GetUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitService_GetUserMergedPullRequestsByRepo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitServiceServer).GetUserMergedPullRequestsByRepo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitService_GetUserMergedPullRequestsByRepo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitServiceServer).GetUserMergedPullRequestsByRepo(ctx, req.(*GetUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitService_GetTotalCodeReviewsByRepo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitServiceServer).GetTotalCodeReviewsByRepo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitService_GetTotalCodeReviewsByRepo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitServiceServer).GetTotalCodeReviewsByRepo(ctx, req.(*GetUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitService_GetUserCodeReviewsByRepo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitServiceServer).GetUserCodeReviewsByRepo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitService_GetUserCodeReviewsByRepo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitServiceServer).GetUserCodeReviewsByRepo(ctx, req.(*GetUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitService_GetLineChangesByRepo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitServiceServer).GetLineChangesByRepo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitService_GetLineChangesByRepo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitServiceServer).GetLineChangesByRepo(ctx, req.(*GetUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitService_GetName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitServiceServer).GetName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitService_GetName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitServiceServer).GetName(ctx, req.(*GetUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitService_GetCompany_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitServiceServer).GetCompany(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitService_GetCompany_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitServiceServer).GetCompany(ctx, req.(*GetUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitService_GetLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitServiceServer).GetLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitService_GetLocation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitServiceServer).GetLocation(ctx, req.(*GetUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitService_GetEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitServiceServer).GetEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitService_GetEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitServiceServer).GetEmail(ctx, req.(*GetUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitService_GetBio_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitServiceServer).GetBio(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitService_GetBio_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitServiceServer).GetBio(ctx, req.(*GetUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitService_GetOrganizations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitServiceServer).GetOrganizations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitService_GetOrganizations_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitServiceServer).GetOrganizations(ctx, req.(*GetUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitService_GetFollowers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsernameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitServiceServer).GetFollowers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitService_GetFollowers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitServiceServer).GetFollowers(ctx, req.(*GetUsernameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitService_GetReadme_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReadmeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitServiceServer).GetReadme(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitService_GetReadme_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitServiceServer).GetReadme(ctx, req.(*GetReadmeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GitService_GetCommits_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCommitsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GitServiceServer).GetCommits(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GitService_GetCommits_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GitServiceServer).GetCommits(ctx, req.(*GetCommitsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GitService_ServiceDesc is the grpc.ServiceDesc for GitService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GitService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "git.GitService",
	HandlerType: (*GitServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDependentRepositoriesByRepo",
			Handler:    _GitService_GetDependentRepositoriesByRepo_Handler,
		},
		{
			MethodName: "GetStarsByRepo",
			Handler:    _GitService_GetStarsByRepo_Handler,
		},
		{
			MethodName: "GetForksByRepo",
			Handler:    _GitService_GetForksByRepo_Handler,
		},
		{
			MethodName: "GetTotalIssuesByRepo",
			Handler:    _GitService_GetTotalIssuesByRepo_Handler,
		},
		{
			MethodName: "GetUserSolvedIssuesByRepo",
			Handler:    _GitService_GetUserSolvedIssuesByRepo_Handler,
		},
		{
			MethodName: "GetTotalPullRequestsByRepo",
			Handler:    _GitService_GetTotalPullRequestsByRepo_Handler,
		},
		{
			MethodName: "GetUserMergedPullRequestsByRepo",
			Handler:    _GitService_GetUserMergedPullRequestsByRepo_Handler,
		},
		{
			MethodName: "GetTotalCodeReviewsByRepo",
			Handler:    _GitService_GetTotalCodeReviewsByRepo_Handler,
		},
		{
			MethodName: "GetUserCodeReviewsByRepo",
			Handler:    _GitService_GetUserCodeReviewsByRepo_Handler,
		},
		{
			MethodName: "GetLineChangesByRepo",
			Handler:    _GitService_GetLineChangesByRepo_Handler,
		},
		{
			MethodName: "GetName",
			Handler:    _GitService_GetName_Handler,
		},
		{
			MethodName: "GetCompany",
			Handler:    _GitService_GetCompany_Handler,
		},
		{
			MethodName: "GetLocation",
			Handler:    _GitService_GetLocation_Handler,
		},
		{
			MethodName: "GetEmail",
			Handler:    _GitService_GetEmail_Handler,
		},
		{
			MethodName: "GetBio",
			Handler:    _GitService_GetBio_Handler,
		},
		{
			MethodName: "GetOrganizations",
			Handler:    _GitService_GetOrganizations_Handler,
		},
		{
			MethodName: "GetFollowers",
			Handler:    _GitService_GetFollowers_Handler,
		},
		{
			MethodName: "GetReadme",
			Handler:    _GitService_GetReadme_Handler,
		},
		{
			MethodName: "GetCommits",
			Handler:    _GitService_GetCommits_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/git.proto",
}
