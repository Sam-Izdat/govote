package govote

import (
	"errors"
)

// PluralityPoll is a poll using a simple plurality method (candidate with most votes wins)
type PluralityPoll struct {
	candidates []string
	ballots    []string
}

// Evaluate poll; returns list of winners as slice of candidate names and
// ranking as slice of CScores
func (p *PluralityPoll) Evaluate() ([]string, []CScore, error) {
	if p.candidates == nil || p.ballots == nil {
		return []string{}, []CScore{}, errors.New("no candidates or no ballots")
	}
	winners, ranks := p.getWinners()
	return winners, ranks, nil
}

// AddBallot submits a ballot to the poll, returns true on success, false on failure
func (p *PluralityPoll) AddBallot(ballot string) bool {
	ok := false
	for _, cv := range p.candidates {
		if cv == ballot {
			ok = true
			break
		}
	}
	if !ok {
		return false
	}
	p.ballots = append(p.ballots, ballot)
	return true
}

func (p *PluralityPoll) getWinners() (winners []string, ranks []CScore) {
	tally := make(map[string]int) // scores keyed by candidate name
	ct := len(p.candidates)       // number of candidates
	for i := 0; i < len(p.ballots); i++ {
		tally[p.ballots[i]]++
	}
	ranks = sortScoresDesc(tally)
	for i := 0; i < ct; i++ {
		if tally[p.candidates[i]] == ranks[0].Score {
			winners = append(winners, p.candidates[i])
		}
	}
	return
}
