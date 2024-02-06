package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	expression1 := "1 + 2 * (3 + 4)"
	expression2 := "(5 / 21) + (3 + 4) * (9 - 10)"
	expression3 := "4 * (5 / 21) + (3 + 4) * (9 - 10)"

	result1, err := evaluateExpression(expression1)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result of", expression1, "is", result1)
	}

	result2, err := evaluateExpression(expression2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result of", expression2, "is", result2)
	}

	result3, err := evaluateExpression(expression3)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result of", expression3, "is", result3)
	}
}

func evaluateExpression(expression string) (float64, error) {
	postfixExpression, err := infixToPostfix(expression)
	fmt.Println("postfixExpression", postfixExpression)
	if err != nil {
		return 0, err
	}

	stack := []float64{}
	tokens := strings.Fields(postfixExpression)
	fmt.Println("tokens", tokens)
	for _, token := range tokens {
		if isNumber(token) {
			number, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, number)
		} else {
			if len(stack) < 2 {
				return 0, fmt.Errorf("Invalid expression")
			}

			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch token {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				stack = append(stack, a/b)
			}
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("Invalid expression")
	}

	return stack[0], nil
}

func infixToPostfix(expression string) (string, error) {
	output := ""
	operators := []rune{}

	for _, char := range expression {
		if unicode.IsSpace(char) {
			continue
		}
		if unicode.IsDigit(char) {
			output += string(char) // добавляем цифру к выходной строке
		} else if char == '(' {
			operators = append(operators, char)
		} else if char == ')' {
			for len(operators) > 0 && operators[len(operators)-1] != '(' {
				output += " " + string(operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = operators[:len(operators)-1] // удаляем '(' из стека
		} else if char == '+' || char == '-' || char == '*' || char == '/' {
			for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(char) {
				output += " " + string(operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, char)
			output += " " // добавляем пробел после оператора
		} else {
			return "", fmt.Errorf("Invalid character in expression")
		}
	}

	for len(operators) > 0 {
		output += " " + string(operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}

	return output, nil
}

func precedence(operator rune) int {
	switch operator {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	default:
		return 0
	}
}

func isNumber(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	return err == nil
}
