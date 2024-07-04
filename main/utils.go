package main

import (
	"fmt"
	"sort"
)

func Var(lit int) uint {
	if lit == 0 {
		panic("文字的值为0！\n")
	}
	if lit > 0 {
		return uint(lit)
	}
	return uint(-lit)
}
func Value(lit int) uint8 {
	if lit > 0 {
		return TRUE
	}
	return FALSE
}
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func Verify(solver *Solver) bool {
	cnt := 0
	for _, clause := range solver.originClauses {
		for _, lit := range clause {
			if Value(lit) == solver.model[Var(lit)] {
				cnt++
				break
			}
		}
	}
	if cnt == solver.nbClauses {
		return true
	}
	return false
}
func PrintClausesByVar(solver *Solver, lit int) {
	v := Var(lit)
	fmt.Printf("---------------------- %d -----------------------------\n", v)
	for _, clause := range solver.GetList(int(v)) {
		temp := make([]int, clause.size)
		copy(temp, clause.literals[:clause.size])
		sort.Ints(temp)
		for _, lit := range temp {
			fmt.Printf("%d ", lit)
		}
		fmt.Print(" | ")
		temp = make([]int, len(clause.literals)-clause.size)
		copy(temp, clause.literals[clause.size:])
		sort.Ints(temp)
		for _, lit := range temp {
			fmt.Printf("%d ", lit)
		}
		fmt.Printf("#%d #%v", clause.size, clause.isTrue)
		fmt.Println()
	}
	fmt.Println("---------------------------------------------------")
	for _, clause := range solver.GetList(-int(v)) {
		temp := make([]int, clause.size)
		copy(temp, clause.literals[:clause.size])
		sort.Ints(temp)
		for _, lit := range temp {
			fmt.Printf("%d ", lit)
		}
		fmt.Print(" | ")
		temp = make([]int, len(clause.literals)-clause.size)
		copy(temp, clause.literals[clause.size:])
		sort.Ints(temp)
		for _, lit := range temp {
			fmt.Printf("%d ", lit)
		}
		fmt.Printf("#%d #%v", clause.size, clause.isTrue)
		fmt.Println()
	}
}

func Check(solver *Solver) {
	for i := 1; i < len(solver.model); i++ {
		if solver.model[i] != UNKNOWN {
			flag := false
			for _, v := range solver.assignStack {
				if Var(v) == uint(i) {
					flag = true
					break
				}
			}
			if !flag {
				panic("变元被赋值，但不在assignStack中")
			}
		}
	}
	for _, v := range solver.assignStack {
		if Value(v) != solver.model[Var(v)] {
			panic("变元的赋值与assignStack中的文字不符")
		}
	}
}
