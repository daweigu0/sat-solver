package main

import "fmt"

func main() {
	nbVars, nbClauses, clauses := ReadCnf("D:/SAT/instances/基准算例/功能测试/sat-20.cnf")
	fmt.Println(nbVars, nbClauses, clauses)
	solver := NewSolver(nbVars, nbClauses, clauses)
	solver.Solve()
}
