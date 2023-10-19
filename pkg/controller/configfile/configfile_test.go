/*
Copyright 2021 The cert-manager Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package configfile

import (
	"fmt"
	"testing"

	"github.com/cert-manager/cert-manager/pkg/util/configfile"
)

func TestFSLoader_LoadMinimal(t *testing.T) {
	const expectedFilename = "/path/to/config/file"
	const kubeConfigPath = "path/to/kubeconfig/file"

	controllerConfig := New()

	loader, err := configfile.NewConfigurationFSLoader(func(filename string) ([]byte, error) {
		if filename != expectedFilename {
			t.Fatalf("unexpected filename %q passed to ReadFile", filename)
			return nil, fmt.Errorf("unexpected filename %q", filename)
		}
		return []byte(fmt.Sprintf(`apiVersion: controller.config.cert-manager.io/v1alpha1
kind: ControllerConfiguration
kubeConfig: %s`, kubeConfigPath)), nil
	}, expectedFilename)
	if err != nil {
		t.Fatal(err)
	}

	if err := loader.Load(controllerConfig); err != nil {
		t.Fatal(err)
	}

	// the config loader will force paths to be 'absolute' if they are provided as relative.
	absKubeConfigPath := "/path/to/config/path/to/kubeconfig/file"
	if controllerConfig.Config.KubeConfig != absKubeConfigPath {
		t.Errorf("expected kubeConfig to be set to %q but got %q", absKubeConfigPath, controllerConfig.Config.KubeConfig)
	}
}

func TestFSLoader_LoadFull(t *testing.T) {
	const expectedFilename = "/path/to/config/file"
	const kubeConfigPath = "path/to/kubeconfig/file"

	controllerConfig := New()

	loader, err := configfile.NewConfigurationFSLoader(func(filename string) ([]byte, error) {
		if filename != expectedFilename {
			t.Fatalf("unexpected filename %q passed to ReadFile", filename)
			return nil, fmt.Errorf("unexpected filename %q", filename)
		}
		return []byte(`apiVersion: controller.config.cert-manager.io/v1alpha1
kind: ControllerConfiguration
apiServerHost: test-api-server-host
kubeConfig: test-kube-config
kubernetesAPIQPS: 123.456
kubernetesAPIBurst: 123456
namespace: test-namespace
clusterResourceNamespace: test-cluster-resource-namespace
leaderElectionConfig:
  enabled: true
  namespace: test-leader-election-namespace
  leaseDuration: 123456s
  renewDeadline: 123456s
  retryPeriod: 123456s
  healthzTimeout: 123456s
acmeHTTP01Config:
  solverImage: helloz
`)
	if err != nil {
		t.Fatal(err)
	}

	if err := loader.Load(controllerConfig); err != nil {
		t.Fatal(err)
	}

	// the config loader will force paths to be 'absolute' if they are provided as relative.
	absKubeConfigPath := "/path/to/config/path/to/kubeconfig/file"
	if controllerConfig.Config.KubeConfig != absKubeConfigPath {
		t.Errorf("expected kubeConfig to be set to %q but got %q", absKubeConfigPath, controllerConfig.Config.KubeConfig)
	}

	if controllerConfig.Config.ACMEHTTP01Config.SolverImage != "hello" {
		t.Errorf("expected hello but got %q", controllerConfig.Config.ACMEHTTP01Config.SolverImage)
	}
}
