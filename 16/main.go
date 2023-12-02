package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	input := aocutil.InputString()
	lines := aocutil.SplitLines(input)

	start := ParseValves(lines)

	part1(start)
}

func part1(start *Valve) {
	var moves Moves
	popMove := func() {
		moves = moves[:len(moves)-1]
	}

	// DFS to find the tunnel path with the highest flow rate.
	// Keep track of which valves have been visited to avoid loops.
	visited := aocutil.NewSet[*Valve](len(start.Tunnels))

	var traverse func(v *Valve, moves Moves) (rate, cost int, newMoves Moves)
	traverse = func(v *Valve, moves Moves) (rate, cost int, newMoves Moves) {
		// What to do: as we traverse down the tunnels, we need to accumulate
		// our flow rate and find the maximum. However, one constraint is that
		// we have to stay within the time limit.
		//
		// Our first return, rate, will be the maximum flow rate we can
		// accumulate for this path.
		//
		// Our second return, cost, will be the cost of the path. The cost is
		// the number of moves we need to make to get to the end of the path.
		//
		// Our third return, newMoves, will be the moves we need to make to
		// get to the end of the path.

		// If we have no more time left, we can't do anything.
		if moves.MinutesLeft() == 0 {
			return 0, 0, moves
		}
	}

	// traverse = func(v *Valve) (rate, cost int) {
	// 	log.Println("visiting", v.ID)
	// 	visited.Add(v)

	// 	var maxRate int
	// 	var curCost int
	// 	var open *Valve

	// 	for _, tunnel := range v.Tunnels {
	// 		if visited.Has(tunnel) {
	// 			continue
	// 		}

	// 		moves = append(moves, MoveToValve{ID: tunnel.ID})
	// 		defer popMove()

	// 		rate, cost := traverse(tunnel)
	// 		if rate > maxRate && cost < moves.MinutesLeft() {
	// 			log.Println("found new max", rate, "from", v.ID, "to", tunnel.ID)
	// 			curCost = cost
	// 			maxRate = rate
	// 			open = tunnel
	// 		}
	// 	}

	// 	if open != nil {
	// 		moves = append(moves, OpenValve{ID: open.ID})
	// 		curCost++
	// 	}

	// 	return maxRate + v.FlowRate, curCost
	// }

	traverse(start)
	fmt.Println("part 1:")
	fmt.Println("  moves:", moves)
	fmt.Println("  pressure:", moves.Pressure(start))
}

// ValveID is a valve ID.
type ValveID string

// Valve is a valve.
type Valve struct {
	ID       ValveID
	FlowRate int
	Tunnels  []*Valve // always sorted flow-rate first, graph-like, need to traverse it
}

var valveLineRe = regexp.MustCompile(
	`Valve ([A-Z]+) has flow rate=(\d+); tunnels? leads? to valves? ([A-Z, ]+)`)

// ParseValves parses the input of valves. The returned valve is the first
// valve, with other parsed valves in its Tunnels map.
func ParseValves(lines []string) *Valve {
	// Parse the valves into a flat map first.
	valves := make(map[ValveID]*Valve, len(lines))
	valveTunnels := make(map[ValveID][]ValveID, len(lines))
	var firstValve *Valve

	for _, line := range lines {
		m := valveLineRe.FindStringSubmatch(line)
		aocutil.Assertf(m != nil, "invalid valve line: %q", line)

		id := ValveID(m[1])
		rate := aocutil.Atoi[int](m[2])
		tunnels := strings.Split(m[3], ", ")

		valve := &Valve{
			ID:       id,
			FlowRate: rate,
			Tunnels:  make([]*Valve, 0, len(tunnels)),
		}

		valves[id] = valve
		valveTunnels[id] = aocutil.Transform(tunnels, func(s string) ValveID { return ValveID(s) })

		if firstValve == nil {
			firstValve = valve
		}
	}

	aocutil.Assertf(firstValve != nil, "no valves found")

	// Start traversing the tunnels.
	var traverse func(*Valve)
	traverse = func(v *Valve) {
		tunnels := valveTunnels[v.ID]
		delete(valveTunnels, v.ID)

		for _, id := range tunnels {
			tunnel := valves[id]
			v.Tunnels = append(v.Tunnels, tunnel)
			traverse(tunnel)
		}

		sort.Slice(v.Tunnels, func(i, j int) bool {
			return v.Tunnels[i].FlowRate > v.Tunnels[j].FlowRate
		})
	}

	traverse(firstValve)
	return firstValve
}

// TunnelTo returns the tunnel to the given valve ID. If the tunnel does not
// exist, it will panic.
func (v *Valve) TunnelTo(id ValveID) *Valve {
	var dst *Valve
	for _, tunnel := range v.Tunnels {
		if tunnel.ID == id {
			dst = tunnel
			break
		}
	}
	aocutil.Assertf(dst != nil, "no tunnel from %q to %q", v.ID, id)
	return dst
}

// Moves is a list of moves.
type Moves []Move

// Move is a move from one valve to another or an opening of a valve. Each move
// takes 1 minute.
type Move interface {
	move()
	fmt.Stringer
}

// MoveToValve moves from the current valve through its tunnels to the next
// valve.
type MoveToValve struct {
	ID ValveID
}

// OpenValve opens a valve.
type OpenValve struct {
	ID ValveID
}

func (MoveToValve) move() {}
func (OpenValve) move()   {}

func (m MoveToValve) String() string { return fmt.Sprintf("move(%q)", m.ID) }
func (o OpenValve) String() string   { return fmt.Sprintf("open(%q)", o.ID) }

// MinutesLeft returns the number of minutes left in the time limit.
func (ms Moves) MinutesLeft() int {
	return 30 - len(ms)
}

// Pressure returns the total eventual pressure of the system given the moves.
func (ms Moves) Pressure(start *Valve) int {
	var total int
	var delta int

	for _, m := range ms {
		switch m := m.(type) {
		case MoveToValve:
			start = start.TunnelTo(m.ID)
		case OpenValve:
			delta += start.FlowRate
		}
		total += delta
	}

	return total
}

// Add returns a new set with the given value added.
func (ms Moves) Add(move Move) Moves {
	return append(ms[:len(ms):len(ms)], move)
}
