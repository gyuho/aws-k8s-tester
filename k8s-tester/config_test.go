package k8s_tester

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestEnv(t *testing.T) {
	cfg := NewDefault()

	os.Setenv("K8S_TESTER_CONFIG_PATH", "test.yaml")
	defer os.Unsetenv("K8S_TESTER_CONFIG_PATH")
	os.Setenv("K8S_TESTER_PROMPT", "false")
	defer os.Unsetenv("K8S_TESTER_PROMPT")
	os.Setenv("K8S_TESTER_CLUSTER_NAME", "hello")
	defer os.Unsetenv("K8S_TESTER_CLUSTER_NAME")
	os.Setenv("K8S_TESTER_CLIENTS", "100")
	defer os.Unsetenv("K8S_TESTER_CLIENTS")
	os.Setenv("K8S_TESTER_CLIENT_TIMEOUT", "100m")
	defer os.Unsetenv("K8S_TESTER_CLIENT_TIMEOUT")
	os.Setenv("K8S_TESTER_KUBECTL_DOWNLOAD_URL", "hello.url")
	defer os.Unsetenv("K8S_TESTER_KUBECTL_DOWNLOAD_URL")
	os.Setenv("K8S_TESTER_KUBECONFIG_PATH", "hello.config")
	defer os.Unsetenv("K8S_TESTER_KUBECONFIG_PATH")
	os.Setenv("K8S_TESTER_KUBECONFIG_CONTEXT", "hello.ctx")
	defer os.Unsetenv("K8S_TESTER_KUBECONFIG_CONTEXT")

	if err := cfg.UpdateFromEnvs(); err != nil {
		t.Fatal(err)
	}

	if cfg.ConfigPath != "test.yaml" {
		t.Fatalf("unexpected cfg.ConfigPath %v", cfg.ConfigPath)
	}
	if cfg.Prompt {
		t.Fatalf("unexpected cfg.Prompt %v", cfg.Prompt)
	}
	if cfg.ClusterName != "hello" {
		t.Fatalf("unexpected cfg.ClusterName %v", cfg.ClusterName)
	}
	if cfg.Clients != 100 {
		t.Fatalf("unexpected cfg.Clients %v", cfg.Clients)
	}
	if cfg.ClientTimeout != 100*time.Minute {
		t.Fatalf("unexpected cfg.ClientTimeout %v", cfg.ClientTimeout)
	}
	if cfg.KubectlDownloadURL != "hello.url" {
		t.Fatalf("unexpected cfg.KubectlDownloadURL %v", cfg.KubectlDownloadURL)
	}
	if cfg.KubeconfigPath != "hello.config" {
		t.Fatalf("unexpected cfg.KubeconfigPath %v", cfg.KubeconfigPath)
	}
	if cfg.KubeconfigContext != "hello.ctx" {
		t.Fatalf("unexpected cfg.KubeconfigContext %v", cfg.KubeconfigContext)
	}
}

