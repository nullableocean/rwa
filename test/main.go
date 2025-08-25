package main

import "fmt"

func main() {
	str := "ab"

	b := repeatedSubstringPattern(str)
	fmt.Println(b)
}

func repeatedSubstringPattern(s string) bool {
	strlen := len(s)
	if strlen <= 1 {
		return false
	}

	sublen := 0
	for {
		sublen++

		if sublen > strlen/2 {
			break
		}

		if strlen%sublen == 0 {
			ok := true
			substr := s[:sublen]
			for i := 0; i < strlen/sublen; i++ {
				sind := i * sublen
				eind := sind + sublen
				if substr != s[sind:eind] {
					ok = false
					break
				}
			}
			if ok {
				return ok
			}
		}
	}

	return false
}
