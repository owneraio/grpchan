// Package grpchan provides an abstraction for a gRPC transport, called a
// Channel. The channel is more general than the concrete *grpc.ClientConn
// and *grpc.Server types provided by gRPC. It allows gRPC over alternate
// substrates and includes sub-packages that provide two such alternatives:
// in-process and HTTP 1.1.
//
// The key types in this package are interfaces that are helpful abstractions
// and expose the same key operations as *grpc.ClientConn and *grpc.Server.
// They are grpchan.Channel and grpchan.ServiceRegistry, respectively. To
// really make use of these interfaces, you also need to use the grpchan
// plugin for protoc. It will generate methods for each RPC service that
// allow you to create client stubs and register server handlers with these
// interfaces. (The code generated by the standard Go plugin for protoc does
// not suffice since it requires the concrete types *grpc.ClientConn and
// *grpc.Server.)
//
// Protoc Plugin
//
// To use the protoc plugin, you need to first build it and make sure its
// location is in your PATH.
//
//   go install github.com/fullstorydev/grpchan/cmd/protoc-gen-grpchan
//   # If necessary, make sure its location is on your path like so:
//   #   export PATH=$PATH:$GOPATH/bin
//
// When you invoke protoc, include a --grpchan_out parameter that indicates
// the same output directory as used for your --go_out parameter. Alongside
// the *.pb.go files generated, the grpchan plugin will also create
// *.pb.grpchan.go files.
//
// The normal Go plugin (when emitted gRPC code) creates client stub factory
// and server registration functions for each RPC service defined in the proto
// source files compiled.
//
//    // A function with this signature is generated, for creating a client
//    // stub that uses the given connection to issue requests.
//    func New<ServerName>Client(cc *grpc.ClientConn) <ServerName>Client {
//        return &<serverName>Client{cc}
//    }
//
//    // And a function with this signature is generated, for registering
//    // server handlers with the given server.
//    func Register<ServerName>Server(s *grpc.Server, srv <ServerName>Server) {
//        s.RegisterService(&_<ServerName>_serviceDesc, srv)
//    }
//
// The grpchan plugins produces similar named methods that accept interfaces:
//
//    func New<ServerName>ChannelClient(ch grpchan.Channel) <ServerName>Client {
//        return &<serverName>ChannelClient{ch}
//    }
//
//    func RegisterHandler<ServerName>(sr grpchan.ServiceRegistry, srv <ServerName>Server) {
//        s.RegisterService(&_<ServerName>_serviceDesc, srv)
//    }
//
// A new transport can then be implemented by just implementing these two
// interfaces, grpchan.Channel for the client side and grpchan.ServiceRegistry
// for the server side.
//
// These alternate methods also work just fine with *grpc.ClientConn and
// *grpc.Server as they implement the necessary interfaces.
//
// Client-Side Channels
//
// The client-side implementation of a transport is done with just the two
// methods in the Channel interface: one for unary RPCs and the other for
// streaming RPCs.
//
// Note that when a unary interceptor is invoked for an RPC on a channel that
// is *not* a *grpc.ClientConn, the parameter of that type will be nil.
//
// Not all client call options will make sense for all transports. This repo
// chooses to ignore call options that do not apply (as opposed to failing
// the RPC or panic'ing). However, several call options are likely important
// to support: those for accessing header and trailer metadata. The peer,
// per-RPC credentials, and message size limits are other options that are
// reasonably straight-forward to apply to other transports. But the other
// options (dealing with on-the-wire encoding, compression, etc) are less
// likely to be meaningful.
//
// Server-Side Service Registries
//
// The server-side implementation of a transport must be able to invoke
// method and stream handlers for a given service implementation. This is done
// by implementing the ServiceRegistry interface. When a service is registered,
// a service description is provided that includes access to method and stream
// handlers. When the transport receives requests for RPC operations, it in
// turn invokes these handlers. For streaming operations, it must also supply
// a grpc.ServerStream implementation, for exchanging messages on the stream.
//
// Note that the server stream's context will need a custom implementation of
// the grpc.ServerTransportStream in it, too. Sadly, this interface is just
// different enough from grpc.ServerStream that they cannot be implemented by
// the same type. This is particularly necessary for unary calls since this is
// how a unary handler dictates what headers and trailers to send back to the
// client.
package grpchan
