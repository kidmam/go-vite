package net

import (
	"errors"
	"fmt"
	"math/rand"
	net2 "net"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/crypto/ed25519"
	"github.com/vitelabs/go-vite/interfaces"
	"github.com/vitelabs/go-vite/net/protos"
	"github.com/vitelabs/go-vite/p2p"
	"github.com/vitelabs/go-vite/p2p/vnode"
)

var errReadTooShort = errors.New("read too short")
var errWriteTooShort = errors.New("write too short")
var errPayloadTooLarge = errors.New("payload too large")
var errHandshakeError = errors.New("sync handshake error")
var errServerNotReady = errors.New("server not ready")
var errIncompleteChunk = errors.New("incomplete chunk")

type syncCode byte

func (s syncCode) code() syncCode {
	return s
}

func (s syncCode) Serialize() ([]byte, error) {
	return nil, nil
}

const (
	syncHandshake syncCode = iota
	syncHandshakeDone
	syncHandshakeErr
	syncRequest
	syncReady // server begin transmit data
	syncMissing
	syncNoAuth
	syncServerError // server error, like open reader failed
	syncQuit
)

type syncMsg interface {
	code() syncCode
	p2p.Serializable
}

type syncCodecer interface {
	read() (syncMsg, error)
	write(msg syncMsg) error
}

func syncMsgParser(code syncCode, payload []byte) (syncMsg, error) {
	switch code {
	case syncHandshake:
		var msg = new(syncHandshakeMsg)
		err := msg.deserialize(payload)
		if err != nil {
			return nil, err
		}
		return msg, nil
	case syncRequest:
		var msg = new(syncRequestMsg)
		err := msg.deserialize(payload)
		if err != nil {
			return nil, err
		}
		return msg, nil
	case syncReady:
		var msg = new(syncReadyMsg)
		err := msg.deserialize(payload)
		if err != nil {
			return nil, err
		}
		return msg, nil
	default:
		return code, nil
	}
}

type syncCodec struct {
	net2.Conn
	builder func(code syncCode, payload []byte) (syncMsg, error)
	buf     [257]byte
	timeout time.Duration
}

// 1 byte code
// 1 byte length
// 0 ~ 255 byte payload
func (s *syncCodec) read() (syncMsg, error) {
	_ = s.Conn.SetReadDeadline(time.Now().Add(s.timeout))

	n, err := s.Conn.Read(s.buf[:2])
	if err != nil {
		return nil, err
	}
	if n != 2 {
		return nil, errReadTooShort
	}

	scode := syncCode(s.buf[0])

	length := s.buf[1]
	n, err = s.Conn.Read(s.buf[:length])
	if err != nil {
		return nil, err
	}
	if n != int(length) {
		return nil, errReadTooShort
	}

	return s.builder(scode, s.buf[:length])
}

func (s *syncCodec) write(msg syncMsg) error {
	_ = s.Conn.SetWriteDeadline(time.Now().Add(s.timeout))

	payload, err := msg.Serialize()
	if err != nil {
		return err
	}

	length := len(payload)

	if length > 255 {
		return errPayloadTooLarge
	}

	s.buf[0] = byte(msg.code())
	s.buf[1] = byte(length)
	copy(s.buf[2:], payload)

	n, err := s.Conn.Write(s.buf[:2+length])
	if err != nil {
		return err
	}
	if n != 2+length {
		return errWriteTooShort
	}

	return nil
}

type syncHandshakeMsg struct {
	key  []byte
	time time.Time
	sign []byte
}

func (s *syncHandshakeMsg) code() syncCode {
	return syncHandshake
}

func (s *syncHandshakeMsg) Serialize() ([]byte, error) {
	pb := &protos.SyncConnHandshake{
		Key:       s.key,
		Timestamp: s.time.Unix(),
		Sign:      s.sign,
	}
	return proto.Marshal(pb)
}

func (s *syncHandshakeMsg) deserialize(data []byte) error {
	pb := &protos.SyncConnHandshake{}
	err := proto.Unmarshal(data, pb)
	if err != nil {
		return err
	}

	s.key = pb.Key
	s.time = time.Unix(pb.Timestamp, 0)
	s.sign = pb.Sign
	return nil
}

type syncRequestMsg struct {
	from, to uint64
}

