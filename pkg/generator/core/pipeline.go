package core

type Pipeline struct {
	steps []Step
}

func (p Pipeline) Execute() error {
	return nil
}
