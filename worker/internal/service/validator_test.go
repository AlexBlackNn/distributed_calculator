package main

import "fmt"

type Stack struct {
	data []string
}

type Calculator struct {
	postfix string
	Stack
}

// IsEmpty: check if stack is empty
func (st *Stack) IsEmpty() bool {
	return len(st.data) == 0
}

// Push a new value onto the stack
func (st *Stack) Push(str string) {
	st.data = append(st.data, str) //Simply append the new value to the end of the stack
}

// Remove top element of stack. Return false if stack is empty.
func (st *Stack) Pop() bool {
	if st.IsEmpty() {
		return false
	} else {
		index := len(st.data) - 1   // Get the index of top most element.
		st.data = (st.data)[:index] // Remove it from the stack by slicing it off.
		return true
	}
}

// Return top element of stack. Return false if stack is empty.
func (st *Stack) Top() string {
	if st.IsEmpty() {
		return ""
	} else {
		index := len(st.data) - 1   // Get the index of top most element.
		element := (st.data)[index] // Index onto the slice and obtain the element.
		return element
	}
}

func IsValidOperation(op string) bool {
	validOperators := map[string]bool{
		"+": true,
		"-": true,
		"*": true,
		"/": true,
		"(": true,
		")": true,
	}

	return validOperators[op]
}

func VerifyExpression(expression string) bool {
	stack := Stack{}
	prevChar := ""
	for _, char := range expression {
		s := string(char)

		if prevChar != "" && isOperator(prevChar) && isOperator(s) {
			return false
		}

		if s == "(" {
			stack.Push(s)
		} else if s == ")" {
			if stack.IsEmpty() || stack.Top() != "(" {
				return false
			}
			stack.Pop()
		}

		prevChar = s
	}

	return stack.IsEmpty()
}

func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/"
}
func main() {
	fmt.Println(VerifyExpression("2+2"))           // true
	fmt.Println(VerifyExpression("3*5+1"))         // true
	fmt.Println(VerifyExpression("2*2+1*(2+1)"))   // true
	fmt.Println(VerifyExpression("2+2-1)"))        // false
	fmt.Println(VerifyExpression("(2+1"))          // false
	fmt.Println(VerifyExpression("2++1"))          // false
	fmt.Println(VerifyExpression("3/+1"))          // false
	fmt.Println(VerifyExpression("5+2("))          // false
	fmt.Println(VerifyExpression("5+++2("))        // false
	fmt.Println(VerifyExpression("5+2-3*(3+1)/2")) // true
	fmt.Println(VerifyExpression("5+2-3(3+1)/2"))  //false

	if !VerifyExpression("1+2") {
		panic("1")
	}
	if !VerifyExpression("1+2*2") {
		panic("2")
	}
	if !VerifyExpression("(1+2)*2") {
		panic("3")
	}
	if !VerifyExpression("(1+2)*2") {
		panic("4")
	}
	if !VerifyExpression("2+(1+2)*2/3*(1-2)") {
		panic("5")
	}
	if !VerifyExpression("6") {
		panic("7")
	}
	if !VerifyExpression("1") {
		panic("8")
	}
}