func (s *syncRequestMsg) code() syncCode {
	return syncRequest
}

func (s *syncRequestMsg) Serialize() ([]byte, error) {
	pb := &protos.ChunkRequest{
		From: s.from,
		To:   s.to,
	}

	return proto.Marshal(pb)
}

func (s *syncRequestMsg) deserialize(data []byte) error {
	pb := &protos.ChunkRequest{}
	err := proto.Unmarshal(data, pb)
	if err != nil {
		return err
	}
	s.from = pb.From
	s.to = pb.To
	return nil
}

type syncReadyMsg struct {
	from, to          uint64
	size              uint64
	prevHash, endHash types.Hash
}

func (s *syncReadyMsg) code() syncCode {
	return syncReady
}

func (s *syncReadyMsg) Serialize() ([]byte, error) {
	pb := &protos.ChunkInfo{
		From:     s.from,
		To:       s.to,
		Size:     s.size,
		PrevHash: s.prevHash.Bytes(),
		EndHash:  s.endHash.Bytes(),
	}

	return proto.Marshal(pb)
}

func (s *syncReadyMsg) deserialize(data []byte) error {
	pb := &protos.ChunkInfo{}
	err := proto.Unmarshal(data, pb)
	if err != nil {
		return err
	}

	s.prevHash, err = types.BytesToHash(pb.PrevHash)
	if err != nil {
		return err
	}

	s.endHash, err = types.BytesToHash(pb.EndHash)
	if err != nil {
		return err
	}

	s.from = pb.From
	s.to = pb.To
	s.size = pb.Size
	return nil
}

type syncConnState byte

const (
	fileConnStateNew syncConnState = iota
	fileConnStateIdle
	fileConnStateBusy
	fileConnStateClosed
)

type SyncConnectionStatus struct {
	Address string `json:"address"`
	Speed   string `json:"speed"`
	Task    string `json:"task"`
}

type syncConnection interface {
	net2.Conn
	syncCodecer
	ID() peerId
	download(from, to uint64) (fatal bool, err error)
	speed() uint64
	state() syncConnState
	status() SyncConnectionStatus
	isBusy() bool
	height() uint64
}

type syncConnectionFactory interface {
	syncConnInitiator
	syncConnReceiver
}

type syncConnReceiver interface {
	receive(conn net2.Conn) (syncConnection, error)
}

type syncConnInitiator interface {
	initiate(conn net2.Conn, peer downloadPeer) (syncConnection, error)
}

type defaultSyncConnectionFactory struct {
	chain      syncCacher
	peers      *peerSet
	privateKey ed25519.PrivateKey
}

func (d *defaultSyncConnectionFactory) makeCodec(conn net2.Conn) syncCodecer {
	return &syncCodec{
		Conn:    conn,
		builder: syncMsgParser,
		timeout: 10 * time.Second,
	}
}

func (d *defaultSyncConnectionFactory) initiate(conn net2.Conn, peer downloadPeer) (syncConnection, error) {
	codec := d.makeCodec(conn)
	err := codec.write(&syncHandshakeMsg{
		key:  d.privateKey.PubByte(),
		time: time.Now(),
		sign: nil,
	})
	if err != nil {
		return nil, err
	}

	msg, err := codec.read()
	if err != nil {
		return nil, err
	}

	if msg.code() != syncHandshakeDone {
		return nil, errHandshakeError
	}

	return &syncConn{
		Conn:         conn,
		syncCodecer:  codec,
		downloadPeer: peer,
		cacher:       d.chain,
	}, nil
}

func (d *defaultSyncConnectionFactory) receive(conn net2.Conn) (syncConnection, error) {
	codec := d.makeCodec(conn)

	msg, err := codec.read()
	if err != nil {
		return nil, err
	}

	if msg.code() != syncHandshake {
		_ = codec.write(syncHandshakeErr)
		return nil, errHandshakeError
	}

	var id peerId
	if h, ok := msg.(*syncHandshakeMsg); ok {
		id, err = vnode.Bytes2NodeID(h.key)
		if err != nil {
			_ = codec.write(syncHandshakeErr)
			return nil, errHandshakeError
		}
	} else {
		_ = codec.write(syncHandshakeErr)
		return nil, errHandshakeError
	}

	p := d.peers.get(id)
	if p == nil {
		_ = codec.write(syncNoAuth)
		return nil, errHandshakeError
	}

	err = codec.write(syncHandshakeDone)
	if err != nil {
		return nil, err
	}

	return &syncConn{
		Conn:         conn,
		syncCodecer:  codec,
		downloadPeer: p,
	}, nil
}

