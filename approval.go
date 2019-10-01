package govote

import (
	"errors"
)

// ApprovalPoll is a poll, where the voters approve one or multiple candidates
type ApprovalPoll struct {
	candidates []string
	ballots    []string
}

// Evaluate poll; returns list of winners as slice of candidate names and
// ranking as slice of CScores
func (p *ApprovalPoll) Evaluate() ([]string, []CScore, error) {
	if p.candidates == nil || p.ballots == nil {
		return []string{}, []CScore{}, errors.New("no candidates or no ballots")
	}
	winners, ranks := p.getWinners()
	return winners, ranks, nil
}

// AddBallot submits a ballot to the poll, returns true on success, false on failure
func (p *ApprovalPoll) AddBallot(ballot []string) bool {
	if hasDuplicates(ballot) {
		return false //sort out invalid votes
	}

	for i := 0; i < len(ballot); i++ {
		p.ballots = append(p.ballots, ballot[i])
	}
	return true
}

func (p *ApprovalPoll) getWinners() (winners []string, ranks []CScore) {
	tally := make(map[string]int) // scores keyed by candidate name

	for i := 0; i < len(p.ballots); i++ {
		tally[p.ballots[i]]++
	}
	ranks = sortScoresDesc(tally)

	ct := len(p.candidates) // number of candidates
	for i := 0; i < ct; i++ {
		if tally[p.candidates[i]] == ranks[0].Score {
			winners = append(winners, p.candidates[i])
		}
	}
	return
}
