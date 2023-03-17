package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ConcurrentStage(done In, in In, stage Stage) Out {
	addedStream := make(Bi)
	go func() {
		defer close(addedStream)
		for i := range stage(in) {
			select {
			case <-done:
				return
			case addedStream <- i:
			}
		}
	}()
	return addedStream
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {

	// Place your code here.
	proxyChannel := in
	for _, stage := range stages {
		proxyChannel = ConcurrentStage(done, proxyChannel, stage)
	}
	return proxyChannel
}
