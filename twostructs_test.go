package twostructs_test

import (
	"testing"
	"time"

	"github.com/davewhat/twostructs"
	"github.com/stretchr/testify/assert"
)

type Epoch int64

type WireModel struct {
	Name            string
	OptionalAddress *string
	TimeInSeconds   Epoch
}

type InternalModel struct {
	FullName string    // field name changes from "Name" to "FullName"
	Address  string    // type changes from *string to string
	Time     time.Time // type changes from Epoch(int64) to time.Time
}

func TestReadme(t *testing.T) {
	wireObj := WireModel{Name: "David", OptionalAddress: nil, TimeInSeconds: 1553878048}
	mapper := twostructs.New()

	// EpochToTimeFunc maps Epoch(int64) to time.Time
	EpochToTimeFunc := func(e Epoch) time.Time {
		return time.Unix(int64(e), 0).UTC()
	}
	mapper.RegisterMappingFunction(EpochToTimeFunc)

	entityObj := InternalModel{}
	err := mapper.Struct(wireObj, &entityObj)

	assert.NoError(t, err)
	assert.Equal(t, "David", entityObj.FullName)
	assert.Equal(t, "2019-03-29 16:47:28 +0000 UTC", entityObj.Time.String())
}

func BenchmarkReadMe(b *testing.B) {
	wireObj := WireModel{Name: "David", OptionalAddress: nil, TimeInSeconds: 1553878048}
	mapper := twostructs.New()

	// EpochToTimeFunc maps Epoch(int64) to time.Time
	EpochToTimeFunc := func(e Epoch) time.Time {
		return time.Unix(int64(e), 0).UTC()
	}
	mapper.RegisterMappingFunction(EpochToTimeFunc)

	entityObj := InternalModel{}
	for i := 0; i < b.N; i++ {
		_ = mapper.Struct(wireObj, &entityObj)
	}
}
