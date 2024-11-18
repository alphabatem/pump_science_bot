package services

import (
	"bytes"
	context2 "context"
	"encoding/json"
	"errors"
	"github.com/alphabatem/common/context"
	"github.com/alphabatem/solana-go/rpc_cached"
	"github.com/gagliardetto/solana-go"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

type SwapRequest struct {
	Wallet   string              `json:"wallet"`
	Action   string              `json:"action"`
	Payload  SwapRequestPayload  `json:"payload"`
	Settings SwapRequestSettings `json:"settings"`
}

type SwapRequestPayload struct {
	Source string  `json:"source"`
	Dest   string  `json:"dest"`
	AToB   bool    `json:"aToB"`
	Amount float64 `json:"amount"`
}

type SwapRequestSettings struct {
	Slippage     int    `json:"slippage"`
	ValidatorTip string `json:"validatorTip"`
	PriorityFee  string `json:"priorityFee"`
}

type SwapResponse struct {
	Signed      bool   `json:"signed"`
	Transaction string `json:"transaction"`
}

type SwapService struct {
	context.DefaultService

	client  *rpc_cached.Client
	hClient *http.Client

	keypair solana.PrivateKey

	validatorTip string
	slippage     int
	inAmountSOL  float64
}

const SWAP_SVC = "swap_svc"

func (svc SwapService) Id() string {
	return SWAP_SVC
}

func (svc *SwapService) Start() error {
	var err error

	svc.validatorTip = os.Getenv("VALIDATOR_TIP")
	svc.inAmountSOL, err = strconv.ParseFloat(os.Getenv("IN_AMOUNT"), 64)
	if err != nil {
		return err
	}

	svc.slippage, err = strconv.Atoi(os.Getenv("SLIPPAGE"))
	if err != nil {
		return err
	}

	return nil
}

func (svc *SwapService) Buy(mint solana.PublicKey, amount float64) error {
	txn, err := svc.swap(mint, true, amount)
	if err != nil {
		return err
	}

	return svc.sendTransaction(txn)
}

func (svc *SwapService) Sell(mint solana.PublicKey, amountPct float64) error {
	if amountPct <= 0 || amountPct > 1 {
		return errors.New("amountPct must be between 0 & 1")
	}

	ata, _, _ := solana.FindAssociatedTokenAddress(svc.keypair.PublicKey(), mint)
	resp, err := svc.client.Raw().GetTokenAccountBalance(context2.TODO(), ata, "processed")
	if err != nil {
		return err
	}

	_tokBal, _ := strconv.Atoi(resp.Value.Amount)
	tokBal := float64(_tokBal) / math.Pow10(int(resp.Value.Decimals))

	if tokBal == 0 {
		return errors.New("no tokens")
	}

	tokBalPct := tokBal * amountPct
	txn, err := svc.swap(mint, false, tokBalPct)
	if err != nil {
		return err
	}

	return svc.sendTransaction(txn)
}

func (svc *SwapService) sendTransaction(txn *solana.Transaction) error {
	sig, err := txn.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(svc.keypair.PublicKey()) {
			return &svc.keypair
		}
		return nil
	})
	if err != nil {
		return err
	}

	_, err = svc.client.Raw().SendTransaction(context2.TODO(), txn)
	if err != nil {
		return err
	}

	log.Printf("TXN Sent: %s", sig)
	return nil
}

func (svc *SwapService) swap(mint solana.PublicKey, aToB bool, inAmount float64) (*solana.Transaction, error) {

	source := solana.SolMint.String()
	dest := solana.SolMint.String()
	if aToB {
		dest = mint.String()
	} else {
		source = mint.String()
	}

	data, err := json.Marshal(&SwapRequest{
		Wallet: svc.keypair.PublicKey().String(),
		Action: "swap",
		Payload: SwapRequestPayload{
			Source: source,
			Dest:   dest,
			AToB:   aToB,
			Amount: inAmount,
		},
		Settings: SwapRequestSettings{
			Slippage:     svc.slippage,
			ValidatorTip: svc.validatorTip,
			PriorityFee:  "market",
		},
	})
	if err != nil {
		return nil, err
	}

	resp, err := svc.hClient.Post("https://gateway.fluxbeam.xyz/bot/actions", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rDat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var swapResp SwapResponse
	err = json.Unmarshal(rDat, &swapResp)
	if err != nil {
		return nil, err
	}

	return solana.TransactionFromBase64(swapResp.Transaction)
}
