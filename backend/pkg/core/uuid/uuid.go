package uuid

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// UUIDSize is the size of the uuid
const UUIDSize = 16

var byteGroups = []int{8, 4, 4, 4, 12}

// NilUUID holds the zero value of UUID: 00000000-0000-0000-0000-000000000000
var NilUUID = OrderedUUID{
	Bytes: [16]byte{},
	Valid: true,
}

// OrderedUUID is an UUIDv1 with some bit shifted for database index performance
type OrderedUUID struct {
	Bytes [UUIDSize]byte
	Valid bool
}

// NewOrderedUUID creates a new UUIDv1 and converts it to an OrderedUUID
//  v1 UUID: aaaabbbb-cccc-dddd-1234-567890123456
//  transposed: ddddcccc-aaaa-bbbb-1234-567890123456
func NewOrderedUUID() OrderedUUID {
	u := uuid.NewV1()

	o := OrderedUUID{}
	o.FromUUIDv1(u)

	return o
}

// FromUUIDv1 populates Byte from uuidV1 and converts it
func (u *OrderedUUID) FromUUIDv1(uuid uuid.UUID) {
	buf := make([]byte, UUIDSize)

	copy(buf[0:2], uuid[6:8])
	copy(buf[2:4], uuid[4:6])
	copy(buf[4:8], uuid[0:4])
	copy(buf[8:], uuid[8:])

	copy(u.Bytes[0:], buf[0:])
	u.Valid = true
}

// String returns the string representation of the UUID
func (u *OrderedUUID) String() string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u.Bytes[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u.Bytes[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u.Bytes[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u.Bytes[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u.Bytes[10:])

	return string(buf)
}

// Value implements the driver.Valuer interface.
func (u OrderedUUID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}

	v := u.Bytes[0:16]
	return v, nil
}

// Scan implements the sql.Scanner interface.
func (u *OrderedUUID) Scan(src interface{}) error {
	u.Valid = false

	if src == nil {
		return nil
	}

	switch src := src.(type) {
	case []byte:
		if len(src) != UUIDSize {
			return fmt.Errorf("uuid: UUID must be exactly 16 bytes long, got %d bytes", len(src))
		}
		copy(u.Bytes[0:], src[0:])
		u.Valid = true

		return nil
	}
	return fmt.Errorf("uuid: cannot convert %T to UUID", src)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// Following formats are supported:
//   "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
//   "6ba7b8109dad11d180b400c04fd430c8"
// Supported UUID text representation follows:
//   uuid := canonical | hashlike | braced | urn
//   plain := canonical | hashlike
//   canonical := 4hexoct '-' 2hexoct '-' 2hexoct '-' 6hexoct
//   hashlike := 12hexoct
func (u *OrderedUUID) UnmarshalText(text []byte) (err error) {
	switch len(text) {
	case 0:
		u.Valid = false
		u.Bytes = [16]byte{}
		return nil
	case 32:
		return u.decodeHashLike(text)
	case 36:
		return u.decodeCanonical(text)
	default:
		return fmt.Errorf("uuid: incorrect UUID length: %s", text)
	}
}

// decodeCanonical decodes UUID string in format
// "6ba7b810-9dad-11d1-80b4-00c04fd430c8".
func (u *OrderedUUID) decodeCanonical(t []byte) (err error) {
	if t[8] != '-' || t[13] != '-' || t[18] != '-' || t[23] != '-' {
		return fmt.Errorf("uuid: incorrect UUID format %s", t)
	}

	src := t[:]
	dst := u.Bytes[:]

	for i, byteGroup := range byteGroups {
		if i > 0 {
			src = src[1:] // skip dash
		}
		_, err = hex.Decode(dst[:byteGroup/2], src[:byteGroup])
		if err != nil {
			return
		}
		src = src[byteGroup:]
		dst = dst[byteGroup/2:]
	}

	u.Valid = true
	return
}

// decodeHashLike decodes UUID string in format
// "6ba7b8109dad11d180b400c04fd430c8".
func (u *OrderedUUID) decodeHashLike(t []byte) (err error) {
	src := t[:]
	dst := u.Bytes[:]

	if _, err = hex.Decode(dst, src); err != nil {
		return err
	}

	u.Valid = true
	return
}

// MarshalText implements the encoding.TextMarshaler interface.
// The encoding is the same as returned by String.
func (u OrderedUUID) MarshalText() (text []byte, err error) {
	if !u.Valid {
		return nil, nil
	}
	text = []byte(u.String())
	return
}

// FromString returns an OrderedUUID parsed from string input.
// Input is expected in a form accepted by UnmarshalText.
func FromString(input string) (u OrderedUUID, err error) {
	err = u.UnmarshalText([]byte(input))
	return
}

// FromBytes returns an OrderedUUID parsed from the byte array.
func FromBytes(input [16]byte) OrderedUUID {
	u := OrderedUUID{}
	copy(u.Bytes[0:], input[0:])

	u.Valid = true
	return u
}
