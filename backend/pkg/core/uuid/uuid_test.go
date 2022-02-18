package uuid

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	uuid "github.com/satori/go.uuid"
)

const nilUUID = "00000000-0000-0000-0000-000000000000"

func TestOrderedUUID(t *testing.T) {
	reset := func() {}
	reset()
	t.Run("Test NewOrderedUUID should return a not null UUID", func(t *testing.T) {
		u := NewOrderedUUID()

		assert.NotEqual(t, nilUUID, u.String())
	})
	reset()
	t.Run("Test FromUUIDv1 reorders data correctly", func(t *testing.T) {
		u := uuid.UUID{}
		data, _ := hex.DecodeString("aaaabbbbccccdddd1234567890123456")
		copy(u[0:], data)

		o := OrderedUUID{}
		o.FromUUIDv1(u)

		assert.Equal(t, "aaaabbbb-cccc-dddd-1234-567890123456", u.String())
		assert.Equal(t, "ddddcccc-aaaa-bbbb-1234-567890123456", o.String())
	})
	reset()
	t.Run("Test Value matches Bytes", func(t *testing.T) {
		u := NewOrderedUUID()
		val, err := u.Value()

		assert.Nil(t, err)
		assert.Equal(t, u.Bytes[:], val)
	})
	reset()
	t.Run("Test Value gives nil value if invalid", func(t *testing.T) {
		u := NewOrderedUUID()
		u.Valid = false
		val, err := u.Value()

		assert.Nil(t, err)
		assert.Nil(t, val)
	})
}
func TestOrderedUUIDDatabase(t *testing.T) {
	reset := func() {}
	reset()
	t.Run("Test Value gives nil value if empty", func (t *testing.T) {
		u := OrderedUUID{}
		val, err := u.Value()

		assert.Nil(t, err)
		assert.Nil(t, val)
	})
	reset()
	t.Run("Test Value gives nil value if empty", func (t *testing.T) {
		u := OrderedUUID{}
		val, err := u.Value()

		assert.Nil(t, err)
		assert.Nil(t, val)
	})
	reset()
	t.Run("Test Scan correctly reads bytes", func (t *testing.T) {
		u := OrderedUUID{}
		val := []byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
		err := u.Scan(val)

		var expectedVal [16]byte
		copy(expectedVal[0:], val)

		assert.Nil(t, err)
		assert.Equal(t, expectedVal, u.Bytes)
		assert.True(t, u.Valid)
	})
	reset()
	t.Run("Test Scan gives invalid uuid on nil value", func (t *testing.T) {
		u := OrderedUUID{}
		err := u.Scan(nil)

		assert.Nil(t, err)
		assert.Equal(t, nilUUID, u.String())
		assert.False(t, u.Valid)
	})
	reset()
	t.Run("Test Scan gives error on invalid type", func (t *testing.T) {
		u := OrderedUUID{}
		err := u.Scan("string")

		assert.NotNil(t, err)
		assert.Equal(t, "uuid: cannot convert string to UUID", err.Error())
		assert.Equal(t, nilUUID, u.String())
		assert.False(t, u.Valid)
	})
	reset()
	t.Run("Test Scan gives error on invalid number of bytes", func (t *testing.T) {
		val := []byte{0x6b, 0xa7, 0xb8, 0x10}

		u := OrderedUUID{}
		err := u.Scan(val)

		assert.NotNil(t, err)
		assert.Equal(t, "uuid: UUID must be exactly 16 bytes long, got 4 bytes", err.Error())
		assert.Equal(t, nilUUID, u.String())
		assert.False(t, u.Valid)

		val = []byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8, 0xc8}
		err = u.Scan(val)
		assert.NotNil(t, err)
		assert.Equal(t, "uuid: UUID must be exactly 16 bytes long, got 17 bytes", err.Error())
		assert.Equal(t, nilUUID, u.String())
		assert.False(t, u.Valid)
	})
}

func TestOrderedUUIDFromString(t *testing.T) {
	b := [16]byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	emptyBytes := [16]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	// Valid uuid
	s1 := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	s2 := "6ba7b8109dad11d180b400c04fd430c8"

	// With non valid hex char
	s3 := "6ba7b810zzzzzzzz80b400c04fd430c8"
	s4 := "6ba7b810-zzzz-zzzz-80b4-00c04fd430c8"

	// Without the dash or not at the right place
	s5 := "6ba7b81019dad111d1180b4100c04fd430c8"
	s6 := "6ba7b8109d-ad-11d1-80b4-00c04fd430c8"

	// Invalid string length
	s7 := "6ba7b810"

	u0, err := FromString("")
	assert.Nil(t, err)
	assert.Equal(t, emptyBytes, u0.Bytes)
	assert.False(t, u0.Valid)

	u1, err := FromString(s1)
	assert.Nil(t, err)
	assert.Equal(t, b, u1.Bytes)
	assert.True(t, u1.Valid)

	u2, err := FromString(s2)
	assert.Nil(t, err)
	assert.Equal(t, b, u2.Bytes)
	assert.True(t, u2.Valid)

	_, err = FromString(s3)
	assert.NotNil(t, err)

	_, err = FromString(s4)
	assert.NotNil(t, err)

	_, err = FromString(s5)
	assert.NotNil(t, err)

	_, err = FromString(s6)
	assert.NotNil(t, err)

	_, err = FromString(s7)
	assert.NotNil(t, err)
}

func TestOrderedUUIMarshalText(t *testing.T) {
	reset := func() {}
	reset()
	t.Run("Test OrderedUUID can be stringified", func (t *testing.T) {
		b := [16]byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
		u := OrderedUUID{Bytes: b, Valid: true}

		expected := []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

		jsonString, err := u.MarshalText()
		assert.Nil(t, err)
		assert.Equal(t, expected, jsonString)
	})
	reset()
	t.Run("Test OrderedUUID gives nil string when invalid", func (t *testing.T) {
		u := OrderedUUID{}

		expected := []byte(nil)

		jsonString, err := u.MarshalText()
		assert.Nil(t, err)
		assert.Equal(t, expected, jsonString)
	})
}

func TestOrderedUUIDFromBytes(t *testing.T) {
	b := [16]byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	u := FromBytes(b)

	assert.Equal(t, b, u.Bytes)
	assert.True(t, u.Valid)
}
