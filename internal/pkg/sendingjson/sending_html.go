package sendingjson

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type SendResponseHTML struct {
	Logger *zap.SugaredLogger
}

func (s *SendResponseHTML) Sending(w http.ResponseWriter, r *http.Request, data any) error {
	html := `
		<!DOCTYPE html>
<html>
<head>
    <title>Таблица данных</title>
    <style>
        table {
            border-collapse: collapse;
            width: 100%;
        }
        
        th, td {
            padding: 8px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        
        th {
            background-color: #f2f2f2;
        }
    </style>
</head>
<body>
    <table>
        <tr>
            <th>order_uid</th>
            <th>track_number</th>
            <th>entry</th>
            <th>delivery</th>
            <th>payment</th>
            <th>items</th>
            <th>locale</th>
            <th>internal_signature</th>
            <th>customer_id</th>
            <th>delivery_service</th>
            <th>shardkey</th>
            <th>sm_id</th>
            <th>date_created</th>
            <th>oof_shard</th>
        </tr>
        <!-- Здесь могут быть строки таблицы с данными -->
    </table>
</body>
</html>
	`

	fmt.Fprint(w, html)
	return nil
}

func NewServiceSendHTML(logger *zap.SugaredLogger) *SendResponseHTML {
	return &SendResponseHTML{Logger: logger}
}