package tui

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"strings"
)

// IsImageFile 判斷指定 Content-Type 或 Key 是否為常見圖片檔
func IsImageFile(contentType, key string) bool {
	ct := strings.ToLower(strings.TrimSpace(contentType))
	if strings.HasPrefix(ct, "image/") {
		return true
	}

	ext := strings.ToLower(filepath.Ext(key))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp", ".ico":
		return true
	}

	return false
}

// RenderImageToTview 將圖片 Data 解碼並轉換為適合 tview (TextView/Modal) 的 half-block ANSI 色彩字串
func RenderImageToTview(imgData []byte, maxCols, maxRows int) (string, error) {
	if len(imgData) == 0 {
		return "", fmt.Errorf("image data is empty")
	}

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return "", fmt.Errorf("decoding image: %w", err)
	}

	bounds := img.Bounds()
	origW := bounds.Dx()
	origH := bounds.Dy()

	if origW == 0 || origH == 0 {
		return "", fmt.Errorf("invalid image dimensions")
	}

	if maxCols <= 0 {
		maxCols = 45
	}
	if maxRows <= 0 {
		maxRows = 15
	}

	// 每個 terminal 字元高度約為寬度的 2 倍，故 targetHeight 像素數以 maxRows * 2 為上限
	targetW := maxCols
	targetH := maxRows * 2

	// 等比例縮放計算
	scaleW := float64(targetW) / float64(origW)
	scaleH := float64(targetH) / float64(origH)

	scale := scaleW
	if scaleH < scale {
		scale = scaleH
	}

	newW := int(float64(origW) * scale)
	newH := int(float64(origH) * scale)

	if newW < 1 {
		newW = 1
	}
	if newH < 1 {
		newH = 1
	}
	// 確保 height 為偶數以利 half-block 成對
	if newH%2 != 0 {
		newH++
	}

	dst := scaleImageNearest(img, newW, newH)

	var sb strings.Builder
	for y := 0; y < newH; y += 2 {
		for x := 0; x < newW; x++ {
			cTop := dst.RGBAAt(x, y)
			cBot := dst.RGBAAt(x, y+1)

			// 處理透明度
			r1, g1, b1 := cTop.R, cTop.G, cTop.B
			if cTop.A < 128 {
				r1, g1, b1 = 0, 0, 0
			}
			r2, g2, b2 := cBot.R, cBot.G, cBot.B
			if cBot.A < 128 {
				r2, g2, b2 = 0, 0, 0
			}

			fmt.Fprintf(&sb, "[#%02x%02x%02x:#%02x%02x%02x]▀", r1, g1, b1, r2, g2, b2)
		}
		sb.WriteString("[-]\n")
	}

	return sb.String(), nil
}

func scaleImageNearest(img image.Image, newW, newH int) *image.RGBA {
	bounds := img.Bounds()
	origW := bounds.Dx()
	origH := bounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))

	for y := 0; y < newH; y++ {
		origY := bounds.Min.Y + (y * origH / newH)
		for x := 0; x < newW; x++ {
			origX := bounds.Min.X + (x * origW / newW)
			r, g, b, a := img.At(origX, origY).RGBA()
			dst.SetRGBA(x, y, color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			})
		}
	}
	return dst
}
