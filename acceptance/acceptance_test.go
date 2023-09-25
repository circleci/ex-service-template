package acceptance

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/circleci/ex/testing/compiler"
)

var (
	apiTestBinary = os.Getenv("API_TEST_BINARY")
)

func TestMain(m *testing.M) {
	status, err := runTests(m)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(status)
}

//nolint:funlen
func runTests(m *testing.M) (int, error) {
	ctx := context.Background()

	p := compiler.NewParallel(2)
	defer p.Cleanup()

	p.Add(compiler.Work{
		Result: &apiTestBinary,
		Name:   "api",
		Target: "..",
		Source: "github.com/circleci/ex-service-template/cmd/api",
	})

	err := p.Run(ctx)
	if err != nil {
		return 0, err
	}

	fmt.Printf("Using 'api' test binary: %q\n", apiTestBinary)

	return m.Run(), nil
}
