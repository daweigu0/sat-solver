package main

import (
	"fmt"
)

func main() {
	nbVars, nbClauses, clauses := ReadCnf("D:/SAT/instances/满足算例/S/problem6-50.cnf")
	//fmt.Println(nbVars, nbClauses, clauses)
	solver := NewSolver(nbVars, nbClauses, clauses)
	result := solver.Solve()
	if result == SAT {
		fmt.Println("SAT")
		if Verify(solver) {
			fmt.Println("verify is success!")
		} else {
			fmt.Println("verify is failure!")
		}
	} else {
		fmt.Println("UNSAT")
	}
}
