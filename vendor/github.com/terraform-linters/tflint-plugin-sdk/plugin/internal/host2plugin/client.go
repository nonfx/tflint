package host2plugin

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-version"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin/internal/fromproto"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin/internal/plugin2host"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin/internal/proto"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin/internal/toproto"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// GRPCClient is a host-side implementation. Host can send requests through the client to plugin's gRPC server.
type GRPCClient struct {
	client proto.RuleSetClient
}

// ClientOpts is an option for initializing a Client.
type ClientOpts struct {
	// Pass server options in the client to bypass grpc
	ServeOpts *[]ServeOpts
}

func NewClient(opts *ClientOpts) (*GRPCClient, error) {
	client := proto.NewRuleSetClient(new(clientConn))

	return &GRPCClient{
		client: client,
	}, nil
}

// RuleSetName returns the name of a plugin.
func (c *GRPCClient) RuleSetName() (string, error) {
	resp, err := c.client.GetName(context.Background(), &proto.GetName_Request{})
	if err != nil {
		return "", fromproto.Error(err)
	}
	return resp.Name, nil
}

// RuleSetVersion returns the version of a plugin.
func (c *GRPCClient) RuleSetVersion() (string, error) {
	resp, err := c.client.GetVersion(context.Background(), &proto.GetVersion_Request{})
	if err != nil {
		return "", fromproto.Error(err)
	}
	return resp.Version, nil
}

// RuleNames returns the list of rule names provided by a plugin.
func (c *GRPCClient) RuleNames() ([]string, error) {
	resp, err := c.client.GetRuleNames(context.Background(), &proto.GetRuleNames_Request{})
	if err != nil {
		return []string{}, fromproto.Error(err)
	}
	return resp.Names, nil
}

// VersionConstraints returns constraints of TFLint versions.
func (c *GRPCClient) VersionConstraints() (version.Constraints, error) {
	resp, err := c.client.GetVersionConstraint(context.Background(), &proto.GetVersionConstraint_Request{})
	if err != nil {
		return nil, fromproto.Error(err)
	}

	if resp.Constraint == "" {
		return version.Constraints{}, nil
	}
	return version.NewConstraint(resp.Constraint)
}

// SDKVersion returns the SDK version.
func (c *GRPCClient) SDKVersion() (*version.Version, error) {
	resp, err := c.client.GetSDKVersion(context.Background(), &proto.GetSDKVersion_Request{})
	if err != nil {
		return nil, fromproto.Error(err)
	}
	return version.NewVersion(resp.Version)
}

// ConfigSchema fetches the config schema from a plugin.
func (c *GRPCClient) ConfigSchema() (*hclext.BodySchema, error) {
	resp, err := c.client.GetConfigSchema(context.Background(), &proto.GetConfigSchema_Request{})
	if err != nil {
		return nil, fromproto.Error(err)
	}
	return fromproto.BodySchema(resp.Schema), nil
}

// ApplyGlobalConfig applies a common config to a plugin.
func (c *GRPCClient) ApplyGlobalConfig(config *tflint.Config) error {
	_, err := c.client.ApplyGlobalConfig(context.Background(), &proto.ApplyGlobalConfig_Request{Config: toproto.Config(config)})
	if err != nil {
		return fromproto.Error(err)
	}
	return nil
}

// ApplyConfig applies the config to a plugin.
func (c *GRPCClient) ApplyConfig(content *hclext.BodyContent, sources map[string][]byte) error {
	_, err := c.client.ApplyConfig(context.Background(), &proto.ApplyConfig_Request{Content: toproto.BodyContent(content, sources)})
	if err != nil {
		return fromproto.Error(err)
	}
	return nil
}

// Check calls its own plugin implementation with an gRPC client that can send
// requests to the host process.
func (c *GRPCClient) Check(runner plugin2host.Server) error {
	// brokerID := c.broker.NextId()
	logger.Debug("starting host-side gRPC server")
	// go c.broker.AcceptAndServe(brokerID, func(opts []grpc.ServerOption) *grpc.Server {
	// 	opts = append(opts, grpc.UnaryInterceptor(interceptor.RequestLogging("plugin2host")))
	// 	server := grpc.NewServer(opts...)
	// 	proto.RegisterRunnerServer(server, &plugin2host.GRPCServer{Impl: runner})
	// 	return server
	// })

	// _, err := c.client.Check(context.Background(), &proto.Check_Request{Runner: brokerID})

	// if err != nil {
	// 	return fromproto.Error(err)
	// }
	return nil
}

// clientConn is a basic implementation of the ClientConnInterface.
type clientConn struct {
	conn *grpc.ClientConn
}

// NewClientConn creates a new clientConn with the given gRPC ClientConn.
func NewClientConn(conn *grpc.ClientConn) grpc.ClientConnInterface {
	return &clientConn{conn: conn}
}

// Invoke performs a unary RPC and returns after the response is received into reply.
func (c *clientConn) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	fmt.Println("Invoke", method, args, reply)
	return c.conn.Invoke(ctx, method, args, reply, opts...)
}

// NewStream begins a streaming RPC.
func (c *clientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	fmt.Println("NewStream", desc, method, opts)
	return c.conn.NewStream(ctx, desc, method, opts...)
}
