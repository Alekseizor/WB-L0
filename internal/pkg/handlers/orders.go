package handlers

import (
	"WB-L0/internal/pkg/repository/orders"
	"WB-L0/internal/pkg/sendingjson"
	"database/sql"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

type OrdersHandler struct {
	OrderRepo orders.OrderInMemoryRepo
	Logger    *zap.SugaredLogger
	Send      sendingjson.ServiceSend
}

func (h *OrdersHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	h.Logger.Info(vars["ID"])
	order, err := h.OrderRepo.GetOrderByID(r.Context(), vars["ID"])
	if err == sql.ErrNoRows {
		h.Logger.Infof("url:%s method:%s error: failed to get order - %v", r.URL.Path, r.Method, err)
		http.Error(w, `this order was not found`, http.StatusBadRequest)
		return
	}
	if err != nil {
		h.Logger.Infof("url:%s method:%s error: failed to get order - %v", r.URL.Path, r.Method, err)
		http.Error(w, `failed to receive order by ID`, http.StatusInternalServerError)
		return
	}
	err = h.Send.Sending(w, r, order)
	if err != nil {
		return
	}
}
