package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
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

type Manifest struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string
	Metadata   struct {
		Namespace   string
		Name        string
		Annotations map[string]string
	}
}

func GetManifests(name, namespace, kubeCtx string) ([]Manifest, error) {
	res, err := getManifests(name, namespace, kubeCtx)
	if err != nil {
		return nil, err
	}
	return parseManifest(res)
}

func parseManifest(res []byte) ([]Manifest, error) {
	manifests := []Manifest{}
	decoder := yaml.NewDecoder(bytes.NewBuffer(res))
	for {
		var m Manifest
		if err := decoder.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("document decode failed: %w", err)
		}
		manifests = append(manifests, m)
	}
	return manifests, nil
}

func getManifests(name, namespace, kubeCtx string) ([]byte, error) {
	args := []string{"get", "manifest", name}
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	if kubeCtx != "" {
		args = append(args, "--kube-context", kubeCtx)
	}
	return execute(args)
}
