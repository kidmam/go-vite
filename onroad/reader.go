package onroad

import (
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/net"
	"github.com/vitelabs/go-vite/producer/producerevent"
	"github.com/vitelabs/go-vite/vm_db"
)

type Pool interface {
	AddDirectAccountBlock(address types.Address, vmAccountBlock *vm_db.VmAccountBlock) error
}

type Producer interface {
	SetAccountEventFunc(func(producerevent.AccountEvent))
}

type Net interface {
	SubscribeSyncStatus(fn func(net.SyncState)) (subId int)
	UnsubscribeSyncStatus(subId int)
	SyncState() net.SyncState
}
