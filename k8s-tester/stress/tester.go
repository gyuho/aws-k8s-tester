// Package stress implements stress tester.
// The rule is "do not parallelize locally, instead parallelize by distributing workers".
// Replace https://github.com/aws/aws-k8s-tester/tree/v1.5.9/eks/stress.
package stress

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-k8s-tester/client"
	k8s_tester "github.com/aws/aws-k8s-tester/k8s-tester/tester"
	"github.com/aws/aws-k8s-tester/utils/latency"
	"github.com/aws/aws-k8s-tester/utils/rand"
	utils_time "github.com/aws/aws-k8s-tester/utils/time"
	"github.com/dustin/go-humanize"
	"github.com/manifoldco/promptui"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	writeRequestsSuccessTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "stress",
			Subsystem: "client",
			Name:      "write_requests_success_total",
			Help:      "Total number of successful write requests.",
		})
	writeRequestsFailureTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "stress",
			Subsystem: "client",
			Name:      "write_requests_failure_total",
			Help:      "Total number of successful write requests.",
		})
	writeRequestLatencyMs = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "stress",
			Subsystem: "client",
			Name:      "write_request_latency_milliseconds",
			Help:      "Bucketed histogram of client-side write request and response latency.",

			// lowest bucket start of upper bound 0.5 ms with factor 2
			// highest bucket start of 0.5 ms * 2^13 == 4.096 sec
			Buckets: prometheus.ExponentialBuckets(0.5, 2, 14),
		})

	updateRequestsSuccessTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "stress",
			Subsystem: "client",
			Name:      "update_requests_success_total",
			Help:      "Total number of successful update requests.",
		})
	updateRequestsFailureTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "stress",
			Subsystem: "client",
			Name:      "update_requests_failure_total",
			Help:      "Total number of successful update requests.",
		})
	updateRequestLatencyMs = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "stress",
			Subsystem: "client",
			Name:      "update_request_latency_milliseconds",
			Help:      "Bucketed histogram of client-side update request and response latency.",

			// lowest bucket start of upper bound 0.5 ms with factor 2
			// highest bucket start of 0.5 ms * 2^13 == 4.096 sec
			Buckets: prometheus.ExponentialBuckets(0.5, 2, 14),
		})

	readRequestsSuccessTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "stress",
			Subsystem: "client",
			Name:      "read_requests_success_total",
			Help:      "Total number of successful read requests.",
		})
	readRequestsFailureTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "stress",
			Subsystem: "client",
			Name:      "read_requests_failure_total",
			Help:      "Total number of successful read requests.",
		})
	readRequestLatencyMs = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "stress",
			Subsystem: "client",
			Name:      "read_request_latency_milliseconds",
			Help:      "Bucketed histogram of client-side read request and response latency.",

			// lowest bucket start of upper bound 0.5 ms with factor 2
			// highest bucket start of 0.5 ms * 2^13 == 4.096 sec
			Buckets: prometheus.ExponentialBuckets(0.5, 2, 14),
		})
)

func init() {
	prometheus.MustRegister(writeRequestsSuccessTotal)
	prometheus.MustRegister(writeRequestsFailureTotal)
	prometheus.MustRegister(writeRequestLatencyMs)

	prometheus.MustRegister(updateRequestsSuccessTotal)
	prometheus.MustRegister(updateRequestsFailureTotal)
	prometheus.MustRegister(updateRequestLatencyMs)

	prometheus.MustRegister(readRequestsSuccessTotal)
	prometheus.MustRegister(readRequestsFailureTotal)
	prometheus.MustRegister(readRequestLatencyMs)
}

