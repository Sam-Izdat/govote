
![govote](doc/logo.png)

...is a Go library for conducting polls using various voting systems

## Install
Grab the package with: 

    $ go get github.com/Sam-Izdat/govote

[![License MIT](http://img.shields.io/badge/license-MIT-red.svg?style=flat-square)](http://opensource.org/licenses/MIT)
[![GoDoc](http://img.shields.io/badge/doc-REFERENCE-blue.svg?style=flat-square)](https://godoc.org/github.com/Sam-Izdat/govote)

## How do I even...
Behold.
```go
package main

import(
    "fmt"
    "github.com/Sam-Izdat/govote"
)

func main() {
    candidates := []string{"Kang", "Kodos"}
    poll, _ := govote.Plurality.New(candidates)
    poll.AddBallot("Kang")
    poll.AddBallot("Kang")
    poll.AddBallot("Kodos")
    fmt.Println(poll.Evaluate())
    // => [Kang] [{Kang 2} {Kodos 1}] <nil>
}
```
There's a more interesting example [used on Wikipedia](http://en.wikipedia.org/wiki/Condorcet_method#Example:_Voting_on_the_location_of_Tennessee.27s_capital) for comparing voting systems, which shows how different methods can yield very different results. 

```go
// Tennessee voting on capital, suppose all voters want it close, blah, blah, blah...
candidates := []string{"Memphis", "Nashville", "Knoxville", "Chattanooga"}
schulze, _ := govote.Schulze.New(candidates)
plurality, _ := govote.Plurality.New(candidates)
runoff, _ := govote.InstantRunoff.New(candidates)
approval, _ := govote.Approval.New(candidates)

ballotMemphis := []string{"Memphis", "Nashville", "Chattanooga", "Knoxville"}
ballotNashville := []string{"Nashville", "Chattanooga", "Knoxville", "Memphis"}
ballotChattanooga := []string{"Chattanooga", "Knoxville", "Nashville", "Memphis"}
ballotKnoxville := []string{"Knoxville", "Chattanooga", "Nashville", "Memphis"}

for i := 0; i < 42; i++ {
    schulze.AddBallot(ballotMemphis)
    plurality.AddBallot(ballotMemphis[0])
    runoff.AddBallot(ballotMemphis)
    approval.AddBallot(ballotMemphis[0:2])
}
for i := 0; i < 26; i++ {
    schulze.AddBallot(ballotNashville)
    plurality.AddBallot(ballotNashville[0])
    runoff.AddBallot(ballotNashville)
    approval.AddBallot(ballotNashville[0:2])
}
for i := 0; i < 15; i++ {
    schulze.AddBallot(ballotChattanooga)
    plurality.AddBallot(ballotChattanooga[0])
    runoff.AddBallot(ballotChattanooga)
    approval.AddBallot(ballotChattanooga[0:2])
}
for i := 0; i < 17; i++ {
    schulze.AddBallot(ballotKnoxville)
    plurality.AddBallot(ballotKnoxville[0])
    runoff.AddBallot(ballotKnoxville)
    approval.AddBallot(ballotKnoxville[0:2])
}

// Schulze scores are a tally of superior strongest-path comparisons for the candidate
fmt.Println("Schulze")
fmt.Println(schulze.Evaluate())
// => [Nashville] [{Nashville 3} {Chattanooga 2} {Knoxville 1} {Memphis 0}] <nil>

// Plurality scores are the number of votes for the candidate
fmt.Println("Plurality")
fmt.Println(plurality.Evaluate())
// => [Memphis] [{Memphis 42} {Nashville 26} {Knoxville 17} {Chattanooga 15}] <nil>

// Instant-runoff returns a slice of rounds, each an ordered slice of candidate scores
fmt.Println("IRV")
fmt.Println(runoff.Evaluate())
// => [Knoxville] [map[Memphis:42 Nashville:26 Chattanooga:15 Knoxville:17] \
// =>   [ [{Memphis 42} {Nashville 26} {Knoxville 17} {Chattanooga 15}] \
// =>     [{Memphis 42} {Knoxville 32} {Nashville 26}] \
// =>     [{Knoxville 58} {Memphis 42}] \
// =>     [{Knoxville 100}] ]
// =>   <nil>

// Approval scores are a tally of approval votes, disregarding their specific ordinal preferences
fmt.Println("Approval")
fmt.Println(approval.Evaluate())
// => [Nashville] [{Nashville 68} {Chattanooga 58} {Memphis 42} {Knoxville 32}] <nil>
```

Keep in mind that multiple winners are possible. 

Condorcet polls may result in a tie and, in the event of a voting paradox, will return all candidates as winners. 

Instant runoff polls, in the event of loser ties at the end of a round, will either eliminate the tied candidates if the sum of their scores is lower than the score of the leader or else elimate a loser at random. If, in the final round, two or more candidates are tied for victory they will all be returned as winners.

Plurality polls will return multiple winner in the event of a tie. 

# What still needs doin'

Voting systems implemented:

- [x] Approval Method
- [x] Instant Runoff Method
- [ ] Minmax Method
- [x] Plurality Method
- [ ] Range Method
- [x] Schulze Method (Condorcet)
- [ ] Chain Method
- [ ] Majority Choice Method

Also need to:

- [ ] Write unit tests

# License

MIT
