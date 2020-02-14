// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulse

import (
	"context"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/stretchr/testify/require"
)

const testPulse = insolar.PulseNumber(42)

func TestContextWith(t *testing.T) {
	t.Parallel()
	ctx := ContextWith(context.Background(), testPulse)
	require.Equal(t, testPulse, ctx.Value(contextKey{}))
}

func TestFromContext(t *testing.T) {
	t.Parallel()
	ctx := context.WithValue(context.Background(), contextKey{}, testPulse)
	result := FromContext(ctx)
	require.Equal(t, testPulse, result)
}
