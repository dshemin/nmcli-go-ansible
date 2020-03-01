//go:generate mockery -name Runner -inpkg -testonly

package nmcli

import (
	"context"
	"fmt"
	"os/exec"
)

// Runner interface for running commands
type Runner interface {
	Run(context.Context, string, ...string) ([]byte, error)
}

type realRunner struct{}

func (realRunner) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, name, args...).CombinedOutput()
}

var runner Runner = realRunner{}

func nmcli(ctx context.Context, args ...string) ([]byte, error) {
	out, err := runner.Run(ctx, "nmcli", args...)

	if err != nil {
		return nil, fmt.Errorf("cannot run `nmcli %#v` output %q: %w", args, string(out), err)
	}
	return out, nil
}
