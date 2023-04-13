package pkg

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/mocyuto/helm-clear/pkg/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigmapOptions struct {
	DryRun        bool
	History       int
	ChartName     string
	Namespace     string
	ConfigmapName string
}

func RunConfigmap(option ConfigmapOptions) error {
	// revision historyを取得
	histories, err := GetHistory(option.ChartName, option.Namespace, option.History)
	if err != nil {
		return fmt.Errorf("%v. cannot get histories", err)
	}
	revisions := make([]int, 0, len(histories))
	for _, h := range histories {
		revisions = append(revisions, h.Revision)
	}
	sort.Ints(revisions)

	// configmapのリストを取得
	cli, err := k8s.NewClient()
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
	removalConfigmaps := make([]string, 0, len(res.Items))
	for _, i := range res.Items {
		if !strings.Contains(i.Name, option.ConfigmapName) {
			continue
		}
		for _, r := range revisions {
			if !strings.Contains(i.Name, strconv.Itoa(r)) {
				removalConfigmaps = append(removalConfigmaps, i.Name)
				break
			}
		}
	}

	// 列挙したものを削除
	if option.DryRun {
		fmt.Println("dry-run mode: configmaps to be removed")
		for _, c := range removalConfigmaps {
			fmt.Println(c)
		}
		return nil
	}
	for _, c := range removalConfigmaps {
		if err := cli.CoreV1().ConfigMaps(option.Namespace).Delete(context.Background(), c, metav1.DeleteOptions{}); err != nil {
			return err
		}
	}
	return nil
}
