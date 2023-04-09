package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func TerminableStage(done In, in In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case input, ok := <-in:
				if !ok {
					return
				}
				out <- input

			case <-done:
				return
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	proxyChannel := in
	for i := range stages {
		proxyChannel = TerminableStage(done, proxyChannel)
		proxyChannel = stages[i](proxyChannel)
	}

	return proxyChannel
}
