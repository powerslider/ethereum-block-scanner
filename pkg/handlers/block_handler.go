package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/powerslider/ethereum-block-scanner/pkg/sdk"
)

// BlockHandler represents an HTTP handler for Ethereum block operations.
type BlockHandler struct {
	Parser *sdk.BlockParser
}

// NewBlockHandler initializes a new instance of BlockHandler.
func NewBlockHandler(parser *sdk.BlockParser) *BlockHandler {
	return &BlockHandler{
		Parser: parser,
	}
}

// GetBlockTransactionsPerAddress godoc
// @Summary Get all transactions for a fixed block range given an address.
// @Description Get all transactions for a fixed block range given an address.
// @Tags blocks
// @Accept  json
// @Produce  json
// @Param address path string true "Address"
// @Router /api/v1/{address}/transactions [get]
func (h *BlockHandler) GetBlockTransactionsPerAddress() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)

		address, ok := vars["address"]
		if !ok {
			rw.WriteHeader(http.StatusBadRequest)
		}

		txs, err := h.Parser.GetTransactions(ctx, address)
		if err != nil {
			log.Printf("Could not get transactions for address %s: %v:", address, err)
			rw.WriteHeader(http.StatusBadRequest)
		}

		handleResponse(ctx, rw, txs)
	}
}

// GetCurrentBlock godoc
// @Summary Get current Ethereum block.
// @Description Get current Ethereum block.
// @Tags blocks
// @Accept  json
// @Produce  json
// @Router /api/v1/block/current [get]
func (h *BlockHandler) GetCurrentBlock() http.HandlerFunc {
	type response struct {
		CurrentBlock int `json:"currentBlock"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		blockNum, err := h.Parser.GetCurrentBlock(ctx)
		if err != nil {
			log.Printf("Could not get current block number: %v:", err)
			rw.WriteHeader(http.StatusBadRequest)
		}

		handleResponse(ctx, rw, response{
			CurrentBlock: blockNum,
		})
	}
}

// SubscribeAddress godoc
// @Summary Subscribe and address to an observer for new inbound/outbound transactions in the latest block.
// @Description Subscribe and address to an observer for new inbound/outbound transactions in the latest block.
// @Tags blocks
// @Accept  json
// @Produce  json
// @Param request body handlers.SubscribeAddress.request true "Address"
// @Router /api/v1/address/subscribe [post]
func (h *BlockHandler) SubscribeAddress() http.HandlerFunc {
	type request struct {
		Address string `json:"address"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var reqBody request

		reqBytes, errReqBytes := io.ReadAll(r.Body)
		errReqUnmarshal := json.Unmarshal(reqBytes, &reqBody)

		errReq := errors.Join(errReqBytes, errReqUnmarshal)
		if errReq != nil {
			log.Println(ctx, "Could not unmarshal request params:", errReq)
			rw.WriteHeader(http.StatusBadRequest)
		}

		subscribed := h.Parser.Subscribe(reqBody.Address)

		handleResponse(ctx, rw, subscribed)
	}
}

func handleResponse(ctx context.Context, rw http.ResponseWriter, resp any) {
	jsonResp, errRespMarshal := json.Marshal(resp)
	_, errRespWrite := rw.Write(jsonResp)

	errResp := errors.Join(errRespMarshal, errRespWrite)
	if errResp != nil {
		log.Println(ctx, "Could not write response:", errResp)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