type syncConn struct {
	net2.Conn
	syncCodecer
	downloadPeer
	busy   int32 // atomic
	st     syncConnState
	t      int64     // timestamp
	_speed uint64    // download speed, byte/s
	chunk  [2]uint64 // task
	closed int32
	cacher syncCacher
	buf    [1024]byte
	failed int32
}

var speedUnits = [...]string{
	" Byte/s",
	" KByte/s",
	" MByte/s",
	" GByte/s",
}

func formatSpeed(s float64) (sf float64, unit int) {
	for unit = 1; unit < len(speedUnits); unit++ {
		if sf = s / 1024.0; sf > 1 {
			s = sf
		} else {
			break
		}
	}

	unit--

	return s, unit
}

func speedToString(s float64) string {
	sf, unit := formatSpeed(s)

	return strconv.FormatFloat(sf, 'f', 2, 64) + speedUnits[unit]
}

func (f *syncConn) status() SyncConnectionStatus {
	st := SyncConnectionStatus{
		Address: f.downloadPeer.ID().Brief() + "@" + f.Conn.RemoteAddr().String(),
		Speed:   speedToString(float64(f._speed)),
		Task:    "",
	}

	if f.isBusy() {
		st.Task = strconv.FormatUint(f.chunk[0], 10) + "-" + strconv.FormatUint(f.chunk[1], 10)
	}

	return st
}

func (f *syncConn) state() syncConnState {
	return f.st
}

func (f *syncConn) fail() bool {
	f.failed++

	return f.failed > 3
}

func (f *syncConn) speed() uint64 {
	return f._speed
}

func (f *syncConn) isBusy() bool {
	return atomic.LoadInt32(&f.busy) == 1
}

func (f *syncConn) download(from, to uint64) (fatal bool, err error) {
	f.setBusy()
	defer f.idle()

	f.chunk = [2]uint64{from, to}

	err = f.write(&syncRequestMsg{
		from: from,
		to:   to,
	})

	if err != nil {
		f.fail()
		return true, err
	}

	var msg syncMsg
	msg, err = f.read()
	if err != nil {
		f.fail()
		return true, err
	}

	if msg.code() != syncReady {
		fatal = f.fail()
		return fatal, errServerNotReady
	}

	chunkInfo := msg.(*syncReadyMsg)

	cache := f.cacher.GetSyncCache()
	segment := interfaces.Segment{
		Bound:    [2]uint64{from, to},
		Hash:     chunkInfo.endHash,
		PrevHash: chunkInfo.prevHash,
	}
	writer, err := cache.NewWriter(segment)
	if err != nil {
		return false, err
	}

	start := time.Now().Unix()
	var nr, nw int
	var total, count uint64
	var rerr, werr error
	_ = f.Conn.SetReadDeadline(time.Now().Add(fileTimeout))
	for {
		count = chunkInfo.size - total
		if count > 1024 {
			count = 1024
		}

		nr, rerr = f.Conn.Read(f.buf[:count])
		total += uint64(nr)

		nw, werr = writer.Write(f.buf[:nr])

		if rerr != nil {
			break
		} else if werr != nil {
			break
		} else if nw != nr {
			werr = errWriteTooShort
			break
		}

		if total == chunkInfo.size {
			break
		}
	}

	_ = writer.Close()

	if rerr != nil {
		fatal = true
	}

	if werr != nil {
		err = fmt.Errorf("failed to write cache %d-%d: %v", from, to, werr)
		_ = cache.Delete(segment)
		return
	}

	if total != chunkInfo.size {
		fatal = true
		err = errIncompleteChunk
		_ = cache.Delete(segment)
		return
	}

	f._speed = total / uint64(time.Now().Unix()-start+1)

	return
}

func (f *syncConn) setBusy() {
	atomic.StoreInt32(&f.busy, 1)
	atomic.StoreInt64(&f.t, time.Now().Unix())
}

func (f *syncConn) idle() {
	atomic.StoreInt32(&f.busy, 0)
	atomic.StoreInt64(&f.t, time.Now().Unix())
}

