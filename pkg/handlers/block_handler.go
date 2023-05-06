package handlers

import (
	"context"
	"encoding/json"
	"errors"
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

func handleResponse(ctx context.Context, rw http.ResponseWriter, resp any) {
	jsonResp, errRespMarshal := json.Marshal(resp)
	_, errRespWrite := rw.Write(jsonResp)

	errResp := errors.Join(errRespMarshal, errRespWrite)
	if errResp != nil {
		log.Println(ctx, "Could not write response:", errResp)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
