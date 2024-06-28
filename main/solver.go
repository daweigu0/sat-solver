package main

import "math"

type Clause struct {
	literals []int
	size     int
	isTrue   bool
}

func NewClause(literals []int) *Clause {
	return &Clause{
		literals: literals,
		size:     len(literals),
		isTrue:   false,
	}
}
func (clause *Clause) GetSize() int {
	return clause.size
}
func (clause *Clause) SetTrue() {
	clause.isTrue = true
}
func (clause *Clause) SetFalse() {
	clause.isTrue = false
}
func (clause *Clause) IsTrue() bool {
	return clause.isTrue
}
func (clause *Clause) Attach(solver *Solver) {
	for _, lit := range clause.literals {
		list := solver.GetList(lit)
		list = append(list, clause)
	}
	solver.clauses = append(solver.clauses, clause)
}
func (clause *Clause) RemoveLit(lit int) int {
	if clause.size == 0 {
		panic("子句长度已经为0，没有文字可以删除！\n")
	}
	if clause.isTrue {
		return clause.size
	}
	v := Var(lit)
	for i := 0; i < clause.size; i++ {
		if v == Var(clause.literals[i]) {
			clause.size--
			clause.literals[i], clause.literals[clause.size] = clause.literals[clause.size], clause.literals[i]
			break
		}
	}
	return clause.size
}

func (clause *Clause) RecoverLit(lit int) {
	v := Var(lit)
	for i := clause.size; i < len(clause.literals); i++ {
		if v == Var(clause.literals[i]) {
			clause.literals[i], clause.literals[clause.size] = clause.literals[clause.size], clause.literals[i]
			clause.size++
			break
		}
	}
}

const (
	UNSAT = false
	SAT   = true
)
const (
	UNKNOWN uint8 = iota
	TRUE
	FALSE
)

type Pair struct {
	idx    int
	isFlip bool
}

func NewPair(idx int, isFlip bool) *Pair {
	return &Pair{
		idx:    idx,
		isFlip: isFlip,
	}
}

type Solver struct {
	nbVars        int
	nbClauses     int
	originClauses [][]int
	clauses       []*Clause
	assignStack   []int
	branchStack   []*Pair
	negList       [][]*Clause
	posList       [][]*Clause
	model         []uint8
	bucket        []int
}

func NewSolver(nbVars, nbClauses int, clauses [][]int) *Solver {
	return &Solver{
		nbVars:        nbVars,
		nbClauses:     nbClauses,
		originClauses: clauses,
		clauses:       make([]*Clause, 0, nbClauses),
		assignStack:   make([]int, 0, nbVars),
		branchStack:   make([]*Pair, 0, nbVars),
		negList:       make([][]*Clause, nbVars+1),
		posList:       make([][]*Clause, nbVars+1),
		model:         make([]uint8, nbVars+1),
		bucket:        make([]int, nbVars+1),
	}
}

func (solver *Solver) GetList(lit int) []*Clause {
	if lit == 0 {
		panic("文字的值为0！\n")
	}
	var list []*Clause
	if lit > 0 {
		list = solver.posList[Var(lit)]
	} else {
		list = solver.negList[Var(lit)]
	}
	return list
}
func (solver *Solver) Solve() bool {
	unit := []int{}
	for _, clause := range solver.originClauses {
		if len(clause) == 0 {
			return UNSAT
		} else if len(clause) == 1 {
			unit = append(unit, clause[0])
		} else {
			for i := 0; i < len(clause)-1; i++ {
				for j := 0; j < len(clause); j++ {
					if clause[i] == clause[j] {
						clause[j] = clause[len(clause)-1]
						clause = clause[:len(clause)-1]
					}
					if clause[i] == -clause[j] {
						goto next
					}
				}
			}
			NewClause(clause).Attach(solver)
		next:
		}
	}
	for _, lit := range unit {
		solver.Push(lit)
		_, ok := solver.UP()
		if ok == false {
			return UNSAT
		}
	}
	for {
		ok := solver.MakeBranch()
		if !ok {
			return SAT
		}
		end, conf := solver.UP()
		if conf {
			ok = solver.Back(end)
			if !ok {
				break
			}
		}
	}
	return UNSAT
}
func (solver *Solver) Push(lit int) {
	v := Var(lit)
	solver.model[v] |= Value(lit)
	solver.assignStack = append(solver.assignStack, lit)
}
func (solver *Solver) Pop(idx int) {
	lit := solver.assignStack[idx]
	v := Var(lit)
	solver.model[v] &= UNKNOWN
	solver.assignStack = solver.assignStack[:idx]
	clauses := solver.GetList(lit)
	for _, clause := range clauses {
		clause.SetFalse()
	}
	clauses = solver.GetList(-lit)
	for _, clause := range clauses {
		clause.RecoverLit(-lit)
	}
}
func (solver *Solver) UP() (int, bool) {
	for i := len(solver.assignStack) - 1; i < len(solver.assignStack); i++ {
		lit := solver.assignStack[i]
		//v := Var(lit)
		clauses := solver.GetList(-lit)
		conf := false
		for _, clause := range clauses {
			size := clause.RemoveLit(-lit)
			if size == 1 {
				solver.Push(clause.literals[0])
			}
			if size == 0 {
				conf = true
			}
		}
		clauses = solver.GetList(lit)
		for _, clause := range clauses {
			clause.SetTrue()
		}
		if conf {
			return i, false
		}
	}
	return len(solver.assignStack), true
}
func (solver *Solver) GetBackLocation() int {
	for i := len(solver.branchStack) - 1; i >= 0; i-- {
		if solver.branchStack[i].isFlip == false {
			solver.branchStack[i].isFlip = true
			solver.branchStack = solver.branchStack[:i+1]
			return solver.branchStack[i].idx
		}
	}
	return -1
}
func (solver *Solver) Back(start int) bool {
	end := solver.GetBackLocation()
	if end == -1 {
		return false
	}
	lit := solver.assignStack[end]
	for i := start; i >= end; i-- {
		solver.Pop(i)
	}
	solver.Push(-lit)
	return true
}

func (solver *Solver) Branch() (uint, bool) {
	minSize := math.MaxInt
	maxCnt := math.MinInt
	ok := false
	branch := uint(0)
	for _, clause := range solver.clauses {
		if clause.isTrue == false {
			minSize = min(minSize, clause.GetSize())
			ok = true
		}
	}
	if !ok {
		return branch, ok
	}
	for _, clause := range solver.clauses {
		if clause.isTrue == false && clause.GetSize() == minSize {
			for i := 0; i < clause.GetSize(); i++ {
				v := Var(clause.literals[i])
				solver.bucket[v]++
				if solver.bucket[v] > maxCnt {
					maxCnt = solver.bucket[v]
					branch = v
				}
			}
		}
	}
	if solver.model[branch] != UNKNOWN {
		panic("分支变元不为自由变元\n")
	}
	return branch, ok
}

func (solver *Solver) MakeBranch() bool {
	branch, ok := solver.Branch()
	if !ok {
		return ok
	}
	solver.Push(int(branch))
	solver.branchStack = append(solver.branchStack, NewPair(len(solver.assignStack)-1, false))
	return ok
}
