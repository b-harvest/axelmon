package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetProxyByVal(valoper string) (string, error) {
	url := fmt.Sprintf("https://axelar-api.bharvest.io/query/snapshot/proxy/%s", valoper)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	proxy := Proxy{}
	err = json.Unmarshal(bodyBytes, &proxy)
	if err != nil {
		return "", err
	}

	return proxy.Result.Address, nil
}
