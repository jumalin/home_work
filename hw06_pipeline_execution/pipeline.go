package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func TerminatableStage(done In, in In, stage Stage) Out {
	addedStream := make(Bi)
	go func() {
		defer close(addedStream)
		for stageResult := range stage(in) {
			select {
			case <-done:
				return
			case addedStream <- stageResult:
			}
		}
	}()
	return addedStream
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	proxyChannel := in
	for i := range stages {
		proxyChannel = TerminatableStage(done, proxyChannel, stages[i])
	}
	return proxyChannel
}
