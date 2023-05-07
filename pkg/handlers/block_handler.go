package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	pkgErrors "github.com/pkg/errors"

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
// @Param blockRange query int false "Block Range" default(0)
// @Router /api/v1/address/{address}/transactions [get]
func (h *BlockHandler) GetBlockTransactionsPerAddress() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		blockRange, err := strconv.Atoi(r.URL.Query().Get("blockRange"))
		if err != nil {
			badRequestError(
				rw,
				pkgErrors.Wrap(err, "invalid query param 'blockRange':"),
			)

			return
		}

		vars := mux.Vars(r)

		address, ok := vars["address"]
		if !ok {
			badRequestError(
				rw,
				errors.New("required path param 'address' is missing"),
			)

			return
		}

		txs, err := h.Parser.GetTransactionsForBlockRange(ctx, address, blockRange)
		if err != nil {
			badRequestError(
				rw,
				pkgErrors.Wrapf(err, "Could not get transactions for address %s", address),
			)

			return
		}

		handleResponse(rw, txs)
	}
}

// GetTransactionsPerSubscriber godoc
// @Summary Get all transactions for a subscribed address.
// @Description Get all transactions for a subscribed address.
// @Tags blocks
// @Accept  json
// @Produce  json
// @Param address path string true "Address"
// @Router /api/v1/subscription/{address}/transactions [get]
func (h *BlockHandler) GetTransactionsPerSubscriber() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		address, ok := vars["address"]
		if !ok {
			badRequestError(
				rw,
				errors.New("required path param 'address' is missing"),
			)

			return
		}

		txs := h.Parser.GetTransactionsPerSubscriber(address)

		handleResponse(rw, txs)
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
			badRequestError(
				rw,
				pkgErrors.Wrap(err, "could not get current block number:"),
			)

			return
		}

		handleResponse(rw, response{
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
		var reqBody request

		reqBytes, errReqBytes := io.ReadAll(r.Body)
		errReqUnmarshal := json.Unmarshal(reqBytes, &reqBody)

		errReq := errors.Join(errReqBytes, errReqUnmarshal)
		if errReq != nil {
			badRequestError(
				rw,
				pkgErrors.Wrap(errReq, "could not unmarshal request params:"),
			)

			return
		}

		subscribed := h.Parser.Subscribe(reqBody.Address)

		handleResponse(rw, subscribed)
	}
}

func handleResponse(rw http.ResponseWriter, resp any) {
	jsonResp, errRespMarshal := json.Marshal(resp)
	_, errRespWrite := rw.Write(jsonResp)

	errResp := errors.Join(errRespMarshal, errRespWrite)
	if errResp != nil {
		http.Error(rw, errResp.Error(), http.StatusInternalServerError)
	}
}

func badRequestError(rw http.ResponseWriter, err error) {
	errBytes, err := json.Marshal(struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
	}{
		Status: http.StatusBadRequest,
		Error:  err.Error(),
	})

	if err == nil {
		http.Error(rw, string(errBytes), http.StatusBadRequest)
	} else {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
