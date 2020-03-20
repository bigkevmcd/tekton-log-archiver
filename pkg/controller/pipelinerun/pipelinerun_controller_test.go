package pipelinerun

import (
	"testing"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	tb "github.com/tektoncd/pipeline/test/builder"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"knative.dev/pkg/apis"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/bigkevmcd/tekton-log-archiver/pkg/archiver"
	ctb "github.com/bigkevmcd/tekton-log-archiver/test/builder"
)

var (
	testNamespace   = "test-namespace"
	pipelineRunName = "test-pipeline-run"
	testToken       = "abcdefghijklmnopqrstuvwxyz12345678901234"
)

var _ reconcile.Reconciler = &ReconcilePipelineRun{}

// TestPipelineRunControllerPendingState runs ReconcilePipelineRun.Reconcile() against a
// fake client that tracks PipelineRun objects.
func TestPipelineRunControllerPendingState(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	pipelineRun := ctb.MakePipelineRunWithResources()
	applyOpts(
		pipelineRun,
		tb.PipelineRunAnnotation(archiver.ArchivableName, "true"),
		tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
			apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionUnknown})))
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

// // TestPipelineRunReconcileWithPreviousPending tests a PipelineRun that
// // we've already sent a pending notification.
// func TestPipelineRunReconcileWithPreviousPending(t *testing.T) {
// 	logf.SetLogger(logf.ZapLogger(true))
// 	pipelineRun := ctb.MakePipelineRunWithResources(
// 		ctb.MakeGitResource("https://github.com/tektoncd/triggers", "master"))
// 	applyOpts(
// 		pipelineRun,
// 		tb.PipelineRunAnnotation(archive.NotifiableName, "true"),
// 		tb.PipelineRunAnnotation(archive.StatusContextName, "test-context"),
// 		tb.PipelineRunAnnotation(archive.StatusDescriptionName, "testing"),
// 		tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
// 			apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionUnknown})))
// 	objs := []runtime.Object{
// 		pipelineRun,
// 		ctb.MakeSecret(archive.SecretName, map[string][]byte{"token": []byte(testToken)}),
// 	}

// 	r, data := makeReconciler(pipelineRun, objs...)

// 	req := reconcile.Request{
// 		NamespacedName: types.NamespacedName{
// 			Name:      pipelineRunName,
// 			Namespace: testNamespace,
// 		},
// 	}
// 	// This runs Reconcile twice.
// 	res, err := r.Reconcile(req)
// 	fatalIfError(t, err, "reconcile: (%v)", err)
// 	if res.Requeue {
// 		t.Fatal("reconcile requeued request")
// 	}
// 	// This cleans out the existing date for the data, because the fake scm
// 	// client updates in-place, so there's no way to know if it received multiple
// 	// pending notifications.
// 	delete(data.Statuses, "master")
// 	res, err = r.Reconcile(req)
// 	fatalIfError(t, err, "reconcile: (%v)", err)
// 	if res.Requeue {
// 		t.Fatal("reconcile requeued request")
// 	}
// 	// There should be no recorded statuses, because the state is still pending
// 	// and the fake client's state was deleted above.
// 	assertNoStatusesRecorded(t, data)
// }

// // TestPipelineRunControllerSuccessState runs ReconcilePipelineRun.Reconcile() against a
// // fake client that tracks PipelineRun objects.
// func TestPipelineRunControllerSuccessState(t *testing.T) {
// 	logf.SetLogger(logf.ZapLogger(true))
// 	pipelineRun := ctb.MakePipelineRunWithResources(
// 		ctb.MakeGitResource("https://github.com/tektoncd/triggers", "master"))
// 	applyOpts(
// 		pipelineRun,
// 		tb.PipelineRunAnnotation(archive.NotifiableName, "true"),
// 		tb.PipelineRunAnnotation(archive.StatusContextName, "test-context"),
// 		tb.PipelineRunAnnotation(archive.StatusDescriptionName, "testing"),
// 		tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
// 			apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionTrue})))
// 	objs := []runtime.Object{
// 		pipelineRun,
// 		ctb.MakeSecret(archive.SecretName, map[string][]byte{"token": []byte(testToken)}),
// 	}
// 	r, data := makeReconciler(pipelineRun, objs...)

