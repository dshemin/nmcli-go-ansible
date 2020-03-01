package nmcli

import (
	"context"
	"testing"

	. "github.com/stretchr/testify/require"
)

func TestNmcli(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		mockRunner := new(MockRunner)
		mockRunner.
			On("Run", ctx, "nmcli", "foo", "bar").
			Once().
			Return([]byte("success"), nil)

		runner = mockRunner

		out, err := nmcli(ctx, "foo", "bar")
		NoError(t, err)
		Equal(t, "success", string(out))
	})
}
