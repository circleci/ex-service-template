package api

import (
	"net/http"
	"testing"

	"github.com/circleci/ex/testing/testcontext"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestAPI_ping(t *testing.T) {
	ctx := testcontext.Background()
	fix := startAPI(ctx, t)

	t.Run("Ping -> pong", func(t *testing.T) {
		m := make(map[string]interface{})
		status := fix.Get(t, "/api/ping", &m)
		assert.Check(t, cmp.Equal(status, http.StatusOK))
		assert.Check(t, cmp.DeepEqual(map[string]interface{}{
			"message": "pong",
		}, m))
	})
}
