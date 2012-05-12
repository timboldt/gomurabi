package main

import (
	"fmt"
	"gomurabi/kingdomstate"
	"math/rand"
	"time"
	"runtime"
	"flag"
	"sync"
)

func doit(wg *sync.WaitGroup, n int, randgen *rand.Rand) {	
	for i := 0; i < n; i++ {
		var ks kingdomstate.KingdomState

		ks.SetupInitialState(randgen)
		for ks.StillInOffice() {
			//ks.PrintSummary()
			ks.TallyUpYear(0, 50, 2000, 10)
		}

	}
	//fmt.Printf("Done %d\n", p)
	wg.Done()
}

func main() {
	var parthreads int
	flag.IntVar(&parthreads, "threads", 1, "# of threads to use")
	flag.Parse()
	
	runtime.GOMAXPROCS(parthreads)	
	fmt.Printf("CPUs=%d\nThreads=%d\n", runtime.NumCPU(), parthreads)
	
	var wg sync.WaitGroup
	n := 10000000/parthreads
	for i:=0; i<parthreads; i++ {
		randgen := rand.New(rand.NewSource(time.Now().UnixNano()))
		wg.Add(1)
		go doit(&wg, n, randgen)
	}
	wg.Wait()
	fmt.Printf("Done\n")
}
