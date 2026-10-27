package main

type P *int

var gp P = new(int)

func f(p *int) P { return p }

func main() { println(f(nil) == nil, gp != nil) }
