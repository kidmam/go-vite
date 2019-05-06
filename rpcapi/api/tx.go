package api

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/consensus"
	"github.com/vitelabs/go-vite/crypto/ed25519"
	"github.com/vitelabs/go-vite/generator"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/net"
	"github.com/vitelabs/go-vite/verifier"
	"github.com/vitelabs/go-vite/vite"
	"github.com/vitelabs/go-vite/vm"
	"github.com/vitelabs/go-vite/vm/quota"
	"github.com/vitelabs/go-vite/vm/util"
	"github.com/vitelabs/go-vite/vm_db"
	"go.uber.org/atomic"
)

type Tx struct {
	vite *vite.Vite
	N    int
}

func NewTxApi(vite *vite.Vite) *Tx {
	tx := &Tx{
		vite: vite,
		N:    5,
	}
	if vite.Producer() == nil {
		return tx
	}
	coinbase := vite.Producer().GetCoinBase()

	manager, err := vite.WalletManager().GetEntropyStoreManager(coinbase.String())
	if err != nil {
		panic(err)
	}

	var fromAddrs []types.Address
	var fromHexPrivKeys []string

	{
		for i := uint32(0); i < uint32(10); i++ {
			_, key, err := manager.DeriveForIndexPath(i)
			if err != nil {
				panic(err)
			}
			binKey, err := key.PrivateKey()
			if err != nil {
				panic(err)
			}

			pubKey, err := key.PublicKey()
			if err != nil {
				panic(err)
			}

			address, err := key.Address()
			if err != nil {
				panic(err)
			}

			fmt.Printf("%s hex public key:%s\n", address, hex.EncodeToString(pubKey))

			hexKey := hex.EncodeToString(binKey)
			fromHexPrivKeys = append(fromHexPrivKeys, hexKey)
			fromAddrs = append(fromAddrs, *address)

		}
	}

	toAddr := types.AddressConsensusGroup
	amount := string("0")

	ss := []string{
		"/cF/JQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAnMxAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"/cF/JQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAnMyAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"/cF/JQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAnMzAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"/cF/JQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAnM0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"/cF/JQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAnM1AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
	}
	var datas [][]byte
	for _, v := range ss {
		byts, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			panic(err)
		}
		datas = append(datas, byts)
	}

	num := atomic.NewUint32(0)

	vite.Consensus().Subscribe(types.SNAPSHOT_GID, "api-auto-send", nil, func(e consensus.Event) {

		if num.Load() > 0 {
			fmt.Printf("something is loading[return].%s\n", time.Now())
			return
		}
		num.Add(1)
		defer num.Sub(1)
		snapshotBlock := vite.Chain().GetLatestSnapshotBlock()
		if snapshotBlock.Height < 10 {
			fmt.Println("latest height must >= 10.")
			return
		}

		state := vite.Net().Status().State
		if state != net.SyncDone {
			fmt.Printf("sync state: %s \n", state)
			return
		}

		for i := 0; i < tx.N; i++ {
			for k, v := range fromAddrs {
				addr := v
				key := fromHexPrivKeys[k]

				block, err := tx.SendTxWithPrivateKey(SendTxWithPrivateKeyParam{
					SelfAddr:     &addr,
					ToAddr:       &toAddr,
					TokenTypeId:  ledger.ViteTokenId,
					PrivateKey:   &key,
					Amount:       &amount,
					Data:         datas[rand.Intn(len(datas))],
					Difficulty:   nil,
					PreBlockHash: nil,
					BlockType:    2,
				})
				if err != nil {
					log.Error(fmt.Sprintf("send block err:%v\n", err))
					return
				} else {
					log.Info(fmt.Sprintf("send block:%s,%s,%s\n", block.AccountAddress, block.Height, block.Hash))
				}
			}
		}
	})

	return tx
}

func (t *Tx) UpdateBenchMark(cnt int) (int, error) {
	t.N = cnt
	return t.N, nil
}

func (t Tx) SendRawTx(block *AccountBlock) error {
	log.Info("SendRawTx")
	if block == nil {
		return errors.New("empty block")
	}

	lb, err := block.LedgerAccountBlock()
	if err != nil {
		return err
	}
	// need to remove Later
	//if len(lb.Data) != 0 && !isBuiltinContracts(lb.ToAddress) {
	//	return ErrorNotSupportAddNot
	//}
	//
	//if len(lb.Data) != 0 && block.BlockType == ledger.BlockTypeReceive {
	//	return ErrorNotSupportRecvAddNote
	//}

	v := verifier.NewVerifier(nil, verifier.NewAccountVerifier(t.vite.Chain(), t.vite.Consensus()))
	latestSb := t.vite.Chain().GetLatestSnapshotBlock()
	if latestSb == nil {
		return errors.New("failed to get latest snapshotBlock")
	}
	result, err := v.VerifyRPCAccBlock(lb, &latestSb.Hash)
	if err != nil {
		return err
	}

	if result != nil {
		return t.vite.Pool().AddDirectAccountBlock(result.AccountBlock.AccountAddress, result)
	} else {
		return errors.New("generator gen an empty block")
	}
	return nil
}

