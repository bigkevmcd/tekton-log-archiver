package archiver

import (
	"testing"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	tb "github.com/tektoncd/pipeline/test/builder"
)

func TestArchivable(t *testing.T) {
	nt := []struct {
		name string
		opts []tb.PipelineRunOp
		want bool
	}{
		{"no labels", nil, false},
		{"no archive label", []tb.PipelineRunOp{tb.PipelineRunAnnotation("testing", "app")}, false},
		{"archive label", []tb.PipelineRunOp{tb.PipelineRunAnnotation(ArchivableName, "true")}, true},
		{"archive label is false", []tb.PipelineRunOp{tb.PipelineRunAnnotation(ArchivableName, "false")}, false},
	}

	for _, tt := range nt {
		r := tb.PipelineRun("test-pipeline-run-with-labels", "foo", tt.opts...)
		if b := Archivable(wrapper{r}); b != tt.want {
			t.Errorf("Archivable() %s got %v, want %v", tt.name, b, tt.want)
		}
	}
}

type wrapper struct {
	*pipelinev1.PipelineRun
}

func (w wrapper) Annotations() map[string]string {
	return w.PipelineRun.Annotations
}