func (f *syncConn) close() error {
	if atomic.CompareAndSwapInt32(&f.closed, 0, 1) {
		return f.Conn.Close()
	}

	return errSyncConnClosed
}

var errSyncConnExist = errors.New("sync connection has exist")
var errSyncConnClosed = errors.New("sync connection has closed")
var errPeerDialing = errors.New("peer is dialing")

type connections []syncConnection

func (fl connections) Len() int {
	return len(fl)
}

func (fl connections) Less(i, j int) bool {
	return fl[i].speed() > fl[j].speed()
}

func (fl connections) Swap(i, j int) {
	fl[i], fl[j] = fl[j], fl[i]
}

func (fl connections) del(i int) connections {
	total := len(fl)
	if i < total {
		copy(fl[i:], fl[i+1:])
		return fl[:total-1]
	}

	return fl
}

type FilePoolStatus struct {
	Connections []SyncConnectionStatus `json:"connections"`
}

type downloadPeer interface {
	ID() peerId
	height() uint64
	fileAddress() string
}

type downloadPeerSet interface {
	pickDownloadPeers(height uint64) (m map[peerId]downloadPeer)
}

type connPool interface {
	addConn(c syncConnection) error
	delConn(c syncConnection)
	chooseSource(to uint64) (downloadPeer, syncConnection, error)
	reset()
	connections() []SyncConnectionStatus
}

type connPoolImpl struct {
	mu    sync.Mutex
	peers downloadPeerSet
	mi    map[peerId]int // value is the index of `connPoolImpl.l`
	l     connections    // connections sort by speed, from fast to slow
}

func newPool(peers downloadPeerSet) *connPoolImpl {
	return &connPoolImpl{
		mi:    make(map[peerId]int),
		peers: peers,
	}
}

func (fp *connPoolImpl) connections() []SyncConnectionStatus {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	cs := make([]SyncConnectionStatus, len(fp.l))

	for i := 0; i < len(fp.l); i++ {
		cs[i] = fp.l[i].status()
	}

	return cs
}

// delete filePeer and connection
func (fp *connPoolImpl) delConn(c syncConnection) {
	_ = c.Close()

	fp.mu.Lock()
	defer fp.mu.Unlock()

	fp.delConnLocked(c.ID())
}

func (fp *connPoolImpl) delConnLocked(id peerId) {
	if i, ok := fp.mi[id]; ok {
		delete(fp.mi, id)

		fp.l = fp.l.del(i)
	}
}

func (fp *connPoolImpl) addConn(c syncConnection) error {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	if _, ok := fp.mi[c.ID()]; ok {
		return errSyncConnExist
	}

	fp.l = append(fp.l, c)
	fp.mi[c.ID()] = len(fp.l) - 1
	return nil
}

// sort list, and update index to map
func (fp *connPoolImpl) sort() {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	fp.sortLocked()
}

func (fp *connPoolImpl) sortLocked() {
	sort.Sort(fp.l)
	for i, c := range fp.l {
		fp.mi[c.ID()] = i
	}
}

// choose the fast fileConn, or create new conn randomly
func (fp *connPoolImpl) chooseSource(to uint64) (downloadPeer, syncConnection, error) {
	peerMap := fp.peers.pickDownloadPeers(to)

	if len(peerMap) == 0 {
		return nil, nil, errNoSuitablePeer
	}

	fp.mu.Lock()
	defer fp.mu.Unlock()

	// only peers without sync connection
	for _, c := range fp.l {
		delete(peerMap, c.ID())
	}

	var createNew bool
	if len(peerMap) > 0 {
		createNew = rand.Intn(10) > 5
	}

	fp.sortLocked()
	for i, c := range fp.l {
		if c.isBusy() || c.height() < to {
			continue
		}

		if len(fp.l)+1 > 3*(i+1) {
			// fast enough
			return nil, c, nil
		}

		if createNew {
			for _, p := range peerMap {
				return p, nil, nil
			}
		} else {
			return nil, c, nil
		}
	}

	for _, p := range peerMap {
		return p, nil, nil
	}

	return nil, nil, nil
}

func (fp *connPoolImpl) reset() {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	fp.mi = make(map[peerId]int)

	for _, c := range fp.l {
		_ = c.Close()
	}

	fp.l = nil
}