func TestEnvAddOnCloudwatchAgent(t *testing.T) {
	cfg := NewDefault()

	os.Setenv("K8S_TESTER_CLUSTER_NAME", "new-name")
	defer os.Unsetenv("K8S_TESTER_CLUSTER_NAME")
	os.Setenv("K8S_TESTER_ADD_ON_CLOUDWATCH_AGENT_ENABLE", "true")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CLOUDWATCH_AGENT_ENABLE")
	os.Setenv("K8S_TESTER_ADD_ON_CLOUDWATCH_AGENT_NAMESPACE", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CLOUDWATCH_AGENT_NAMESPACE")

	if err := cfg.UpdateFromEnvs(); err != nil {
		t.Fatal(err)
	}
	if err := cfg.ValidateAndSetDefaults(); err != nil {
		t.Fatal(err)
	}
	os.RemoveAll(cfg.ConfigPath)

	if cfg.ClusterName != "new-name" {
		t.Fatalf("unexpected cfg.ClusterName %v", cfg.AddOnCloudwatchAgent.Enable)
	}
	if !cfg.AddOnCloudwatchAgent.Enable {
		t.Fatalf("unexpected cfg.AddOnCloudwatchAgent.Enable %v", cfg.AddOnCloudwatchAgent.Enable)
	}
	if cfg.AddOnCloudwatchAgent.Namespace != "hello" {
		t.Fatalf("unexpected cfg.AddOnCloudwatchAgent.Namespace %v", cfg.AddOnCloudwatchAgent.Namespace)
	}
	// should be overwritten since it's a read-only field
	if cfg.AddOnCloudwatchAgent.ClusterName != "new-name" {
		t.Fatalf("unexpected cfg.AddOnCloudwatchAgent.ClusterName %v", cfg.AddOnCloudwatchAgent.ClusterName)
	}
}

func TestEnvAddOnFluentBit(t *testing.T) {
	cfg := NewDefault()

	os.Setenv("K8S_TESTER_ADD_ON_FLUENT_BIT_ENABLE", "true")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_FLUENT_BIT_ENABLE")
	os.Setenv("K8S_TESTER_ADD_ON_FLUENT_BIT_NAMESPACE", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_FLUENT_BIT_NAMESPACE")

	if err := cfg.UpdateFromEnvs(); err != nil {
		t.Fatal(err)
	}

	if !cfg.AddOnFluentBit.Enable {
		t.Fatalf("unexpected cfg.AddOnFluentBit.Enable %v", cfg.AddOnFluentBit.Enable)
	}
	if cfg.AddOnFluentBit.Namespace != "hello" {
		t.Fatalf("unexpected cfg.AddOnFluentBit.Namespace %v", cfg.AddOnFluentBit.Namespace)
	}
}

func TestEnvAddOnMetricsServer(t *testing.T) {
	cfg := NewDefault()

	os.Setenv("K8S_TESTER_ADD_ON_METRICS_SERVER_ENABLE", "true")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_METRICS_SERVER_ENABLE")

	if err := cfg.UpdateFromEnvs(); err != nil {
		t.Fatal(err)
	}

	if !cfg.AddOnMetricsServer.Enable {
		t.Fatalf("unexpected cfg.AddOnMetricsServer.Enable %v", cfg.AddOnMetricsServer.Enable)
	}
}

func TestEnvAddOnConformance(t *testing.T) {
	cfg := NewDefault()

	os.Setenv("K8S_TESTER_ADD_ON_CONFORMANCE_ENABLE", "true")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFORMANCE_ENABLE")
	os.Setenv("K8S_TESTER_ADD_ON_CONFORMANCE_NAMESPACE", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFORMANCE_NAMESPACE")
	os.Setenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_PATH", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_PATH")
	os.Setenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_DOWNLOAD_URL", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_DOWNLOAD_URL")
	os.Setenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RUN_TIMEOUT", "1h")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RUN_TIMEOUT")
	os.Setenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_DELETE_TIMEOUT", "1h")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_DELETE_TIMEOUT")
	os.Setenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RUN_MODE", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RUN_MODE")
	os.Setenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RUN_E2E_FOCUS", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RUN_E2E_FOCUS")
	os.Setenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RUN_E2E_SKIP", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RUN_E2E_SKIP")
	os.Setenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RUN_KUBE_CONFORMANCE_IMAGE", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RUN_KUBE_CONFORMANCE_IMAGE")
	os.Setenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RUN_E2E_REPO_CONFIG", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RUN_E2E_REPO_CONFIG")
	os.Setenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RUN_IMAGE", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RUN_IMAGE")
	os.Setenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RUN_SYSTEMD_LOGS_IMAGE", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RUN_SYSTEMD_LOGS_IMAGE")
	os.Setenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RESULTS_TAR_GZ_PATH", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RESULTS_TAR_GZ_PATH")
	os.Setenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RESULTS_E2E_LOG_PATH", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RESULTS_E2E_LOG_PATH")
	os.Setenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RESULTS_JUNIT_XML_PATH", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RESULTS_JUNIT_XML_PATH")
	os.Setenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RESULTS_OUTPUT_DIR", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFORMANCE_SONOBUOY_RESULTS_OUTPUT_DIR")

	if err := cfg.UpdateFromEnvs(); err != nil {
		t.Fatal(err)
	}

	if !cfg.AddOnConformance.Enable {
		t.Fatalf("unexpected cfg.AddOnConformance.Enable %v", cfg.AddOnConformance.Enable)
	}
	if cfg.AddOnConformance.Namespace != "hello" {
		t.Fatalf("unexpected cfg.AddOnConformance.Namespace %v", cfg.AddOnConformance.Namespace)
	}
	if cfg.AddOnConformance.SonobuoyPath != "hello" {
		t.Fatalf("unexpected cfg.AddOnConformance.SonobuoyPath %v", cfg.AddOnConformance.SonobuoyPath)
	}
	if cfg.AddOnConformance.SonobuoyDownloadURL != "hello" {
		t.Fatalf("unexpected cfg.AddOnConformance.SonobuoyDownloadURL %v", cfg.AddOnConformance.SonobuoyDownloadURL)
	}
	if cfg.AddOnConformance.SonobuoyRunTimeout != time.Hour {
		t.Fatalf("unexpected cfg.AddOnConformance.SonobuoyRunTimeout %v", cfg.AddOnConformance.SonobuoyRunTimeout)
	}
	if cfg.AddOnConformance.SonobuoyDeleteTimeout != time.Hour {
		t.Fatalf("unexpected cfg.AddOnConformance.SonobuoyDeleteTimeout %v", cfg.AddOnConformance.SonobuoyDeleteTimeout)
	}
	if cfg.AddOnConformance.SonobuoyRunMode != "hello" {
		t.Fatalf("unexpected cfg.AddOnConformance.SonobuoyRunMode %v", cfg.AddOnConformance.SonobuoyRunMode)
	}
	if cfg.AddOnConformance.SonobuoyRunE2EFocus != "hello" {
		t.Fatalf("unexpected cfg.AddOnConformance.SonobuoyRunE2EFocus %v", cfg.AddOnConformance.SonobuoyRunE2EFocus)
	}
	if cfg.AddOnConformance.SonobuoyRunE2ESkip != "hello" {
		t.Fatalf("unexpected cfg.AddOnConformance.SonobuoyRunE2ESkip %v", cfg.AddOnConformance.SonobuoyRunE2ESkip)
	}
	if cfg.AddOnConformance.SonobuoyRunKubeConformanceImage != "hello" {
		t.Fatalf("unexpected cfg.AddOnConformance.SonobuoyRunKubeConformanceImage %v", cfg.AddOnConformance.SonobuoyRunKubeConformanceImage)
	}
	if cfg.AddOnConformance.SonobuoyRunE2ERepoConfig != "hello" {
		t.Fatalf("unexpected cfg.AddOnConformance.SonobuoyRunE2ERepoConfig %v", cfg.AddOnConformance.SonobuoyRunE2ERepoConfig)
	}
	if cfg.AddOnConformance.SonobuoyRunImage != "hello" {
		t.Fatalf("unexpected cfg.AddOnConformance.SonobuoyRunImage %v", cfg.AddOnConformance.SonobuoyRunImage)
	}
	if cfg.AddOnConformance.SonobuoyRunSystemdLogsImage != "hello" {
		t.Fatalf("unexpected cfg.AddOnConformance.SonobuoyRunSystemdLogsImage %v", cfg.AddOnConformance.SonobuoyRunSystemdLogsImage)
	}
	if cfg.AddOnConformance.SonobuoyResultsTarGzPath != "hello" {
		t.Fatalf("unexpected cfg.AddOnConformance.SonobuoyResultsTarGzPath %v", cfg.AddOnConformance.SonobuoyResultsTarGzPath)
	}
	if cfg.AddOnConformance.SonobuoyResultsE2ELogPath != "hello" {
		t.Fatalf("unexpected cfg.AddOnConformance.SonobuoyResultsE2ELogPath %v", cfg.AddOnConformance.SonobuoyResultsE2ELogPath)
	}
	if cfg.AddOnConformance.SonobuoyResultsJunitXMLPath != "hello" {
		t.Fatalf("unexpected cfg.AddOnConformance.SonobuoyResultsJunitXMLPath %v", cfg.AddOnConformance.SonobuoyResultsJunitXMLPath)
	}
	if cfg.AddOnConformance.SonobuoyResultsOutputDir != "hello" {
		t.Fatalf("unexpected cfg.AddOnConformance.SonobuoyResultsOutputDir %v", cfg.AddOnConformance.SonobuoyResultsOutputDir)
	}
}

