package pkg

import (
	"encoding/json"
	"strconv"
	"strings"
)

const notFoundError = "not found"

type History struct {
	Revision   int    `json:"revision"`
	Updated    string `json:"updated"`
	Status     string `json:"status"`
	Chart      string `json:"chart"`
	AppVersion string `json:"app_version"`
}

func GetHistory(name, namespace, kubeCtx string, count int) ([]History, error) {
	res, err := getHistory(name, namespace, kubeCtx, count)
	if err != nil {
		if strings.Contains(err.Error(), notFoundError) {
			return nil, nil
		}
		return nil, err
	}
	var histories []History
	if err := json.Unmarshal(res, &histories); err != nil {
		return nil, err
	}
	return histories, nil
}
func getHistory(name, namespace, kubeCtx string, count int) ([]byte, error) {
	var max int
	if count != 0 {
		max = count
	}
	args := []string{"history", name, "--max", strconv.Itoa(max), "--output", "json"}
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	if kubeCtx != "" {
		args = append(args, "--kube-context", kubeCtx)
	}
	return execute(args)
}
