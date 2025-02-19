package qrcode

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/yeqown/go-qrcode/matrix"
)

type ImageOption interface {
	apply(o *outputImageOptions)
}

// funcOption wraps a function that modifies outputImageOptions into an
// implementation of the ImageOption interface.
type funcOption struct {
	f func(oo *outputImageOptions)
}

func (fo *funcOption) apply(oo *outputImageOptions) {
	fo.f(oo)
}

func newFuncDialOption(f func(oo *outputImageOptions)) *funcOption {
	return &funcOption{
		f: f,
	}
}

// WithBgColor background color
func WithBgColor(c color.Color) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.bgColor = c
		_stateToRGBA[matrix.StateFalse] = oo.bgColor
	})
}

// WithBgColorRGBHex background color
func WithBgColorRGBHex(hex string) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.bgColor = hexToRGBA(hex)
		_stateToRGBA[matrix.StateFalse] = oo.bgColor
	})
}

// WithFgColor QR color
func WithFgColor(c color.Color) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.qrColor = c
		_stateToRGBA[matrix.StateTrue] = oo.qrColor
		_stateToRGBA[matrix.StateFinder] = oo.qrColor
	})
}

// WithFgColorRGBHex Hex string to set QR Color
func WithFgColorRGBHex(hex string) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.qrColor = hexToRGBA(hex)
		_stateToRGBA[matrix.StateTrue] = oo.qrColor
		_stateToRGBA[matrix.StateFinder] = oo.qrColor
	})
}

// WithLogoImage image should only has 1/5 width of QRCode at most
func WithLogoImage(img image.Image) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.logo = img
	})
}

// WithLogoImageFileJPEG load image from file, jpeg is required.
// image should only has 1/5 width of QRCode at most
func WithLogoImageFileJPEG(f string) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		fd, err := os.Open(f)
		if err != nil {
			fmt.Printf("could not open file(%s), error=%v\n", f, err)
			return
		}

		img, err := jpeg.Decode(fd)
		if err != nil {
			fmt.Printf("could not open file(%s), error=%v\n", f, err)
			return
		}

		oo.logo = img
	})
}

// WithLogoImageFilePNG load image from file, PNG is required.
// image should only has 1/5 width of QRCode at most
func WithLogoImageFilePNG(f string) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		fd, err := os.Open(f)
		if err != nil {
			fmt.Printf("Open file(%s) failed: %v\n", f, err)
			return
		}

		img, err := png.Decode(fd)
		if err != nil {
			fmt.Printf("Decode file(%s) as PNG failed: %v\n", f, err)
			return
		}

		oo.logo = img
	})
}

// WithQRWidth specify width of each qr block
func WithQRWidth(width uint8) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.qrWidth = int(width)
	})
}

// WithCircleShape use circle shape as rectangle(default)
func WithCircleShape() ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.shape = _shapeCircle
	})
}

// WithCustomShape use custom shape as rectangle(default)
func WithCustomShape(shape IShape) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		oo.shape = shape
	})
}

// WithBuiltinImageEncoder option includes: JPEG_FORMAT as default, PNG_FORMAT.
// This works like WithBuiltinImageEncoder, the different between them is
// formatTyp is enumerated in (JPEG_FORMAT, PNG_FORMAT)
func WithBuiltinImageEncoder(format formatTyp) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		var encoder ImageEncoder
		switch format {
		case JPEG_FORMAT:
			encoder = jpegEncoder{}
		case PNG_FORMAT:
			encoder = pngEncoder{}
		default:
			panic("Not supported file format")
		}

		oo.imageEncoder = encoder
	})
}

// WithBuiltinImageEncoder to use custom image encoder to encode image.Image into
// io.Writer
func WithCustomImageEncoder(encoder ImageEncoder) ImageOption {
	return newFuncDialOption(func(oo *outputImageOptions) {
		if encoder == nil {
			return
		}

		oo.imageEncoder = encoder
	})
}