func TestEnvAddOnCSIEBS(t *testing.T) {
	cfg := NewDefault()

	os.Setenv("K8S_TESTER_ADD_ON_CSI_EBS_ENABLE", "true")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CSI_EBS_ENABLE")

	if err := cfg.UpdateFromEnvs(); err != nil {
		t.Fatal(err)
	}

	if !cfg.AddOnCSIEBS.Enable {
		t.Fatalf("unexpected cfg.AddOnCSIEBS.Enable %v", cfg.AddOnCSIEBS.Enable)
	}
}

func TestEnvAddOnKubernetesDashboard(t *testing.T) {
	cfg := NewDefault()

	os.Setenv("K8S_TESTER_ADD_ON_KUBERNETES_DASHBOARD_ENABLE", "true")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_KUBERNETES_DASHBOARD_ENABLE")
	os.Setenv("K8S_TESTER_ADD_ON_KUBERNETES_DASHBOARD_MINIMUM_NODES", "10")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_KUBERNETES_DASHBOARD_MINIMUM_NODES")

	if err := cfg.UpdateFromEnvs(); err != nil {
		t.Fatal(err)
	}

	if !cfg.AddOnKubernetesDashboard.Enable {
		t.Fatalf("unexpected cfg.AddOnKubernetesDashboard.Enable %v", cfg.AddOnKubernetesDashboard.Enable)
	}
	if cfg.AddOnKubernetesDashboard.MinimumNodes != 10 {
		t.Fatalf("unexpected cfg.AddOnKubernetesDashboard.MinimumNodes %v", cfg.AddOnKubernetesDashboard.MinimumNodes)
	}
}

