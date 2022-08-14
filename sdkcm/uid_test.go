package sdkcm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUID(t *testing.T) {
	for _, c := range []struct {
		uid    UID
		expect string
	}{
		{uid: NewUID(1, 1, 1), expect: "268697601"},
		{uid: NewUID(2, 1, 1), expect: "537133057"},
		{uid: NewUID(3, 1, 1), expect: "805568513"},
		{uid: NewUID(4, 1, 1), expect: "1074003969"},
		{uid: NewUID(5, 1, 1), expect: "1342439425"},
		{uid: NewUID(6, 1, 1), expect: "1610874881"},
		{uid: NewUID(7, 1, 1), expect: "1879310337"},
		{uid: NewUID(8, 1, 1), expect: "2147745793"},
		{uid: NewUID(9, 1, 1), expect: "2416181249"},
		{uid: NewUID(10, 1, 1), expect: "2684616705"},
	} {
		actual := c.uid.String()
		assert.Equal(t, c.expect, actual, "should be equal")
	}
}

func TestDecomposeUID(t *testing.T) {
	for _, c := range []struct {
		uid    string
		expect UID
	}{
		{expect: NewUID(1, 1, 1), uid: "268697601"},
		{expect: NewUID(2, 1, 1), uid: "537133057"},
		{expect: NewUID(3, 1, 1), uid: "805568513"},
		{expect: NewUID(4, 1, 1), uid: "1074003969"},
	} {
		actual, err := DecomposeUID(c.uid)
		assert.Nil(t, err, "must be nil")
		assert.Equal(t, c.expect.GetLocalID(), actual.GetLocalID(), "should be equal")
		assert.Equal(t, c.expect.GetObjectType(), actual.GetObjectType(), "should be equal")
		assert.Equal(t, c.expect.GetShardID(), actual.GetShardID(), "should be equal")
	}

	wrongFormat := "abc"
	_, err := DecomposeUID(wrongFormat)
	assert.NotNil(t, err, "should be an error")
}
