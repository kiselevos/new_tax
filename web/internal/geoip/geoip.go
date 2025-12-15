package geoip

import (
	"encoding/csv"
	"io"
	"log/slog"
	"net"
	"os"
	"sort"
	"strconv"
)

type IPv4Range struct {
	From   uint32
	To     uint32
	Region string
}

type DB struct {
	Ranges []IPv4Range
}

func NewEmpty() *DB {
	return &DB{
		Ranges: nil,
	}
}

func LoadFromCSV(path string) (*DB, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("failed to close file", "err", err)
		}
	}()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	const maxIPv4 = uint64(4294967295)

	var ranges []IPv4Range

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(rec) < 6 {
			// пропускаем кривые строки
			continue
		}

		ipFromStr := rec[0]
		ipToStr := rec[1]
		region := rec[4]

		ipFrom, err := strconv.ParseUint(ipFromStr, 10, 64)
		if err != nil {
			continue
		}
		ipTo, err := strconv.ParseUint(ipToStr, 10, 64)
		if err != nil {
			continue
		}

		// оставляем только честные IPv4
		if ipFrom > maxIPv4 || ipTo > maxIPv4 {
			continue
		}

		ranges = append(ranges, IPv4Range{
			From:   uint32(ipFrom),
			To:     uint32(ipTo),
			Region: region,
		})
	}

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].From < ranges[j].From
	})

	return &DB{Ranges: ranges}, nil
}

func (db *DB) findRegionUint32(ip uint32) string {

	ranges := db.Ranges
	lo, hi := 0, len(ranges)-1

	for lo <= hi {
		mid := (lo + hi) / 2
		r := ranges[mid]

		if ip < r.From {
			hi = mid - 1
		} else if ip > r.To {
			lo = mid + 1
		} else {
			return r.Region
		}
	}

	return "Unknown"
}

func ipv4ToUint32(ipStr string) uint32 {
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		return 0
	}

	return uint32(ip[0])<<24 |
		uint32(ip[1])<<16 |
		uint32(ip[2])<<8 |
		uint32(ip[3])
}

func (db *DB) LookupRegion(ipStr string) string {
	if db == nil || len(db.Ranges) == 0 {
		return "Unknown"
	}

	ipUint := ipv4ToUint32(ipStr)

	if ipUint == 0 {
		return "Unknown"
	}

	return db.findRegionUint32(ipUint)
}
