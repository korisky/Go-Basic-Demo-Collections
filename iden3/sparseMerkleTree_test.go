package idenDemo

import (
	"context"
	"fmt"
	"github.com/iden3/go-merkletree-sql"
	"github.com/iden3/go-merkletree-sql/db/memory"
	"math/big"
	"testing"
)

func Test_SparseMerkleTree(t *testing.T) {

	ctx := context.Background()

	// Tree storage
	store := memory.NewMemoryStorage()

	// generate a new MerkleTree with 32 levels
	mt, _ := merkletree.NewMerkleTree(ctx, store, 32)

	// add a leaf to the tree with idx:1 & val:10
	idx1 := big.NewInt(1)
	val1 := big.NewInt(10)
	_ = mt.Add(ctx, idx1, val1)

	// add another leaf
	idx2 := big.NewInt(2)
	val2 := big.NewInt(15)
	_ = mt.Add(ctx, idx2, val2)

	// proof of membership of a leaf with idx 1 (just send your idx into it)
	proof, pr_val, _ := mt.GenerateProof(ctx, idx1, mt.Root())
	fmt.Println("Proof of membership:", proof.Existence)
	fmt.Println("Val corresponding to the queried idx:", pr_val)

	proof_2, _, _ := mt.GenerateProof(ctx, big.NewInt(23), mt.Root())
	fmt.Println("Proof of membership:", proof_2.Existence)
}
