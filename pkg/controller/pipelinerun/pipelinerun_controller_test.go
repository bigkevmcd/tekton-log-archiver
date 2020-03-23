package pipelinerun

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"knative.dev/pkg/apis"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	tb "github.com/tektoncd/pipeline/test/builder"

	"github.com/bigkevmcd/tekton-log-archiver/pkg/archiver"
	ctb "github.com/bigkevmcd/tekton-log-archiver/test/builder"
)

var (
	testNamespace   = "test-namespace"
	pipelineRunName = "test-pipeline-run"
)

var _ reconcile.Reconciler = &ReconcilePipelineRun{}

func TestPipelineRunControllerUploadsContent(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	pipelineRun := ctb.MakePipelineRunWithResources()
	applyOpts(
		pipelineRun,
		tb.PipelineRunAnnotation(archiver.ArchivableName, "true"),
		tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
			apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionTrue}),
			tb.PipelineRunTaskRunsStatus("my-test-taskrun", &pipelinev1.PipelineRunTaskRunStatus{
				PipelineTaskName: "task-1",
				Status:           &pipelinev1.TaskRunStatus{},
			})))

	objs := []runtime.Object{
		pipelineRun,
	}
	r := makeReconciler(pipelineRun, objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      pipelineRunName,
			Namespace: testNamespace,
		},
	}
	res, err := r.Reconcile(req)
	fatalIfError(t, err, "reconcile: (%v)", err)
	if res.Requeue {
		t.Fatal("reconcile requeued request")
	}
}

func applyOpts(pr *pipelinev1.PipelineRun, opts ...tb.PipelineRunOp) {
	for _, o := range opts {
		o(pr)
	}
}

func makeReconciler(pr *pipelinev1.PipelineRun, objs ...runtime.Object) *ReconcilePipelineRun {
	s := scheme.Scheme
	s.AddKnownTypes(pipelinev1.SchemeGroupVersion, pr)
	cl := fake.NewFakeClient(objs...)
	return &ReconcilePipelineRun{
		client: cl,
		scheme: s,
	}
}

func fatalIfError(t *testing.T, err error, format string, a ...interface{}) {
	if err != nil {
		t.Fatalf(format, a...)
	}
}
