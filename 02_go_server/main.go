package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/open-policy-agent/opa/rego"
)

func main() {
	ctx := context.Background()

	// Prepare Rego query
	query, err := rego.New(
		rego.Query("data.authz.allow"),
		rego.Load([]string{"policy/"}, func(abspath string, info os.FileInfo, depth int) bool {
			return !strings.HasSuffix(abspath, "_test.rego")
		}),
	).PrepareForEval(ctx)

	if err != nil {
		log.Fatalf("Failed to prepare rego query: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		user := r.Header.Get("X-User")
		input := map[string]any{
			"method": r.Method,
			"path":   r.URL.Path,
			"user":   user,
		}

		results, err := query.Eval(ctx, rego.EvalInput(input))
		if err != nil {
			http.Error(w, "Evaluation error", http.StatusInternalServerError)
			return
		}

		if len(results) > 0 && results[0].Expressions[0].Value == true {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Access Allowed")
		} else {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, "Access Denied")
		}
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
