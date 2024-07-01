package main

import (
	"fmt"
	"math"
)

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
func (clause *Clause) SetTrue(lit int) {
	if clause.isTrue {
		return
	}
	ok := false
	for i := 0; i < clause.size; i++ {
		if clause.literals[i] == lit {
			clause.literals[0], clause.literals[i] = clause.literals[i], clause.literals[0]
			ok = true
			break
		}
	}
	if !ok {
		panic("赋值为真的文字不在子句中")
	}
	clause.isTrue = true
}
func (clause *Clause) SetFalse(lit int) {
	if lit != clause.literals[0] {
		return
	}
	clause.isTrue = false
}
func (clause *Clause) IsTrue() bool {
	return clause.isTrue
}
func (clause *Clause) Attach(solver *Solver) {
	for _, lit := range clause.literals {
		if lit > 0 {
			solver.posList[Var(lit)] = append(solver.posList[Var(lit)], clause)
		} else {
			solver.negList[Var(lit)] = append(solver.negList[Var(lit)], clause)
		}
	}
	solver.clauses = append(solver.clauses, clause)
}

func (clause *Clause) RemoveLit(lit int) (int, bool) {
	if clause.size == 0 {
		panic("子句长度已经为0，没有文字可以删除！\n")
	}
	if clause.isTrue {
		return clause.size, false
	}
	v := Var(lit)
	ok := false
	for i := 0; i < clause.size; i++ {
		if v == Var(clause.literals[i]) {
			clause.size--
			clause.literals[i], clause.literals[clause.size] = clause.literals[clause.size], clause.literals[i]
			ok = true
			break
		}
	}
	if !ok {
		panic(fmt.Sprintf("子句中没有可以删除的文字%d", lit))
	}
	return clause.size, true
}

func (clause *Clause) RecoverLit(lit int) {
	if clause.isTrue {
		return
	}
	for i := clause.size; i < len(clause.literals); i++ {
		if lit == clause.literals[i] {
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
	solver := &Solver{
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
	for i := range solver.negList {
		solver.negList[i] = []*Clause{}
	}
	for i := range solver.posList {
		solver.posList[i] = []*Clause{}
	}
	return solver
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
	var unit []int
	for _, clause := range solver.originClauses {
		if len(clause) == 0 {
			return UNSAT
		} else if len(clause) == 1 {
			unit = append(unit, clause[0])
		} else {
			for i := 0; i < len(clause)-1; i++ {
				for j := i + 1; j < len(clause); j++ {
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
		if ok == true {
			return UNSAT
		}
	}
	for {
		end, conf := solver.UP()
		if conf {
			ok := solver.Back(end)
			if !ok {
				break
			}
			continue
		}
		ok := solver.MakeBranch()
		if !ok {
			return SAT
		}
	}
	return UNSAT
}
func (solver *Solver) Push(lit int) {
	v := Var(lit)
	if solver.model[v] != UNKNOWN {
		return
	}
	solver.model[v] |= Value(lit)
	solver.assignStack = append(solver.assignStack, lit)
}
func (solver *Solver) Pop(tail int) {
	lit := solver.assignStack[tail]
	v := Var(lit)
	solver.model[v] &= UNKNOWN
	solver.assignStack = solver.assignStack[:tail]
	clauses := solver.GetList(lit)
	//fmt.Println("========= 回溯前 =========")
	//PrintClausesByVar(solver, lit)
	for _, clause := range clauses {
		for i := 0; i < clause.size; i++ {
			clause.SetFalse(lit)
		}
	}
	clauses = solver.GetList(-lit)
	for _, clause := range clauses {
		clause.RecoverLit(-lit)
	}
	//fmt.Println("========= 回溯后 =========")
	//PrintClausesByVar(solver, lit)
}

// UP
//
//	@Description: 单子句传播
//	@receiver solver
//	@return int 最后传播的单子句在assignStack中的下标
//	@return bool 有冲突时为真
func (solver *Solver) UP() (int, bool) {
	for i := len(solver.assignStack) - 1; i < len(solver.assignStack) && i >= 0; i++ {
		//Check(solver)
		lit := solver.assignStack[i]
		//v := Var(lit)
		clauses := solver.GetList(-lit)
		conf := false
		for _, clause := range clauses {
			size, ok := clause.RemoveLit(-lit)
			if ok && size == 1 {
				solver.Push(clause.literals[0])
			}
			if ok && size == 0 {
				conf = true
			}
		}
		clauses = solver.GetList(lit)
		for _, clause := range clauses {
			clause.SetTrue(lit)
		}
		if conf {
			return i, conf
		}
	}
	return len(solver.assignStack) - 1, false
}
func (solver *Solver) GetBackHead() (int, bool) {
	for i := len(solver.branchStack) - 1; i >= 0; i-- {
		if solver.branchStack[i].isFlip == false {
			solver.branchStack[i].isFlip = true
			solver.branchStack = solver.branchStack[:i+1]
			return solver.branchStack[i].idx, true
		}
	}
	return -1, false
}
func (solver *Solver) Back(tail int) bool {
	head, ok := solver.GetBackHead()
	if !ok {
		return false
	}
	for i := len(solver.assignStack) - 1; i > tail; i-- {
		lit := solver.assignStack[i]
		v := Var(lit)
		solver.model[v] &= UNKNOWN
	}
	lit := solver.assignStack[head]
	for i := tail; i >= head; i-- {
		solver.Pop(i)
	}
	//PrintClausesByVar(solver, lit)
	solver.Push(-lit)
	//Check(solver)
	return true
}

// Branch
//
//	@Description: moms策略挑选分支变元
//	@receiver solver
//	@return uint
//	@return bool 如果全部子句已经满足，则为假
func (solver *Solver) Branch() (uint, bool) {
	minSize := math.MaxInt
	maxCnt := math.MinInt
	ok := false
	branch := uint(0)
	for _, clause := range solver.clauses {
		if !clause.isTrue {
			minSize = min(minSize, clause.GetSize())
			ok = true
		}
	}
	if !ok {
		return branch, ok
	}
	for _, clause := range solver.clauses {
		if !clause.isTrue && clause.GetSize() == minSize {
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
	if branch == 0 {
		panic("分支变元为0!\n")
	}
	if solver.model[branch] != UNKNOWN {
		panic("分支变元不为自由变元\n")
	}
	for i := range solver.bucket {
		solver.bucket[i] = 0
	}
	return branch, ok
}

// MakeBranch
//
//	@Description: 分支
//	@receiver solver
//	@return bool 成功分支返回true；分支失败返回false，表明算例已经满足
func (solver *Solver) MakeBranch() bool {
	branch, ok := solver.Branch()
	if !ok {
		return ok
	}
	solver.Push(int(branch))
	solver.branchStack = append(solver.branchStack, NewPair(len(solver.assignStack)-1, false))
	return ok
}
