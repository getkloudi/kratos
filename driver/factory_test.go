// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/ory/x/configx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/kratos/driver"
	"github.com/ory/kratos/driver/config"
)

func TestDriverNew(t *testing.T) {
	ctx := context.Background()
	r, err := driver.New(
		context.Background(),
		os.Stderr,
		driver.WithConfigOptions(
			configx.WithValue(config.ViperKeyDSN, config.DefaultSQLiteMemoryDSN),
			configx.SkipValidation(),
		))
	require.NoError(t, err)

	assert.EqualValues(t, config.DefaultSQLiteMemoryDSN, r.Config().DSN(ctx))
	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	t.Cleanup(cancel)
	require.NoError(t, r.Persister().Ping(pingCtx))

	assert.NotEqual(t, uuid.Nil.String(), r.Persister().NetworkID(context.Background()).String())

	n, err := r.Persister().DetermineNetwork(context.Background())
	require.NoError(t, err)
	assert.Equal(t, r.Persister().NetworkID(context.Background()), n.ID)
}
