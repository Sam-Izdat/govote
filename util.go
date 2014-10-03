package govote
 
import (
    "sort"
    "math/rand"
    "time"
)

// returns lesser value of a and b
func min(a, b int) (r int) {
    if a < b { r = a } else { r = b }
    return
 }

// returns greater value of a and b
func max(a, b int) (r int) {
    if a > b { r = a } else { r = b }
    return
}

// CPair represents two candidates' indices for pairwise comparison
type CPair struct { a, b int }

// CScore represents a candidate and the candidate's score
type CScore struct { 
    Name string 
    Score int 
}

// sortScoresAsc sorts a string-key/int-value map in ascending order,
// by value returning a slice of CScores
func sortScoresAsc (m map[string]int) (res []CScore) {
    vs := newValSorter(m)
    vs.sortAsc()
    for i := 0; i < len(vs.keys); i++ {
        res = append(res, CScore{vs.keys[i], vs.vals[i]})
    }
    return
}

// sortScoresDesc sorts a string-key/int-value map in descending order,
// by value returning a slice of CScores
func sortScoresDesc (m map[string]int) (res []CScore) {
    vs := newValSorter(m)
    vs.sortDesc()
    for i := 0; i < len(vs.keys); i++ {
        res = append(res, CScore{vs.keys[i], vs.vals[i]})
    }
    return
}

type valSorter struct {
    keys []string
    vals []int
}
 
func newValSorter(m map[string]int) *valSorter {
    vs := &valSorter{
        keys: make([]string, 0, len(m)),
        vals: make([]int, 0, len(m)),
    }
    for k, v := range m {
        vs.keys = append(vs.keys, k)
        vs.vals = append(vs.vals, v)
    }
    return vs
}
 
func (vs *valSorter) sortAsc() {
    sort.Sort(vs)
}

func (vs *valSorter) sortDesc() {
    sort.Sort(sort.Reverse(vs))
} 

func (vs *valSorter) Len() int           { return len(vs.vals) }
func (vs *valSorter) Less(i, j int) bool { return vs.vals[i] < vs.vals[j] }
func (vs *valSorter) Swap(i, j int) {
    vs.vals[i], vs.vals[j] = vs.vals[j], vs.vals[i]
    vs.keys[i], vs.keys[j] = vs.keys[j], vs.keys[i]
}

const maxUint = ^uint(0) 
const minUint = 0 
const maxInt = int(maxUint >> 1) 
const minInt = -maxInt - 1

// Returns integer between 0 and n
func randIntn(n int) int {
    if n == 0 { return 0 }
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    return int( r.Intn(n) )
}

// Removes duplicates from string slice
func removeDuplicates(xs *[]string) {
    found := make(map[string]bool)
    j := 0
    for i, x := range *xs {
        if !found[x] {
            found[x] = true
            (*xs)[j] = (*xs)[i]
            j++
        }
    }
    *xs = (*xs)[:j]
}