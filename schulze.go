package govote

import (
   "errors"
)

// SchulzePoll is a Condorcet poll using the Schulze method 
type SchulzePoll struct {
   candidates     []string
   ballots        [][]string
   PP             map[CPair]int // pairwise preferences
   SP             map[CPair]int // strongest path
}

// Evaluate poll; returns list of winners as slice of candidate names and
// ranking as slice of CScores
func (p *SchulzePoll) Evaluate() ([]string, []CScore, error) {
   if p.candidates == nil || p.ballots == nil { 
      return []string{}, []CScore{}, errors.New("no candidates or no ballots") 
   }
   winners, ranks := p.getWinners()
   return winners, ranks, nil
}

// AddBallot submits a ballot to the poll, returns true on success, false on failure
func (p *SchulzePoll) AddBallot(ballot []string) bool {
   removeDuplicates(&ballot)
   if len(ballot) == 0 { return false }
   var ok bool
   for _, bv := range ballot {
      ok = false
      for _, cv := range p.candidates {
         if cv == bv { 
            ok = true 
            break
         }
      }
      if !ok { return false }
   }
   p.ballots = append(p.ballots, ballot)
   return true
}

// Returns the rank of a candidate on a specific ballot or -1 if unranked
func (p SchulzePoll) getBallotRank(idx, i int) int {
   ballot := p.ballots[idx]
   for k, v := range ballot {
      if v == p.candidates[i] { return k }
   }
   return -1
}

// Returns number of voters strictly prefering candidate i to j
func (p SchulzePoll) comparePref(i, j int) int {
   count := 0
   for k := range p.ballots {
      ri, rj := p.getBallotRank(k, i), p.getBallotRank(k, j)
      if ri == -1 || rj == -1 { continue } 
      if ri < rj { count++ }
   }
   return count
}

func (p *SchulzePoll) getWinners() (winners []string, ranks []CScore) {
   p.PP = make(map[CPair]int)                   // pairwise preferences
   p.SP = make(map[CPair]int)                   // strongest paths
   won := make([]bool, len(p.candidates))       // winners marked true by index
   ct := len(p.candidates)                      // number of candidates
   tally := make(map[string]int)                // scores keyed by candidate name

   // compute pairwise preferences
   for i := 0; i < ct; i++ {
      for j := 0; j < ct; j++ {
         if i != j {
            p.PP[CPair{i, j}] = p.comparePref(i, j)
         }
      }
   }

   // compute strongest paths
   for i := 0; i < ct; i++ {
      for j := 0; j < ct; j++ {
         if i != j {
            if p.PP[CPair{i, j}] > p.PP[CPair{j, i}] {
               p.SP[CPair{i, j}] = p.PP[CPair{i, j}]
            } else {
               p.SP[CPair{i, j}] = 0
            }
         }
      }
   }

   for i := 0; i < ct; i++ {
      for j := 0; j < ct; j++ {
         if i != j {
            for k := 0; k < ct; k++ {
               if i != k && j != k {
                  var(
                     SPjk = p.SP[CPair{j, k}]
                     SPji = p.SP[CPair{j, i}]
                     SPik = p.SP[CPair{i, k}]
                  ) 
                  p.SP[CPair{j, k}] = max(SPjk, min(SPji, SPik))
               }
            }
         }
      }
   }

   // mark winners
   for i := 0; i < ct; i++ {
      won[i] = true
   }

   for i := 0; i < ct; i++ {
      for j := 0; j < ct; j++ {
         if i != j {
            if p.SP[CPair{j, i}] > p.SP[CPair{i, j}] {
               won[i] = false
               tally[p.candidates[i]] += 0   // creates entry for Condorcet loser
            } else {
               tally[p.candidates[i]]++      // candidate preferred to x others
            }
         }
      }
   }
   ranks = sortScoresDesc(tally)

   // make winners list
   for i := 0; i < ct; i++ {
      if won[i] == true { 
         winners = append(winners, p.candidates[i])
      }
   }
   return
}