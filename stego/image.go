package stego

import (
	"image"
	"io"
	"sync"
)

const (
	COLOR_RGBA = iota
	COLOR_GRAY
)

type Image struct {
	rwm       sync.RWMutex
	img       image.Image
	colortype int
	entropy   uint32
	stegosize uint32
	pix       []uint8
}

func NewImage(img image.Image) *Image {
	if img == nil {
		return nil
	}
	r := img.Bounds()
	sz := r.Dx() * r.Dy()
	var bpp int
	steg := &Image{}
	switch img.(type) {
	case *image.RGBA:
		steg.colortype = COLOR_RGBA
		steg.img = img
		pix := img.(*image.RGBA).Pix
		steg.pix = pix
		for i, y := 0, r.Min.Y; y < r.Max.Y; y++ {
			for x := r.Min.X; x < r.Max.X; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				pix[i], pix[i+1], pix[i+2] = uint8(r), uint8(g), uint8(b)
				i += 4
			}
		}
		bpp = 3
	case *image.Gray:
		steg.colortype = COLOR_GRAY
		steg.img = img
		pix := img.(*image.Gray).Pix
		steg.pix = pix
		for i, y := 0, r.Min.Y; y < r.Max.Y; y++ {
			for x := r.Min.X; x < r.Max.X; x++ {
				g, _, _, _ := img.At(x, y).RGBA()
				pix[i] = uint8(g >> 8) // this is probably horribly wrong
				i++
			}
		}
		bpp = 1
	default:
		panic("color model not supported")
	}
	steg.entropy = uint32(sz >> 3 * bpp)
	return steg
}

func (img *Image) Entropy() uint32 {
	return img.entropy
}

func (img *Image) Reader() io.Reader {
	img.rwm.RLock()
	defer img.rwm.RUnlock()
	r := &imageReader{img: img}
	switch img.colortype {
	case COLOR_RGBA:
		r.index = 42
	case COLOR_GRAY:
		r.index = 32
	}
	r.img.stegosize = r.img.readSize()
	return r
}

func (img *Image) Writer() io.Writer {
	img.rwm.RLock()
	defer img.rwm.RUnlock()
	w := &imageWriter{img: img}
	switch img.colortype {
	case COLOR_RGBA:
		w.index = 42
	case COLOR_GRAY:
		w.index = 32
	default:
		panic("color model not supported")
	}
	w.img.stegosize = 0
	return w
}

func (img *Image) Export() image.Image {
	return img.img
}

func (img *Image) writeSize(n uint32) {
	switch img.colortype {
	case COLOR_RGBA:
		img.writeSizeRGBA(n)
	case COLOR_GRAY:
		img.writeSizeGray(n)
	default:
		panic("color model not supported")
	}
}

func (img *Image) readSize() (n uint32) {
	switch img.colortype {
	case COLOR_RGBA:
		n = img.readSizeRGBA()
	case COLOR_GRAY:
		n = img.readSizeGray()
	default:
		panic("color model not supported")
	}
	return
}

func (img *Image) writeSizeRGBA(n uint32) {
	pix := img.pix
	pix[0] = pix[0]&^1 | uint8(n&1)
	pix[1] = pix[1]&^1 | uint8(n>>1&1)
	pix[2] = pix[2]&^1 | uint8(n>>2&1)

	pix[4] = pix[4]&^1 | uint8(n>>3&1)
	pix[5] = pix[5]&^1 | uint8(n>>4&1)
	pix[6] = pix[6]&^1 | uint8(n>>5&1)

	pix[8] = pix[8]&^1 | uint8(n>>6&1)
	pix[9] = pix[9]&^1 | uint8(n>>7&1)
	pix[10] = pix[10]&^1 | uint8(n>>8&1)

	pix[12] = pix[12]&^1 | uint8(n>>9&1)
	pix[13] = pix[13]&^1 | uint8(n>>10&1)
	pix[14] = pix[14]&^1 | uint8(n>>11&1)

	pix[16] = pix[16]&^1 | uint8(n>>12&1)
	pix[17] = pix[17]&^1 | uint8(n>>13&1)
	pix[18] = pix[18]&^1 | uint8(n>>14&1)

	pix[20] = pix[20]&^1 | uint8(n>>15&1)
	pix[21] = pix[21]&^1 | uint8(n>>16&1)
	pix[22] = pix[22]&^1 | uint8(n>>17&1)

	pix[24] = pix[24]&^1 | uint8(n>>18&1)
	pix[25] = pix[25]&^1 | uint8(n>>19&1)
	pix[26] = pix[26]&^1 | uint8(n>>20&1)

	pix[28] = pix[28]&^1 | uint8(n>>21&1)
	pix[29] = pix[29]&^1 | uint8(n>>22&1)
	pix[30] = pix[30]&^1 | uint8(n>>23&1)

	pix[32] = pix[32]&^1 | uint8(n>>24&1)
	pix[33] = pix[33]&^1 | uint8(n>>25&1)
	pix[34] = pix[34]&^1 | uint8(n>>26&1)

	pix[36] = pix[36]&^1 | uint8(n>>27&1)
	pix[37] = pix[37]&^1 | uint8(n>>28&1)
	pix[38] = pix[38]&^1 | uint8(n>>29&1)

	pix[40] = pix[40]&^1 | uint8(n>>30&1)
	pix[41] = pix[41]&^1 | uint8(n>>31&1)
}

