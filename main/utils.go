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