func TestEnvAddOnPHPApache(t *testing.T) {
	cfg := NewDefault()

	os.Setenv("K8S_TESTER_ADD_ON_PHP_APACHE_ENABLE", "true")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_PHP_APACHE_ENABLE")
	os.Setenv("K8S_TESTER_ADD_ON_PHP_APACHE_MINIMUM_NODES", "100")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_PHP_APACHE_MINIMUM_NODES")
	os.Setenv("K8S_TESTER_ADD_ON_PHP_APACHE_NAMESPACE", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_PHP_APACHE_NAMESPACE")
	os.Setenv("K8S_TESTER_ADD_ON_PHP_APACHE_REPOSITORY_PARTITION", "aws")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_PHP_APACHE_REPOSITORY_PARTITION")
	os.Setenv("K8S_TESTER_ADD_ON_PHP_APACHE_REPOSITORY_ACCOUNT_ID", "123")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_PHP_APACHE_REPOSITORY_ACCOUNT_ID")
	os.Setenv("K8S_TESTER_ADD_ON_PHP_APACHE_REPOSITORY_REGION", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_PHP_APACHE_REPOSITORY_REGION")
	os.Setenv("K8S_TESTER_ADD_ON_PHP_APACHE_REPOSITORY_NAME", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_PHP_APACHE_REPOSITORY_NAME")
	os.Setenv("K8S_TESTER_ADD_ON_PHP_APACHE_REPOSITORY_IMAGE_TAG", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_PHP_APACHE_REPOSITORY_IMAGE_TAG")
	os.Setenv("K8S_TESTER_ADD_ON_PHP_APACHE_DEPLOYMENT_NODE_SELECTOR", `{"a":"b","c":"d"}`)
	defer os.Unsetenv("K8S_TESTER_ADD_ON_PHP_APACHE_DEPLOYMENT_NODE_SELECTOR")

	if err := cfg.UpdateFromEnvs(); err != nil {
		t.Fatal(err)
	}

	if !cfg.AddOnPHPApache.Enable {
		t.Fatalf("unexpected cfg.AddOnPHPApache.Enable %v", cfg.AddOnPHPApache.Enable)
	}
	if cfg.AddOnPHPApache.MinimumNodes != 100 {
		t.Fatalf("unexpected cfg.AddOnPHPApache.MinimumNodes %v", cfg.AddOnPHPApache.MinimumNodes)
	}
	if cfg.AddOnPHPApache.Namespace != "hello" {
		t.Fatalf("unexpected cfg.AddOnPHPApache.Namespace %v", cfg.AddOnPHPApache.Namespace)
	}
	if cfg.AddOnPHPApache.RepositoryPartition != "aws" {
		t.Fatalf("unexpected cfg.AddOnPHPApache.RepositoryPartition %v", cfg.AddOnPHPApache.RepositoryPartition)
	}
	if cfg.AddOnPHPApache.RepositoryAccountID != "123" {
		t.Fatalf("unexpected cfg.AddOnPHPApache.RepositoryAccountID %v", cfg.AddOnPHPApache.RepositoryAccountID)
	}
	if cfg.AddOnPHPApache.RepositoryRegion != "hello" {
		t.Fatalf("unexpected cfg.AddOnPHPApache.RepositoryRegion %v", cfg.AddOnPHPApache.RepositoryRegion)
	}
	if cfg.AddOnPHPApache.RepositoryName != "hello" {
		t.Fatalf("unexpected cfg.AddOnPHPApache.RepositoryName %v", cfg.AddOnPHPApache.RepositoryName)
	}
	if cfg.AddOnPHPApache.RepositoryImageTag != "hello" {
		t.Fatalf("unexpected cfg.AddOnPHPApache.RepositoryImageTag %v", cfg.AddOnPHPApache.RepositoryImageTag)
	}
	if !reflect.DeepEqual(cfg.AddOnPHPApache.DeploymentNodeSelector, map[string]string{"a": "b", "c": "d"}) {
		t.Fatalf("unexpected cfg.AddOnPHPApache.DeploymentNodeSelector %v", cfg.AddOnPHPApache.DeploymentNodeSelector)
	}
}

func TestEnvAddOnNLBHelloWorld(t *testing.T) {
	cfg := NewDefault()

	os.Setenv("K8S_TESTER_CONFIG_PATH", "test.yaml")
	defer os.Unsetenv("K8S_TESTER_CONFIG_PATH")
	os.Setenv("K8S_TESTER_ADD_ON_NLB_HELLO_WORLD_ENABLE", "true")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_NLB_HELLO_WORLD_ENABLE")
	os.Setenv("K8S_TESTER_ADD_ON_NLB_HELLO_WORLD_MINIMUM_NODES", "100")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_NLB_HELLO_WORLD_MINIMUM_NODES")
	os.Setenv("K8S_TESTER_ADD_ON_NLB_HELLO_WORLD_NAMESPACE", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_NLB_HELLO_WORLD_NAMESPACE")
	os.Setenv("K8S_TESTER_ADD_ON_NLB_HELLO_WORLD_DEPLOYMENT_NODE_SELECTOR", `{"a":"b","c":"d"}`)
	defer os.Unsetenv("K8S_TESTER_ADD_ON_NLB_HELLO_WORLD_DEPLOYMENT_NODE_SELECTOR")

	if err := cfg.UpdateFromEnvs(); err != nil {
		t.Fatal(err)
	}

	if cfg.ConfigPath != "test.yaml" {
		t.Fatalf("unexpected cfg.ConfigPath %v", cfg.ConfigPath)
	}
	if !cfg.AddOnNLBHelloWorld.Enable {
		t.Fatalf("unexpected cfg.AddOnNLBHelloWorld.Enable %v", cfg.AddOnNLBHelloWorld.Enable)
	}
	if cfg.AddOnNLBHelloWorld.MinimumNodes != 100 {
		t.Fatalf("unexpected cfg.AddOnNLBHelloWorld.MinimumNodes %v", cfg.AddOnNLBHelloWorld.MinimumNodes)
	}
	if cfg.AddOnNLBHelloWorld.Namespace != "hello" {
		t.Fatalf("unexpected cfg.AddOnNLBHelloWorld.Namespace %v", cfg.AddOnNLBHelloWorld.Namespace)
	}
	if !reflect.DeepEqual(cfg.AddOnNLBHelloWorld.DeploymentNodeSelector, map[string]string{"a": "b", "c": "d"}) {
		t.Fatalf("unexpected cfg.AddOnNLBHelloWorld.DeploymentNodeSelector %v", cfg.AddOnNLBHelloWorld.DeploymentNodeSelector)
	}
}

func TestEnvAddOnJobsPi(t *testing.T) {
	cfg := NewDefault()

	os.Setenv("K8S_TESTER_ADD_ON_JOBS_PI_ENABLE", "true")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_JOBS_PI_ENABLE")
	os.Setenv("K8S_TESTER_ADD_ON_JOBS_PI_MINIMUM_NODES", "100")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_JOBS_PI_MINIMUM_NODES")
	os.Setenv("K8S_TESTER_ADD_ON_JOBS_PI_NAMESPACE", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_JOBS_PI_NAMESPACE")
	os.Setenv("K8S_TESTER_ADD_ON_JOBS_PI_COMPLETES", `222`)
	defer os.Unsetenv("K8S_TESTER_ADD_ON_JOBS_PI_COMPLETES")
	os.Setenv("K8S_TESTER_ADD_ON_JOBS_PI_PARALLELS", `333`)
	defer os.Unsetenv("K8S_TESTER_ADD_ON_JOBS_PI_PARALLELS")

	if err := cfg.UpdateFromEnvs(); err != nil {
		t.Fatal(err)
	}

	if !cfg.AddOnJobsPi.Enable {
		t.Fatalf("unexpected cfg.AddOnJobsPi.Enable %v", cfg.AddOnJobsPi.Enable)
	}
	if cfg.AddOnJobsPi.MinimumNodes != 100 {
		t.Fatalf("unexpected cfg.AddOnJobsPi.MinimumNodes %v", cfg.AddOnJobsPi.MinimumNodes)
	}
	if cfg.AddOnJobsPi.Namespace != "hello" {
		t.Fatalf("unexpected cfg.AddOnJobsPi.Namespace %v", cfg.AddOnJobsPi.Namespace)
	}
	if cfg.AddOnJobsPi.Completes != 222 {
		t.Fatalf("unexpected cfg.AddOnJobsPi.Completes %v", cfg.AddOnJobsPi.Completes)
	}
	if cfg.AddOnJobsPi.Parallels != 333 {
		t.Fatalf("unexpected cfg.AddOnJobsPi.Parallels %v", cfg.AddOnJobsPi.Parallels)
	}
}

func TestEnvAddOnJobsEcho(t *testing.T) {
	cfg := NewDefault()

	os.Setenv("K8S_TESTER_ADD_ON_JOBS_ECHO_ENABLE", "true")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_JOBS_ECHO_ENABLE")
	os.Setenv("K8S_TESTER_ADD_ON_JOBS_ECHO_MINIMUM_NODES", "100")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_JOBS_ECHO_MINIMUM_NODES")
	os.Setenv("K8S_TESTER_ADD_ON_JOBS_ECHO_NAMESPACE", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_JOBS_ECHO_NAMESPACE")
	os.Setenv("K8S_TESTER_ADD_ON_JOBS_ECHO_REPOSITORY_BUSYBOX_PARTITION", "aws")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_JOBS_ECHO_REPOSITORY_BUSYBOX_PARTITION")
	os.Setenv("K8S_TESTER_ADD_ON_JOBS_ECHO_REPOSITORY_BUSYBOX_ACCOUNT_ID", "123")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_JOBS_ECHO_REPOSITORY_BUSYBOX_ACCOUNT_ID")
	os.Setenv("K8S_TESTER_ADD_ON_JOBS_ECHO_REPOSITORY_BUSYBOX_REGION", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_JOBS_ECHO_REPOSITORY_BUSYBOX_REGION")
	os.Setenv("K8S_TESTER_ADD_ON_JOBS_ECHO_REPOSITORY_BUSYBOX_NAME", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_JOBS_ECHO_REPOSITORY_BUSYBOX_NAME")
	os.Setenv("K8S_TESTER_ADD_ON_JOBS_ECHO_REPOSITORY_BUSYBOX_IMAGE_TAG", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_JOBS_ECHO_REPOSITORY_BUSYBOX_IMAGE_TAG")
	os.Setenv("K8S_TESTER_ADD_ON_JOBS_ECHO_COMPLETES", `222`)
	defer os.Unsetenv("K8S_TESTER_ADD_ON_JOBS_ECHO_COMPLETES")
	os.Setenv("K8S_TESTER_ADD_ON_JOBS_ECHO_PARALLELS", `333`)
	defer os.Unsetenv("K8S_TESTER_ADD_ON_JOBS_ECHO_PARALLELS")
	os.Setenv("K8S_TESTER_ADD_ON_JOBS_ECHO_ECHO_SIZE", `555`)
	defer os.Unsetenv("K8S_TESTER_ADD_ON_JOBS_ECHO_ECHO_SIZE")

	if err := cfg.UpdateFromEnvs(); err != nil {
		t.Fatal(err)
	}

	if !cfg.AddOnJobsEcho.Enable {
		t.Fatalf("unexpected cfg.AddOnJobsEcho.Enable %v", cfg.AddOnJobsEcho.Enable)
	}
	if cfg.AddOnJobsEcho.MinimumNodes != 100 {
		t.Fatalf("unexpected cfg.AddOnJobsEcho.MinimumNodes %v", cfg.AddOnJobsEcho.MinimumNodes)
	}
	if cfg.AddOnJobsEcho.Namespace != "hello" {
		t.Fatalf("unexpected cfg.AddOnJobsEcho.Namespace %v", cfg.AddOnJobsEcho.Namespace)
	}
	if cfg.AddOnJobsEcho.RepositoryBusyboxPartition != "aws" {
		t.Fatalf("unexpected cfg.AddOnJobsEcho.RepositoryBusyboxPartition %v", cfg.AddOnJobsEcho.RepositoryBusyboxPartition)
	}
	if cfg.AddOnJobsEcho.RepositoryBusyboxAccountID != "123" {
		t.Fatalf("unexpected cfg.AddOnJobsEcho.RepositoryBusyboxAccountID %v", cfg.AddOnJobsEcho.RepositoryBusyboxAccountID)
	}
	if cfg.AddOnJobsEcho.RepositoryBusyboxRegion != "hello" {
		t.Fatalf("unexpected cfg.AddOnJobsEcho.RepositoryBusyboxRegion %v", cfg.AddOnJobsEcho.RepositoryBusyboxRegion)
	}
	if cfg.AddOnJobsEcho.RepositoryBusyboxName != "hello" {
		t.Fatalf("unexpected cfg.AddOnJobsEcho.RepositoryBusyboxName %v", cfg.AddOnJobsEcho.RepositoryBusyboxName)
	}
	if cfg.AddOnJobsEcho.RepositoryBusyboxImageTag != "hello" {
		t.Fatalf("unexpected cfg.AddOnJobsEcho.RepositoryBusyboxImageTag %v", cfg.AddOnJobsEcho.RepositoryBusyboxImageTag)
	}
	if cfg.AddOnJobsEcho.Completes != 222 {
		t.Fatalf("unexpected cfg.AddOnJobsEcho.Completes %v", cfg.AddOnJobsEcho.Completes)
	}
	if cfg.AddOnJobsEcho.Parallels != 333 {
		t.Fatalf("unexpected cfg.AddOnJobsEcho.Parallels %v", cfg.AddOnJobsEcho.Parallels)
	}
	if cfg.AddOnJobsEcho.EchoSize != 555 {
		t.Fatalf("unexpected cfg.AddOnJobsEcho.EchoSize %v", cfg.AddOnJobsEcho.EchoSize)
	}
}

func TestEnvAddOnCronJobsEcho(t *testing.T) {
	cfg := NewDefault()

	os.Setenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_ENABLE", "true")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_ENABLE")
	os.Setenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_MINIMUM_NODES", "100")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_MINIMUM_NODES")
	os.Setenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_NAMESPACE", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_NAMESPACE")
	os.Setenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_REPOSITORY_BUSYBOX_PARTITION", "aws")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_REPOSITORY_BUSYBOX_PARTITION")
	os.Setenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_REPOSITORY_BUSYBOX_ACCOUNT_ID", "123")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_REPOSITORY_BUSYBOX_ACCOUNT_ID")
	os.Setenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_REPOSITORY_BUSYBOX_REGION", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_REPOSITORY_BUSYBOX_REGION")
	os.Setenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_REPOSITORY_BUSYBOX_NAME", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_REPOSITORY_BUSYBOX_NAME")
	os.Setenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_REPOSITORY_BUSYBOX_IMAGE_TAG", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_REPOSITORY_BUSYBOX_IMAGE_TAG")
	os.Setenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_COMPLETES", `222`)
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_COMPLETES")
	os.Setenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_PARALLELS", `333`)
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_PARALLELS")
	os.Setenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_ECHO_SIZE", `555`)
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_ECHO_SIZE")
	os.Setenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_SCHEDULE", `*/10 */10 * * *`)
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_SCHEDULE")
	os.Setenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_SUCCESSFUL_JOBS_HISTORY_LIMIT", "55555")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_SUCCESSFUL_JOBS_HISTORY_LIMIT")
	os.Setenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_FAILED_JOBS_HISTORY_LIMIT", "77777")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CRON_JOBS_ECHO_FAILED_JOBS_HISTORY_LIMIT")

	if err := cfg.UpdateFromEnvs(); err != nil {
		t.Fatal(err)
	}

	if !cfg.AddOnCronJobsEcho.Enable {
		t.Fatalf("unexpected cfg.AddOnCronJobsEcho.Enable %v", cfg.AddOnCronJobsEcho.Enable)
	}
	if cfg.AddOnCronJobsEcho.MinimumNodes != 100 {
		t.Fatalf("unexpected cfg.AddOnCronJobsEcho.MinimumNodes %v", cfg.AddOnCronJobsEcho.MinimumNodes)
	}
	if cfg.AddOnCronJobsEcho.Namespace != "hello" {
		t.Fatalf("unexpected cfg.AddOnCronJobsEcho.Namespace %v", cfg.AddOnCronJobsEcho.Namespace)
	}
	if cfg.AddOnCronJobsEcho.RepositoryBusyboxPartition != "aws" {
		t.Fatalf("unexpected cfg.AddOnCronJobsEcho.RepositoryBusyboxPartition %v", cfg.AddOnCronJobsEcho.RepositoryBusyboxPartition)
	}
	if cfg.AddOnCronJobsEcho.RepositoryBusyboxAccountID != "123" {
		t.Fatalf("unexpected cfg.AddOnCronJobsEcho.RepositoryBusyboxAccountID %v", cfg.AddOnCronJobsEcho.RepositoryBusyboxAccountID)
	}
	if cfg.AddOnCronJobsEcho.RepositoryBusyboxRegion != "hello" {
		t.Fatalf("unexpected cfg.AddOnCronJobsEcho.RepositoryBusyboxRegion %v", cfg.AddOnCronJobsEcho.RepositoryBusyboxRegion)
	}
	if cfg.AddOnCronJobsEcho.RepositoryBusyboxName != "hello" {
		t.Fatalf("unexpected cfg.AddOnCronJobsEcho.RepositoryBusyboxName %v", cfg.AddOnCronJobsEcho.RepositoryBusyboxName)
	}
	if cfg.AddOnCronJobsEcho.RepositoryBusyboxImageTag != "hello" {
		t.Fatalf("unexpected cfg.AddOnCronJobsEcho.RepositoryBusyboxImageTag %v", cfg.AddOnCronJobsEcho.RepositoryBusyboxImageTag)
	}
	if cfg.AddOnCronJobsEcho.Completes != 222 {
		t.Fatalf("unexpected cfg.AddOnCronJobsEcho.Completes %v", cfg.AddOnCronJobsEcho.Completes)
	}
	if cfg.AddOnCronJobsEcho.Parallels != 333 {
		t.Fatalf("unexpected cfg.AddOnCronJobsEcho.Parallels %v", cfg.AddOnCronJobsEcho.Parallels)
	}
	if cfg.AddOnCronJobsEcho.EchoSize != 555 {
		t.Fatalf("unexpected cfg.AddOnCronJobsEcho.EchoSize %v", cfg.AddOnCronJobsEcho.EchoSize)
	}
	if cfg.AddOnCronJobsEcho.Schedule != "*/10 */10 * * *" {
		t.Fatalf("unexpected cfg.AddOnCronJobsEcho.Schedule %v", cfg.AddOnCronJobsEcho.Schedule)
	}
	if cfg.AddOnCronJobsEcho.SuccessfulJobsHistoryLimit != 55555 {
		t.Fatalf("unexpected cfg.AddOnCronJobsEcho.SuccessfulJobsHistoryLimit %v", cfg.AddOnCronJobsEcho.SuccessfulJobsHistoryLimit)
	}
	if cfg.AddOnCronJobsEcho.FailedJobsHistoryLimit != 77777 {
		t.Fatalf("unexpected cfg.AddOnCronJobsEcho.FailedJobsHistoryLimit %v", cfg.AddOnCronJobsEcho.FailedJobsHistoryLimit)
	}
}

