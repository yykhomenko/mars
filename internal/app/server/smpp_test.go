package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {

	params := ParseTLVStatus("id:186960018316110 sub:001 dlvrd:001 submit date:1908190140 done date:1908190140 stat:DELIVRD err:000 text:Test")

	assert.Equal(t, "186960018316110", params["id"])
}
