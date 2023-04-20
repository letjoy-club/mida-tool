package areacode

import "strings"

type AreaCode string

const codeLen = 6

func CodeFrom(code string) AreaCode {
	length := len(code)
	if length > codeLen {
		return AreaCode(code[:codeLen])
	} else if length < codeLen {
		return AreaCode(code + strings.Repeat("0", codeLen-length))
	}
	return AreaCode(code)
}

func (a AreaCode) String() string {
	return string(a)
}

func (a AreaCode) FuzzyQuery() string {
	if len(a) != codeLen {
		return ""
	}
	depth := a.Depth()
	// 已经是最小单位，没法进行模糊查找
	if depth == 2 {
		return a.String()
	}
	prefix := string(a[:depth*2])
	prefix = strings.ReplaceAll(prefix, "%", " ")
	return prefix + "%"
}

func (a AreaCode) Parent() AreaCode {
	if len(a) <= 1 {
		return AreaCode(string(a) + strings.Repeat("0", codeLen-len(a)))
	}
	depth := a.Depth()
	return AreaCode(string(a[:depth*2]) + strings.Repeat("0", codeLen-depth*2))
}

func (a AreaCode) checkDepth(i int) bool {
	if len(a) >= i*2 {
		return a[i*2] != '0' || a[i*2+1] != '0'
	}
	return false
}

func (a AreaCode) Depth() int {
	for i := 2; i >= 0; i-- {
		if a.checkDepth(i) {
			return i
		}
	}
	return 0
}