func TestEnvAddOnCSRs(t *testing.T) {
	cfg := NewDefault()

	os.Setenv("K8S_TESTER_CLIENTS", "100")
	defer os.Unsetenv("K8S_TESTER_CLIENTS")
	os.Setenv("K8S_TESTER_ADD_ON_CSRS_ENABLE", "true")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CSRS_ENABLE")
	os.Setenv("K8S_TESTER_ADD_ON_CSRS_MINIMUM_NODES", "100")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CSRS_MINIMUM_NODES")
	os.Setenv("K8S_TESTER_ADD_ON_CSRS_OBJECTS", `222`)
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CSRS_OBJECTS")
	os.Setenv("K8S_TESTER_ADD_ON_CSRS_INITIAL_REQUEST_CONDITION_TYPE", "Approved")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CSRS_INITIAL_REQUEST_CONDITION_TYPE")

	if err := cfg.UpdateFromEnvs(); err != nil {
		t.Fatal(err)
	}

	if cfg.Clients != 100 {
		t.Fatalf("unexpected cfg.Clients %v", cfg.Clients)
	}
	if !cfg.AddOnCSRs.Enable {
		t.Fatalf("unexpected cfg.AddOnCSRs.Enable %v", cfg.AddOnCSRs.Enable)
	}
	if cfg.AddOnCSRs.MinimumNodes != 100 {
		t.Fatalf("unexpected cfg.AddOnCSRs.MinimumNodes %v", cfg.AddOnCSRs.MinimumNodes)
	}
	if cfg.AddOnCSRs.Objects != 222 {
		t.Fatalf("unexpected cfg.AddOnCSRs.Objects %v", cfg.AddOnCSRs.Objects)
	}
	if cfg.AddOnCSRs.InitialRequestConditionType != "Approved" {
		t.Fatalf("unexpected cfg.AddOnCSRs.InitialRequestConditionType %v", cfg.AddOnCSRs.InitialRequestConditionType)
	}
}

