// Code generated by protoc-gen-go. DO NOT EDIT.
// source: vitepb/message.proto

package vitepb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type StatusMsg struct {
	NetID                uint64   `protobuf:"varint,1,opt,name=NetID,proto3" json:"NetID,omitempty"`
	Version              uint64   `protobuf:"varint,2,opt,name=Version,proto3" json:"Version,omitempty"`
	Height               uint64   `protobuf:"varint,3,opt,name=Height,proto3" json:"Height,omitempty"`
	CurrentBlock         []byte   `protobuf:"bytes,4,opt,name=CurrentBlock,proto3" json:"CurrentBlock,omitempty"`
	GenesisBlock         []byte   `protobuf:"bytes,5,opt,name=GenesisBlock,proto3" json:"GenesisBlock,omitempty"`
	Port                 uint32   `protobuf:"varint,6,opt,name=Port,proto3" json:"Port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusMsg) Reset()         { *m = StatusMsg{} }
func (m *StatusMsg) String() string { return proto.CompactTextString(m) }
func (*StatusMsg) ProtoMessage()    {}
func (*StatusMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{0}
}

func (m *StatusMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusMsg.Unmarshal(m, b)
}
func (m *StatusMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusMsg.Marshal(b, m, deterministic)
}
func (m *StatusMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusMsg.Merge(m, src)
}
func (m *StatusMsg) XXX_Size() int {
	return xxx_messageInfo_StatusMsg.Size(m)
}
func (m *StatusMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusMsg.DiscardUnknown(m)
}

var xxx_messageInfo_StatusMsg proto.InternalMessageInfo

func (m *StatusMsg) GetNetID() uint64 {
	if m != nil {
		return m.NetID
	}
	return 0
}

func (m *StatusMsg) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *StatusMsg) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *StatusMsg) GetCurrentBlock() []byte {
	if m != nil {
		return m.CurrentBlock
	}
	return nil
}

func (m *StatusMsg) GetGenesisBlock() []byte {
	if m != nil {
		return m.GenesisBlock
	}
	return nil
}

func (m *StatusMsg) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

type GetSnapshotBlocksMsg struct {
	Origin               []byte   `protobuf:"bytes,1,opt,name=Origin,proto3" json:"Origin,omitempty"`
	Count                uint64   `protobuf:"varint,2,opt,name=Count,proto3" json:"Count,omitempty"`
	Forward              bool     `protobuf:"varint,3,opt,name=Forward,proto3" json:"Forward,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSnapshotBlocksMsg) Reset()         { *m = GetSnapshotBlocksMsg{} }
func (m *GetSnapshotBlocksMsg) String() string { return proto.CompactTextString(m) }
func (*GetSnapshotBlocksMsg) ProtoMessage()    {}
func (*GetSnapshotBlocksMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{1}
}

func (m *GetSnapshotBlocksMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSnapshotBlocksMsg.Unmarshal(m, b)
}
func (m *GetSnapshotBlocksMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSnapshotBlocksMsg.Marshal(b, m, deterministic)
}
func (m *GetSnapshotBlocksMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSnapshotBlocksMsg.Merge(m, src)
}
func (m *GetSnapshotBlocksMsg) XXX_Size() int {
	return xxx_messageInfo_GetSnapshotBlocksMsg.Size(m)
}
func (m *GetSnapshotBlocksMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSnapshotBlocksMsg.DiscardUnknown(m)
}

var xxx_messageInfo_GetSnapshotBlocksMsg proto.InternalMessageInfo

func (m *GetSnapshotBlocksMsg) GetOrigin() []byte {
	if m != nil {
		return m.Origin
	}
	return nil
}

func (m *GetSnapshotBlocksMsg) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *GetSnapshotBlocksMsg) GetForward() bool {
	if m != nil {
		return m.Forward
	}
	return false
}

