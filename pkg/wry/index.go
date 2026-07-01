package wry

import (
	"encoding/binary"
)

// maxSearchIter caps binary-search iterations. log2 of any plausible index
// count is well under this, so it only ever fires on corrupt/misaligned headers
// (where the r-l==entryLen termination could otherwise never be reached).
const maxSearchIter = 128

// NOTE: the receiver type parameter is deliberately named uint32/uint64 (as in
// the original), so uint32(...)/uint64(...) below convert to the parameter type
// rather than the builtin. The explicit conversions bridge the builtin values
// returned by encoding/binary and Bytes3ToUint32.

func (db *IPDB[uint32]) SearchIndexV4(ip uint32) uint32 {
	ipLen := uint32(db.IPLen)
	entryLen := uint32(db.OffLen) + uint32(db.IPLen)
	dataLen := uint32(len(db.Data))

	// Validate header-derived offsets before indexing into Data. Checks are
	// written to avoid uint overflow (idx+entryLen could wrap past dataLen).
	if ipLen < 4 || db.OffLen < 3 || db.IdxStart > db.IdxEnd ||
		entryLen > dataLen || db.IdxEnd > dataLen-entryLen {
		return 0
	}

	l, r := db.IdxStart, db.IdxEnd
	for iter := 0; iter < maxSearchIter; iter++ {
		mid := (r-l)/entryLen/2*entryLen + l
		if mid > dataLen-entryLen {
			return 0
		}
		buf := db.Data[mid : mid+entryLen]
		ipc := uint32(binary.LittleEndian.Uint32(buf[:4]))

		if r-l == entryLen {
			if ip >= uint32(binary.LittleEndian.Uint32(db.Data[r:r+4])) {
				buf = db.Data[r : r+entryLen]
			}
			return uint32(Bytes3ToUint32(buf[ipLen:entryLen]))
		}

		if ipc > ip {
			r = mid
		} else if ipc < ip {
			if l == mid { // no forward progress: avoid an infinite loop
				return uint32(Bytes3ToUint32(buf[ipLen:entryLen]))
			}
			l = mid
		} else {
			return uint32(Bytes3ToUint32(buf[ipLen:entryLen]))
		}
	}
	return 0
}

func (db *IPDB[uint64]) SearchIndexV6(ip uint64) uint32 {
	ipLen := uint64(db.IPLen)
	entryLen := uint64(db.OffLen) + uint64(db.IPLen)
	dataLen := uint64(len(db.Data))

	// Validate header-derived offsets before indexing into Data. Checks are
	// written to avoid uint overflow (idx+entryLen could wrap past dataLen).
	if ipLen < 8 || db.OffLen < 3 || db.IdxStart > db.IdxEnd ||
		entryLen > dataLen || db.IdxEnd > dataLen-entryLen {
		return 0
	}

	l, r := db.IdxStart, db.IdxEnd
	for iter := 0; iter < maxSearchIter; iter++ {
		mid := (r-l)/entryLen/2*entryLen + l
		if mid > dataLen-entryLen {
			return 0
		}
		buf := db.Data[mid : mid+entryLen]
		ipc := uint64(binary.LittleEndian.Uint64(buf[:8]))

		if r-l == entryLen {
			if ip >= uint64(binary.LittleEndian.Uint64(db.Data[r:r+8])) {
				buf = db.Data[r : r+entryLen]
			}
			return Bytes3ToUint32(buf[ipLen:entryLen])
		}

		if ipc > ip {
			r = mid
		} else if ipc < ip {
			if l == mid { // no forward progress: avoid an infinite loop
				return Bytes3ToUint32(buf[ipLen:entryLen])
			}
			l = mid
		} else {
			return Bytes3ToUint32(buf[ipLen:entryLen])
		}
	}
	return 0
}