func TestEnvAddOnConfigmaps(t *testing.T) {
	cfg := NewDefault()

	os.Setenv("K8S_TESTER_CLIENTS", "100")
	defer os.Unsetenv("K8S_TESTER_CLIENTS")
	os.Setenv("K8S_TESTER_ADD_ON_CONFIGMAPS_ENABLE", "true")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFIGMAPS_ENABLE")
	os.Setenv("K8S_TESTER_ADD_ON_CONFIGMAPS_MINIMUM_NODES", "100")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFIGMAPS_MINIMUM_NODES")
	os.Setenv("K8S_TESTER_ADD_ON_CONFIGMAPS_NAMESPACE", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFIGMAPS_NAMESPACE")
	os.Setenv("K8S_TESTER_ADD_ON_CONFIGMAPS_OBJECTS", `222`)
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFIGMAPS_OBJECTS")
	os.Setenv("K8S_TESTER_ADD_ON_CONFIGMAPS_OBJECT_SIZE", `333`)
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CONFIGMAPS_OBJECT_SIZE")

	if err := cfg.UpdateFromEnvs(); err != nil {
		t.Fatal(err)
	}

	if cfg.Clients != 100 {
		t.Fatalf("unexpected cfg.Clients %v", cfg.Clients)
	}
	if !cfg.AddOnConfigmaps.Enable {
		t.Fatalf("unexpected cfg.AddOnConfigmaps.Enable %v", cfg.AddOnConfigmaps.Enable)
	}
	if cfg.AddOnConfigmaps.MinimumNodes != 100 {
		t.Fatalf("unexpected cfg.AddOnConfigmaps.MinimumNodes %v", cfg.AddOnConfigmaps.MinimumNodes)
	}
	if cfg.AddOnConfigmaps.Namespace != "hello" {
		t.Fatalf("unexpected cfg.AddOnConfigmaps.Namespace %v", cfg.AddOnConfigmaps.Namespace)
	}
	if cfg.AddOnConfigmaps.Objects != 222 {
		t.Fatalf("unexpected cfg.AddOnConfigmaps.Objects %v", cfg.AddOnConfigmaps.Objects)
	}
	if cfg.AddOnConfigmaps.ObjectSize != 333 {
		t.Fatalf("unexpected cfg.AddOnConfigmaps.ObjectSize %v", cfg.AddOnConfigmaps.ObjectSize)
	}
}

func TestEnvAddOnSecrets(t *testing.T) {
	cfg := NewDefault()

	os.Setenv("K8S_TESTER_CLIENTS", "100")
	defer os.Unsetenv("K8S_TESTER_CLIENTS")
	os.Setenv("K8S_TESTER_ADD_ON_SECRETS_ENABLE", "true")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_SECRETS_ENABLE")
	os.Setenv("K8S_TESTER_ADD_ON_SECRETS_MINIMUM_NODES", "100")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_SECRETS_MINIMUM_NODES")
	os.Setenv("K8S_TESTER_ADD_ON_SECRETS_NAMESPACE", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_SECRETS_NAMESPACE")
	os.Setenv("K8S_TESTER_ADD_ON_SECRETS_OBJECTS", `222`)
	defer os.Unsetenv("K8S_TESTER_ADD_ON_SECRETS_OBJECTS")
	os.Setenv("K8S_TESTER_ADD_ON_SECRETS_OBJECT_SIZE", `333`)
	defer os.Unsetenv("K8S_TESTER_ADD_ON_SECRETS_OBJECT_SIZE")

	if err := cfg.UpdateFromEnvs(); err != nil {
		t.Fatal(err)
	}

	if cfg.Clients != 100 {
		t.Fatalf("unexpected cfg.Clients %v", cfg.Clients)
	}
	if !cfg.AddOnSecrets.Enable {
		t.Fatalf("unexpected cfg.AddOnSecrets.Enable %v", cfg.AddOnSecrets.Enable)
	}
	if cfg.AddOnSecrets.MinimumNodes != 100 {
		t.Fatalf("unexpected cfg.AddOnSecrets.MinimumNodes %v", cfg.AddOnSecrets.MinimumNodes)
	}
	if cfg.AddOnSecrets.Namespace != "hello" {
		t.Fatalf("unexpected cfg.AddOnSecrets.Namespace %v", cfg.AddOnSecrets.Namespace)
	}
	if cfg.AddOnSecrets.Objects != 222 {
		t.Fatalf("unexpected cfg.AddOnSecrets.Objects %v", cfg.AddOnSecrets.Objects)
	}
	if cfg.AddOnSecrets.ObjectSize != 333 {
		t.Fatalf("unexpected cfg.AddOnSecrets.ObjectSize %v", cfg.AddOnSecrets.ObjectSize)
	}
}

func TestEnvAddOnClusterloader(t *testing.T) {
	cfg := NewDefault()

	os.Setenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_ENABLE", "true")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_ENABLE")
	os.Setenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_CLUSTERLOADER_PATH", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_CLUSTERLOADER_PATH")
	os.Setenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_CLUSTERLOADER_DOWNLOAD_URL", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_CLUSTERLOADER_DOWNLOAD_URL")
	os.Setenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_RUNS", "10")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_RUNS")
	os.Setenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_RUN_TIMEOUT", "10h")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_RUN_TIMEOUT")
	os.Setenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_TEST_CONFIG_PATH", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_TEST_CONFIG_PATH")
	os.Setenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_REPORT_DIR", "hello")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_REPORT_DIR")
	os.Setenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_RUN_FROM_CLUSTER", "true")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_RUN_FROM_CLUSTER")
	os.Setenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_NODES", "100")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_NODES")
	os.Setenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_TEST_OVERRIDE_NODES_PER_NAMESPACE", "100")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_TEST_OVERRIDE_NODES_PER_NAMESPACE")
	os.Setenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_TEST_OVERRIDE_PODS_PER_NODE", "100")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_TEST_OVERRIDE_PODS_PER_NODE")
	os.Setenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_TEST_OVERRIDE_MEDIUM_STATEFUL_SETS_PER_NAMESPACE", "5000")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_TEST_OVERRIDE_MEDIUM_STATEFUL_SETS_PER_NAMESPACE")
	os.Setenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_TEST_OVERRIDE_CL2_SCHEDULER_THROUGHPUT_THRESHOLD", "1000")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_TEST_OVERRIDE_CL2_SCHEDULER_THROUGHPUT_THRESHOLD")
	os.Setenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_TEST_OVERRIDE_ENABLE_SYSTEM_POD_METRICS", "true")
	defer os.Unsetenv("K8S_TESTER_ADD_ON_CLUSTERLOADER_TEST_OVERRIDE_ENABLE_SYSTEM_POD_METRICS")

	if err := cfg.UpdateFromEnvs(); err != nil {
		t.Fatal(err)
	}

	if !cfg.AddOnClusterloader.Enable {
		t.Fatalf("unexpected cfg.AddOnClusterloader.Enable %v", cfg.AddOnClusterloader.Enable)
	}
	if cfg.AddOnClusterloader.ClusterloaderPath != "hello" {
		t.Fatalf("unexpected cfg.AddOnClusterloader.ClusterloaderPath %v", cfg.AddOnClusterloader.ClusterloaderPath)
	}
	if cfg.AddOnClusterloader.ClusterloaderDownloadURL != "hello" {
		t.Fatalf("unexpected cfg.AddOnClusterloader.ClusterloaderDownloadURL %v", cfg.AddOnClusterloader.ClusterloaderDownloadURL)
	}
	if cfg.AddOnClusterloader.Runs != 10 {
		t.Fatalf("unexpected cfg.AddOnClusterloader.Runs %v", cfg.AddOnClusterloader.Runs)
	}
	if cfg.AddOnClusterloader.RunTimeout != 10*time.Hour {
		t.Fatalf("unexpected cfg.AddOnClusterloader.RunTimeout %v", cfg.AddOnClusterloader.RunTimeout)
	}
	if cfg.AddOnClusterloader.TestConfigPath != "hello" {
		t.Fatalf("unexpected cfg.AddOnClusterloader.TestConfigPath %v", cfg.AddOnClusterloader.TestConfigPath)
	}
	if cfg.AddOnClusterloader.ReportDir != "hello" {
		t.Fatalf("unexpected cfg.AddOnClusterloader.ReportDir %v", cfg.AddOnClusterloader.ReportDir)
	}
	if !cfg.AddOnClusterloader.RunFromCluster {
		t.Fatalf("unexpected cfg.AddOnClusterloader.RunFromCluster %v", cfg.AddOnClusterloader.RunFromCluster)
	}
	if cfg.AddOnClusterloader.Nodes != 100 {
		t.Fatalf("unexpected cfg.AddOnClusterloader.Nodes %v", cfg.AddOnClusterloader.Nodes)
	}
	if cfg.AddOnClusterloader.TestOverride.NodesPerNamespace != 100 {
		t.Fatalf("unexpected cfg.AddOnClusterloader.TestOverride.NodesPerNamespace %v", cfg.AddOnClusterloader.TestOverride.NodesPerNamespace)
	}
	if cfg.AddOnClusterloader.TestOverride.PodsPerNode != 100 {
		t.Fatalf("unexpected cfg.AddOnClusterloader.TestOverride.PodsPerNode %v", cfg.AddOnClusterloader.TestOverride.PodsPerNode)
	}
	if cfg.AddOnClusterloader.TestOverride.MediumStatefulSetsPerNamespace != 5000 {
		t.Fatalf("unexpected cfg.AddOnClusterloader.TestOverride.MediumStatefulSetsPerNamespace %v", cfg.AddOnClusterloader.TestOverride.MediumStatefulSetsPerNamespace)
	}
	if cfg.AddOnClusterloader.TestOverride.CL2SchedulerThroughputThreshold != 1000 {
		t.Fatalf("unexpected cfg.AddOnClusterloader.TestOverride.CL2SchedulerThroughputThreshold %v", cfg.AddOnClusterloader.TestOverride.CL2SchedulerThroughputThreshold)
	}
	if !cfg.AddOnClusterloader.TestOverride.EnableSystemPodMetrics {
		t.Fatalf("unexpected cfg.AddOnClusterloader.TestOverride.EnableSystemPodMetrics %v", cfg.AddOnClusterloader.TestOverride.EnableSystemPodMetrics)
	}
}
