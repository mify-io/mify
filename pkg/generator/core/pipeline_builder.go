package core

type PipelineBuilder struct {
	steps []Step
}

func NewPipelineBuilder() *PipelineBuilder {
	return &PipelineBuilder{
		steps: make([]Step, 0),
	}
}

func (pb *PipelineBuilder) Register(step Step) *PipelineBuilder {
	pb.steps = append(pb.steps, step)
	return pb
}

func (pb *PipelineBuilder) Build() Pipeline {
	return Pipeline{
		steps: pb.steps,
	}
}
