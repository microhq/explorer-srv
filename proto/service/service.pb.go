// Code generated by protoc-gen-go.
// source: github.com/micro/explorer-srv/proto/service/service.proto
// DO NOT EDIT!

/*
Package service is a generated protocol buffer package.

It is generated from these files:
	github.com/micro/explorer-srv/proto/service/service.proto

It has these top-level messages:
	Service
	Version
	API
	Endpoint
	Dependency
	Source
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
	CreateVersionRequest
	CreateVersionResponse
	DeleteVersionRequest
	DeleteVersionResponse
	ReadVersionRequest
	ReadVersionResponse
	UpdateVersionRequest
	UpdateVersionResponse
	SearchVersionRequest
	SearchVersionResponse
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
	Url         string            `protobuf:"bytes,7,opt,name=url" json:"url,omitempty"`
	Readme      string            `protobuf:"bytes,8,opt,name=readme" json:"readme,omitempty"`
	Metadata    map[string]string `protobuf:"bytes,9,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
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

type Version struct {
	Id           string            `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	ServiceId    string            `protobuf:"bytes,2,opt,name=serviceId" json:"serviceId,omitempty"`
	Version      string            `protobuf:"bytes,3,opt,name=version" json:"version,omitempty"`
	Created      int64             `protobuf:"varint,4,opt,name=created" json:"created,omitempty"`
	Updated      int64             `protobuf:"varint,5,opt,name=updated" json:"updated,omitempty"`
	Api          *API              `protobuf:"bytes,6,opt,name=api" json:"api,omitempty"`
	Sources      []*Source         `protobuf:"bytes,7,rep,name=sources" json:"sources,omitempty"`
	Dependencies []*Dependency     `protobuf:"bytes,8,rep,name=dependencies" json:"dependencies,omitempty"`
	Metadata     map[string]string `protobuf:"bytes,9,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Version) Reset()         { *m = Version{} }
func (m *Version) String() string { return proto.CompactTextString(m) }
func (*Version) ProtoMessage()    {}

func (m *Version) GetApi() *API {
	if m != nil {
		return m.Api
	}
	return nil
}

func (m *Version) GetSources() []*Source {
	if m != nil {
		return m.Sources
	}
	return nil
}

func (m *Version) GetDependencies() []*Dependency {
	if m != nil {
		return m.Dependencies
	}
	return nil
}

func (m *Version) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type API struct {
	Endpoints []*Endpoint       `protobuf:"bytes,1,rep,name=endpoints" json:"endpoints,omitempty"`
	Metadata  map[string]string `protobuf:"bytes,2,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *API) Reset()         { *m = API{} }
func (m *API) String() string { return proto.CompactTextString(m) }
func (*API) ProtoMessage()    {}

func (m *API) GetEndpoints() []*Endpoint {
	if m != nil {
		return m.Endpoints
	}
	return nil
}

func (m *API) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type Endpoint struct {
	Name     string            `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Request  map[string]string `protobuf:"bytes,2,rep,name=request" json:"request,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Response map[string]string `protobuf:"bytes,3,rep,name=response" json:"response,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Metadata map[string]string `protobuf:"bytes,4,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Endpoint) Reset()         { *m = Endpoint{} }
func (m *Endpoint) String() string { return proto.CompactTextString(m) }
func (*Endpoint) ProtoMessage()    {}

func (m *Endpoint) GetRequest() map[string]string {
	if m != nil {
		return m.Request
	}
	return nil
}

func (m *Endpoint) GetResponse() map[string]string {
	if m != nil {
		return m.Response
	}
	return nil
}

func (m *Endpoint) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type Dependency struct {
	Name     string            `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Type     string            `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Metadata map[string]string `protobuf:"bytes,3,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Dependency) Reset()         { *m = Dependency{} }
func (m *Dependency) String() string { return proto.CompactTextString(m) }
func (*Dependency) ProtoMessage()    {}

func (m *Dependency) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type Source struct {
	Name     string            `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Type     string            `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Metadata map[string]string `protobuf:"bytes,3,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Source) Reset()         { *m = Source{} }
func (m *Source) String() string { return proto.CompactTextString(m) }
func (*Source) ProtoMessage()    {}

func (m *Source) GetMetadata() map[string]string {
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
	Order  string `protobuf:"bytes,5,opt,name=order" json:"order,omitempty"`
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

// Version request/response
type CreateVersionRequest struct {
	Version *Version `protobuf:"bytes,1,opt,name=version" json:"version,omitempty"`
}

func (m *CreateVersionRequest) Reset()         { *m = CreateVersionRequest{} }
func (m *CreateVersionRequest) String() string { return proto.CompactTextString(m) }
func (*CreateVersionRequest) ProtoMessage()    {}

func (m *CreateVersionRequest) GetVersion() *Version {
	if m != nil {
		return m.Version
	}
	return nil
}

type CreateVersionResponse struct {
}

func (m *CreateVersionResponse) Reset()         { *m = CreateVersionResponse{} }
func (m *CreateVersionResponse) String() string { return proto.CompactTextString(m) }
func (*CreateVersionResponse) ProtoMessage()    {}

type DeleteVersionRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *DeleteVersionRequest) Reset()         { *m = DeleteVersionRequest{} }
func (m *DeleteVersionRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteVersionRequest) ProtoMessage()    {}

type DeleteVersionResponse struct {
}

func (m *DeleteVersionResponse) Reset()         { *m = DeleteVersionResponse{} }
func (m *DeleteVersionResponse) String() string { return proto.CompactTextString(m) }
func (*DeleteVersionResponse) ProtoMessage()    {}

type ReadVersionRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *ReadVersionRequest) Reset()         { *m = ReadVersionRequest{} }
func (m *ReadVersionRequest) String() string { return proto.CompactTextString(m) }
func (*ReadVersionRequest) ProtoMessage()    {}

type ReadVersionResponse struct {
	Version *Version `protobuf:"bytes,1,opt,name=version" json:"version,omitempty"`
}

func (m *ReadVersionResponse) Reset()         { *m = ReadVersionResponse{} }
func (m *ReadVersionResponse) String() string { return proto.CompactTextString(m) }
func (*ReadVersionResponse) ProtoMessage()    {}

func (m *ReadVersionResponse) GetVersion() *Version {
	if m != nil {
		return m.Version
	}
	return nil
}

type UpdateVersionRequest struct {
	Version *Version `protobuf:"bytes,1,opt,name=version" json:"version,omitempty"`
}

func (m *UpdateVersionRequest) Reset()         { *m = UpdateVersionRequest{} }
func (m *UpdateVersionRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateVersionRequest) ProtoMessage()    {}

func (m *UpdateVersionRequest) GetVersion() *Version {
	if m != nil {
		return m.Version
	}
	return nil
}

type UpdateVersionResponse struct {
}

func (m *UpdateVersionResponse) Reset()         { *m = UpdateVersionResponse{} }
func (m *UpdateVersionResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateVersionResponse) ProtoMessage()    {}

type SearchVersionRequest struct {
	ServiceId string `protobuf:"bytes,1,opt,name=serviceId" json:"serviceId,omitempty"`
	Version   string `protobuf:"bytes,2,opt,name=version" json:"version,omitempty"`
	Limit     int64  `protobuf:"varint,3,opt,name=limit" json:"limit,omitempty"`
	Offset    int64  `protobuf:"varint,4,opt,name=offset" json:"offset,omitempty"`
}

func (m *SearchVersionRequest) Reset()         { *m = SearchVersionRequest{} }
func (m *SearchVersionRequest) String() string { return proto.CompactTextString(m) }
func (*SearchVersionRequest) ProtoMessage()    {}

type SearchVersionResponse struct {
	Versions []*Version `protobuf:"bytes,1,rep,name=versions" json:"versions,omitempty"`
}

func (m *SearchVersionResponse) Reset()         { *m = SearchVersionResponse{} }
func (m *SearchVersionResponse) String() string { return proto.CompactTextString(m) }
func (*SearchVersionResponse) ProtoMessage()    {}

func (m *SearchVersionResponse) GetVersions() []*Version {
	if m != nil {
		return m.Versions
	}
	return nil
}

func init() {
}
