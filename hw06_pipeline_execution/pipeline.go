package hw06pipelineexecution

import "log"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// если не передали ни одного этапа обработки - возвращаем входной канал in
	if len(stages) == 0 {
		return in
	}

	// начинаем с in. далее в цикле канал будет подменяться
	// после выполнения последнего в списке этапа (stage) будет содержать конечные результаты работы pipeline
	out := in

	for _, stage := range stages {
		// создаем новый канал
		funcIn := make(Bi)

		go func(funcIn Bi, out Out) {
			// горутина единственный писатель в этот канал. поэтому закрываем при выходе из горутины
			defer func() {
				close(funcIn)
				// этот блок нужен в основном для обработки случаев отмены
				// нужен для того, чтобы дать возможность писателю в out сбросить значение в канал, докрутить
				// итерацию цикла внутри себя и на следующей итерации зафиксировать закрытие канала
				// без этого блока текущая горутина завершит работу и никто не сможет считать значение из out
				// писатель в out в этом случае зависнет и не сможет увидеть закрытие своего входящего канала
				for v := range out {
					log.Printf("Сброс из канала %v", v)
				}
			}()

			for {
				select {
				case <-done:
					return
				case v, ok := <-out:
					if !ok {
						return
					}
					funcIn <- v
				}
			}
		}(funcIn, out)

		// получаем канал с исходящими для текущего этапа данными
		// канал заполняемся асинхронно. возврат из stage(...) практически сразу
		out = stage(funcIn)
	}

	return out
}
