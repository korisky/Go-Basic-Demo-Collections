package idenDemo

import (
	"encoding/json"
	"fmt"
	core "github.com/iden3/go-iden3-core"
	"math/big"
	"testing"
	"time"
)

func Test_GenericClaim(t *testing.T) {

	// set claim expiration date
	expireTime := time.Date(1361, 3, 22, 0, 44, 48, 0, time.UTC)

	// set schema
	ageSchema, _ := core.NewSchemaHashFromHex("2e2d1c11ad3e500de68d7ce16a0a559e")

	// define data slots
	birthday := big.NewInt(199960424)
	documentType := big.NewInt(1)

	// set revocation nonce
	revocationNonce := uint64(1909830690)

	// set ID of the claim subject
	id, _ := core.IDFromString("113TCVw5KMeMp99Qdvub9Mssfz7krL9jWNvbdB7Fd2")

	// create claim
	claim, _ := core.NewClaim(ageSchema, core.WithExpirationDate(expireTime),
		core.WithRevocationNonce(revocationNonce), core.WithIndexID(id),
		core.WithIndexDataInts(birthday, documentType))

	// transform claim from bytes array to json
	claimToMarshal, _ := json.Marshal(claim)
	fmt.Println(string(claimToMarshal))
}
