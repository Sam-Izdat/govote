package govote

import "testing"

func createResults() map[string]int {
	return map[string]int{
		"Most":  5,
		"Least": 0,
		"One":   1,
		"Two":   2,
	}
}

func Test_min(t *testing.T) {
	a, b := 5, 3
	result := min(a, b)

	if b != result {
		t.Errorf("%v should be %v", result, b)
	}
}

func Test_min_reverse(t *testing.T) {
	a, b := 3, 5
	result := min(a, b)

	if a != result {
		t.Errorf("%v should be %v", result, b)
	}
}

func Benchmark_min(b *testing.B) {
	for i := 0; i < b.N; i++ {
		min(i, b.N)
	}
}

func Test_max(t *testing.T) {
	a, b := 5, 3
	result := max(a, b)

	if a != result {
		t.Errorf("%v should be %v", result, a)
	}
}

func Test_max_reverse(t *testing.T) {
	a, b := 3, 5
	result := max(a, b)

	if b != result {
		t.Errorf("%v should be %v", result, b)
	}
}

func Benchmark_max(b *testing.B) {
	for i := 0; i < b.N; i++ {
		max(i, b.N)
	}
}

func Test_sortScoresAsc(t *testing.T) {
	m := createResults()
	results := sortScoresAsc(m)
	expected := []int{0, 1, 2, 5}

	for i, k := range results {
		if expected[i] != k.Score {
			t.Errorf("Expected %v to be %v at index %v", k, expected[i], i)
		}
	}
}

func Benchmark_sortScoresAsc(b *testing.B) {
	m := createResults()

	for i := 0; i < b.N; i++ {
		sortScoresAsc(m)
	}
}

func Test_sortScoresDesc(t *testing.T) {
	m := createResults()
	results := sortScoresDesc(m)
	expected := []int{5, 2, 1, 0}

	for i, k := range results {
		if expected[i] != k.Score {
			t.Errorf("Expected %v to be %v but was %v", k, expected[i], k.Score)
		}
	}
}

func Benchmark_sortScoresDesc(b *testing.B) {
	m := createResults()

	for i := 0; i < b.N; i++ {
		sortScoresDesc(m)
	}
}

func Test_randIntn_PassZero(t *testing.T) {
	result := randIntn(0)

	if result != 0 {
		t.Error("Result should be 0 when 0 is provided")
	}
}

func Test_randIntn_PassNumber(t *testing.T) {
	result := randIntn(10)

	if result < 0 || result > 10 {
		t.Error("Result is outside the bounds given")
	}
}

func Benchmark_randIntn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randIntn(i)
	}
}

func Test_removeDuplicates(t *testing.T) {
	candidates := []string{"Memphis", "Nashville", "Knoxville", "Chattanooga", "Memphis"}

	removeDuplicates(&candidates)

	if len(candidates) != 4 {
		t.Error("Duplicates not removed")
	}
}

func Benchmark_removeDuplicates(b *testing.B) {
	for i := 0; i < b.N; i++ {
		candidates := []string{"Memphis", "Nashville", "Knoxville", "Chattanooga", "Memphis"}

		removeDuplicates(&candidates)
	}
}
