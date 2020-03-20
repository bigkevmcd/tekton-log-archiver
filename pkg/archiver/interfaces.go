package archiver

type stateGetter interface {
	RunState() State
}

type annotationsGetter interface {
	Annotations() map[string]string
}

type trackableResource interface {
	stateGetter
	annotationsGetter
}
