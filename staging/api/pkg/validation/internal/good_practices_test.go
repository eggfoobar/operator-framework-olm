package internal

import (
	"github.com/operator-framework/api/pkg/manifests"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ValidateGoodPractices(t *testing.T) {
	bundleWithDeploymentSpecEmpty, _ := manifests.GetBundleFromDir("./testdata/valid_bundle")
	bundleWithDeploymentSpecEmpty.CSV.Spec.InstallStrategy.StrategySpec.DeploymentSpecs = nil

	type args struct {
		bundleDir string
		bundle    *manifests.Bundle
	}
	tests := []struct {
		name        string
		args        args
		wantError   bool
		wantWarning bool
		errStrings  []string
		warnStrings []string
	}{
		{
			name: "should pass successfully when the resource request is set for " +
				"all containers defined in the bundle",
			args: args{
				bundleDir: "./testdata/valid_bundle",
			},
		},
		{
			name:        "should raise an waring when the resource request is NOT set for any of the containers defined in the bundle",
			wantWarning: true,
			warnStrings: []string{"Warning: Value memcached-operator.v0.0.1: unable to find the resource requests for the container kube-rbac-proxy. It is recommended to ensure the resource request for CPU and Memory. Be aware that for some clusters configurations it is required to specify requests or limits for those values. Otherwise, the system or quota may reject Pod creation. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/",
				"Warning: Value memcached-operator.v0.0.1: unable to find the resource requests for the container manager. It is recommended to ensure the resource request for CPU and Memory. Be aware that for some clusters configurations it is required to specify requests or limits for those values. Otherwise, the system or quota may reject Pod creation. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/"},
			args: args{
				bundleDir: "./testdata/valid_bundle_v1",
			},
		},
		{
			name:      "should fail when the bundle is nil",
			wantError: true,
			args: args{
				bundle: nil,
			},
			errStrings: []string{"Error: : Bundle is nil"},
		},
		{
			name:      "should fail when the bundle csv is nil",
			wantError: true,
			args: args{
				bundle: &manifests.Bundle{CSV: nil, Name: "test"},
			},
			errStrings: []string{"Error: Value test: Bundle csv is nil"},
		},
		{
			name:      "should fail when the csv.Spec.InstallStrategy.StrategySpec.DeploymentSpecs is nil",
			wantError: true,
			args: args{
				bundle: bundleWithDeploymentSpecEmpty,
			},
			errStrings: []string{"Error: Value etcdoperator.v0.9.4: unable to find a deployment to install in the CSV"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if len(tt.args.bundleDir) > 0 {
				tt.args.bundle, err = manifests.GetBundleFromDir(tt.args.bundleDir)
				require.NoError(t, err)
			}
			results := validateGoodPracticesFrom(tt.args.bundle)
			require.Equal(t, tt.wantWarning, len(results.Warnings) > 0)
			if tt.wantWarning {
				require.Equal(t, len(tt.warnStrings), len(results.Warnings))
				for _, w := range results.Warnings {
					wString := w.Error()
					require.Contains(t, tt.warnStrings, wString)
				}
			}

			require.Equal(t, tt.wantError, len(results.Errors) > 0)
			if tt.wantError {
				require.Equal(t, len(tt.errStrings), len(results.Errors))
				for _, err := range results.Errors {
					errString := err.Error()
					require.Contains(t, tt.errStrings, errString)
				}
			}
		})
	}
}
