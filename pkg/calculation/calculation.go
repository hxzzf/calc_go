package calculation

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

func Calc(expression string) (float64, error) {
	if strings.TrimSpace(expression) == "" {
		return 0, errors.New("empty expression")
	}

	tokens, err := tokenize(expression)
	if err != nil {
		return 0, err
	}

	if isOperator(tokens[0]) || isOperator(tokens[len(tokens)-1]) {
		return 0, errors.New("invalid expression")
	}

	for i := 0; i < len(tokens)-1; i++ {
		if isOperator(tokens[i]) && isOperator(tokens[i+1]) {
			return 0, errors.New("invalid expression: consecutive operators")
		}
	}

	return evaluate(tokens)
}

func isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/"
}

func tokenize(expression string) ([]string, error) {
	tokens := make([]string, 0, len(expression))
	var currentToken strings.Builder
	var hasDecimalPoint bool

	for _, char := range expression {
		if unicode.IsSpace(char) {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
				hasDecimalPoint = false
			}
			continue
		}

		if unicode.IsDigit(char) {
			currentToken.WriteRune(char)
		} else if char == '.' {
			if hasDecimalPoint {
				return nil, fmt.Errorf("invalid token")
			}
			hasDecimalPoint = true
			currentToken.WriteRune(char)
		} else if unicode.IsLetter(char) {
			return nil, fmt.Errorf("invalid token")
		} else {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
				hasDecimalPoint = false
			}
			tokens = append(tokens, string(char))
		}
	}

	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens, nil
}

func evaluate(tokens []string) (float64, error) {
	numbers := make([]float64, 0, len(tokens)/2+1)
	operators := make([]string, 0, len(tokens)/2)

	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(token) {
				if err := applyOperation(&numbers, &operators); err != nil {
					return 0, err
				}
			}
			operators = append(operators, token)
		case "(":
			operators = append(operators, token)
		case ")":
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				if err := applyOperation(&numbers, &operators); err != nil {
					return 0, err
				}
			}
			if len(operators) == 0 || operators[len(operators)-1] != "(" {
				return 0, errors.New("mismatched parentheses")
			}
			operators = operators[:len(operators)-1]
		default:
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid token: %s", token)
			}
			numbers = append(numbers, num)
		}
	}

	for len(operators) > 0 {
		if operators[len(operators)-1] == "(" {
			return 0, errors.New("mismatched parentheses")
		}
		if err := applyOperation(&numbers, &operators); err != nil {
			return 0, err
		}
	}

	if len(numbers) != 1 {
		return 0, errors.New("invalid expression")
	}

	return numbers[0], nil
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

func applyOperation(numbers *[]float64, operators *[]string) error {
	if len(*numbers) < 2 {
		return errors.New("not enough operands")
	}
	b, a := (*numbers)[len(*numbers)-1], (*numbers)[len(*numbers)-2]
	*numbers = (*numbers)[:len(*numbers)-2]

	op := (*operators)[len(*operators)-1]
	*operators = (*operators)[:len(*operators)-1]

	var result float64
	switch op {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
		if result > 1.797693e+308 || result < -1.797693e+308 || result != result {
			return errors.New("numeric overflow")
		}
	case "/":
		if b == 0 {
			return errors.New("division by zero")
		}
		result = a / b
	}

	if math.IsInf(result, 0) || math.IsNaN(result) {
		return errors.New("numeric overflow")
	}

	*numbers = append(*numbers, result)
	return nil
}
