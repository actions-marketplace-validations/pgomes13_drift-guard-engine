package schema

// GRPCSchema is the normalized representation of a protobuf/gRPC schema.
type GRPCSchema struct {
	Services []GRPCService
	Messages []GRPCMessage
}

// GRPCService represents a gRPC service definition.
type GRPCService struct {
	Name string
	RPCs []GRPCRPC
}

// GRPCRPC represents a single RPC method within a service.
type GRPCRPC struct {
	Name              string
	RequestType       string
	ResponseType      string
	ClientStreaming    bool
	ServerStreaming    bool
}

// GRPCMessage represents a protobuf message definition.
type GRPCMessage struct {
	Name   string
	Fields []GRPCField
}

// GRPCField represents a single field within a protobuf message.
type GRPCField struct {
	Name     string
	Type     string
	Number   int  // field number
	Repeated bool
	Optional bool
}