func (img *Image) writeSizeGray(n uint32) {
	pix := img.pix
	pix[0] = pix[0]&^1 | uint8(n&1)
	pix[1] = pix[1]&^1 | uint8(n>>1&1)
	pix[2] = pix[2]&^1 | uint8(n>>2&1)
	pix[3] = pix[3]&^1 | uint8(n>>3&1)
	pix[4] = pix[4]&^1 | uint8(n>>4&1)
	pix[5] = pix[5]&^1 | uint8(n>>5&1)
	pix[6] = pix[6]&^1 | uint8(n>>6&1)
	pix[7] = pix[7]&^1 | uint8(n>>7&1)

	pix[8] = pix[8]&^1 | uint8(n>>8&1)
	pix[9] = pix[9]&^1 | uint8(n>>9&1)
	pix[10] = pix[10]&^1 | uint8(n>>10&1)
	pix[11] = pix[11]&^1 | uint8(n>>11&1)
	pix[12] = pix[12]&^1 | uint8(n>>12&1)
	pix[13] = pix[13]&^1 | uint8(n>>13&1)
	pix[14] = pix[14]&^1 | uint8(n>>14&1)
	pix[15] = pix[15]&^1 | uint8(n>>15&1)

	pix[16] = pix[16]&^1 | uint8(n>>16&1)
	pix[17] = pix[17]&^1 | uint8(n>>17&1)
	pix[18] = pix[18]&^1 | uint8(n>>18&1)
	pix[19] = pix[19]&^1 | uint8(n>>19&1)
	pix[20] = pix[20]&^1 | uint8(n>>20&1)
	pix[21] = pix[21]&^1 | uint8(n>>21&1)
	pix[22] = pix[22]&^1 | uint8(n>>22&1)
	pix[23] = pix[23]&^1 | uint8(n>>23&1)

	pix[24] = pix[24]&^1 | uint8(n>>24&1)
	pix[25] = pix[25]&^1 | uint8(n>>25&1)
	pix[26] = pix[26]&^1 | uint8(n>>26&1)
	pix[27] = pix[27]&^1 | uint8(n>>27&1)
	pix[28] = pix[28]&^1 | uint8(n>>28&1)
	pix[29] = pix[29]&^1 | uint8(n>>29&1)
	pix[30] = pix[30]&^1 | uint8(n>>30&1)
	pix[31] = pix[31]&^1 | uint8(n>>31&1)
}

func (img *Image) readSizeRGBA() (n uint32) {
	pix := img.pix
	n = uint32(pix[0] & 1)
	n |= uint32(pix[1]&1) << 1
	n |= uint32(pix[2]&1) << 2

	n |= uint32(pix[4]&1) << 3
	n |= uint32(pix[5]&1) << 4
	n |= uint32(pix[6]&1) << 5

	n |= uint32(pix[8]&1) << 6
	n |= uint32(pix[9]&1) << 7
	n |= uint32(pix[10]&1) << 8

	n |= uint32(pix[12]&1) << 9
	n |= uint32(pix[13]&1) << 10
	n |= uint32(pix[14]&1) << 11

	n |= uint32(pix[16]&1) << 12
	n |= uint32(pix[17]&1) << 13
	n |= uint32(pix[18]&1) << 14

	n |= uint32(pix[20]&1) << 15
	n |= uint32(pix[21]&1) << 16
	n |= uint32(pix[22]&1) << 17

	n |= uint32(pix[24]&1) << 18
	n |= uint32(pix[25]&1) << 19
	n |= uint32(pix[26]&1) << 20

	n |= uint32(pix[28]&1) << 21
	n |= uint32(pix[29]&1) << 22
	n |= uint32(pix[30]&1) << 23

	n |= uint32(pix[32]&1) << 24
	n |= uint32(pix[33]&1) << 25
	n |= uint32(pix[34]&1) << 26

	n |= uint32(pix[36]&1) << 27
	n |= uint32(pix[37]&1) << 28
	n |= uint32(pix[38]&1) << 29

	n |= uint32(pix[40]&1) << 30
	n |= uint32(pix[41]&1) << 31
	return
}

