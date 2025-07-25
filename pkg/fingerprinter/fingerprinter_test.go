package fingerprinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSHA1FromStrings(t *testing.T) {
	actual := SHA1FromStrings("some", "random", "text")
	assert.Equal(t, "cb7d13eaca9a402d9c45035dd4766562671a73d2", actual)
}
