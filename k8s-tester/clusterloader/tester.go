// Package clusterloader installs clusterloader.
// Replace https://github.com/aws/aws-k8s-tester/tree/v1.5.9/eks/cluster-loader.
package clusterloader

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-k8s-tester/client"
	k8s_tester "github.com/aws/aws-k8s-tester/k8s-tester/tester"
	"github.com/aws/aws-k8s-tester/utils/file"
	"github.com/manifoldco/promptui"
	"go.uber.org/zap"
)

// Config defines parameters for Kubernetes clusterloader tests.
type Config struct {
	Enable bool `json:"enable"`
	Prompt bool `json:"-"`

	Stopc     chan struct{} `json:"-"`
	Logger    *zap.Logger   `json:"-"`
	LogWriter io.Writer     `json:"-"`
	Client    client.Client `json:"-"`

	// MinimumNodes is the minimum number of Kubernetes nodes required for installing this addon.
	MinimumNodes int `json:"minimum_nodes"`

	// ClusterloaderPath is the path to download the "clusterloader".
	ClusterloaderPath string `json:"clusterloader_path"`
	// ClusterloaderDownloadURL is the download URL to download "clusterloader" binary from.
	ClusterloaderDownloadURL string `json:"clusterloader_download_url"`

	Provider string `json:"provider"`

	// Runs is the number of "clusterloader2" runs back-to-back.
	Runs int `json:"runs"`
	// RunTimeout is the timeout for the total test runs.
	RunTimeout       time.Duration `json:"run_timeout"`
	RunTimeoutString string        `json:"run_timeout_string" read-only:"true"`

	// TestConfigPath is the clusterloader2 test configuration file.
	// Must be located along with other configuration files.
	// e.g. ${HOME}/go/src/k8s.io/perf-tests/clusterloader2/testing/load/config.yaml
	// ref. https://github.com/kubernetes/perf-tests/blob/master/clusterloader2/testing/load/config.yaml
	// Set via "--testconfig" flag.
	TestConfigPath string `json:"test_config_path"`
	// ReportDir is the clusterloader2 test report output directory.
	// Set via "--report-dir" flag.
	ReportDir string `json:"report_dir"`

	// RunFromCluster is set 'true' to override KUBECONFIG set in "Client" field.
	// If "false", instead pass Client.Config().KubeconfigPath to "--kubeconfig" flag.
	// Set via "--run-from-cluster" flag.
	// ref. https://github.com/kubernetes/perf-tests/pull/1295
	RunFromCluster bool `json:"run_from_cluster"`
	// Nodes is the number of nodes.
	// Set via "--nodes" flag.
	Nodes int `json:"nodes"`

	// PodStartupLatency is the result of clusterloader runs.
	PodStartupLatency PerfData `json:"pod_startup_latency" read-only:"true"`

	// TestOverride defines "testoverrides" flag values.
	// Set via "--testoverrides" flag.
	// See https://github.com/kubernetes/perf-tests/tree/master/clusterloader2/testing/overrides for more.
	// ref. https://github.com/kubernetes/perf-tests/pull/1345
	TestOverride *TestOverride `json:"test_override"`
}

func (cfg *Config) ValidateAndSetDefaults() error {
	if cfg.ClusterloaderPath == "" {
		cfg.ClusterloaderPath = DefaultClusterloaderPath()
	}
	if cfg.ClusterloaderDownloadURL == "" {
		cfg.ClusterloaderDownloadURL = DefaultClusterloaderDownloadURL()
	}

	if cfg.Runs == 0 {
		return fmt.Errorf("invalid Runs %d", cfg.Runs)
	}
	if cfg.RunTimeout == time.Duration(0) {
		cfg.RunTimeout = DefaultRunTimeout
	}
	cfg.RunTimeoutString = cfg.RunTimeout.String()

	if !file.Exist(cfg.TestConfigPath) {
		return fmt.Errorf("TestConfigPath %q does not exist", cfg.TestConfigPath)
	}

	if cfg.Nodes == 0 {
		cfg.Nodes = cfg.MinimumNodes
	}

	return nil
}

