package main

import (
	"log"

	"github.com/diamondburned/aoc-2022/aocutil"
)

func main() {
	var nums []int

	input := aocutil.InputString()
	lines := aocutil.SplitLines(input)
	for _, line := range lines {
		nums = append(nums, aocutil.Atoi[int](line))
	}

	part1(nums)
}

func part1(input []int) {
	type numPair struct{ i, v int }
	nums := make([]numPair, len(input))
	for i, v := range input {
		nums[i] = numPair{i, v}
	}

	mov := func(i int) {
		if i == 0 {
			return
		}

		// Move number at i backwards or forwards i times.
		// We'll do this the smart way by calculating the new index using the
		// modulo operator.
		var j int
		if i > 0 {
			// i is positive, so we're moving forwards. We can just do basic
			// modulo for this one.
			j = i % len(nums)
		} else {
			// i is negative, so we're moving backwards. We need to do some
			// extra work to make it work.
			j = len(nums) - (i % len(nums))
		}

		// Shift the numbers.
		nums = append(nums[j:], nums[:j]...)
	}

	// Do the thing.
	for i, n := range input {
		mov(nums[i])
		log.Println("moving by", nums[i])
	}
}