// 	req := reconcile.Request{
// 		NamespacedName: types.NamespacedName{
// 			Name:      pipelineRunName,
// 			Namespace: testNamespace,
// 		},
// 	}
// 	res, err := r.Reconcile(req)
// 	fatalIfError(t, err, "reconcile: (%v)", err)
// 	if res.Requeue {
// 		t.Fatal("reconcile requeued request")
// 	}
// 	wanted := &scm.Status{State: scm.StateSuccess, Label: "test-context", Desc: "testing", Target: ""}
// 	status := data.Statuses["master"][0]
// 	if !reflect.DeepEqual(status, wanted) {
// 		t.Fatalf("log-archive notification got %#v, wanted %#v\n", status, wanted)
// 	}
// }

// // TestPipelineRunControllerFailedState runs ReconcilePipelineRun.Reconcile() against a
// // fake client that tracks PipelineRun objects.
// func TestPipelineRunControllerFailedState(t *testing.T) {
// 	logf.SetLogger(logf.ZapLogger(true))
// 	pipelineRun := ctb.MakePipelineRunWithResources(
// 		ctb.MakeGitResource("https://github.com/tektoncd/triggers", "master"))
// 	applyOpts(
// 		pipelineRun,
// 		tb.PipelineRunAnnotation(archive.NotifiableName, "true"),
// 		tb.PipelineRunAnnotation(archive.StatusContextName, "test-context"),
// 		tb.PipelineRunAnnotation(archive.StatusDescriptionName, "testing"),
// 		tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
// 			apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionFalse})))
// 	objs := []runtime.Object{
// 		pipelineRun,
// 		ctb.MakeSecret(archive.SecretName, map[string][]byte{"token": []byte(testToken)}),
// 	}
// 	r, data := makeReconciler(pipelineRun, objs...)

// 	req := reconcile.Request{
// 		NamespacedName: types.NamespacedName{
// 			Name:      pipelineRunName,
// 			Namespace: testNamespace,
// 		},
// 	}
// 	res, err := r.Reconcile(req)
// 	fatalIfError(t, err, "reconcile: (%v)", err)
// 	if res.Requeue {
// 		t.Fatal("reconcile requeued request")
// 	}
// 	wanted := &scm.Status{State: scm.StateFailure, Label: "test-context", Desc: "testing", Target: ""}
// 	status := data.Statuses["master"][0]
// 	if !reflect.DeepEqual(status, wanted) {
// 		t.Fatalf("log-archive notification got %#v, wanted %#v\n", status, wanted)
// 	}
// }

// // TestPipelineRunReconcileWithNoGitCredentials tests a non-archivable
// // PipelineRun.
// func TestPipelineRunReconcileNonNotifiable(t *testing.T) {
// 	logf.SetLogger(logf.ZapLogger(true))
// 	pipelineRun := ctb.MakePipelineRunWithResources(
// 		ctb.MakeGitResource("https://github.com/tektoncd/triggers", "master"))
// 	applyOpts(
// 		pipelineRun,
// 		tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
// 			apis.Condition{Type: apis.ConditionSucceeded, Status: corev1.ConditionUnknown})))
// 	objs := []runtime.Object{
// 		pipelineRun,
// 		ctb.MakeSecret(archive.SecretName, map[string][]byte{"token": []byte(testToken)}),
// 	}
// 	r, data := makeReconciler(pipelineRun, objs...)

// 	req := reconcile.Request{
// 		NamespacedName: types.NamespacedName{
// 			Name:      pipelineRunName,
// 			Namespace: testNamespace,
// 		},
// 	}
// 	res, err := r.Reconcile(req)
// 	fatalIfError(t, err, "reconcile: (%v)", err)
// 	if res.Requeue {
// 		t.Fatal("reconcile requeued request")
// 	}
// 	assertNoStatusesRecorded(t, data)
// }

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
