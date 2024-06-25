package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Place your code here.
	stageIn := in
	for _, stage := range stages {
		out := stage(stageIn)
		stageIn = out
	}

	select {
	case <-done:
		close(in)

	}

	return nil
}