func (img *Image) readSizeGray() (n uint32) {
	pix := img.pix
	n = uint32(pix[0] & 1)
	n |= uint32(pix[1]&1) << 1
	n |= uint32(pix[2]&1) << 2
	n |= uint32(pix[3]&1) << 3
	n |= uint32(pix[4]&1) << 4
	n |= uint32(pix[5]&1) << 5
	n |= uint32(pix[6]&1) << 6
	n |= uint32(pix[7]&1) << 7

	n |= uint32(pix[8]&1) << 8
	n |= uint32(pix[9]&1) << 9
	n |= uint32(pix[10]&1) << 10
	n |= uint32(pix[11]&1) << 11
	n |= uint32(pix[12]&1) << 12
	n |= uint32(pix[13]&1) << 13
	n |= uint32(pix[14]&1) << 14
	n |= uint32(pix[15]&1) << 15

	n |= uint32(pix[16]&1) << 16
	n |= uint32(pix[17]&1) << 17
	n |= uint32(pix[18]&1) << 18
	n |= uint32(pix[19]&1) << 19
	n |= uint32(pix[20]&1) << 20
	n |= uint32(pix[21]&1) << 21
	n |= uint32(pix[22]&1) << 22
	n |= uint32(pix[23]&1) << 23

	n |= uint32(pix[24]&1) << 24
	n |= uint32(pix[25]&1) << 25
	n |= uint32(pix[26]&1) << 26
	n |= uint32(pix[27]&1) << 27
	n |= uint32(pix[28]&1) << 28
	n |= uint32(pix[29]&1) << 29
	n |= uint32(pix[30]&1) << 30
	n |= uint32(pix[31]&1) << 31
	return
}

// -----------------------------------------------------------------------------

type imageReader struct {
	img       *Image
	index     int // read index
	bytesread uint32
}

func (r *imageReader) Read(buf []byte) (n int, err error) {
	r.img.rwm.RLock()
	defer r.img.rwm.RUnlock()
	switch r.img.colortype {
	case COLOR_RGBA:
		n, err = r.readRGBA(buf)
	case COLOR_GRAY:
		n, err = r.readGray(buf)
	default:
		panic("color model not supported")
	}
	return
}

func (r *imageReader) readRGBA(buf []byte) (int, error) {
	i := r.index
	pix := r.img.pix
	sz := r.img.stegosize - r.bytesread
	var err error
	if sz > uint32(len(buf)) {
		sz = uint32(len(buf))
	} else {
		err = io.EOF
	}
	for j := uint32(0); j < sz; j++ {
		var b byte
		switch i & 3 {
		case 3:
			i++
			fallthrough
		case 0:
			b = pix[i] & 1
			b |= pix[i+1] & 1 << 1
			b |= pix[i+2] & 1 << 2
			b |= pix[i+4] & 1 << 3
			b |= pix[i+5] & 1 << 4
			b |= pix[i+6] & 1 << 5
			b |= pix[i+8] & 1 << 6
			b |= pix[i+9] & 1 << 7
			buf[j] = b
			i += 10
		case 1:
			b = pix[i] & 1
			b |= pix[i+1] & 1 << 1
			b |= pix[i+3] & 1 << 2
			b |= pix[i+4] & 1 << 3
			b |= pix[i+5] & 1 << 4
			b |= pix[i+7] & 1 << 5
			b |= pix[i+8] & 1 << 6
			b |= pix[i+9] & 1 << 7
			buf[j] = b
			i += 10
		case 2:
			b = pix[i] & 1
			b |= pix[i+2] & 1 << 1
			b |= pix[i+3] & 1 << 2
			b |= pix[i+4] & 1 << 3
			b |= pix[i+6] & 1 << 4
			b |= pix[i+7] & 1 << 5
			b |= pix[i+8] & 1 << 6
			b |= pix[i+10] & 1 << 7
			buf[j] = b
			i += 11
		}
		buf[j] = b
	}
	r.index = i
	r.bytesread += sz
	return int(sz), err
}

