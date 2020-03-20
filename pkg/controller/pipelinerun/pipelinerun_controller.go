package pipelinerun

import (
	"bytes"
	"context"
	"fmt"
	"io"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/bigkevmcd/tekton-log-archiver/pkg/archiver"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

var log = logf.Log.WithName("controller_pipelinerun")

// Add creates a new PipelineRun Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	clientset, err := newClientSet()
	if err != nil {
		return err
	}
	return add(mgr, newReconciler(mgr, clientset))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, clientset *kubernetes.Clientset) reconcile.Reconciler {
	return &ReconcilePipelineRun{
		client:    mgr.GetClient(),
		scheme:    mgr.GetScheme(),
		clientset: clientset,
	}
}

func newClientSet() (*kubernetes.Clientset, error) {
	clusterConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get in cluster config: %w", err)
	}
	kubeClient, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create KubeClient: %w", err)
	}
	return kubeClient, nil
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("pipelinerun-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &pipelinev1.PipelineRun{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	return nil
}

// ReconcilePipelineRun reconciles a PipelineRun object
type ReconcilePipelineRun struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client    client.Client
	scheme    *runtime.Scheme
	clientset *kubernetes.Clientset
}

// Reconcile reads that state of the cluster for a PipelineRun object and makes changes based on the state read
// and what is in the PipelineRun.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcilePipelineRun) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling PipelineRun")
	ctx := context.Background()

	// Fetch the PipelineRun instance
	pipelineRun := &pipelinev1.PipelineRun{}
	err := r.client.Get(ctx, request.NamespacedName, pipelineRun)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	w := wrap(pipelineRun)
	if !archiver.Archivable(w) {
		reqLogger.Info("not a notifiable pipeline run")
		return reconcile.Result{}, nil
	}

	state := w.RunState()
	if !state.Complete() {
		return reconcile.Result{}, nil

	}

	for _, tr := range pipelineRun.Status.TaskRuns {
		if err != nil {
			continue
		}
		logs, err := logsForPod(ctx, request.Namespace, tr.Status.PodName, r.clientset)
		if err != nil {
			return reconcile.Result{}, err
		}
		log.Info(fmt.Sprintf("KEVIN!!! output from logs:\n%s\n", logs))
	}
	reqLogger.Info("archived logs")
	return reconcile.Result{}, nil
}

func logsForPod(ctx context.Context, ns, name string, c *kubernetes.Clientset) (string, error) {
	podLogOpts := corev1.PodLogOptions{}
	req := c.CoreV1().Pods(ns).GetLogs(name, &podLogOpts)
	podLogs, err := req.Stream()
	if err != nil {
		return "", fmt.Errorf("error in opening stream: %w", err)
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", fmt.Errorf("error in copy information from podLogs to buf: %w", err)
	}
	str := buf.String()
	return str, nil
}