const (
	DefaultMinimumNodes int = 1

	DefaultRuns       = 2
	DefaultRunTimeout = 30 * time.Minute

	DefaultRunFromCluster = false
	DefaultNodes          = 10
)

func NewDefault() *Config {
	return &Config{
		Enable:       false,
		Prompt:       false,
		MinimumNodes: DefaultMinimumNodes,

		ClusterloaderPath:        DefaultClusterloaderPath(),
		ClusterloaderDownloadURL: DefaultClusterloaderDownloadURL(),

		Provider: DefaultProvider,

		Runs:       DefaultRuns,
		RunTimeout: DefaultRunTimeout,

		RunFromCluster: DefaultRunFromCluster,
		Nodes:          DefaultNodes,

		TestOverride: newDefaultTestOverride(),
	}
}

func New(cfg *Config) k8s_tester.Tester {
	return &tester{
		cfg: cfg,
	}
}

type tester struct {
	cfg *Config
}

var pkgName = path.Base(reflect.TypeOf(tester{}).PkgPath())

func Env() string {
	return "ADD_ON_" + strings.ToUpper(strings.Replace(pkgName, "-", "_", -1))
}

func EnvTestOverride() string {
	return Env() + "_TEST_OVERRIDE"
}

func (ts *tester) Name() string { return pkgName }

func (ts *tester) Enabled() bool { return ts.cfg.Enable }

func (ts *tester) Apply() error {
	if ok := ts.runPrompt("apply"); !ok {
		return errors.New("cancelled")
	}

	if nodes, err := client.ListNodes(ts.cfg.Client.KubernetesClient()); len(nodes) < ts.cfg.MinimumNodes || err != nil {
		return fmt.Errorf("failed to validate minimum nodes requirement %d (nodes %v, error %v)", ts.cfg.MinimumNodes, len(nodes), err)
	}

	if err := installClusterloader(ts.cfg.Logger, ts.cfg.ClusterloaderPath, ts.cfg.ClusterloaderDownloadURL); err != nil {
		return err
	}

	if err := os.MkdirAll(ts.cfg.ReportDir, 0700); err != nil {
		return err
	}
	if err := file.IsDirWriteable(ts.cfg.ReportDir); err != nil {
		return err
	}
	ts.cfg.Logger.Info("created report dir", zap.String("dir", ts.cfg.ReportDir))

	if err := ts.cfg.TestOverride.Sync(); err != nil {
		return err
	}
	ts.cfg.Logger.Info("wrote test override file", zap.String("path", ts.cfg.TestOverride.Path))

	return nil
}

func (ts *tester) Delete() error {
	if ok := ts.runPrompt("delete"); !ok {
		return errors.New("cancelled")
	}

	var errs []string

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}

func (ts *tester) runPrompt(action string) (ok bool) {
	if ts.cfg.Prompt {
		msg := fmt.Sprintf("Ready to %q resources, should we continue?", action)
		prompt := promptui.Select{
			Label: msg,
			Items: []string{
				"No, cancel it!",
				fmt.Sprintf("Yes, let's %q!", action),
			},
		}
		idx, answer, err := prompt.Run()
		if err != nil {
			panic(err)
		}
		if idx != 1 {
			fmt.Printf("cancelled %q [index %d, answer %q]\n", action, idx, answer)
			return false
		}
	}
	return true
}

func (ts *tester) getCL2Args() (args []string) {
	args = []string{
		ts.cfg.ClusterloaderPath,
		"--logtostderr",     // log to standard error instead of files (default true)
		"--alsologtostderr", // log to standard error as well as files
		"--testconfig=" + ts.cfg.TestConfigPath,
		"--testoverrides=" + ts.cfg.TestOverride.Path,
		"--report-dir=" + ts.cfg.ReportDir,
		"--nodes=" + fmt.Sprintf("%d", ts.cfg.Nodes),
		"--provider=" + ts.cfg.Provider,
	}
	if ts.cfg.RunFromCluster {
		// ref. https://github.com/kubernetes/perf-tests/pull/1295
		args = append(args, "--run-from-cluster=true")
	} else {
		args = append(args, "--kubeconfig="+ts.cfg.Client.Config().KubeconfigPath)
	}
	return args
}