type Config struct {
	Enable bool `json:"enable"`
	Prompt bool `json:"-"`

	Stopc     chan struct{} `json:"-"`
	Logger    *zap.Logger   `json:"-"`
	LogWriter io.Writer     `json:"-"`
	Client    client.Client `json:"-"`

	// MinimumNodes is the minimum number of Kubernetes nodes required for installing this addon.
	MinimumNodes int `json:"minimum_nodes"`
	// Namespace to create test resources.
	Namespace string `json:"namespace"`

	// RunTimeout is the duration of stress runs.
	// After timeout, it stops all stress requests.
	RunTimeout       time.Duration `json:"run_timeout"`
	RunTimeoutString string        `json:"run_timeout_string" read-only:"true"`

	// CreateObjects is the desired number of objects to create.
	CreateObjects int `json:"create_objects"`
	// CreateObjectSize is the size in bytes per object.
	CreateObjectSize int `json:"create_object_size"`

	LatencySummaryCreates latency.Summary `json:"latency_summary_creates" read-only:"true"`
	LatencySummaryUpdates latency.Summary `json:"latency_summary_updates" read-only:"true"`
	LatencySummaryReads   latency.Summary `json:"latency_summary_reads" read-only:"true"`
}

func (cfg *Config) ValidateAndSetDefaults() error {
	if cfg.MinimumNodes == 0 {
		cfg.MinimumNodes = DefaultMinimumNodes
	}
	if cfg.Namespace == "" {
		return errors.New("empty Namespace")
	}

	return nil
}

const (
	DefaultMinimumNodes int = 1
	DefaultRunTimeout       = time.Minute

	DefaultCreateObjects    int = 10
	DefaultCreateObjectSize int = 10 * 1024 // 10 KB

	// writes total 300 MB data to etcd
	// CreateObjects: 1000,
	// CreateObjectSize: 300000, // 0.3 MB
)

func NewDefault() *Config {
	return &Config{
		Enable:           false,
		Prompt:           false,
		MinimumNodes:     DefaultMinimumNodes,
		Namespace:        pkgName + "-" + rand.String(10) + "-" + utils_time.GetTS(10),
		RunTimeout:       DefaultRunTimeout,
		RunTimeoutString: DefaultRunTimeout.String(),
		CreateObjects:    DefaultCreateObjects,
		CreateObjectSize: DefaultCreateObjectSize,
	}
}

func New(cfg *Config) k8s_tester.Tester {
	return &tester{
		cfg:            cfg,
		donec:          make(chan struct{}),
		donecCloseOnce: new(sync.Once),
	}
}

type tester struct {
	cfg            *Config
	donec          chan struct{}
	donecCloseOnce *sync.Once
}

var pkgName = path.Base(reflect.TypeOf(tester{}).PkgPath())

func Env() string {
	return "ADD_ON_" + strings.ToUpper(strings.Replace(pkgName, "-", "_", -1))
}

