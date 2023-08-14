package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		bi := make(Bi)
		go func(_bi Bi, _out Out) {
			defer close(_bi)
			for {
				select {
				case <-done:
					return
				case value, ok := <-_out:
					if !ok {
						return
					}
					_bi <- value
				}
			}
		}(bi, out)
		out = stage(bi)
	}

	return out
}
