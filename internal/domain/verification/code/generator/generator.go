package generator

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
)

const (
	codeLength    = 6
	hardnessRatio = 0.66
	digitsLimit   = 9 + 1
)

type CodeGenerator struct{}

func New() *CodeGenerator {
	return &CodeGenerator{}
}

func (g *CodeGenerator) Generate() string {
	uniqueNumbersAmount := int(math.Ceil(float64(codeLength) * hardnessRatio))

	numbersToUse := make([]int, 0, uniqueNumbersAmount)
	for range uniqueNumbersAmount {
		numbersToUse = append(numbersToUse, rand.Intn(digitsLimit))
	}

	var sb strings.Builder
	sb.Grow(codeLength)

	for range codeLength {
		chat := numbersToUse[rand.Intn(len(numbersToUse))]
		sb.WriteString(strconv.Itoa(chat))
	}

	return sb.String()
}
