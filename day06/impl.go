package day06

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var tokenRE = regexp.MustCompile(`\S+`)

func Round1(path string, verbose bool) (int64, error) {
	table, err := LoadTable(path)
	if err != nil {
		return 0, err
	}
	var result int64
	for col := range table[0] {
		var exp MathExpression
		for _, row := range table {
			cell := row[col]
			switch cell {
			case "+":
				exp.operator = sum
			case "*":
				exp.operator = multiply
			default:
				o, err := strconv.ParseInt(cell, 10, 64)
				if err != nil {
					return 0, err
				}
				exp.operands = append(exp.operands, o)
			}
		}
		result += exp.Evaluate()
	}
	return result, nil
}

func Round2(path string, verbose bool) (int64, error) {
	m, err := LoadTextMatrix(path)
	if err != nil {
		return 0, err
	}
	exprs, err := extractExpressionsFromMatrix(m, verbose)
	if err != nil {
		return 0, err
	}
	var sum int64
	for _, expr := range exprs {
		sum += expr.Evaluate()
	}
	return sum, nil
}

func LoadTable(path string) (table [][]string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := tokenRE.FindAllString(line, -1)
		if len(table) > 0 && len(table[0]) != len(row) {
			return nil, fmt.Errorf("line has unexpected length: %q", line)
		}
		table = append(table, row)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return table, nil
}

func LoadTextMatrix(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var matrix [][]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, "")
		if len(matrix) > 0 && len(matrix[0]) != len(row) {
			return nil, fmt.Errorf("row has unexpected length: %q", row)
		}
		matrix = append(matrix, row)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return matrix, nil
}

func extractExpressionsFromMatrix(m [][]string, verbose bool) ([]*MathExpression, error) {
	var exprs []*MathExpression
	operatorRow := len(m) - 1
	for col := 0; col < len(m[0]); col++ {
		// Scan to the next expression. It has the operator in the bottom row.
		var expr MathExpression
		switch op := m[operatorRow][col]; op {
		case " ":
			continue
		case "+":
			if verbose {
				fmt.Println(op)
			}
			expr.operator = sum
		case "*":
			if verbose {
				fmt.Println(op)
			}
			expr.operator = multiply
		default:
			return nil, fmt.Errorf("unexpected operator: %q", op)
		}
		// Parse the operands.
		for ; col < len(m[0]); col++ {
			var buf strings.Builder
			for row := 0; row < operatorRow; row++ {
				ch := m[row][col]
				if ch != " " {
					buf.WriteString(ch)
				}
			}
			if buf.Len() == 0 {
				break
			}
			operandStr := buf.String()
			if verbose {
				fmt.Println(operandStr)
			}
			operand, err := strconv.ParseInt(operandStr, 10, 64)
			if err != nil {
				return nil, err
			}
			expr.operands = append(expr.operands, operand)
		}
		exprs = append(exprs, &expr)
	}
	return exprs, nil
}

type MathExpression struct {
	operands []int64
	operator func([]int64) int64
}

func (m *MathExpression) Evaluate() int64 {
	return m.operator(m.operands)
}

func sum(operands []int64) int64 {
	var sum int64
	for _, o := range operands {
		sum += o
	}
	return sum
}

func multiply(operands []int64) int64 {
	var product int64 = 1
	for _, o := range operands {
		product *= o
	}
	return product
}
