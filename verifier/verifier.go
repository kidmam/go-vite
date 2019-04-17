package verifier

import (
	"errors"
	"fmt"

	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/log15"
	"github.com/vitelabs/go-vite/vm_db"
)

type VerifyResult int

type verifier struct {
	Sv  *SnapshotVerifier
	Av  *AccountVerifier
	log log15.Logger
}

func NewVerifier(sv *SnapshotVerifier, av *AccountVerifier) Verifier {
	return &verifier{
		Sv:  sv,
		Av:  av,
		log: log15.New("module", "verifier"),
	}
}

func (v *verifier) VerifyNetSb(block *ledger.SnapshotBlock) error {
	return v.Sv.VerifyNetSb(block)
}

func (v *verifier) VerifyNetAb(block *ledger.AccountBlock) error {
	// fixme 1. makesure genesis and initial-balance blocks don't need to check, return nil
	//2. VerifyHash
	if err := v.Av.verifyHash(block); err != nil {
		return err
	}
	//3. VerifySignature
	if v.Av.chain.IsGenesisAccountBlock(block.Hash) {
		return nil
	}
	if err := v.Av.verifySignature(block); err != nil {
		return err
	}
	return nil
}

func (v *verifier) VerifyPoolAccBlock(block *ledger.AccountBlock, snapshotHash *types.Hash) (*AccBlockPendingTask, *vm_db.VmAccountBlock, error) {
	eLog := v.log.New("method", "VerifyPoolAccBlock")

	detail := fmt.Sprintf("sbHash=%v; block:addr=%v height=%v hash=%v; ", snapshotHash, block.AccountAddress, block.Height, block.Hash)
	if block.IsReceiveBlock() {
		detail += fmt.Sprintf("fromHash=%v;", block.FromBlockHash)
	}

	verifyResult, task, err := v.Av.verifyReferred(block)
	if err != nil {
		eLog.Error(err.Error(), "d", detail)
	}
	switch verifyResult {
	case PENDING:
		return task, nil, nil
	case SUCCESS:
		blocks, err := v.Av.vmVerify(block, snapshotHash)
		if err != nil {
			eLog.Error(err.Error(), "d", detail)
			return nil, nil, err
		}
		return nil, blocks, nil
	default:
		return nil, nil, err
	}
}

func (v *verifier) VerifyRPCAccBlock(block *ledger.AccountBlock, snapshotHash *types.Hash) (*vm_db.VmAccountBlock, error) {
	log := v.log.New("method", "VerifyRPCAccBlock")

	detail := fmt.Sprintf("sbHash:%v, addr:%v, height:%v, hash:%v", snapshotHash, block.AccountAddress, block.Height, block.Hash)
	if block.IsReceiveBlock() {
		detail += fmt.Sprintf(",fromH:%v", block.FromBlockHash)
	}

	if verifyResult, task, err := v.Av.verifyReferred(block); verifyResult != SUCCESS {
		if err != nil {
			log.Error(err.Error(), "d", detail)
			return nil, err
		}
		log.Error("verify block failed, pending for:"+task.pendingHashListToStr(), "d", detail)
		return nil, ErrVerifyRPCBlockPendingState
	}

	vmBlock, err := v.Av.vmVerify(block, snapshotHash)
	if err != nil {
		log.Error(err.Error(), "d", detail)
		return nil, err
	}
	return vmBlock, nil
}

func (v *verifier) VerifyAccBlockHash(block *ledger.AccountBlock) error {
	return v.Av.verifyHash(block)
}

func (v *verifier) VerifyAccBlockSignature(block *ledger.AccountBlock) error {
	accType, err := v.Av.verifyAccAddress(block)
	if err != nil {
		return err
	}
	if accType == AccountTypeNotSure {
		return ErrVerifyAccountTypeNotSure
	}
	isGeneralAddr := isAccTypeGeneral(accType)
	if !isGeneralAddr && block.IsSendBlock() {
		if len(block.Signature) != 0 || len(block.PublicKey) != 0 {
			return errors.New("signature and publicKey of the contract's send must be nil")
		}
	} else {
		if err := v.Av.verifySignature(block); err != nil {
			return err
		}
	}
	return nil
}

func (v *verifier) VerifyAccBlockNonce(block *ledger.AccountBlock) error {
	accType, err := v.Av.verifyAccAddress(block)
	if err != nil {
		return err
	}
	if accType == AccountTypeNotSure {
		return ErrVerifyAccountTypeNotSure
	}
	return v.Av.verifyNonce(block, isAccTypeGeneral(accType))
}

func (v *verifier) VerifyAccBlockProducerLegality(block *ledger.AccountBlock) error {
	accType, err := v.Av.verifyAccAddress(block)
	if err != nil {
		return err
	}
	if accType == AccountTypeNotSure {
		return ErrVerifyAccountTypeNotSure
	}
	return v.Av.verifyProducerLegality(block, isAccTypeGeneral(accType))
}

func (v *verifier) GetSnapshotVerifier() *SnapshotVerifier {
	return v.Sv
}
