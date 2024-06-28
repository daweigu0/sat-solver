package main

type Clause struct {
	literals []int
	size     int
}

func NewClause(literals []int) *Clause {
	return &Clause{
		literals: literals,
		size:     len(literals),
	}
}
func (clause *Clause) GetSize() int {
	return clause.size
}
func (clause *Clause) Attach(solver *Solver) {
	for _, lit := range clause.literals {
		list := solver.GetList(lit)
		list = append(list, clause)
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

type Solver struct {
	nbVars      int
	nbClauses   int
	clauses     [][]int
	assignStack []int
	negList     [][]*Clause
	posList     [][]*Clause
	model       []uint8
}

func NewSolver(nbVars, nbClauses int, clauses [][]int) *Solver {
	return &Solver{
		nbVars:      nbVars,
		nbClauses:   nbClauses,
		clauses:     clauses,
		assignStack: make([]int, 0, nbVars<<1),
		negList:     make([][]*Clause, nbVars+1),
		posList:     make([][]*Clause, nbVars+1),
		model:       make([]uint8, nbVars+1),
	}
}
func (solver *Solver) SetClauses(clauses [][]int) {
	solver.clauses = clauses
}
func Var(lit int) uint {
	if lit == 0 {
		panic("文字的值为0！")
	}
	if lit > 0 {
		return uint(lit)
	}
	return uint(-lit)
}
func (solver *Solver) GetList(lit int) []*Clause {
	if lit == 0 {
		panic("文字的值为0！")
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
	for _, clause := range solver.clauses {
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
	return false
}