func (r *imageReader) readGray(buf []byte) (int, error) {
	i := r.index
	pix := r.img.pix
	sz := r.img.stegosize - r.bytesread
	var err error
	if sz > uint32(len(buf)) {
		sz = uint32(len(buf))
	} else {
		err = io.EOF
	}
	for j := uint32(0); j < sz; j++ {
		b := pix[i] & 1
		b |= pix[i+1] & 1 << 1
		b |= pix[i+2] & 1 << 2
		b |= pix[i+3] & 1 << 3
		b |= pix[i+4] & 1 << 4
		b |= pix[i+5] & 1 << 5
		b |= pix[i+6] & 1 << 6
		b |= pix[i+7] & 1 << 7
		i += 8
		buf[j] = b
		sz -= 8
	}
	r.index = i
	r.bytesread += sz
	return int(sz), err
}

// -----------------------------------------------------------------------------

type imageWriter struct {
	img   *Image
	index int // read index
}

func (w *imageWriter) Write(buf []byte) (n int, err error) {
	w.img.rwm.Lock()
	defer w.img.rwm.Unlock()
	switch w.img.colortype {
	case COLOR_RGBA:
		n, err = w.writeRGBA(buf)
	case COLOR_GRAY:
		n, err = w.writeGray(buf)
	default:
		panic("color model not supported")
	}
	if n > 0 {
		w.img.stegosize += uint32(n)
		w.img.writeSize(w.img.stegosize)
	}
	return
}

func (w *imageWriter) writeRGBA(buf []byte) (int, error) {
	i := w.index
	pix := w.img.pix
	sz := w.img.entropy - w.img.stegosize
	var err error
	if sz > uint32(len(buf)) {
		sz = uint32(len(buf))
	} else {
		err = io.EOF
	}
	for j := uint32(0); j < sz; j++ {
		b := buf[j]
		switch i & 3 {
		case 3:
			i++
			fallthrough
		case 0:
			pix[i] = pix[i]&^1 | b&1
			pix[i+1] = pix[i+1]&^1 | b>>1&1
			pix[i+2] = pix[i+2]&^1 | b>>2&1
			pix[i+4] = pix[i+4]&^1 | b>>3&1
			pix[i+5] = pix[i+5]&^1 | b>>4&1
			pix[i+6] = pix[i+6]&^1 | b>>5&1
			pix[i+8] = pix[i+8]&^1 | b>>6&1
			pix[i+9] = pix[i+9]&^1 | b>>7&1
			i += 10
		case 1:
			pix[i] = pix[i]&^1 | b&1
			pix[i+1] = pix[i+1]&^1 | b>>1&1
			pix[i+3] = pix[i+3]&^1 | b>>2&1
			pix[i+4] = pix[i+4]&^1 | b>>3&1
			pix[i+5] = pix[i+5]&^1 | b>>4&1
			pix[i+7] = pix[i+7]&^1 | b>>5&1
			pix[i+8] = pix[i+8]&^1 | b>>6&1
			pix[i+9] = pix[i+9]&^1 | b>>7&1
			i += 10
		case 2:
			pix[i] = pix[i]&^1 | b&1
			pix[i+2] = pix[i+2]&^1 | b>>1&1
			pix[i+3] = pix[i+3]&^1 | b>>2&1
			pix[i+4] = pix[i+4]&^1 | b>>3&1
			pix[i+6] = pix[i+6]&^1 | b>>4&1
			pix[i+7] = pix[i+7]&^1 | b>>5&1
			pix[i+8] = pix[i+8]&^1 | b>>6&1
			pix[i+10] = pix[i+10]&^1 | b>>7&1
			i += 11
		}
	}
	w.index = i
	return int(sz), err
}

func (w *imageWriter) writeGray(buf []byte) (int, error) {
	i := w.index
	pix := w.img.pix
	sz := w.img.entropy - w.img.stegosize
	var err error
	if sz > uint32(len(buf)) {
		sz = uint32(len(buf))
	} else {
		err = io.EOF
	}
	for j := uint32(0); j < sz; j++ {
		b := buf[j]
		pix[i] = pix[i]&^1 | b&1
		pix[i+1] = pix[i+1]&^1 | b>>1&1
		pix[i+2] = pix[i+2]&^1 | b>>2&1
		pix[i+3] = pix[i+3]&^1 | b>>3&1
		pix[i+4] = pix[i+4]&^1 | b>>4&1
		pix[i+5] = pix[i+5]&^1 | b>>5&1
		pix[i+6] = pix[i+6]&^1 | b>>6&1
		pix[i+7] = pix[i+7]&^1 | b>>7&1
		i += 8
	}
	w.index = i
	return int(sz), err
}
