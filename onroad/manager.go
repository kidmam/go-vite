package onroad

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/vitelabs/go-vite/chain"
	"github.com/vitelabs/go-vite/common"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/generator"
	"github.com/vitelabs/go-vite/log15"
	"github.com/vitelabs/go-vite/net"
	"github.com/vitelabs/go-vite/onroad/pool"
	"github.com/vitelabs/go-vite/producer/producerevent"
	"github.com/vitelabs/go-vite/wallet"
)

var (
	slog           = log15.New("module", "onroad")
	ErrNotSyncDone = errors.New("network synchronization is not complete")

	DefaultContractGidList = []types.Gid{types.DELEGATE_GID}
)

type Manager struct {
	net      Net
	producer Producer
	wallet   *wallet.Manager

	pool      Pool
	chain     chain.Chain
	consensus generator.Consensus

	contractWorkers     map[types.Gid]*ContractWorker
	newContractListener sync.Map //map[types.Gid]func(address types.Address)

	onRoadPools sync.Map //map[types.Gid]contract_pool.OnRoadPool

	unlockLid   int
	netStateLid int

	lastProducerAccEvent *producerevent.AccountStartEvent

	log log15.Logger
}

func NewManager(net Net, pool Pool, producer Producer, consensus generator.Consensus, wallet *wallet.Manager) *Manager {
	m := &Manager{
		net:             net,
		producer:        producer,
		wallet:          wallet,
		pool:            pool,
		consensus:       consensus,
		contractWorkers: make(map[types.Gid]*ContractWorker),
		log:             slog.New("w", "manager"),
	}
	return m
}

func (manager *Manager) Init(chain chain.Chain) {
	manager.chain = chain
	for _, gid := range DefaultContractGidList {
		manager.prepareOnRoadPool(gid)
	}
}

func (manager *Manager) Start() {
	manager.netStateLid = manager.Net().SubscribeSyncStatus(manager.netStateChangedFunc)
	if manager.producer != nil {
		manager.producer.SetAccountEventFunc(manager.producerStartEventFunc)
	}
	manager.Chain().Register(manager)
}

func (manager *Manager) Stop() {
	manager.log.Info("Close")
	manager.Net().UnsubscribeSyncStatus(manager.netStateLid)
	manager.wallet.RemoveUnlockChangeChannel(manager.unlockLid)
	if manager.producer != nil {
		manager.Producer().SetAccountEventFunc(nil)
	}
	manager.Chain().UnRegister(manager)
	manager.stopAllWorks()
	manager.log.Info("Close end")
}

func (manager *Manager) Close() error {
	return nil
}

func (manager *Manager) prepareOnRoadPool(gid types.Gid) {
	orPool, exist := manager.onRoadPools.Load(gid)
	manager.log.Info(fmt.Sprintf("prepareOnRoadPool"), "gid", gid, "exist", exist, "orPool", orPool)
	if !exist || orPool == nil {
		manager.onRoadPools.Store(gid, onroad_pool.NewContractOnRoadPool(gid, manager.chain))
		return
	}
}

func (manager *Manager) netStateChangedFunc(state net.SyncState) {
	manager.log.Info("receive chain net event", "state_bak", state)
	common.Go(func() {
		if state == net.SyncDone {
			manager.resumeContractWorks()
		} else {
			manager.stopAllWorks()
		}
	})
}

func (manager *Manager) producerStartEventFunc(accevent producerevent.AccountEvent) {
	netstate := manager.Net().SyncState()
	manager.log.Info("producerStartEventFunc receive event", "netstate", netstate)
	if netstate != net.SyncDone {
		manager.log.Error(ErrNotSyncDone.Error())
		return
	}

	event, ok := accevent.(producerevent.AccountStartEvent)
	if !ok {
		manager.log.Info("producerStartEventFunc not support this event")
		return
	}

	if !manager.wallet.GlobalCheckAddrUnlock(event.Address) {
		manager.log.Error("receive chain right event but address locked", "event", event)
		return
	}

	manager.lastProducerAccEvent = &event

	w, found := manager.contractWorkers[event.Gid]
	if !found {
		w = NewContractWorker(manager)
		manager.contractWorkers[event.Gid] = w
	}

	nowTime := time.Now()
	if nowTime.After(event.Stime) && nowTime.Before(event.Etime) {
		w.Start(event)
		time.AfterFunc(event.Etime.Sub(nowTime), func() {
			w.Stop()
		})
	} else {
		w.Stop()
	}
}

func (manager *Manager) stopAllWorks() {
	manager.log.Info("stopAllWorks called")
	var wg = sync.WaitGroup{}
	for _, v := range manager.contractWorkers {
		wg.Add(1)
		common.Go(func() {
			v.Stop()
			wg.Done()
		})
	}
	wg.Wait()
	manager.log.Info("stopAllWorks end")
}

func (manager *Manager) resumeContractWorks() {
	manager.log.Info("resumeContractWorks")
	if manager.lastProducerAccEvent != nil {
		nowTime := time.Now()
		if nowTime.After(manager.lastProducerAccEvent.Stime) && nowTime.Before(manager.lastProducerAccEvent.Etime) {
			cw, ok := manager.contractWorkers[manager.lastProducerAccEvent.Gid]
			if ok {
				manager.log.Info("resumeContractWorks found an cw need to resume", "gid", manager.lastProducerAccEvent.Gid)
				cw.Start(*manager.lastProducerAccEvent)
				time.AfterFunc(manager.lastProducerAccEvent.Etime.Sub(nowTime), func() {
					cw.Stop()
				})
			}
		}
	}
	manager.log.Info("end resumeContractWorks")
}

func (manager Manager) Chain() chain.Chain {
	return manager.chain
}

func (manager Manager) Net() Net {
	return manager.net
}

func (manager Manager) Producer() Producer {
	return manager.producer
}

func (manager Manager) Consensus() generator.Consensus {
	return manager.consensus
}
