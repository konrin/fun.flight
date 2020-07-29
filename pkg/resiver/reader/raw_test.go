package reader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAverage(t *testing.T) {
	reader := RawReader{}

	m, err := reader.ReadManual("*5D42438E9C7BD1;")
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, m.Data, "5D42438E9C7BD1")
	assert.Equal(t, m.MlatAt, int64(0))
	assert.Equal(t, m.RSSI, 0)
}

func TestAverage1(t *testing.T) {
	reader := RawReader{}

	m, err := reader.ReadWithMlat("@00097012B8868D4242DCEA447864DD3C0879780D;")
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, m.Data, "8D4242DCEA447864DD3C0879780D")
	assert.Equal(t, m.MlatAt, int64(40534980742))
	assert.Equal(t, m.RSSI, 0)
}
