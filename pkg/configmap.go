package pkg

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	v1 "k8s.io/api/core/v1"

	"github.com/mocyuto/helm-clear/pkg/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigmapOptions struct {
	DryRun      bool
	History     int
	ChartName   string
	Namespace   string
	KubeContext string
}

func RunConfigmap(option ConfigmapOptions) error {
	// revision historyを取得
	histories, err := GetHistory(option.ChartName, option.Namespace, option.KubeContext, option.History)
	if err != nil {
		return fmt.Errorf("%v. cannot get histories", err)
	}
	revisions := make([]int, 0, len(histories))
	for _, h := range histories {
		revisions = append(revisions, h.Revision)
	}
	sort.Ints(revisions)

	// 最新のrevisionのconfigmapを取得
	manifests, err := GetManifests(option.ChartName, option.Namespace, option.KubeContext)
	if err != nil {
		return fmt.Errorf("%v. cannot get manifests", err)
	}
	latestConfigMaps := []string{}
	for _, m := range manifests {
		if m.Kind == "ConfigMap" {
			latestConfigMaps = append(latestConfigMaps, m.Metadata.Name)
		}
	}

	// configmapのリストを取得
	cli, err := k8s.NewClient(option.KubeContext)
	if err != nil {
		return err
	}
	res, err := cli.CoreV1().ConfigMaps(option.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	if len(res.Items) == 0 {
		return nil
	}
	removals := removalConfigMaps(res.Items, revisions, latestConfigMaps)

	// 列挙したものを削除
	if option.DryRun {
		fmt.Println("dry-run mode: configmaps to be removed")
		for _, c := range removals {
			fmt.Println(c)
		}
		return nil
	}
	for _, c := range removals {
		if err := cli.CoreV1().ConfigMaps(option.Namespace).Delete(context.Background(), c, metav1.DeleteOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func removalConfigMaps(items []v1.ConfigMap, revisions []int, latestConfigName []string) []string {
	latestRevision := revisions[len(revisions)-1]
	removals := make([]string, 0, len(items))
	for _, c := range latestConfigName {
		keys := strings.Split(c, strconv.Itoa(latestRevision))
		if len(keys) != 2 {
			continue
		}
		for _, item := range items {
			i := strings.Index(item.Name, keys[0])
			j := strings.LastIndex(item.Name, keys[1])
			hasKey := i != -1 && j != -1 && i < j
			if !hasKey {
				continue
			}
			ok := true
			for _, r := range revisions {
				if strings.Contains(item.Name, strconv.Itoa(r)) {
					ok = false
					break
				}
			}
			if ok {
				removals = append(removals, item.Name)
			}
		}
	}
	return removals
}
