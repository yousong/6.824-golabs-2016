package mapreduce

import (
	"fmt"
)

// schedule starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mr.files)
		nios = mr.nReduce
	case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

	ctask := make(chan int, ntasks)
	cdone := make(chan bool, 4)
	if mr.workersChannel == nil {
		var bufsiz int
		if ntasks >= nios {
			bufsiz = ntasks
		} else {
			bufsiz = nios
		}
		mr.workersChannel = make(chan string, bufsiz)
		go func() {
			for {
				w := <-mr.registerChannel
				mr.workersChannel <- w
			}
		}()
	}
	go func() {
		for i := 0; i < ntasks; i++ {
			ctask <- i
		}
	}()

	go func() {
		for {
			t, ok := <-ctask
			if ok == false {
				return
			}
			w := <-mr.workersChannel

			go func() {
				var res bool
				args := new(DoTaskArgs)
				args.JobName = mr.jobName
				args.TaskNumber = t
				args.NumOtherPhase = nios
				args.Phase = phase
				switch phase {
				case mapPhase:
					args.File = mr.files[t]
					res = call(w, "Worker.DoTask", args, new(struct{}))
				case reducePhase:
					res = call(w, "Worker.DoTask", args, new(struct{}))
				}
				if res == true {
					mr.workersChannel <- w
					cdone <- res
				} else {
					ctask <- t
				}
			}()
		}
	}()

	var doneCnt = 0
	for {
		res := <-cdone
		if res == true {
			doneCnt += 1
			if doneCnt == ntasks {
				close(ctask)
				break
			}
		}
	}

	fmt.Printf("Schedule: %v phase done\n", phase)
}
