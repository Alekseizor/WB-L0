package handlers

import (
	"WB-L0/internal/pkg/repository/orders"
	"WB-L0/internal/pkg/sendingjson"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

type OrdersHandler struct {
	OrderRepo orders.OrderRepo //тут будем юзать ин мемори
	Logger    *zap.SugaredLogger
	Send      sendingjson.ServiceSend
}

func (h *OrdersHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vars["ID"]

}
