package idenDemo

import (
	"fmt"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"testing"
)

// Baby Jubjub is the elliptic curve that used in Iden3, which is designed to work efficiently with zkSNARKs
// https://docs.iden3.io/publications/pdfs/Baby-Jubjub.pdf
func Test_BabyJubGen(t *testing.T) {

	// generate private key randomly
	babyJubjubPrivKey := babyjub.NewRandPrivKey()

	// generate pub key from priKey
	babyJubjubPubKey := babyJubjubPrivKey.Public()

	fmt.Println(babyJubjubPubKey.String())
}
