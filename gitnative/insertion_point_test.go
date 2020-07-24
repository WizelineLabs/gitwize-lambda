package gitnative

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetInsertionPoint(t *testing.T) {
	count, err := GetInsertionPoint("./", "661e0a91146aff0f53f248b49c8722c30355beaf")
	assert.Equal(t, 6, count)
	assert.Nil(t, err)
	count, err = GetInsertionPoint("./", "e25da7a76548ab56172e0f985835fb078cc920fd")
	assert.Equal(t, 0, count)
	assert.NotNil(t, err) // grep return exit code 1 if not found
	count, err = GetInsertionPoint("./", "invalidhash")
	assert.NotNil(t, err)
}