type SnapshotBlocksMsg struct {
	Blocks               []*SnapshotBlock `protobuf:"bytes,1,rep,name=blocks,proto3" json:"blocks,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *SnapshotBlocksMsg) Reset()         { *m = SnapshotBlocksMsg{} }
func (m *SnapshotBlocksMsg) String() string { return proto.CompactTextString(m) }
func (*SnapshotBlocksMsg) ProtoMessage()    {}
func (*SnapshotBlocksMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{2}
}

func (m *SnapshotBlocksMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SnapshotBlocksMsg.Unmarshal(m, b)
}
func (m *SnapshotBlocksMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SnapshotBlocksMsg.Marshal(b, m, deterministic)
}
func (m *SnapshotBlocksMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SnapshotBlocksMsg.Merge(m, src)
}
func (m *SnapshotBlocksMsg) XXX_Size() int {
	return xxx_messageInfo_SnapshotBlocksMsg.Size(m)
}
func (m *SnapshotBlocksMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_SnapshotBlocksMsg.DiscardUnknown(m)
}

var xxx_messageInfo_SnapshotBlocksMsg proto.InternalMessageInfo

func (m *SnapshotBlocksMsg) GetBlocks() []*SnapshotBlock {
	if m != nil {
		return m.Blocks
	}
	return nil
}

type GetAccountBlocksMsg struct {
	Origin               []byte   `protobuf:"bytes,1,opt,name=Origin,proto3" json:"Origin,omitempty"`
	Count                uint64   `protobuf:"varint,2,opt,name=Count,proto3" json:"Count,omitempty"`
	Forward              bool     `protobuf:"varint,3,opt,name=Forward,proto3" json:"Forward,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAccountBlocksMsg) Reset()         { *m = GetAccountBlocksMsg{} }
func (m *GetAccountBlocksMsg) String() string { return proto.CompactTextString(m) }
func (*GetAccountBlocksMsg) ProtoMessage()    {}
func (*GetAccountBlocksMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{3}
}

func (m *GetAccountBlocksMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAccountBlocksMsg.Unmarshal(m, b)
}
func (m *GetAccountBlocksMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAccountBlocksMsg.Marshal(b, m, deterministic)
}
func (m *GetAccountBlocksMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAccountBlocksMsg.Merge(m, src)
}
func (m *GetAccountBlocksMsg) XXX_Size() int {
	return xxx_messageInfo_GetAccountBlocksMsg.Size(m)
}
func (m *GetAccountBlocksMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAccountBlocksMsg.DiscardUnknown(m)
}

var xxx_messageInfo_GetAccountBlocksMsg proto.InternalMessageInfo

func (m *GetAccountBlocksMsg) GetOrigin() []byte {
	if m != nil {
		return m.Origin
	}
	return nil
}

func (m *GetAccountBlocksMsg) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *GetAccountBlocksMsg) GetForward() bool {
	if m != nil {
		return m.Forward
	}
	return false
}

type AccountBlocksMsg struct {
	Blocks               []*AccountBlock `protobuf:"bytes,3,rep,name=blocks,proto3" json:"blocks,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *AccountBlocksMsg) Reset()         { *m = AccountBlocksMsg{} }
func (m *AccountBlocksMsg) String() string { return proto.CompactTextString(m) }
func (*AccountBlocksMsg) ProtoMessage()    {}
func (*AccountBlocksMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{4}
}

func (m *AccountBlocksMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountBlocksMsg.Unmarshal(m, b)
}
func (m *AccountBlocksMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountBlocksMsg.Marshal(b, m, deterministic)
}
func (m *AccountBlocksMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountBlocksMsg.Merge(m, src)
}
func (m *AccountBlocksMsg) XXX_Size() int {
	return xxx_messageInfo_AccountBlocksMsg.Size(m)
}
func (m *AccountBlocksMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountBlocksMsg.DiscardUnknown(m)
}

var xxx_messageInfo_AccountBlocksMsg proto.InternalMessageInfo

func (m *AccountBlocksMsg) GetBlocks() []*AccountBlock {
	if m != nil {
		return m.Blocks
	}
	return nil
}

// version 2
type BlockID struct {
	Hash                 []byte   `protobuf:"bytes,1,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Height               uint64   `protobuf:"varint,2,opt,name=Height,proto3" json:"Height,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockID) Reset()         { *m = BlockID{} }
