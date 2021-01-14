package handlers

import (
	"encoding/json"
	"net/http"

	auth "../authorization"
	node "../node"
	pay "../payment"
	util "../util"
)

var connPath string = "./conn.conf"

// var connPath string = "/home/jtn/payment_core/conn.conf"

type respReturn struct {
	Data   interface{} `json:"data,omitempty"`
	Status respStatus  `json:"status,omitempty"`
}
type respStatus struct {
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	reqService, err := auth.RegProcessRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		url, err := getFuncUrl("REG")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		pwd, err := getFuncUrl("PWD")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := reqService.ProcessMethod(url, pwd)
		processResponse(resp, err, w)
	}
}
func Payment(w http.ResponseWriter, r *http.Request) {
	reqService, err := pay.PayProcessRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		url, err := getFuncUrl("PAY")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := reqService.PayMethod(url)
		processResponse(resp, err, w)
	}
}
func Trade(w http.ResponseWriter, r *http.Request) {
	reqService, err := pay.TradeProcessRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		url, err := getFuncUrl("PAY")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := reqService.TradeMethod(url)
		processResponse(resp, err, w)
	}
}
func Transfer(w http.ResponseWriter, r *http.Request) {
	reqService, err := pay.TransferProcessRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		url, err := getFuncUrl("PAY")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := reqService.TransferMethod(url)
		processResponse(resp, err, w)
	}
}
func Offline(w http.ResponseWriter, r *http.Request) {
	reqService, err := pay.OfflineProcessRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		url, err := getFuncUrl("PAY")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := reqService.OfflineMethod(url)
		processResponse(resp, err, w)
	}
}
func Nolimit(w http.ResponseWriter, r *http.Request) {
	reqService, err := pay.NolimitProcessRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		url, err := getFuncUrl("PAY")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := reqService.NolimitMethod(url)
		processResponse(resp, err, w)
	}
}

func processResponse(respService interface{}, err error, w http.ResponseWriter) {
	var resp respReturn
	if err != nil {
		resp.Data = new(interface{})
		resp.Status.Code = "500"
		resp.Status.Message = err.Error()
	} else {
		resp.Data = respService
		resp.Status.Code = "0"
		resp.Status.Message = "success"
	}
	resp.Status.Timestamp = util.Timestamp()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
	}
}
func getFuncUrl(funcName string) (string, error) {
	var conn string = ""
	fileData, err := util.GetParamFromPath(connPath)
	if err != nil {
		return conn, err
	}
	connInfo, err := fileData.GetConnInfo()
	if err != nil {
		return conn, err
	}

	switch funcName {
	case "GPE":
		client := node.Init(connInfo.GPE)
		connect := client.GPEConnect()
		conn = connect.URL()
	case "DDMX":
		client := node.Init(connInfo.DDMX)
		connect := client.DDMXConnect()
		conn = connect.URL()
	case "BTC":
		client := node.Init(connInfo.BTC)
		connect := client.BTCConnect()
		conn = connect.URL()
	case "USDT":
		client := node.Init(connInfo.USDT)
		connect := client.USDTConnect()
		conn = connect.URL()
	case "ETH":
		client := node.Init(connInfo.ETH)
		connect := client.ETHConnect()
		conn = connect.URL()
	case "REG":
		client := node.Init(connInfo.REG)
		connect := client.RegConnect()
		conn = connect.URL()
	case "PAY":
		client := node.Init(connInfo.PAY)
		connect := client.PayConnect()
		conn = connect.URL()
	case "PWD":
		client := node.Init(connInfo.PWD)
		connect := client.PwdConnect()
		conn = connect.URL()
	}
	return conn, nil
}
