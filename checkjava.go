package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func check(a string) (s1 []string) {
	var cmd *exec.Cmd
	cmd = exec.Command("ls", a)
	test, _ := cmd.Output()
	b := string(test)
	s1 = strings.Fields(b)
	fmt.Printf("%T,%v\n", s1, s1)
	return s1
}

func check2() (s2 []string) {
	var cmd2 *exec.Cmd
	cmd2 = exec.Command("/bin/bash", "/data/go/test.sh")
	test1, _ := cmd2.Output()
	y := string(test1)
	s2 = strings.Fields(y)
	fmt.Printf("%T,%v\n", s2, s2)
	return s2
}

func intersect(s1, s2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range s1 {
		m[v]++
	}

	for _, v := range s2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

func difference(s1, s2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	inter := intersect(s1, s2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range s1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}

func restart(s4 []string) {
	var cmd3 *exec.Cmd
	j := "/data/domains/"
	for _, i := range s4 {
		i = j + i + ".sh" + "  restart"
		cmd3 = exec.Command("/bin/bash", i)
		err := cmd3.Run()
		if err != nil {
			fmt.Println("Execute Command failed:" + err.Error())
			return
		}
	}
}

func main() {
	a := check("/data/domains")
	b := check2()
	c := difference(a, b)
	fmt.Printf("%T,%v\n", c, c)
	restart(c)	
}