func (m *BlockID) String() string { return proto.CompactTextString(m) }
func (*BlockID) ProtoMessage()    {}
func (*BlockID) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{5}
}

func (m *BlockID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockID.Unmarshal(m, b)
}
func (m *BlockID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockID.Marshal(b, m, deterministic)
}
func (m *BlockID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockID.Merge(m, src)
}
func (m *BlockID) XXX_Size() int {
	return xxx_messageInfo_BlockID.Size(m)
}
func (m *BlockID) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockID.DiscardUnknown(m)
}

var xxx_messageInfo_BlockID proto.InternalMessageInfo

func (m *BlockID) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *BlockID) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

type Segment struct {
	From                 *BlockID `protobuf:"bytes,1,opt,name=From,proto3" json:"From,omitempty"`
	To                   *BlockID `protobuf:"bytes,2,opt,name=To,proto3" json:"To,omitempty"`
	Step                 uint64   `protobuf:"varint,3,opt,name=Step,proto3" json:"Step,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Segment) Reset()         { *m = Segment{} }
func (m *Segment) String() string { return proto.CompactTextString(m) }
func (*Segment) ProtoMessage()    {}
func (*Segment) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{6}
}

func (m *Segment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Segment.Unmarshal(m, b)
}
func (m *Segment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Segment.Marshal(b, m, deterministic)
}
func (m *Segment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Segment.Merge(m, src)
}
func (m *Segment) XXX_Size() int {
	return xxx_messageInfo_Segment.Size(m)
}
func (m *Segment) XXX_DiscardUnknown() {
	xxx_messageInfo_Segment.DiscardUnknown(m)
}

var xxx_messageInfo_Segment proto.InternalMessageInfo

func (m *Segment) GetFrom() *BlockID {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *Segment) GetTo() *BlockID {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *Segment) GetStep() uint64 {
	if m != nil {
		return m.Step
	}
	return 0
}

type AccountSegment struct {
	Address              []byte   `protobuf:"bytes,1,opt,name=Address,proto3" json:"Address,omitempty"`
	Segment              *Segment `protobuf:"bytes,2,opt,name=Segment,proto3" json:"Segment,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AccountSegment) Reset()         { *m = AccountSegment{} }
func (m *AccountSegment) String() string { return proto.CompactTextString(m) }
func (*AccountSegment) ProtoMessage()    {}
func (*AccountSegment) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{7}
}

func (m *AccountSegment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountSegment.Unmarshal(m, b)
}
func (m *AccountSegment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountSegment.Marshal(b, m, deterministic)
}
func (m *AccountSegment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountSegment.Merge(m, src)
}
func (m *AccountSegment) XXX_Size() int {
	return xxx_messageInfo_AccountSegment.Size(m)
}
func (m *AccountSegment) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountSegment.DiscardUnknown(m)
}

var xxx_messageInfo_AccountSegment proto.InternalMessageInfo

func (m *AccountSegment) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *AccountSegment) GetSegment() *Segment {
	if m != nil {
		return m.Segment
	}
	return nil
}