func EnvRepository() string {
	return Env() + "_REPOSITORY"
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

	if err := client.CreateNamespace(ts.cfg.Logger, ts.cfg.Client.KubernetesClient(), ts.cfg.Namespace); err != nil {
		return err
	}

	latenciesWritesCh := make(chan latency.Durations)
	latenciesUpdatesCh := make(chan latency.Durations)
	latenciesReadsCh := make(chan latency.Durations)
	go func() {
		latenciesWritesCh <- ts.startWrites()
	}()
	go func() {
		latenciesUpdatesCh <- ts.startUpdates()
	}()
	go func() {
		latenciesReadsCh <- ts.startReads()
	}()
	select {
	case <-ts.cfg.Stopc:
		ts.cfg.Logger.Warn("all stopped")
		return nil
	case <-time.After(ts.cfg.RunTimeout):
		ts.cfg.Logger.Info("run timeout, closing done channel")
		ts.donecCloseOnce.Do(func() {
			close(ts.donec)
		})
	}
	var latenciesWrites latency.Durations
	select {
	case <-ts.cfg.Stopc:
		ts.cfg.Logger.Warn("stopped while waiting for writes results")
		return nil
	case latenciesWrites = <-latenciesWritesCh:
	}
	var latenciesUpdates latency.Durations
	select {
	case <-ts.cfg.Stopc:
		ts.cfg.Logger.Warn("stopped while waiting for writes results")
		return nil
	case latenciesUpdates = <-latenciesUpdatesCh:
	}
	var latenciesReads latency.Durations
	select {
	case <-ts.cfg.Stopc:
		ts.cfg.Logger.Warn("stopped while waiting for writes results")
		return nil
	case latenciesReads = <-latenciesReadsCh:
	}

	ts.cfg.Logger.Info("sorting write latency results", zap.Int("total-data-points", latenciesWrites.Len()))
	now := time.Now()
	sort.Sort(latenciesWrites)
	ts.cfg.Logger.Info("sorted write latency results", zap.Int("total-data-points", latenciesWrites.Len()), zap.String("took", time.Since(now).String()))

	ts.cfg.Logger.Info("sorting update latency results", zap.Int("total-data-points", latenciesUpdates.Len()))
	now = time.Now()
	sort.Sort(latenciesUpdates)
	ts.cfg.Logger.Info("sorted update latency results", zap.Int("total-data-points", latenciesUpdates.Len()), zap.String("took", time.Since(now).String()))

	ts.cfg.Logger.Info("sorting read latency results", zap.Int("total-data-points", latenciesReads.Len()))
	now = time.Now()
	sort.Sort(latenciesReads)
	ts.cfg.Logger.Info("sorted read latency results", zap.Int("total-data-points", latenciesReads.Len()), zap.String("took", time.Since(now).String()))

	testID := time.Now().UTC().Format(time.RFC3339Nano)

	ts.cfg.LatencySummaryCreates.TestID = testID
	ts.cfg.LatencySummaryCreates.P50 = latenciesWrites.PickP50()
	ts.cfg.LatencySummaryCreates.P90 = latenciesWrites.PickP90()
	ts.cfg.LatencySummaryCreates.P99 = latenciesWrites.PickP99()
	ts.cfg.LatencySummaryCreates.P999 = latenciesWrites.PickP999()
	ts.cfg.LatencySummaryCreates.P9999 = latenciesWrites.PickP9999()

	ts.cfg.LatencySummaryUpdates.TestID = testID
	ts.cfg.LatencySummaryUpdates.P50 = latenciesUpdates.PickP50()
	ts.cfg.LatencySummaryUpdates.P90 = latenciesUpdates.PickP90()
	ts.cfg.LatencySummaryUpdates.P99 = latenciesUpdates.PickP99()
	ts.cfg.LatencySummaryUpdates.P999 = latenciesUpdates.PickP999()
	ts.cfg.LatencySummaryUpdates.P9999 = latenciesUpdates.PickP9999()

	ts.cfg.LatencySummaryReads.TestID = testID
	ts.cfg.LatencySummaryReads.P50 = latenciesReads.PickP50()
	ts.cfg.LatencySummaryReads.P90 = latenciesReads.PickP90()
	ts.cfg.LatencySummaryReads.P99 = latenciesReads.PickP99()
	ts.cfg.LatencySummaryReads.P999 = latenciesReads.PickP999()
	ts.cfg.LatencySummaryReads.P9999 = latenciesReads.PickP9999()

	// https://pkg.go.dev/github.com/prometheus/client_golang/prometheus?tab=doc#Gatherer
	mfs, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		ts.cfg.Logger.Warn("failed to gather prometheus metrics", zap.Error(err))
		return err
	}
	for _, mf := range mfs {
		if mf == nil {
			continue
		}

		switch *mf.Name {
		case "stress_client_write_requests_success_total":
			gg := mf.Metric[0].GetGauge()
			ts.cfg.LatencySummaryCreates.SuccessTotal = gg.GetValue()
		case "stress_client_write_requests_failure_total":
			gg := mf.Metric[0].GetGauge()
			ts.cfg.LatencySummaryCreates.FailureTotal = gg.GetValue()
		case "stress_client_write_request_latency_milliseconds":
			ts.cfg.LatencySummaryCreates.Histogram, err = latency.ParseHistogram("milliseconds", mf.Metric[0].GetHistogram())
			if err != nil {
				return err
			}

		case "stress_client_update_requests_success_total":
			gg := mf.Metric[0].GetGauge()
			ts.cfg.LatencySummaryUpdates.SuccessTotal = gg.GetValue()
		case "stress_client_update_requests_failure_total":
			gg := mf.Metric[0].GetGauge()
			ts.cfg.LatencySummaryUpdates.FailureTotal = gg.GetValue()
		case "stress_client_update_request_latency_milliseconds":
			ts.cfg.LatencySummaryUpdates.Histogram, err = latency.ParseHistogram("milliseconds", mf.Metric[0].GetHistogram())
			if err != nil {
				return err
			}

		case "stresser_client_read_requests_success_total":
			gg := mf.Metric[0].GetGauge()
			ts.cfg.LatencySummaryReads.SuccessTotal = gg.GetValue()
		case "stresser_client_read_requests_failure_total":
			gg := mf.Metric[0].GetGauge()
			ts.cfg.LatencySummaryReads.FailureTotal = gg.GetValue()
		case "stresser_client_read_request_latency_milliseconds":
			ts.cfg.LatencySummaryReads.Histogram, err = latency.ParseHistogram("milliseconds", mf.Metric[0].GetHistogram())
			if err != nil {
				return err
			}
		}
	}
	fmt.Fprintf(ts.cfg.LogWriter, "\n\nLatencySummary:\n%s\n", ts.cfg.LatencySummaryCreates.Table())

	return nil
}

