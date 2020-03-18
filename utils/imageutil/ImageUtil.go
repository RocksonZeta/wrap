package imageutil

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"math"
	"os"
	"strings"

	"github.com/disintegration/imaging"
)

type ImageTransform struct {
	Raw    image.Image
	Im     image.Image
	Format imaging.Format
	// RawWidth, RawHeight int
}

func NewImageFileTransform(filename string) (*ImageTransform, error) {
	format, err := imaging.FormatFromFilename(filename)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewImageTransform(f, format)
}
func NewImageTransform(file io.Reader, format imaging.Format) (*ImageTransform, error) {
	t := &ImageTransform{Format: format}
	im, err := imaging.Decode(file)
	if nil != err {
		return nil, err
	}
	t.Raw = im
	t.Im = im
	return t, nil
}
func NewImageTransformExt(file io.Reader, ext string) (*ImageTransform, error) {
	i := strings.LastIndex(ext, ".")
	if i != -1 {
		ext = ext[i+1:]
	}
	format, err := imaging.FormatFromExtension(strings.ToLower(ext))
	if err != nil {
		return nil, err
	}
	return NewImageTransform(file, format)
}

func (t *ImageTransform) ImClone() *ImageTransform {
	return &ImageTransform{
		Raw:    t.Im,
		Im:     t.Im,
		Format: t.Format,
	}
}
func (t *ImageTransform) CropRect(rect image.Rectangle) *ImageTransform {
	t.Im = imaging.Crop(t.Raw, rect)
	return t
}
func (t *ImageTransform) Crop(x, y, w, h int) *ImageTransform {
	t.Im = imaging.Crop(t.Raw, image.Rect(x, y, x+w, y+h))
	return t
}

//Crop4 crop x, y, w, h
func (t *ImageTransform) Crop4(crop []int) *ImageTransform {
	return t.Crop(crop[0], crop[1], crop[2], crop[3])
}

func (t *ImageTransform) Resize(w, h int) *ImageTransform {

	t.Im = imaging.Resize(t.Raw, w, h, imaging.Lanczos)
	return t
}

//ResizeMax 不改变压缩比，选择最大的压缩比
func (t *ImageTransform) ResizeMax(maxWidth, maxHeight int) bool {
	ow, oh := t.Size()
	if ow < maxWidth && oh < maxHeight {
		return false
	}
	rw := -1.0
	rh := -1.0
	if maxHeight > 0 {
		rh = float64(oh) / float64(maxHeight)
	}
	if maxWidth > 0 {
		rw = float64(ow) / float64(maxWidth)
	}
	nw, nh := maxWidth, maxHeight
	if rw > rh {
		nh = int(float64(oh) / rw)
	} else {
		nw = int(float64(ow) / rh)
	}

	t.Resize(nw, nh)
	return true
}
func (t *ImageTransform) ResizeKeepRatio(w, h int) *ImageTransform {
	ow, oh := t.Size()
	nw, nh, x, y := GetGoodResize(ow, oh, w, h)
	t.Im = imaging.Resize(t.Raw, nw, nh, imaging.Lanczos)
	t.Crop(x, y, w, h)
	return t
}

//GetResizeParam return : resize and crop
func GetGoodResize(ow, oh, w, h int) (int, int, int, int) {
	rw := float64(w) / float64(ow)
	rh := float64(h) / float64(oh)
	x, y := 0, 0
	nw, nh := w, h
	if math.Abs(rw-1) < math.Abs(rh-1) { //选择压缩比最小的 ，
		nh = int(float64(oh) * rw)
		y = (h - nh) / 2
		if y < 0 {
			y = -y
		}
	} else {
		nw = int(float64(ow) * rh)
		x = (w - nw) / 2
		if x < 0 {
			x = -x
		}
	}

	return nw, nh, x, y
}
func (t *ImageTransform) Buffer() (*bytes.Buffer, error) {
	buff := new(bytes.Buffer)
	err := imaging.Encode(buff, t.Im, t.Format, imaging.PNGCompressionLevel(png.BestCompression), imaging.JPEGQuality(90))
	if nil != err {
		return nil, err
	}
	return buff, nil
}
func (t *ImageTransform) Save(filePath string) error {
	return imaging.Save(t.Im, filePath, imaging.PNGCompressionLevel(png.BestCompression), imaging.JPEGQuality(90))
}
func (t *ImageTransform) Size() (int, int) {
	return t.Im.Bounds().Dx(), t.Im.Bounds().Dy()
}
func (t *ImageTransform) Write(filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if nil != err {
		return err
	}
	defer f.Close()
	src, err := t.Buffer()
	if nil != err {
		return err
	}
	io.Copy(f, src)
	return nil
}

// func GetImageFormat(filename string) (imaging.Format, error) {
// 	formats := map[string]imaging.Format{
// 		".jpg":  imaging.JPEG,
// 		".jpeg": imaging.JPEG,
// 		".png":  imaging.PNG,
// 		".tif":  imaging.TIFF,
// 		".tiff": imaging.TIFF,
// 		".bmp":  imaging.BMP,
// 		".gif":  imaging.GIF,
// 	}

// 	ext := strings.ToLower(filepath.Ext(filename))
// 	f, ok := formats[ext]
// 	if !ok {
// 		return 0, imaging.ErrUnsupportedFormat
// 	}
// 	return f, nil
// }
