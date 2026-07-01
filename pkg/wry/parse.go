package wry

const (
	// RedirectMode1 [IP][0x01][国家和地区信息的绝对偏移地址]
	RedirectMode1 = 0x01
	// RedirectMode2 [IP][0x02][信息的绝对偏移][...] or [IP][国家][...]
	RedirectMode2 = 0x02
)

// maxRedirectDepth bounds redirect-mode-1 chains so a cyclic/self-referencing
// redirect in a malformed database can't recurse forever and overflow the stack.
const maxRedirectDepth = 16

// Parse reads the record at offset into the Reader's Result, following redirects.
func (r *Reader) Parse(offset uint32) {
	r.parse(offset, 0)
}

func (r *Reader) parse(offset uint32, depth int) {
	if r.err != nil {
		return
	}
	if depth > maxRedirectDepth {
		r.fail()
		return
	}
	if offset != 0 {
		r.seekAbs(offset)
	}

	switch r.readMode() {
	case RedirectMode1:
		r.readOffset(true)
		r.parse(0, depth+1)
	case RedirectMode2:
		r.Result.Country = r.parseRedMode2()
		r.Result.Area = r.readArea()
	default:
		r.seekBack()
		r.Result.Country = r.readString(true)
		r.Result.Area = r.readArea()
	}
}

func (r *Reader) parseRedMode2() string {
	r.readOffset(true)
	str := r.readString(false)
	r.seekBack()
	return str
}

func (r *Reader) readArea() string {
	mode := r.readMode()
	if mode == RedirectMode1 || mode == RedirectMode2 {
		offset := r.readOffset(true)
		if offset == 0 {
			return ""
		}
	} else {
		r.seekBack()
	}
	return r.readString(false)
}
