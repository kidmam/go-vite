package net

import (
	"fmt"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/log15"
	"github.com/vitelabs/go-vite/monitor"
	"github.com/vitelabs/go-vite/p2p"
	"github.com/vitelabs/go-vite/vite/net/message"
	"sync/atomic"
)

type fetcher struct {
	filter   Filter
	peers    *peerSet
	receiver Receiver
	pool     RequestPool
	ready    int32 // atomic
	log      log15.Logger
}

func newFetcher(filter Filter, peers *peerSet, receiver Receiver, pool RequestPool) *fetcher {
	return &fetcher{
		filter:   filter,
		peers:    peers,
		receiver: receiver,
		pool:     pool,
		log:      log15.New("module", "net/fetcher"),
	}
}

func (f *fetcher) ID() string {
	return "fetcher"
}

func (f *fetcher) Cmds() []cmd {
	return []cmd{SnapshotBlocksCode, AccountBlocksCode}
}

func (f *fetcher) Handle(msg *p2p.Msg, sender *Peer) error {
	switch cmd(msg.Cmd) {
	case SnapshotBlocksCode:
		bs := new(message.SnapshotBlocks)
		err := bs.Deserialize(msg.Payload)
		if err != nil {
			return err
		}

		f.receiver.ReceiveSnapshotBlocks(bs.Blocks)
	case AccountBlocksCode:
		bs := new(message.AccountBlocks)
		err := bs.Deserialize(msg.Payload)
		if err != nil {
			return err
		}

		f.receiver.ReceiveAccountBlocks(bs.Blocks)
	}

	return nil
}

func (f *fetcher) FetchSnapshotBlocks(start types.Hash, count uint64) {
	monitor.LogEvent("net/fetch", "s")

	// been suppressed
	if f.filter.hold(start) {
		f.log.Warn(fmt.Sprintf("fetch suppressed: %s %d", start, count))
		return
	}

	if atomic.LoadInt32(&f.ready) == 0 {
		f.log.Warn("not ready")
		return
	}

	m := &message.GetSnapshotBlocks{
		From:    ledger.HashHeight{Hash: start},
		Count:   count,
		Forward: true,
	}

	p := f.peers.BestPeer()
	if p != nil {
		id := f.pool.MsgID()
		err := p.Send(GetSnapshotBlocksCode, id, m)
		if err != nil {
			f.log.Error(fmt.Sprintf("send %s to %s error: %v", GetSnapshotBlocksCode, p, err))
		} else {
			f.log.Info(fmt.Sprintf("send %s to %s done", GetSnapshotBlocksCode, p))
		}
	} else {
		f.log.Error(errNoPeer.Error())
	}
}

func (f *fetcher) FetchAccountBlocks(start types.Hash, count uint64, address *types.Address) {
	monitor.LogEvent("net/fetch", "a")

	// been suppressed
	if f.filter.hold(start) {
		f.log.Warn(fmt.Sprintf("fetch suppressed: %s %d", start, count))
		return
	}

	if atomic.LoadInt32(&f.ready) == 0 {
		f.log.Warn("not ready")
		return
	}

	addr := NULL_ADDRESS
	if address != nil {
		addr = *address
	}
	m := &message.GetAccountBlocks{
		Address: addr,
		From: ledger.HashHeight{
			Hash: start,
		},
		Count:   count,
		Forward: true,
	}

	p := f.peers.BestPeer()
	if p != nil {
		id := f.pool.MsgID()
		err := p.Send(GetAccountBlocksCode, id, m)
		if err != nil {
			f.log.Error(fmt.Sprintf("send %s to %s error: %v", GetAccountBlocksCode, p, err))
		} else {
			f.log.Info(fmt.Sprintf("send %s to %s done", GetAccountBlocksCode, p))
		}
	} else {
		f.log.Error(errNoPeer.Error())
	}
}

func (f *fetcher) listen(st SyncState) {
	if st == Syncdone || st == SyncDownloaded {
		f.log.Info(fmt.Sprintf("ready: %s", st))

		atomic.StoreInt32(&f.ready, 1)
	}
}
