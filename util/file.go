package util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type filedata struct {
	dataBytes []byte
}
type ConnInfo struct {
	BTC  []string `json:"btc"`
	ETH  []string `json:"eth"`
	DDMX []string `json:"ddmx"`
	GPE  []string `json:"gpe"`
	USDT []string `json:"usdt"`
	REG  []string `json:"reg"`
	PAY  []string `json:"pay"`
	PWD  []string `json:"pwd"`
}

func GetParamFromPath(path string) (*filedata, error) {
	if path == "" {
		return nil, errors.New("path required.")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	byteVal, _ := ioutil.ReadAll(f)

	return &filedata{
		dataBytes: byteVal,
	}, nil
}
func (fi *filedata) GetConnInfo() (*ConnInfo, error) {
	var conn ConnInfo
	err := json.Unmarshal(fi.dataBytes, &conn)
	if err != nil {
		return nil, err
	}
	return &conn, nil
}
