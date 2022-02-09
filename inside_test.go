package think

import (
	"context"
	"testing"

	"github.com/hiank/think/db"
	"gotest.tools/v3/assert"
)

func TestOptions(t *testing.T) {
	dopt := makeOptions()
	assert.Equal(t, len(dopt.mdialer), 0)
	assert.Equal(t, len(dopt.mdopts), 0)
	assert.Equal(t, dopt.natsUrl, "")
	assert.Equal(t, dopt.todo, context.TODO())

	dopt = makeOptions(
		WithMemKvDialer(FuncKvDialer(func(c context.Context, do ...db.DialOption) (db.KvDB, error) { return nil, nil })),
		WithDiskvDialer(FuncKvDialer(func(c context.Context, do ...db.DialOption) (db.KvDB, error) { return nil, nil })),
		WithNatsUrl("natsurl"),
		WithTODO(context.Background()),
	)
	assert.Equal(t, len(dopt.mdialer), 2)
	assert.Equal(t, len(dopt.mdopts), 2)
	assert.Equal(t, dopt.natsUrl, "natsurl")
	assert.Equal(t, dopt.todo, context.Background())
}
