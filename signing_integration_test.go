// Optional live Testnet checks: set MERIDIAN_INTEGRATION=1 (and MERIDIAN_PK if required).
// The legacy ETHEREAL_INTEGRATION / ETHEREAL_PK names are still accepted.
package etherealRest

import (
	"cmp"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	abi "github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func TestLive_orderSigning_againstTestnet(t *testing.T) {
	if cmp.Or(os.Getenv("MERIDIAN_INTEGRATION"), os.Getenv("ETHEREAL_INTEGRATION")) != "1" {
		t.Skip("set MERIDIAN_INTEGRATION=1 to run live Testnet signing checks")
	}
	pk := cmp.Or(os.Getenv("MERIDIAN_PK"), os.Getenv("ETHEREAL_PK"))
	if pk == "" {
		pk = "0bb5d63b84421e1268dda020818ae30cf26e7f10e321fb820a8aa69216dea92a"
	}

	order := Order{
		Sender:     "0xdeadbeef00000000000000000000000000000000",
		Subaccount: "0x123456789abcde00000000000000000000000000000000000000000000000000",
		Quantity:   "1",
		Price:      "3000",
		ReduceOnly: false,
		Side:       BUY,
		EngineType: PERPETUAL,
		OnchainID:  2,
		Nonce:      "1764897077655477722",
		SignedAt:   int64(1764897077),
	}

	ctx := context.Background()
	client, err := NewClient(ctx, pk, Testnet)
	if err != nil {
		t.Fatal(err)
	}

	domainHashString, err := client.InitDomain(ctx)
	if err != nil {
		t.Fatal(err)
	}
	expectedDomainHash := "f45e2c25657b6ad7b4835073c1c65bec7f6b973d6064b491dd248315296e7d6a"
	if domainHashString != expectedDomainHash {
		t.Fatalf("domain hash drift (API may have changed EIP-712 domain): got %s want %s", domainHashString, expectedDomainHash)
	}

	msg, err := order.ToMessage()
	if err != nil {
		t.Fatal(err)
	}
	messageHash, err := client.GetTypes().HashStruct("TradeOrder", msg)
	if err != nil {
		t.Fatal(err)
	}

	signer := &integrationTestSigner{pk: mustIntegrationKey(t, pk), types: client.GetTypes()}

	signature, err := Sign(&order, "TradeOrder", signer)
	if err != nil {
		t.Fatal(err)
	}

	domainBytes, err := reverseHexIntegration(domainHashString)
	if err != nil {
		t.Fatal(err)
	}
	fullHash := MakeFullHash(domainBytes, messageHash)
	t.Logf("message hash: %s", common.Bytes2Hex(messageHash))
	t.Logf("full hash: %s", common.Bytes2Hex(fullHash))

	expectedSignature := "0xbc76ecd2bcd2bab0c6dc989d238bcd0216c4b574b2de0d0141ae205aabdd53532206441c3593e34de96cbe2e1bff501095a73cb2f2dcb05d1e2b03f2db6b24461c"
	if signature != expectedSignature {
		t.Fatalf("signature drift: got %s want %s", signature, expectedSignature)
	}

	payload := SignedMessage[*Order]{Data: &order, Signature: signature}
	payloadJSON, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("signed payload:\n%s", string(payloadJSON))
}

type integrationTestSigner struct {
	pk    *ecdsa.PrivateKey
	types *abi.TypedData
}

func mustIntegrationKey(t *testing.T, pk string) *ecdsa.PrivateKey {
	t.Helper()
	k, err := crypto.HexToECDSA(strings.TrimPrefix(pk, "0x"))
	if err != nil {
		t.Fatal(err)
	}
	return k
}

func (s *integrationTestSigner) GetPk() *ecdsa.PrivateKey {
	return s.pk
}

func (s *integrationTestSigner) GetTypes() *abi.TypedData {
	return s.types
}

func reverseHexIntegration(s string) ([]byte, error) {
	clean := strings.TrimPrefix(s, "0x")
	return hex.DecodeString(clean)
}
