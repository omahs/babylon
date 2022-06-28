package keeper

import (
	"errors"
	"github.com/babylonchain/babylon/x/checkpointing/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CheckpointsState struct {
	cdc         codec.BinaryCodec
	checkpoints sdk.KVStore
}

func (k Keeper) CheckpointsState(ctx sdk.Context) CheckpointsState {
	// Build the CheckpointsState storage
	store := ctx.KVStore(k.storeKey)
	return CheckpointsState{
		cdc:         k.cdc,
		checkpoints: prefix.NewStore(store, types.CheckpointsPrefix),
	}
}

// CreateRawCkptWithMeta inserts the raw checkpoint with meta into the storage by its epoch number
// a new checkpoint is created with the status of UNCEHCKPOINTED
func (cs CheckpointsState) CreateRawCkptWithMeta(ckpt *types.RawCheckpoint) {
	// save concrete ckpt object
	ckptWithMeta := types.NewCheckpointWithMeta(ckpt, types.Uncheckpointed)
	cs.checkpoints.Set(types.CkptsObjectKey(ckpt.EpochNum), types.CkptWithMetaToBytes(cs.cdc, ckptWithMeta))
}

// GetRawCkptWithMeta retrieves a raw checkpoint with meta by its epoch number
func (cs CheckpointsState) GetRawCkptWithMeta(epoch uint64) (*types.RawCheckpointWithMeta, error) {
	ckptsKey := types.CkptsObjectKey(epoch)
	rawBytes := cs.checkpoints.Get(ckptsKey)
	if rawBytes == nil {
		return nil, types.ErrCkptDoesNotExist.Wrap("no raw checkpoint with provided epoch")
	}

	return types.BytesToCkptWithMeta(cs.cdc, rawBytes)
}

// GetRawCkptsWithMetaByStatus retrieves raw checkpoints with meta by their status by the descending order of epoch
func (cs CheckpointsState) GetRawCkptsWithMetaByStatus(status types.CheckpointStatus, f func(sig *types.RawCheckpointWithMeta) bool) error {
	store := prefix.NewStore(cs.checkpoints, types.CkptsObjectPrefix)
	iter := store.ReverseIterator(nil, nil)
	defer iter.Close()

	// the iterator starts from the highest epoch number
	// once it gets to an epoch where the status is CONFIRMED,
	// all the lower epochs will be CONFIRMED
	for ; iter.Valid(); iter.Next() {
		ckptBytes := iter.Value()
		ckptWithMeta, err := types.BytesToCkptWithMeta(cs.cdc, ckptBytes)
		if err != nil {
			return err
		}
		// the loop can end if the current status is CONFIRMED but the requested status is not CONFIRMED
		if status != types.Confirmed && ckptWithMeta.Status == types.Confirmed {
			return nil
		}
		if ckptWithMeta.Status != status {
			continue
		}
		stop := f(ckptWithMeta)
		if stop {
			return nil
		}
	}
	return nil
}

// UpdateCkptStatus updates the checkpoint's status
func (cs CheckpointsState) UpdateCkptStatus(ckpt *types.RawCheckpoint, status types.CheckpointStatus) error {
	ckptWithMeta, err := cs.GetRawCkptWithMeta(ckpt.EpochNum)
	if err != nil {
		// the checkpoint should exist
		return err
	}
	if !ckptWithMeta.Ckpt.Hash().Equals(ckpt.Hash()) {
		return errors.New("hash not the same with existing checkpoint")
	}
	ckptWithMeta.Status = status
	cs.checkpoints.Set(sdk.Uint64ToBigEndian(ckpt.EpochNum), types.CkptWithMetaToBytes(cs.cdc, ckptWithMeta))

	return nil
}
