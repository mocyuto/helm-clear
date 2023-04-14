package pkg

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "k8s.io/api/core/v1"

	"github.com/google/go-cmp/cmp"
)

func Test_ParseManifest(t *testing.T) {
	text := `---
# Source: api/templates/serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: api
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::123123123:role/stg-api
---
# Source: -api/templates/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: api-1
  annotations:
    helm.sh/resource-policy: "keep"
data:
  config.yaml: |-
    db:
      user: user
      password: password
`
	want := []Manifest{
		{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
			Metadata: struct {
				Namespace   string
				Name        string
				Annotations map[string]string
			}{
				Namespace:   "",
				Name:        "api",
				Annotations: map[string]string{"eks.amazonaws.com/role-arn": "arn:aws:iam::123123123:role/stg-api"},
			},
		},
		{
			APIVersion: "v1",
			Kind:       "ConfigMap",
			Metadata: struct {
				Namespace   string
				Name        string
				Annotations map[string]string
			}{
				Namespace:   "",
				Name:        "api-1",
				Annotations: map[string]string{"helm.sh/resource-policy": "keep"},
			},
		},
	}
	m, err := parseManifest([]byte(text))
	if err != nil {
		t.Fatal(err)
	}
	if d := cmp.Diff(m, want); d != "" {
		t.Fatalf("diff: %s", d)
	}
}

func Test_removalConfigMaps(t *testing.T) {
	type args struct {
		manifests        []v1.ConfigMap
		revisions        []int
		latestConfigName []string
	}
	tests := map[string]struct {
		args args
		want []string
	}{
		"test1": {
			args: args{
				manifests: []v1.ConfigMap{
					{ObjectMeta: metav1.ObjectMeta{Name: "api-1"}},
					{ObjectMeta: metav1.ObjectMeta{Name: "api-2"}},
					{ObjectMeta: metav1.ObjectMeta{Name: "api-3"}},
					{ObjectMeta: metav1.ObjectMeta{Name: "api-4"}},
				},
				revisions:        []int{3, 4},
				latestConfigName: []string{"api-4"},
			},
			want: []string{"api-1", "api-2"},
		},
		"test2": {
			args: args{
				manifests: []v1.ConfigMap{
					{ObjectMeta: metav1.ObjectMeta{Name: "api-1-dev"}},
					{ObjectMeta: metav1.ObjectMeta{Name: "api-2-dev"}},
					{ObjectMeta: metav1.ObjectMeta{Name: "api-3-dev"}},
					{ObjectMeta: metav1.ObjectMeta{Name: "api-4-dev"}},
					{ObjectMeta: metav1.ObjectMeta{Name: "api-4-prd"}},
				},
				revisions:        []int{3, 4},
				latestConfigName: []string{"api-4-dev"},
			},
			want: []string{"api-1-dev", "api-2-dev"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := removalConfigMaps(tt.args.manifests, tt.args.revisions, tt.args.latestConfigName)
			if d := cmp.Diff(got, tt.want); d != "" {
				t.Errorf("diff: %s", d)
			}
		})
	}

}
