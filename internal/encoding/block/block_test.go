// Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file

package block

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlock_Types(t *testing.T) {
	o, err := ioutil.ReadFile(smallFile)
	assert.NotEmpty(t, o)
	assert.NoError(t, err)

	b, err := FromOrc("test", o)
	assert.NoError(t, err)

	schema := b.Schema()
	assert.Equal(t, 5, len(schema))
	assert.Contains(t, schema, "boolean1")
	assert.Contains(t, schema, "double1")
	assert.Contains(t, schema, "int1")
	assert.Contains(t, schema, "long1")
	assert.Contains(t, schema, "string1")

	{
		result, err := b.Select([]string{"int1"})
		assert.NoError(t, err)
		assert.Equal(t, 2, result["int1"].IntegerData.Count())
	}

	{
		result, err := b.Select([]string{"string1"})
		assert.NoError(t, err)
		assert.Equal(t, 2, result["string1"].VarcharData.Count())
	}

	{
		result, err := b.Select([]string{"long1"})
		assert.NoError(t, err)
		assert.Equal(t, 2, result["long1"].BigintData.Count())
	}
}

// BenchmarkBlockRead/read-8         	     190	   6141327 ns/op	23054922 B/op	      11 allocs/op
func BenchmarkBlockRead(b *testing.B) {
	o, err := ioutil.ReadFile(testFile)
	noerror(err)

	blk, err := FromOrc("test", o)
	noerror(err)

	// 122MB uncompressed
	// 13MB snappy compressed
	buf, err := blk.Encode()
	noerror(err)

	columns := []string{"_col5"}
	b.Run("read", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, _ = Read(buf, columns)
		}
	})
}

// BenchmarkFrom/orc-8         	    7519	    140205 ns/op	  454838 B/op	    1143 allocs/op
// BenchmarkFrom/batch-8       	  138270	      8714 ns/op	    7651 B/op	     100 allocs/op
func BenchmarkFrom(b *testing.B) {
	orc, err := ioutil.ReadFile(smallFile)
	noerror(err)

	b.Run("orc", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, err = FromOrc("test", orc)
			noerror(err)
		}
	})

	b.Run("batch", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_, err = FromBatchBy(testBatch, "d")
			noerror(err)
		}
	})

}

func noerror(err error) {
	if err != nil {
		panic(err)
	}
}
