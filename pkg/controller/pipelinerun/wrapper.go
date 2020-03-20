package pipelinerun

import (
	"github.com/bigkevmcd/tekton-log-archiver/pkg/archiver"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

type pipelineRunWrapper struct {
	*pipelinev1.PipelineRun
}

func wrap(pr *pipelinev1.PipelineRun) pipelineRunWrapper {
	return pipelineRunWrapper{pr}
}

// RunState returns whether or not a PipelineRun was successful or
// not.
func (p pipelineRunWrapper) RunState() archiver.State {
	return archiver.ConditionsToState(p.Status.Conditions)
}

func (p pipelineRunWrapper) Annotations() map[string]string {
	return p.PipelineRun.Annotations
}