func (t Tx) SendTxWithPrivateKey(param SendTxWithPrivateKeyParam) (*AccountBlock, error) {

	if param.Amount == nil {
		return nil, errors.New("amount is nil")
	}

	if param.SelfAddr == nil {
		return nil, errors.New("selfAddr is nil")
	}

	if param.ToAddr == nil && param.BlockType != ledger.BlockTypeSendCreate {
		return nil, errors.New("toAddr is nil")
	}

	if param.PrivateKey == nil {
		return nil, errors.New("privateKey is nil")
	}

	var d *big.Int = nil
	if param.Difficulty != nil {
		t, ok := new(big.Int).SetString(*param.Difficulty, 10)
		if !ok {
			return nil, ErrStrToBigInt
		}
		d = t
	}

	amount, ok := new(big.Int).SetString(*param.Amount, 10)
	if !ok {
		return nil, ErrStrToBigInt
	}

	var blockType byte
	if param.BlockType > 0 {
		blockType = param.BlockType
	} else {
		blockType = ledger.BlockTypeSendCall
	}

	msg := &generator.IncomingMessage{
		BlockType:      blockType,
		AccountAddress: *param.SelfAddr,
		ToAddress:      param.ToAddr,
		TokenId:        &param.TokenTypeId,
		Amount:         amount,
		Fee:            nil,
		Data:           param.Data,
		Difficulty:     d,
	}

	addrState, err := generator.GetAddressStateForGenerator(t.vite.Chain(), &msg.AccountAddress)
	if err != nil || addrState == nil {
		return nil, errors.New(fmt.Sprintf("failed to get addr state for generator, err:%v", err))
	}
	g, e := generator.NewGenerator(t.vite.Chain(), t.vite.Consensus(), msg.AccountAddress, addrState.LatestSnapshotHash, addrState.LatestAccountHash)
	if e != nil {
		return nil, e
	}
	result, e := g.GenerateWithMessage(msg, &msg.AccountAddress, func(addr types.Address, data []byte) (signedData, pubkey []byte, err error) {
		var privkey ed25519.PrivateKey
		privkey, e := ed25519.HexToPrivateKey(*param.PrivateKey)
		if e != nil {
			return nil, nil, e
		}
		signData := ed25519.Sign(privkey, data)
		pubkey = privkey.PubByte()
		return signData, pubkey, nil
	})
	if e != nil {
		return nil, e
	}
	if result.Err != nil {
		return nil, result.Err
	}
	if result.VmBlock != nil {
		if err := t.vite.Pool().AddDirectAccountBlock(msg.AccountAddress, result.VmBlock); err != nil {
			return nil, err
		}
		return ledgerToRpcBlock(result.VmBlock.AccountBlock, t.vite.Chain())
	} else {
		return nil, errors.New("generator gen an empty block")
	}
}

type SendTxWithPrivateKeyParam struct {
	SelfAddr     *types.Address    `json:"selfAddr"`
	ToAddr       *types.Address    `json:"toAddr"`
	TokenTypeId  types.TokenTypeId `json:"tokenTypeId"`
	PrivateKey   *string           `json:"privateKey"` //hex16
	Amount       *string           `json:"amount"`
	Data         []byte            `json:"data"` //base64
	Difficulty   *string           `json:"difficulty,omitempty"`
	PreBlockHash *types.Hash       `json:"preBlockHash,omitempty"`
	BlockType    byte              `json:"blockType"`
}

type CalcPoWDifficultyParam struct {
	SelfAddr types.Address `json:"selfAddr"`
	PrevHash types.Hash    `json:"prevHash"`

	BlockType byte           `json:"blockType"`
	ToAddr    *types.Address `json:"toAddr"`
	Data      []byte         `json:"data"`

	UsePledgeQuota bool `json:"usePledgeQuota"`
}

type CalcPoWDifficultyResult struct {
	QuotaRequired uint64 `json:"quota"`
	Difficulty    string `json:"difficulty"`
}

func (t Tx) CalcPoWDifficulty(param CalcPoWDifficultyParam) (result *CalcPoWDifficultyResult, err error) {
	latestBlock, err := t.vite.Chain().GetLatestAccountBlock(param.SelfAddr)
	if err != nil {
		return nil, err
	}
	if (latestBlock == nil && !param.PrevHash.IsZero()) ||
		(latestBlock != nil && latestBlock.Hash != param.PrevHash) {
		return nil, util.ErrChainForked
	}
	// get quota required
	block := &ledger.AccountBlock{
		BlockType:      param.BlockType,
		AccountAddress: param.SelfAddr,
		PrevHash:       param.PrevHash,
		Data:           param.Data,
	}
	if param.ToAddr != nil {
		block.ToAddress = *param.ToAddr
	} else if param.BlockType == ledger.BlockTypeSendCall {
		return nil, errors.New("toAddr is nil")
	}
	sb := t.vite.Chain().GetLatestSnapshotBlock()
	db, err := vm_db.NewVmDb(t.vite.Chain(), &param.SelfAddr, &sb.Hash, &param.PrevHash)
	if err != nil {
		return nil, err
	}
	quotaRequired, err := vm.GasRequiredForBlock(db, block)
	if err != nil {
		return nil, err
	}

	// get current quota
	var pledgeAmount *big.Int
	var q types.Quota
	if param.UsePledgeQuota {
		pledgeAmount, err = t.vite.Chain().GetPledgeBeneficialAmount(param.SelfAddr)
		if err != nil {
			return nil, err
		}
		q, err := quota.GetPledgeQuota(db, param.SelfAddr, pledgeAmount)
		if err != nil {
			return nil, err
		}
		if q.Current() >= quotaRequired {
			return &CalcPoWDifficultyResult{quotaRequired, ""}, nil
		}
	} else {
		pledgeAmount = big.NewInt(0)
		q = types.NewQuota(0, 0, 0, 0)
	}
	// calc difficulty if current quota is not enough
	canPoW, err := quota.CanPoW(db, block.AccountAddress)
	if err != nil {
		return nil, err
	}
	if !canPoW {
		return nil, util.ErrCalcPoWTwice
	}
	d, err := quota.CalcPoWDifficulty(quotaRequired, q)
	if err != nil {
		return nil, err
	}
	return &CalcPoWDifficultyResult{quotaRequired, d.String()}, nil
}
