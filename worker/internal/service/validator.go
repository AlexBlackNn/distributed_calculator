package service

type Validator struct {
}

func (*Validator) VerifyExpression(expression string) bool {
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
