package main

import (
	"fmt"
	"strconv"
	"unicode"
)

type Stack []string

// IsEmpty: check if stack is empty
func (st *Stack) IsEmpty() bool {
	return len(*st) == 0
}

// Push a new value onto the stack
func (st *Stack) Push(str string) {
	*st = append(*st, str) //Simply append the new value to the end of the stack
}

// Remove top element of stack. Return false if stack is empty.
func (st *Stack) Pop() bool {
	if st.IsEmpty() {
		return false
	} else {
		index := len(*st) - 1 // Get the index of top most element.
		*st = (*st)[:index]   // Remove it from the stack by slicing it off.
		return true
	}
}

// Return top element of stack. Return false if stack is empty.
func (st *Stack) Top() string {
	if st.IsEmpty() {
		return ""
	} else {
		index := len(*st) - 1   // Get the index of top most element.
		element := (*st)[index] // Index onto the slice and obtain the element.
		return element
	}
}

// Function to return precedence of operators
func prec(s string) int {
	if s == "^" {
		return 3
	} else if (s == "/") || (s == "*") {
		return 2
	} else if (s == "+") || (s == "-") {
		return 1
	} else {
		return -1
	}
}

func infixToPostfix(infix string) string {
	var sta Stack
	var postfix string
	var operand string
	for _, char := range infix {
		opchar := string(char)
		// if scanned character is operand, add it to output string
		if unicode.IsDigit(char) {
			operand += opchar
		} else {
			if operand != "" {
				postfix += operand + " "
				operand = "" // reset operand
			}
			if char == '(' {
				sta.Push(opchar)
			} else if char == ')' {
				for sta.Top() != "(" {
					postfix += sta.Top() + " "
					sta.Pop()
				}
				sta.Pop()
			} else {
				for !sta.IsEmpty() && prec(opchar) <= prec(sta.Top()) {
					postfix += sta.Top() + " "
					sta.Pop()
				}
				sta.Push(opchar)
			}
		}
	}
	if operand != "" {
		postfix += operand + " "
	}
	// Pop all the remaining elements from the stack
	for !sta.IsEmpty() {
		postfix += sta.Top() + " "
		sta.Pop()
	}
	return postfix
}

func evaluatePostfix(postfix string) int {
	sta := Stack{}
	var fullNum string
	for _, char := range postfix {

		str := string(char)
		if unicode.IsDigit(char) {
			fullNum += str
			continue
		} else if string(char) == " " {
			if fullNum != "" {
				sta.Push(fullNum)
				fullNum = ""
			}
		} else if str == "+" || str == "-" || str == "*" || str == "/" {
			op2, _ := strconv.Atoi(sta.Top())
			sta.Pop()
			op1, _ := strconv.Atoi(sta.Top())
			sta.Pop()
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
			sta.Push(strconv.Itoa(res))
		}
	}
	result, _ := strconv.Atoi(sta.Top())
	return result
}

func main() {
	infix := "(11-1)/2+(22+11)*2/2"
	postfix := infixToPostfix(infix)
	fmt.Println(postfix)

	result := evaluatePostfix(postfix)
	fmt.Println(result)
}
