package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"time"

	"github.com/open-policy-agent/opa/rego"
)

type TestCase struct {
	Name           string
	InputGrants    []string
	ExpectedResult []string
}

func main() {
	// 1. Read data.json to construct the Store
	dataBytes, err := os.ReadFile("policy/data.json")
	if err != nil {
		panic(err)
	}

	var dataMap map[string]any
	if err := json.Unmarshal(dataBytes, &dataMap); err != nil {
		panic(err)
	}

	ctx := context.Background()

	// Prepare Rego object
	r := rego.New(
		rego.Query("data.authz.abilities"),
		rego.Load([]string{"./policy/"}, nil),
	)

	// 3. Prepare for evaluation (Compilation)
	startPrepare := time.Now()
	pq, err := r.PrepareForEval(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Preparation/Compilation took: %v\n", time.Since(startPrepare))

	// Define Test Cases
	testCases := []TestCase{
		{
			Name:           "Empty Case",
			InputGrants:    []string{},
			ExpectedResult: []string{},
		},
		{
			Name:           "Single Ability Match (ability_1)",
			InputGrants:    []string{"grant_1", "grant_2"},
			ExpectedResult: []string{"ability_1"},
		},
		{
			Name:           "Multiple Abilities Match (ability_1, ability_2)",
			InputGrants:    []string{"grant_1", "grant_2", "grant_3", "grant_4", "grant_5"},
			ExpectedResult: []string{"ability_1", "ability_2"},
		},
		{
			Name:           "Partial Match (Missing grant_2 for ability_1)",
			InputGrants:    []string{"grant_1"},
			ExpectedResult: []string{},
		},
		{
			Name:           "Superset Match (ability_1 + unused)",
			InputGrants:    []string{"grant_1", "grant_2", "unused_grant"},
			ExpectedResult: []string{"ability_1"},
		},
		{
			Name:           "Generated Ability Match (gen_ability_0)",
			InputGrants:    []string{"gen_grant_0_qjh"},
			ExpectedResult: []string{"gen_ability_0"},
		},
	}

	// 4. Run Evaluation Benchmark
	numIterations := 10
	fmt.Printf("\nRunning benchmark with %d iterations per test case...\n", numIterations)

	for _, tc := range testCases {
		fmt.Printf("\nTest Case: %s\n", tc.Name)
		var totalDuration time.Duration
		var lastResult []string
		var evalErr error

		input := map[string]any{
			"grants": tc.InputGrants,
		}

		for i := 0; i < numIterations; i++ {
			start := time.Now()
			rs, err := pq.Eval(ctx, rego.EvalInput(input))
			duration := time.Since(start)
			totalDuration += duration

			if err != nil {
				evalErr = err
				break
			}

			if len(rs) > 0 && len(rs[0].Expressions) > 0 {
				val := rs[0].Expressions[0].Value
				// OPA returns []any, need to convert to []string for comparison
				if valSlice, ok := val.([]any); ok {
					var resultStr []string
					for _, v := range valSlice {
						resultStr = append(resultStr, v.(string))
					}
					lastResult = resultStr
				} else {
					lastResult = []string{}
				}
			} else {
				lastResult = []string{}
			}
		}

		if evalErr != nil {
			fmt.Printf("ERROR: %v\n", evalErr)
			continue
		}

		// Validation
		if validateResult(lastResult, tc.ExpectedResult) {
			avg := totalDuration / time.Duration(numIterations)
			fmt.Printf("PASS | Average Latency: %v\n", avg)
		} else {
			fmt.Printf("FAIL | Expected: %v, Got: %v\n", tc.ExpectedResult, lastResult)
		}
	}
}

func validateResult(actual, expected []string) bool {
	if len(actual) != len(expected) {
		return false
	}

	// Sort both slices to ensure order-independent comparison
	a := make([]string, len(actual))
	copy(a, actual)
	sort.Strings(a)

	e := make([]string, len(expected))
	copy(e, expected)
	sort.Strings(e)

	return reflect.DeepEqual(a, e)
}
