package calculator

import (
	"strconv"
	"time"
	"unicode"
	transport "worker/internal/transport"
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

func (st *Calculator) InfixToPostfix(message transport.RequestMessage) {
	st.Stack = Stack{}
	st.postfix = ""
	var operand string
	for _, char := range message.Operation {
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

func (st *Calculator) EvaluatePostfix(message transport.RequestMessage) int {
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
				time.Sleep(time.Duration(message.MessageExectutionTime.PlusOperationExecutionTime) * time.Millisecond)
			case "-":
				res = op1 - op2
				time.Sleep(time.Duration(message.MessageExectutionTime.MinusOperationExecutionTime) * time.Millisecond)
			case "*":
				res = op1 * op2
				time.Sleep(time.Duration(message.MessageExectutionTime.MultiplicationOperationExecutionTime) * time.Millisecond)
			case "/":
				res = op1 / op2
				time.Sleep(time.Duration(message.MessageExectutionTime.DivisionOperationExecutionTime) * time.Millisecond)
			}
			st.Push(strconv.Itoa(res))
		}
	}
	result, _ := strconv.Atoi(st.Top())
	return result
}

func (st *Calculator) Start(requestMessage transport.RequestMessage) int {
	st.InfixToPostfix(requestMessage)
	return st.EvaluatePostfix(requestMessage)
}

func New() *Calculator {
	stack := Stack{}
	return &Calculator{"", stack}
}
