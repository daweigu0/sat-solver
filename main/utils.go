package main

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
