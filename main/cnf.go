package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func ReadCnf(path string) (int, int, [][]int) {
	var (
		nbVars    int
		nbClauses int
		clauses   [][]int
	)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("打开cnf文件失败 ", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nbVars, nbClauses, clauses
			}
			fmt.Println("读取cnf时发生错误 ", err)

		}
		if line[0] == 'c' {
			continue
		} else if line[0] == 'p' {
			_, err := fmt.Sscanf(line, "p cnf %d %d", &nbVars, &nbClauses)
			if err != nil {
				fmt.Println("读取变元和子句数量时发生错误 ", err)
			}
		} else {
			ss := strings.Split(line[:len(line)-1], " ")
			clause := make([]int, 0, len(ss)-1)
			for _, v := range ss {
				if v != "0" {
					lit, err := strconv.Atoi(v)
					if err != nil {
						fmt.Println("读取子句时发生错误")
					}
					clause = append(clause, lit)
				}
			}
			clauses = append(clauses, clause)
		}
	}
}