type MultiAccountSegment struct {
	Segments             []*AccountSegment `protobuf:"bytes,1,rep,name=Segments,proto3" json:"Segments,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *MultiAccountSegment) Reset()         { *m = MultiAccountSegment{} }
func (m *MultiAccountSegment) String() string { return proto.CompactTextString(m) }
func (*MultiAccountSegment) ProtoMessage()    {}
func (*MultiAccountSegment) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{8}
}

func (m *MultiAccountSegment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MultiAccountSegment.Unmarshal(m, b)
}
func (m *MultiAccountSegment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MultiAccountSegment.Marshal(b, m, deterministic)
}
func (m *MultiAccountSegment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MultiAccountSegment.Merge(m, src)
}
func (m *MultiAccountSegment) XXX_Size() int {
	return xxx_messageInfo_MultiAccountSegment.Size(m)
}
func (m *MultiAccountSegment) XXX_DiscardUnknown() {
	xxx_messageInfo_MultiAccountSegment.DiscardUnknown(m)
}

var xxx_messageInfo_MultiAccountSegment proto.InternalMessageInfo

func (m *MultiAccountSegment) GetSegments() []*AccountSegment {
	if m != nil {
		return m.Segments
	}
	return nil
}

type FileList struct {
	Files                []string `protobuf:"bytes,1,rep,name=Files,proto3" json:"Files,omitempty"`
	Start                uint64   `protobuf:"varint,2,opt,name=Start,proto3" json:"Start,omitempty"`
	End                  uint64   `protobuf:"varint,3,opt,name=End,proto3" json:"End,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FileList) Reset()         { *m = FileList{} }
func (m *FileList) String() string { return proto.CompactTextString(m) }
func (*FileList) ProtoMessage()    {}
func (*FileList) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{9}
}

func (m *FileList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FileList.Unmarshal(m, b)
}
func (m *FileList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FileList.Marshal(b, m, deterministic)
}
func (m *FileList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileList.Merge(m, src)
}
func (m *FileList) XXX_Size() int {
	return xxx_messageInfo_FileList.Size(m)
}
func (m *FileList) XXX_DiscardUnknown() {
	xxx_messageInfo_FileList.DiscardUnknown(m)
}

var xxx_messageInfo_FileList proto.InternalMessageInfo

func (m *FileList) GetFiles() []string {
	if m != nil {
		return m.Files
	}
	return nil
}

func (m *FileList) GetStart() uint64 {
	if m != nil {
		return m.Start
	}
	return 0
}

func (m *FileList) GetEnd() uint64 {
	if m != nil {
		return m.End
	}
	return 0
}

