// internal/types/point.go
package types

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"math"
)

// Point represents a geographic point for database operations
type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Scan implements the sql.Scanner interface for reading from database
func (p *Point) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		if lat, lon, ok := parseWKT(v); ok {
			p.Lat = lat
			p.Lng = lon
			return nil
		}
		if lat, lon, ok := parseEWKB(v); ok {
			p.Lat = lat
			p.Lng = lon
			return nil
		}
		return fmt.Errorf("cannot parse point from string: %s", v)
	case []byte:
		str := string(v)
		if lat, lon, ok := parseEWKB(str); ok {
			p.Lat = lat
			p.Lng = lon
			return nil
		}
		return fmt.Errorf("cannot parse point from bytes")
	default:
		return fmt.Errorf("cannot scan %T into Point", value)
	}
}

// Value implements the driver.Valuer interface for writing to database
func (p Point) Value() (driver.Value, error) {
	return fmt.Sprintf("SRID=4326;POINT(%f %f)", p.Lng, p.Lat), nil
}

// parseWKT parses Well-Known Text format: "POINT(longitude latitude)"
func parseWKT(wkt string) (lat, lon float64, ok bool) {
	var longitude, latitude float64
	n, err := fmt.Sscanf(wkt, "POINT(%f %f)", &longitude, &latitude)
	if err != nil || n != 2 {
		return 0, 0, false
	}
	return latitude, longitude, true
}

// parseEWKB parses Extended Well-Known Binary format from PostGIS
func parseEWKB(ewkb string) (lat, lon float64, ok bool) {
	// Remove any prefix like \x if present
	if len(ewkb) > 2 && ewkb[:2] == "\\x" {
		ewkb = ewkb[2:]
	}

	// Decode hex string to bytes
	data, err := hex.DecodeString(ewkb)
	if err != nil {
		return 0, 0, false
	}

	if len(data) < 25 { // Minimum size for a point with SRID
		return 0, 0, false
	}

	// Parse EWKB format
	// Byte 0: Endianness (01 = little endian, 00 = big endian)
	// Bytes 1-4: Geometry type (with SRID flag)
	// Bytes 5-8: SRID
	// Bytes 9-16: X coordinate (longitude)
	// Bytes 17-24: Y coordinate (latitude)

	endian := data[0]
	var x, y float64

	if endian == 1 { // Little endian
		// Skip geometry type (4 bytes) and SRID (4 bytes)
		x = math.Float64frombits(uint64(data[9]) | uint64(data[10])<<8 | uint64(data[11])<<16 | uint64(data[12])<<24 |
			uint64(data[13])<<32 | uint64(data[14])<<40 | uint64(data[15])<<48 | uint64(data[16])<<56)
		y = math.Float64frombits(uint64(data[17]) | uint64(data[18])<<8 | uint64(data[19])<<16 | uint64(data[20])<<24 |
			uint64(data[21])<<32 | uint64(data[22])<<40 | uint64(data[23])<<48 | uint64(data[24])<<56)
	} else { // Big endian
		x = math.Float64frombits(uint64(data[16]) | uint64(data[15])<<8 | uint64(data[14])<<16 | uint64(data[13])<<24 |
			uint64(data[12])<<32 | uint64(data[11])<<40 | uint64(data[10])<<48 | uint64(data[9])<<56)
		y = math.Float64frombits(uint64(data[24]) | uint64(data[23])<<8 | uint64(data[22])<<16 | uint64(data[21])<<24 |
			uint64(data[20])<<32 | uint64(data[19])<<40 | uint64(data[18])<<48 | uint64(data[17])<<56)
	}

	return y, x, true // Return latitude, longitude
}
