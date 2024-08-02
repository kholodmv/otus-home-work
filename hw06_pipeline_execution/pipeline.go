package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, terminate In, stages ...Stage) Out {
	for _, stage := range stages {
		in = runStage(stage, in, terminate)
	}
	return in
}

func runStage(stage Stage, in In, terminate In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		stageOut := stage(in)

		for {
			select {
			case v, ok := <-stageOut:
				if !ok {
					return
				}
				select {
				case out <- v:
				case <-terminate:
					return
				}
			case <-terminate:
				return
			}
		}
	}()
	return out
}