type SubLedger struct {
	SBlocks              []*SnapshotBlock `protobuf:"bytes,1,rep,name=SBlocks,proto3" json:"SBlocks,omitempty"`
	ABlocks              []*AccountBlock  `protobuf:"bytes,2,rep,name=ABlocks,proto3" json:"ABlocks,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *SubLedger) Reset()         { *m = SubLedger{} }
func (m *SubLedger) String() string { return proto.CompactTextString(m) }
func (*SubLedger) ProtoMessage()    {}
func (*SubLedger) Descriptor() ([]byte, []int) {
	return fileDescriptor_2a6a8486deb9ab39, []int{10}
}

func (m *SubLedger) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubLedger.Unmarshal(m, b)
}
func (m *SubLedger) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubLedger.Marshal(b, m, deterministic)
}
func (m *SubLedger) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubLedger.Merge(m, src)
}
func (m *SubLedger) XXX_Size() int {
	return xxx_messageInfo_SubLedger.Size(m)
}
func (m *SubLedger) XXX_DiscardUnknown() {
	xxx_messageInfo_SubLedger.DiscardUnknown(m)
}

var xxx_messageInfo_SubLedger proto.InternalMessageInfo

func (m *SubLedger) GetSBlocks() []*SnapshotBlock {
	if m != nil {
		return m.SBlocks
	}
	return nil
}

func (m *SubLedger) GetABlocks() []*AccountBlock {
	if m != nil {
		return m.ABlocks
	}
	return nil
}

func init() {
	proto.RegisterType((*StatusMsg)(nil), "vitepb.StatusMsg")
	proto.RegisterType((*GetSnapshotBlocksMsg)(nil), "vitepb.GetSnapshotBlocksMsg")
	proto.RegisterType((*SnapshotBlocksMsg)(nil), "vitepb.SnapshotBlocksMsg")
	proto.RegisterType((*GetAccountBlocksMsg)(nil), "vitepb.GetAccountBlocksMsg")
	proto.RegisterType((*AccountBlocksMsg)(nil), "vitepb.AccountBlocksMsg")
	proto.RegisterType((*BlockID)(nil), "vitepb.BlockID")
	proto.RegisterType((*Segment)(nil), "vitepb.Segment")
	proto.RegisterType((*AccountSegment)(nil), "vitepb.AccountSegment")
	proto.RegisterType((*MultiAccountSegment)(nil), "vitepb.MultiAccountSegment")
	proto.RegisterType((*FileList)(nil), "vitepb.FileList")
	proto.RegisterType((*SubLedger)(nil), "vitepb.SubLedger")
}

func init() { proto.RegisterFile("vitepb/message.proto", fileDescriptor_2a6a8486deb9ab39) }

var fileDescriptor_2a6a8486deb9ab39 = []byte{
	// 495 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x53, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0x95, 0x9d, 0xd4, 0x6e, 0xa6, 0x01, 0xca, 0x36, 0x44, 0x56, 0x38, 0x10, 0x2d, 0x97, 0x20,
	0x41, 0x2a, 0x05, 0x71, 0x27, 0x69, 0xc9, 0x87, 0xd4, 0x02, 0xb2, 0x0b, 0x37, 0x40, 0x4e, 0x3c,
	0x72, 0x2c, 0x12, 0x6f, 0xb4, 0xbb, 0x86, 0x9f, 0xc5, 0x5f, 0x44, 0xfb, 0x65, 0x25, 0x29, 0x42,
	0x1c, 0xb8, 0xcd, 0x9b, 0x79, 0x3b, 0xf3, 0x66, 0x9e, 0x0d, 0x9d, 0x1f, 0x85, 0xc4, 0xdd, 0xf2,
	0x72, 0x8b, 0x42, 0xa4, 0x39, 0x0e, 0x77, 0x9c, 0x49, 0x46, 0x02, 0x93, 0xed, 0xf5, 0x6c, 0x35,
	0x5d, 0xad, 0x58, 0x55, 0xca, 0x6f, 0xcb, 0x0d, 0x5b, 0x7d, 0x37, 0x9c, 0xde, 0x53, 0x5b, 0x13,
	0x65, 0xba, 0x13, 0x6b, 0x76, 0x50, 0xa4, 0xbf, 0x3c, 0x68, 0x25, 0x32, 0x95, 0x95, 0xb8, 0x15,
	0x39, 0xe9, 0xc0, 0xc9, 0x7b, 0x94, 0x8b, 0xeb, 0xc8, 0xeb, 0x7b, 0x83, 0x66, 0x6c, 0x00, 0x89,
	0x20, 0xfc, 0x8c, 0x5c, 0x14, 0xac, 0x8c, 0x7c, 0x9d, 0x77, 0x90, 0x74, 0x21, 0x98, 0x63, 0x91,
	0xaf, 0x65, 0xd4, 0xd0, 0x05, 0x8b, 0x08, 0x85, 0xf6, 0x55, 0xc5, 0x39, 0x96, 0x72, 0xa2, 0x66,
	0x45, 0xcd, 0xbe, 0x37, 0x68, 0xc7, 0x07, 0x39, 0xc5, 0x99, 0x61, 0x89, 0xa2, 0x10, 0x86, 0x73,
	0x62, 0x38, 0xfb, 0x39, 0x42, 0xa0, 0xf9, 0x91, 0x71, 0x19, 0x05, 0x7d, 0x6f, 0xf0, 0x20, 0xd6,
	0x31, 0xfd, 0x0a, 0x9d, 0x19, 0xca, 0xc4, 0x2e, 0xa3, 0x79, 0x5a, 0x7b, 0x17, 0x82, 0x0f, 0xbc,
	0xc8, 0x8b, 0x52, 0x8b, 0x6f, 0xc7, 0x16, 0xa9, 0x9d, 0xae, 0xd4, 0x4d, 0xac, 0x76, 0x03, 0xd4,
	0x4e, 0x53, 0xc6, 0x7f, 0xa6, 0x3c, 0xd3, 0xd2, 0x4f, 0x63, 0x07, 0xe9, 0x04, 0x1e, 0xdf, 0x6f,
	0xfe, 0x0a, 0x02, 0x7d, 0x35, 0x11, 0x79, 0xfd, 0xc6, 0xe0, 0x6c, 0xf4, 0x64, 0x68, 0x8e, 0x3a,
	0x3c, 0xa0, 0xc6, 0x96, 0x44, 0xbf, 0xc0, 0xc5, 0x0c, 0xe5, 0xd8, 0x98, 0xf1, 0xff, 0x25, 0xbe,
	0x85, 0xf3, 0x7b, 0xbd, 0x5f, 0xd6, 0x0a, 0x1b, 0x5a, 0x61, 0xc7, 0x29, 0xdc, 0x67, 0xd6, 0x02,
	0xdf, 0x40, 0xa8, 0x13, 0x8b, 0x6b, 0x75, 0xe3, 0x79, 0x2a, 0xd6, 0x56, 0x92, 0x8e, 0xf7, 0x7c,
	0xf5, 0xf7, 0x7d, 0xa5, 0x2b, 0x08, 0x13, 0xcc, 0xb7, 0x58, 0x4a, 0xf2, 0x1c, 0x9a, 0x53, 0xce,
	0xb6, 0xfa, 0xd9, 0xd9, 0xe8, 0x91, 0x9b, 0x66, 0xbb, 0xc6, 0xba, 0x48, 0x9e, 0x81, 0x7f, 0xc7,
	0x74, 0x8f, 0x3f, 0x50, 0xfc, 0x3b, 0xa6, 0x86, 0x27, 0x12, 0x77, 0xf6, 0xf3, 0xd1, 0x31, 0xfd,
	0x04, 0x0f, 0xad, 0x66, 0x37, 0x2b, 0x82, 0x70, 0x9c, 0x65, 0x1c, 0x85, 0xb0, 0x2a, 0x1d, 0x24,
	0x2f, 0x6a, 0x41, 0xc7, 0x53, 0x6c, 0x3a, 0x76, 0x75, 0xba, 0x80, 0x8b, 0xdb, 0x6a, 0x23, 0x8b,
	0xa3, 0xde, 0x23, 0x38, 0xb5, 0xa1, 0xf3, 0xb6, 0x7b, 0x74, 0x39, 0xd7, 0xa9, 0xe6, 0xd1, 0x39,
	0x9c, 0x4e, 0x8b, 0x0d, 0xde, 0x14, 0x42, 0x2a, 0xef, 0x54, 0x6c, 0x1e, 0xb7, 0x62, 0x03, 0x54,
	0x36, 0x91, 0x29, 0xaf, 0x1d, 0xd5, 0x80, 0x9c, 0x43, 0xe3, 0x5d, 0x99, 0xd9, 0x65, 0x55, 0x48,
	0x37, 0xd0, 0x4a, 0xaa, 0xe5, 0x0d, 0x66, 0x39, 0x72, 0x72, 0x09, 0x61, 0x32, 0xf9, 0x87, 0xaf,
	0xcc, 0xb1, 0xc8, 0x10, 0xc2, 0xb1, 0x7d, 0xe0, 0xff, 0xc5, 0x74, 0x47, 0x5a, 0x06, 0xfa, 0x9f,
	0x7f, 0xfd, 0x3b, 0x00, 0x00, 0xff, 0xff, 0xda, 0x25, 0x07, 0x94, 0x4c, 0x04, 0x00, 0x00,
}