func (ts *tester) Delete() error {
	if ok := ts.runPrompt("delete"); !ok {
		return errors.New("cancelled")
	}

	ts.donecCloseOnce.Do(func() {
		close(ts.donec)
	})

	var errs []string

	if err := client.DeleteNamespaceAndWait(
		ts.cfg.Logger,
		ts.cfg.Client.KubernetesClient(),
		ts.cfg.Namespace,
		client.DefaultNamespaceDeletionInterval,
		client.DefaultNamespaceDeletionTimeout,
		client.WithForceDelete(true),
	); err != nil {
		errs = append(errs, fmt.Sprintf("failed to delete namespace (%v)", err))
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}

func (ts *tester) runPrompt(action string) (ok bool) {
	if ts.cfg.Prompt {
		msg := fmt.Sprintf("Ready to %q resources for the namespace %q, should we continue?", action, ts.cfg.Namespace)
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

func (ts *tester) startWrites() (latencies latency.Durations) {
	if ts.cfg.CreateObjectSize == 0 {
		ts.cfg.Logger.Info("skipping writes, object size 0")
		return latencies
	}

	ts.cfg.Logger.Info("writing", zap.Int("objects", ts.cfg.CreateObjects), zap.String("object-size", humanize.Bytes(uint64(ts.cfg.CreateObjectSize))))
	latencies = make(latency.Durations, 0, 20000)

	val := rand.String(ts.cfg.CreateObjectSize)
	for i := 0; i < ts.cfg.CreateObjects; i++ {
		select {
		case <-ts.cfg.Stopc:
			ts.cfg.Logger.Warn("writes stopped")
			return
		case <-ts.donec:
			ts.cfg.Logger.Info("writes done")
			return
		default:
		}

		key := fmt.Sprintf("configmap%d%s", i, rand.String(7))

		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.Client.Config().ClientTimeout)
		_, err := ts.cfg.Client.KubernetesClient().
			CoreV1().
			ConfigMaps(ts.cfg.Namespace).
			Create(ctx, &core_v1.ConfigMap{
				TypeMeta: meta_v1.TypeMeta{
					APIVersion: "v1",
					Kind:       "ConfigMap",
				},
				ObjectMeta: meta_v1.ObjectMeta{
					Name:      key,
					Namespace: ts.cfg.Namespace,
					Labels: map[string]string{
						"name": key,
					},
				},
				Data: map[string]string{key: val},
			}, meta_v1.CreateOptions{})
		cancel()
		took := time.Since(start)
		tookMS := float64(took / time.Millisecond)
		writeRequestLatencyMs.Observe(tookMS)
		latencies = append(latencies, took)
		if err != nil {
			writeRequestsFailureTotal.Inc()
			ts.cfg.Logger.Warn("write configmap failed", zap.String("namespace", ts.cfg.Namespace), zap.Error(err))
		} else {
			writeRequestsSuccessTotal.Inc()
			if i%20 == 0 {
				ts.cfg.Logger.Info("wrote configmap", zap.Int("iteration", i), zap.String("namespace", ts.cfg.Namespace))
			}
		}
	}
	return latencies
}

func (ts *tester) startUpdates() (latencies latency.Durations) {
	if ts.cfg.CreateObjectSize == 0 {
		ts.cfg.Logger.Info("skipping updates, object size 0")
		return latencies
	}

	ts.cfg.Logger.Info("updating", zap.Int("objects", ts.cfg.CreateObjects), zap.String("object-size", humanize.Bytes(uint64(ts.cfg.CreateObjectSize))))
	latencies = make(latency.Durations, 0, 20000)

	val := rand.String(ts.cfg.CreateObjectSize)
	for i := 0; i < ts.cfg.CreateObjects; i++ {
		select {
		case <-ts.cfg.Stopc:
			ts.cfg.Logger.Warn("updates stopped")
			return
		case <-ts.donec:
			ts.cfg.Logger.Info("updates done")
			return
		default:
		}

		key := fmt.Sprintf("configmap%d%s", i, rand.String(7))

		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.Client.Config().ClientTimeout)
		_, err := ts.cfg.Client.KubernetesClient().
			CoreV1().
			ConfigMaps(ts.cfg.Namespace).
			Create(ctx, &core_v1.ConfigMap{
				TypeMeta: meta_v1.TypeMeta{
					APIVersion: "v1",
					Kind:       "ConfigMap",
				},
				ObjectMeta: meta_v1.ObjectMeta{
					Name:      key,
					Namespace: ts.cfg.Namespace,
					Labels: map[string]string{
						"name": key,
					},
				},
				Data: map[string]string{key: val},
			}, meta_v1.CreateOptions{})
		cancel()
		took := time.Since(start)
		tookMS := float64(took / time.Millisecond)
		updateRequestLatencyMs.Observe(tookMS)
		latencies = append(latencies, took)
		if err != nil {
			updateRequestsFailureTotal.Inc()
			ts.cfg.Logger.Warn("update configmap failed", zap.String("namespace", ts.cfg.Namespace), zap.Error(err))
		} else {
			updateRequestsSuccessTotal.Inc()
			if i%20 == 0 {
				ts.cfg.Logger.Info("updated configmap", zap.Int("iteration", i), zap.String("namespace", ts.cfg.Namespace))
			}
		}
	}
	return latencies
}

func (ts *tester) startReads() (latencies latency.Durations) {
	if ts.cfg.CreateObjectSize == 0 {
		ts.cfg.Logger.Info("skipping reads, object size 0")
		return latencies
	}

	ts.cfg.Logger.Info("reading", zap.Int("objects", ts.cfg.CreateObjects))
	latencies = make(latency.Durations, 0, 20000)

	val := rand.String(ts.cfg.CreateObjectSize)
	for i := 0; i < ts.cfg.CreateObjects; i++ {
		select {
		case <-ts.cfg.Stopc:
			ts.cfg.Logger.Warn("reads stopped")
			return
		case <-ts.donec:
			ts.cfg.Logger.Info("reads done")
			return
		default:
		}

		key := fmt.Sprintf("configmap%d%s", i, rand.String(7))

		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.Client.Config().ClientTimeout)
		_, err := ts.cfg.Client.KubernetesClient().
			CoreV1().
			ConfigMaps(ts.cfg.Namespace).
			Create(ctx, &core_v1.ConfigMap{
				TypeMeta: meta_v1.TypeMeta{
					APIVersion: "v1",
					Kind:       "ConfigMap",
				},
				ObjectMeta: meta_v1.ObjectMeta{
					Name:      key,
					Namespace: ts.cfg.Namespace,
					Labels: map[string]string{
						"name": key,
					},
				},
				Data: map[string]string{key: val},
			}, meta_v1.CreateOptions{})
		cancel()
		took := time.Since(start)
		tookMS := float64(took / time.Millisecond)
		readRequestLatencyMs.Observe(tookMS)
		latencies = append(latencies, took)
		if err != nil {
			readRequestsFailureTotal.Inc()
			ts.cfg.Logger.Warn("read configmap failed", zap.String("namespace", ts.cfg.Namespace), zap.Error(err))
		} else {
			readRequestsSuccessTotal.Inc()
			if i%20 == 0 {
				ts.cfg.Logger.Info("read configmap", zap.Int("iteration", i), zap.String("namespace", ts.cfg.Namespace))
			}
		}
	}
	return latencies
}
