package main

import (
	"fmt"
	"time"
)

func main() {

	nbVars, nbClauses, clauses := ReadCnf("D:/SAT/instances/SAT_training_instance_test/middle/4blocks.cnf")
	//fmt.Println(nbVars, nbClauses, clauses)
	solver := NewSolver(nbVars, nbClauses, clauses)
	start := time.Now()
	result := solver.Solve()
	duration := time.Since(start)
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
	fmt.Printf("time(sec)：%v", float64(duration)/1e9)
}
