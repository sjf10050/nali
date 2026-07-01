package wry

import (
	"encoding/binary"
	"testing"
)

func TestBytes3ToUint32(t *testing.T) {
	if got := Bytes3ToUint32([]byte{0x01, 0x02, 0x03}); got != 0x030201 {
		t.Fatalf("Bytes3ToUint32 = %#x, want 0x030201", got)
	}
}

func TestReaderMissingNulTerminator(t *testing.T) {
	// No 0x00 in the buffer: readString must flag an error, not panic or
	// slice with a length of uint32(-1).
	r := NewReader([]byte{0x41, 0x42, 0x43})
	if s := r.readString(true); s != "" {
		t.Fatalf("readString = %q, want empty", s)
	}
	if r.Err() == nil {
		t.Fatal("expected error for missing NUL terminator")
	}
}

func TestReaderOffsetOutOfRange(t *testing.T) {
	r := NewReader([]byte{0x01, 0x02, 0x03})
	r.Parse(1 << 20) // offset far beyond the buffer
	if r.Err() == nil {
		t.Fatal("expected error for out-of-range offset")
	}
}

func TestParseCyclicRedirectTerminates(t *testing.T) {
	// [0x01][00 00 00] is redirect-mode-1 with a target offset of 0, i.e. a
	// self-referencing cycle. The depth guard must stop it instead of
	// overflowing the stack.
	r := NewReader([]byte{0x01, 0x00, 0x00, 0x00})
	r.Parse(0)
	if r.Err() == nil {
		t.Fatal("expected error from exceeding redirect depth")
	}
}

func buildIndexV4(entries []struct{ ip, off uint32 }) ([]byte, uint32, uint32) {
	const entryLen = 7
	start := uint32(8)
	end := start + uint32(len(entries)-1)*entryLen
	data := make([]byte, int(end)+entryLen+16)
	binary.LittleEndian.PutUint32(data[0:4], start)
	binary.LittleEndian.PutUint32(data[4:8], end)
	for i, e := range entries {
		base := int(start) + i*entryLen
		binary.LittleEndian.PutUint32(data[base:base+4], e.ip)
		data[base+4] = byte(e.off)
		data[base+5] = byte(e.off >> 8)
		data[base+6] = byte(e.off >> 16)
	}
	return data, start, end
}

func TestSearchIndexV4(t *testing.T) {
	entries := []struct{ ip, off uint32 }{
		{0x01000000, 100},
		{0x02000000, 200},
		{0x03000000, 300},
	}
	data, start, end := buildIndexV4(entries)
	db := IPDB[uint32]{Data: data, OffLen: 3, IPLen: 4, IdxStart: start, IdxEnd: end}

	cases := []struct {
		ip   uint32
		want uint32
	}{
		{0x02000000, 200}, // exact match
		{0x09000000, 300}, // above all -> last entry
		{0x01500000, 100}, // between first and second -> lower bound
	}
	for _, c := range cases {
		if got := db.SearchIndexV4(c.ip); got != c.want {
			t.Errorf("SearchIndexV4(%#x) = %d, want %d", c.ip, got, c.want)
		}
	}
}

func TestSearchIndexV4CorruptHeaderNoPanic(t *testing.T) {
	// IdxEnd well past the buffer, and misaligned: must return 0, not panic
	// or loop forever.
	db := IPDB[uint32]{Data: make([]byte, 32), OffLen: 3, IPLen: 4, IdxStart: 8, IdxEnd: 1 << 20}
	if got := db.SearchIndexV4(0x01020304); got != 0 {
		t.Fatalf("corrupt header should yield 0, got %d", got)
	}
}

func FuzzParse(f *testing.F) {
	f.Add([]byte{0x01, 0x00, 0x00, 0x00}, uint32(0))
	f.Add([]byte{0x02, 0x05, 0x00, 0x00, 0x41, 0x00}, uint32(0))
	f.Add([]byte("hello\x00world\x00"), uint32(3))
	f.Fuzz(func(t *testing.T, data []byte, offset uint32) {
		r := NewReader(data)
		r.Parse(offset) // must never panic
		_ = r.Result.String()
		_ = r.Err()
	})
}

func FuzzSearchIndexV4(f *testing.F) {
	f.Add(make([]byte, 32), uint32(0x01020304))
	f.Fuzz(func(t *testing.T, data []byte, ip uint32) {
		if len(data) < 8 {
			return
		}
		db := IPDB[uint32]{
			Data:     data,
			OffLen:   3,
			IPLen:    4,
			IdxStart: binary.LittleEndian.Uint32(data[0:4]),
			IdxEnd:   binary.LittleEndian.Uint32(data[4:8]),
		}
		_ = db.SearchIndexV4(ip) // must never panic or hang (iteration-capped)
	})
}

func FuzzSearchIndexV6(f *testing.F) {
	f.Add(make([]byte, 64), uint64(0x0102030405060708))
	f.Fuzz(func(t *testing.T, data []byte, ip uint64) {
		if len(data) < 16 {
			return
		}
		db := IPDB[uint64]{
			Data:     data,
			OffLen:   3,
			IPLen:    8,
			IdxStart: binary.LittleEndian.Uint64(data[0:8]),
			IdxEnd:   binary.LittleEndian.Uint64(data[8:16]),
		}
		_ = db.SearchIndexV6(ip) // must never panic or hang (iteration-capped)
	})
}
