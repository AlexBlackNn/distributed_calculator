package calculator

import (
	"fmt"
	"strconv"
	"unicode"
)

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

// Function to return precedence of operators
func prec(s string) int {
	if (s == "/") || (s == "*") {
		return 2
	} else if (s == "+") || (s == "-") {
		return 1
	} else {
		return -1
	}
}

func (st *Calculator) InfixToPostfix(infix string) {
	var operand string
	for _, char := range infix {
		opchar := string(char)
		// if scanned character is operand, add it to output string
		if unicode.IsDigit(char) {
			operand += opchar
		} else {
			if operand != "" {
				st.postfix += operand + " "
				operand = "" // reset operand
			}
			if char == '(' {
				st.Push(opchar)
			} else if char == ')' {
				for st.Top() != "(" {
					st.postfix += st.Top() + " "
					st.Pop()
				}
				st.Pop()
			} else {
				for !st.IsEmpty() && prec(opchar) <= prec(st.Top()) {
					st.postfix += st.Top() + " "
					st.Pop()
				}
				st.Push(opchar)
			}
		}
	}
	if operand != "" {
		st.postfix += operand + " "
	}
	// Pop all the remaining elements from the stack
	for !st.IsEmpty() {
		st.postfix += st.Top() + " "
		st.Pop()
	}
}

func (st *Calculator) EvaluatePostfix() int {
	var fullNum string
	for _, char := range st.postfix {
		str := string(char)
		if unicode.IsDigit(char) {
			fullNum += str
			continue
		} else if string(char) == " " {
			if fullNum != "" {
				st.Push(fullNum)
				fullNum = ""
			}
		} else if str == "+" || str == "-" || str == "*" || str == "/" {
			op2, _ := strconv.Atoi(st.Top())
			st.Pop()
			op1, _ := strconv.Atoi(st.Top())
			st.Pop()
			var res int
			switch str {
			case "+":
				res = op1 + op2
			case "-":
				res = op1 - op2
			case "*":
				res = op1 * op2
			case "/":
				res = op1 / op2
			}
			st.Push(strconv.Itoa(res))
		}
	}
	result, _ := strconv.Atoi(st.Top())
	return result
}

func (st *Calculator) Start(infix string) int {
	st.InfixToPostfix(infix)
	return st.EvaluatePostfix()
}

func New() *Calculator {
	stack := Stack{}
	return &Calculator{"", stack}
}

func main() {
	calculator := New()
	infix := "(11-1)/2+1*(22+11)*2/2"
	result := calculator.Start(infix)
	fmt.Println(result)
}
