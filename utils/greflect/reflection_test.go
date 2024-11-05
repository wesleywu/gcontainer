package greflect

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type demoStruct struct {
	name  string
	value int
}

func Test_Instance(t *testing.T) {
	var (
		s1  *demoStruct
		s2  demoStruct
		s3  []*demoStruct
		s4  []demoStruct
		s5  map[string]*demoStruct
		s6  map[string]demoStruct
		s7  *string
		s8  string
		s9  *bool
		s10 bool
		s11 chan *demoStruct
		s12 chan demoStruct
	)
	s1 = Instance[*demoStruct]()
	assert.NotNil(t, s1)
	assert.Equal(t, "", s1.name)
	assert.Equal(t, 0, s1.value)

	s2 = Instance[demoStruct]()
	assert.NotNil(t, s2)
	assert.Equal(t, "", s2.name)
	assert.Equal(t, 0, s2.value)

	s3 = Instance[[]*demoStruct]()
	assert.NotNil(t, s3)
	assert.Equal(t, 0, len(s3))

	s4 = Instance[[]demoStruct]()
	assert.NotNil(t, s4)
	assert.Equal(t, 0, len(s4))

	s5 = Instance[map[string]*demoStruct]()
	assert.NotNil(t, s5)
	assert.Equal(t, 0, len(s5))

	s6 = Instance[map[string]demoStruct]()
	assert.NotNil(t, s6)
	assert.Equal(t, 0, len(s6))

	s7 = Instance[*string]()
	assert.NotNil(t, s7)
	assert.Equal(t, "", *s7)

	s8 = Instance[string]()
	assert.NotNil(t, s8)
	assert.Equal(t, "", s8)

	s9 = Instance[*bool]()
	assert.NotNil(t, s9)
	assert.Equal(t, false, *s9)

	s10 = Instance[bool]()
	assert.NotNil(t, s10)
	assert.Equal(t, false, s10)

	s11 = Instance[chan *demoStruct]()
	assert.NotNil(t, s11)
	go func() {
		s11 <- &demoStruct{
			name:  "test",
			value: 123,
		}
	}()
	d1 := <-s11
	assert.NotNil(t, d1)
	assert.Equal(t, "test", d1.name)
	assert.Equal(t, 123, d1.value)

	s12 = Instance[chan demoStruct]()
	assert.NotNil(t, s12)
	go func() {
		s12 <- demoStruct{
			name:  "test",
			value: 123,
		}
	}()
	d2 := <-s12
	assert.NotNil(t, d2)
	assert.Equal(t, "test", d2.name)
	assert.Equal(t, 123, d2.value)
}

func Test_InstanceOf(t *testing.T) {
	var (
		s1  *demoStruct
		s2  demoStruct
		s3  []*demoStruct
		s4  []demoStruct
		s5  map[string]*demoStruct
		s6  map[string]demoStruct
		s7  *string
		s8  string
		s9  *bool
		s10 bool
		s11 chan *demoStruct
		s12 chan demoStruct
	)
	s1 = InstanceOf(reflect.TypeOf(s1)).(*demoStruct)
	assert.NotNil(t, s1)
	assert.Equal(t, "", s1.name)
	assert.Equal(t, 0, s1.value)

	s2 = InstanceOf(reflect.TypeOf(s2)).(demoStruct)
	assert.NotNil(t, s2)
	assert.Equal(t, "", s2.name)
	assert.Equal(t, 0, s2.value)

	s3 = InstanceOf(reflect.TypeOf(s3)).([]*demoStruct)
	assert.NotNil(t, s3)
	assert.Equal(t, 0, len(s3))

	s4 = InstanceOf(reflect.TypeOf(s4)).([]demoStruct)
	assert.NotNil(t, s4)
	assert.Equal(t, 0, len(s4))

	s5 = InstanceOf(reflect.TypeOf(s5)).(map[string]*demoStruct)
	assert.NotNil(t, s5)
	assert.Equal(t, 0, len(s5))

	s6 = InstanceOf(reflect.TypeOf(s6)).(map[string]demoStruct)
	assert.NotNil(t, s6)
	assert.Equal(t, 0, len(s6))

	s7 = InstanceOf(reflect.TypeOf(s7)).(*string)
	assert.NotNil(t, s7)
	assert.Equal(t, "", *s7)

	s8 = InstanceOf(reflect.TypeOf(s8)).(string)
	assert.NotNil(t, s8)
	assert.Equal(t, "", s8)

	s9 = InstanceOf(reflect.TypeOf(s9)).(*bool)
	assert.NotNil(t, s9)
	assert.Equal(t, false, *s9)

	s10 = InstanceOf(reflect.TypeOf(s10)).(bool)
	assert.NotNil(t, s10)
	assert.Equal(t, false, s10)

	s11 = InstanceOf(reflect.TypeOf(s11)).(chan *demoStruct)
	assert.NotNil(t, s11)
	go func() {
		s11 <- &demoStruct{
			name:  "test",
			value: 123,
		}
	}()
	d1 := <-s11
	assert.NotNil(t, d1)
	assert.Equal(t, "test", d1.name)
	assert.Equal(t, 123, d1.value)

	s12 = InstanceOf(reflect.TypeOf(s12)).(chan demoStruct)
	assert.NotNil(t, s12)
	go func() {
		s12 <- demoStruct{
			name:  "test",
			value: 123,
		}
	}()
	d2 := <-s12
	assert.NotNil(t, d2)
	assert.Equal(t, "test", d2.name)
	assert.Equal(t, 123, d2.value)
}

func Test_NilOrEmpty(t *testing.T) {
	var (
		s1  *demoStruct
		s2  demoStruct
		s3  []*demoStruct
		s4  []demoStruct
		s5  map[string]*demoStruct
		s6  map[string]demoStruct
		s7  *string
		s8  string
		s9  *bool
		s10 bool
		s11 chan *demoStruct
		s12 chan demoStruct
	)
	s1 = NilOrEmpty[*demoStruct]()
	assert.Nil(t, s1)

	s2 = NilOrEmpty[demoStruct]()
	assert.NotNil(t, s2)
	assert.Equal(t, "", s2.name)
	assert.Equal(t, 0, s2.value)

	s3 = NilOrEmpty[[]*demoStruct]()
	assert.Nil(t, s3)

	s4 = NilOrEmpty[[]demoStruct]()
	assert.Nil(t, s4)

	s5 = NilOrEmpty[map[string]*demoStruct]()
	assert.Nil(t, s5)

	s6 = NilOrEmpty[map[string]demoStruct]()
	assert.Nil(t, s6)

	s7 = NilOrEmpty[*string]()
	assert.Nil(t, s7)

	s8 = NilOrEmpty[string]()
	assert.NotNil(t, s8)
	assert.Equal(t, "", s8)

	s9 = NilOrEmpty[*bool]()
	assert.Nil(t, s9)

	s10 = NilOrEmpty[bool]()
	assert.NotNil(t, s10)
	assert.Equal(t, false, s10)

	s11 = NilOrEmpty[chan *demoStruct]()
	assert.Nil(t, s11)

	s12 = NilOrEmpty[chan demoStruct]()
	assert.Nil(t, s12)
}
