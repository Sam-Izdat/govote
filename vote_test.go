package govote

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"
)

func democrats() map[string]int {
	return map[string]int{
		"Biden":        29,
		"Sanders":      20,
		"Harris":       5,
		"O'Rourke":     4,
		"Warren":       4,
		"Booker":       3,
		"Delaney":      3,
		"Klobuchar":    3,
		"Castro":       1,
		"Gabbard":      1,
		"Gillibrand":   1,
		"Hickenlooper": 1,
		"Inslee":       1,
		"Yang":         1,
		"Buttigieg":    0,
		"Williamson":   0,
	}
}

func candidateList(weights map[string]int) []string {
	out := make([]string, 0, len(weights))
	for name := range weights {
		out = append(out, name)
	}
	return out
}

func weighted(weights map[string]int) []*string {
	i := 0
	length := 100 + len(weights)
	result := make([]*string, 0, length)
	for name, v := range weights {
	REDO:
		n := name // BUG: loop capture?
		result = append(result, &n)
		i++
		v--
		if v >= 0 {
			goto REDO
		}
	}
	return result
}

func randomBallot(random *rand.Rand, choices []*string, min, max int) []string {
	if min > max {
		panic("min > max")
	}
	if min < 1 {
		panic("min < 1")
	}
	if max > len(choices) {
		panic("max > len(choices)")
	}

	count := random.Intn(max-min+1) + min
	votes := make(map[string]bool, count)
	output := make([]string, count, count)
	for i := 0; i < count; i++ {
	RETRY: // suboptimal but better than building a giant decision tree in memory since we don't expect more than a handful of votes per
		index := random.Intn(len(choices))
		vote := *choices[index]
		if exists := votes[vote]; exists {
			goto RETRY
		}
		votes[vote] = true
		output[i] = vote
	}
	return output
}

// always store the result to a package level variable
// so the compiler cannot eliminate the Benchmark itself.
var resultSchulze []string

func BenchmarkSchulze1Procs(b *testing.B) {
	benchmarkSchulze(1, b)
}

func BenchmarkSchulze2Procs(b *testing.B) {
	benchmarkSchulze(2, b)
}

func BenchmarkSchulzeGOMAXPROCS(b *testing.B) {
	benchmarkSchulze(runtime.GOMAXPROCS(0), b)
}

func BenchmarkSchulze2xGOMAXPROCS(b *testing.B) {
	benchmarkSchulze(2*runtime.GOMAXPROCS(0), b)
}

var benchmarkSchulzeSetupOnce sync.Once
var benchmarkSchulzePoll SchulzePoll

func benchmarkSchulzeSetup() {
	const population = 10000

	random := rand.New(rand.NewSource(2019))

	dems := democrats()
	candidates := candidateList(dems)
	deck := weighted(dems)
	poll, _ := Schulze.New(candidates)

	length := runtime.GOMAXPROCS(0)
	c := make(chan []string, length)
	quitters := make([]chan bool, length)
	for i := 0; i < length; i++ {
		quitters[i] = make(chan bool)
		s := rand.NewSource(random.Int63())
		random := rand.New(s)
		go func(quit chan bool) {
		RETRY:
			ballot := randomBallot(random, deck, 1, len(candidates))
			select {
			case <-quit:
			case c <- ballot:
				goto RETRY
			}
		}(quitters[i])
	}
	for i := 0; i < population; i++ {
		ballot := <-c
		poll.AddBallot(ballot)
	}
	for i := 0; i < length; i++ {
		quitters[i] <- true
	}
	close(c)
	benchmarkSchulzePoll = poll
}

func benchmarkSchulze(numProcs int, b *testing.B) {
	benchmarkSchulzeSetupOnce.Do(benchmarkSchulzeSetup)

	oldnumProcs := runtime.GOMAXPROCS(numProcs)
	defer runtime.GOMAXPROCS(oldnumProcs)
	// run the benchmark
	for n := 0; n < b.N; n++ {
		result, _, err := benchmarkSchulzePoll.Evaluate()
		resultSchulze = result
		if err != err {
			b.Fatalf("Evaluate returned error %v", err)
		}
	}
}
