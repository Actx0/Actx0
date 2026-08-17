// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdValid(t *testing.T) {
	assert.True(t, Id("8faf1e0b-bcf9-4c51-b54f-608949473754").Valid())
	assert.True(t, BotUserId.Valid())
	assert.False(t, Id("").Valid())
	assert.False(t, Id("not-a-uuid").Valid())
}

func TestIdValue(t *testing.T) {
	value, err := Id("8faf1e0b-bcf9-4c51-b54f-608949473754").Value()
	assert.NoError(t, err)
	assert.Equal(t, "8faf1e0b-bcf9-4c51-b54f-608949473754", value)

	value, err = Id("").Value()
	assert.NoError(t, err)
	assert.Nil(t, value)
}
