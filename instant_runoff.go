package govote

import (
   "errors"
 )


type IRVPoll struct {
   candidates     []string
   ballots        [][]string
}

// Evaluate poll; returns list of winners as slice of candidate names and
// a slice of name-key score-value maps representing separate rounds
func (p *IRVPoll) Evaluate() ([]string, []map[string]int, error) {
   if p.candidates == nil || p.ballots == nil { 
      return []string{}, []map[string]int{}, errors.New("no candidates or no ballots") 
   }
   winners, rounds := p.getWinners()
   return winners, rounds, nil
}

// AddBallot submits a ballot to the poll, returns true on success, false on failure
func (p *IRVPoll) AddBallot(ballot []string) bool {
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

func (p *IRVPoll) runRound(elim map[string]bool) map[string]int {
   tally := make(map[string]int)                // scores keyed by candidate name
   for _, v := range p.ballots {
      for i := 0; i < len(v); i++ {
         if !elim[v[i]] {
            tally[v[i]]++
            break
         }
      }
   }
   return tally
}

func (p *IRVPoll) getWinners() (winners []string, rounds []map[string]int) {
   elim := make(map[string]bool)                // eliminated candidates
   tally := make(map[string]int)                // scores keyed by candidate name
   ct := len(p.candidates)                      // number of candidates

   for i := 0; i < ct; i++ {
      elim[p.candidates[i]] = false
   }

   for true {
      tally = p.runRound(elim)
      rounds = append(rounds, tally)

      // Figure out the lowest and highest score
      min, max := maxInt, 0
      for k, v := range tally {
         if !elim[k] && v < min { min = v }
         if !elim[k] && v > max { max = v }
      }

      if min == max { // victory or tie for the win
         break;
      }

      targets := []string{} // targets for elimination
      for k, v := range tally {
         if v == min { targets = append(targets, k) }
      }

      if len(targets) == 1 {
         elim[targets[0]] = true
      } else { // There's tied losers      
         for _, v := range tally {
            // Is their combined score less than the winner's?
            // If so, throw them all out.
            if v > min * len(targets) {
               for _, c := range targets {
                  elim[c] = true
               }
               targets = []string{}
               break
            } else { // otherwise, pick one at random
               mo := targets[randIntn(len(targets))]
               elim[mo] = true
            }
         }
      }
   }

   for i := 0; i < ct; i++ {
      if !elim[p.candidates[i]]  {
         winners = append(winners, p.candidates[i])
      }
   }

   return
}