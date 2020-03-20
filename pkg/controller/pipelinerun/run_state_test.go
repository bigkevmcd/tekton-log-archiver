package pipelinerun

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"

	"github.com/bigkevmcd/tekton-log-archiver/pkg/archiver"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	tb "github.com/tektoncd/pipeline/test/builder"
)

func TestGetPipelineRunStatus(t *testing.T) {
	statusTests := []struct {
		conditionType   apis.ConditionType
		conditionStatus corev1.ConditionStatus
		want            archiver.State
	}{
		{apis.ConditionSucceeded, corev1.ConditionTrue, archiver.Successful},
		{apis.ConditionSucceeded, corev1.ConditionUnknown, archiver.Pending},
		{apis.ConditionSucceeded, corev1.ConditionFalse, archiver.Failed},
	}

	for _, tt := range statusTests {
		w := pipelineRunWrapper{makePipelineRunWithCondition(tt.conditionType, tt.conditionStatus)}
		s := w.RunState()
		if s != tt.want {
			t.Errorf("RunState(%s) got %v, want %v", tt.conditionStatus, s, tt.want)
		}
	}
}

func makePipelineRunWithCondition(s apis.ConditionType, c corev1.ConditionStatus) *pipelinev1.PipelineRun {
	return tb.PipelineRun(pipelineRunName, testNamespace, tb.PipelineRunSpec(
		"tomatoes",
	), tb.PipelineRunStatus(tb.PipelineRunStatusCondition(
		apis.Condition{Type: s, Status: c}),
		tb.PipelineRunTaskRunsStatus("trname", &pipelinev1.PipelineRunTaskRunStatus{
			PipelineTaskName: "task-1",
		})))
}
