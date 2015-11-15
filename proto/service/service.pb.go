// Code generated by protoc-gen-go.
// source: github.com/myodc/explorer-srv/proto/service/service.proto
// DO NOT EDIT!

/*
Package service is a generated protocol buffer package.

It is generated from these files:
	github.com/myodc/explorer-srv/proto/service/service.proto

It has these top-level messages:
	Service
	CreateRequest
	CreateResponse
	DeleteRequest
	DeleteResponse
	ReadRequest
	ReadResponse
	UpdateRequest
	UpdateResponse
	SearchRequest
	SearchResponse
*/
package service

import proto "github.com/golang/protobuf/proto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

type Service struct {
	Id          string            `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name        string            `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Owner       string            `protobuf:"bytes,3,opt,name=owner" json:"owner,omitempty"`
	Description string            `protobuf:"bytes,4,opt,name=description" json:"description,omitempty"`
	Created     int64             `protobuf:"varint,5,opt,name=created" json:"created,omitempty"`
	Updated     int64             `protobuf:"varint,6,opt,name=updated" json:"updated,omitempty"`
	Metadata    map[string]string `protobuf:"bytes,7,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Service) Reset()         { *m = Service{} }
func (m *Service) String() string { return proto.CompactTextString(m) }
func (*Service) ProtoMessage()    {}

func (m *Service) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type CreateRequest struct {
	Service *Service `protobuf:"bytes,1,opt,name=service" json:"service,omitempty"`
}

func (m *CreateRequest) Reset()         { *m = CreateRequest{} }
func (m *CreateRequest) String() string { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()    {}

func (m *CreateRequest) GetService() *Service {
	if m != nil {
		return m.Service
	}
	return nil
}

type CreateResponse struct {
}

func (m *CreateResponse) Reset()         { *m = CreateResponse{} }
func (m *CreateResponse) String() string { return proto.CompactTextString(m) }
func (*CreateResponse) ProtoMessage()    {}

type DeleteRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *DeleteRequest) Reset()         { *m = DeleteRequest{} }
func (m *DeleteRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteRequest) ProtoMessage()    {}

type DeleteResponse struct {
}

func (m *DeleteResponse) Reset()         { *m = DeleteResponse{} }
func (m *DeleteResponse) String() string { return proto.CompactTextString(m) }
func (*DeleteResponse) ProtoMessage()    {}

type ReadRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *ReadRequest) Reset()         { *m = ReadRequest{} }
func (m *ReadRequest) String() string { return proto.CompactTextString(m) }
func (*ReadRequest) ProtoMessage()    {}

type ReadResponse struct {
	Service *Service `protobuf:"bytes,1,opt,name=service" json:"service,omitempty"`
}

func (m *ReadResponse) Reset()         { *m = ReadResponse{} }
func (m *ReadResponse) String() string { return proto.CompactTextString(m) }
func (*ReadResponse) ProtoMessage()    {}

func (m *ReadResponse) GetService() *Service {
	if m != nil {
		return m.Service
	}
	return nil
}

type UpdateRequest struct {
	Service *Service `protobuf:"bytes,1,opt,name=service" json:"service,omitempty"`
}

func (m *UpdateRequest) Reset()         { *m = UpdateRequest{} }
func (m *UpdateRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateRequest) ProtoMessage()    {}

func (m *UpdateRequest) GetService() *Service {
	if m != nil {
		return m.Service
	}
	return nil
}

type UpdateResponse struct {
}

func (m *UpdateResponse) Reset()         { *m = UpdateResponse{} }
func (m *UpdateResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateResponse) ProtoMessage()    {}

type SearchRequest struct {
	Name   string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Owner  string `protobuf:"bytes,2,opt,name=owner" json:"owner,omitempty"`
	Limit  int64  `protobuf:"varint,3,opt,name=limit" json:"limit,omitempty"`
	Offset int64  `protobuf:"varint,4,opt,name=offset" json:"offset,omitempty"`
}

func (m *SearchRequest) Reset()         { *m = SearchRequest{} }
func (m *SearchRequest) String() string { return proto.CompactTextString(m) }
func (*SearchRequest) ProtoMessage()    {}

type SearchResponse struct {
	Services []*Service `protobuf:"bytes,1,rep,name=services" json:"services,omitempty"`
}

func (m *SearchResponse) Reset()         { *m = SearchResponse{} }
func (m *SearchResponse) String() string { return proto.CompactTextString(m) }
func (*SearchResponse) ProtoMessage()    {}

func (m *SearchResponse) GetServices() []*Service {
	if m != nil {
		return m.Services
	}
	return nil
}

func init() {
}
