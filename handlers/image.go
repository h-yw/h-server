package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"os"

	// "image/draw"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/image/colornames"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type GenSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}
type GenQuery struct {
	Title string  `json:"title"`
	Logo  string  `json:"logo"`
	Date  string  `json:"date"`
	Size  GenSize `json:"size"`
	Bg    string  `json:"string"`
}

type GenImage struct {
	Params GenQuery `json:"params"`
}

func NewGenImage() *GenImage {
	return &GenImage{}
}

func ImageHandler(c *gin.Context) {
	gen := NewGenImage()
	gen.getQuery(c.Request.URL.Query())
	fmt.Println("json=====>", gen.Params)
	img := gen.createImage()
	c.Header("Content-Type", "image/png")
	err := png.Encode(c.Writer, img)
	if err != nil {
		c.JSON(201, gin.H{
			"error": "图片生成异常",
		})
	}
	// c.JSON(200, gin.H{
	// 	"title": "Home",
	// 	"query": query,
	// })
}

func (gen *GenImage) getQuery(val url.Values) GenQuery {
	// var query GenQuery
	var res GenQuery
	jsonStr := val.Get("jsonStr")
	err := json.Unmarshal([]byte(jsonStr), &res)
	if err != nil {
		fmt.Println("jsoonStr error", err)
	}
	// fmt.Println("jsonStr====>", res.Logo)
	// if ; len(jsonStr) > 0 {

	// }
	// if title := val.Get("title"); len(title) > 0 {
	// 	gen.Params.Title = title
	// }
	// if logo := val.Get("logo"); len(logo) > 0 {
	// 	gen.Params.Logo = logo
	// }
	// if date := val.Get("date"); len(date) > 0 {
	// 	gen.Params.Date = date
	// }
	// if sizeStr := val.Get("size"); len(sizeStr) > 0 {
	// 	sizeArr := strings.Split(sizeStr, "x")
	// 	defaultSize := GenSize{
	// 		Width:  0,
	// 		Height: 0,
	// 	}
	// 	if len(sizeArr) != 2 {
	// 		gen.Params.Size = defaultSize
	// 	} else {
	// 		w, werr := strconv.Atoi(sizeArr[0])
	// 		h, herr := strconv.Atoi(sizeArr[1])
	// 		fmt.Println("size===>", w, h)
	// 		if werr == nil && herr == nil {
	// 			gen.Params.Size = GenSize{
	// 				Width:  w,
	// 				Height: h,
	// 			}
	// 		}
	// 	}
	// }
	// if bg := val.Get("bg"); len(bg) > 0 {
	// 	// val,_:=parseBackground(bg)
	// }
	// if logo := val.Get("logo"); isValidURL(logo) {
	// 	gen.Params.Logo = logo
	// }
	gen.Params = res
	return res
}

func (gen *GenImage) createImage() *image.RGBA {
	params := gen.Params
	fmt.Println("text====>", params.Title)
	// gen.params.Date
	file, _ := os.Open("assets/image/twitter.png")
	defer file.Close()
	twitter, _, _ := image.Decode(file)
	// bgColor := color.RGBA{
	// 	R: 255,
	// 	G: 255,
	// 	B: 255,
	// 	A: 255,
	// }
	bg := image.NewRGBA(image.Rect(
		0, 0,
		gen.Params.Size.Width, gen.Params.Size.Height))
	draw.Draw(bg, bg.Bounds(), twitter, image.Point{}, draw.Over)

	// netImg, err := downloadImg(gen.Params.Logo)
	// if err != nil {
	// 	panic("error====>")
	// }
	// dstRect := image.Rect(16, 240, 64+16, 240+64)
	// scaleImg := resizeImage(netImg, 64, 64)
	// draw.Draw(bg, dstRect, scaleImg, image.Point{}, draw.Over)

	maxWidth := 400
	// 计算文字的起始绘制位置，使文字居中
	// assets\fonts\LXGWWenKaiMono-Regular.ttf
	face, err := loadFont("assets/fonts/LXGWWenKaiMono-Regular.ttf", 32)
	if err != nil {
		fmt.Println("laod font error")
	}
	// textHeight := len(wordWrap(gen.Params.Title, face, maxWidth)) * face.Metrics().Height.Ceil()
	startX := (params.Size.Width - maxWidth) / 2
	startY := 100 //(params.Size.Height-textHeight)/2 + face.Metrics().Ascent.Ceil()
	drawText(bg, params.Title, startX, startY, maxWidth, colornames.Black, face)
	drawText(bg, params.Title, startX+1, startY+1, maxWidth, color.RGBA{0, 0, 0, 100}, face)
	drawText(bg, params.Title, startX, startY+1, maxWidth, color.RGBA{0, 0, 0, 100}, face)
	face2, _ := loadFont("assets/fonts/LXGWBright-Italic.ttf", 24)
	sx := params.Size.Width - face.Metrics().Ascent.Ceil()*len(params.Date)
	sy := params.Size.Height - face2.Metrics().Height.Ceil()
	drawText(bg, params.Date, sx, sy, maxWidth, color.RGBA{0, 0, 0, 200}, face2)
	// 边框
	return bg
}

func parseRGB(rgbStr string) ([]int, error) {
	parts := strings.Split(rgbStr, ",")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid RGB format")
	}
	r, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid R value: %v", err)
	}

	g, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid G value: %v", err)
	}

	b, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid B value: %v", err)
	}

	return []int{r, g, b}, nil
}

// 解析 bg 参数
func parseBackground(bg string) (interface{}, error) {
	if isValid := isValidURL(bg); !isValid {
		// 处理图片链接
		return bg, nil
	} else {
		// 处理 RGB 字符串
		return parseRGB(bg)
	}
}

// 校验 URL 是否合法
func isValidURL(rawURL string) bool {
	_, err := url.ParseRequestURI(rawURL)
	return err == nil
}

func downloadImg(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	return img, nil
}

// 调整图片大小
func resizeImage(src image.Image, newWidth, newHeight int) image.Image {
	// 创建目标图像
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// 使用 `draw.CatmullRom` 缩放图像
	draw.CatmullRom.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	return dst
}

// 计算文字宽度并进行自动换行
func wordWrap(text string, face font.Face, maxWidth int) []string {
	var lines []string
	var currentLine strings.Builder
	var currentWidth fixed.Int26_6

	for _, word := range strings.Split(text, "") {
		wordWidth := font.MeasureString(face, word)
		if currentWidth+wordWidth >= fixed.I(maxWidth) {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentWidth = 0
		}
		// if currentLine.Len() > 0 {
		// 	spaceWidth := font.MeasureString(face, " ")
		// 	currentLine.WriteString(" ")
		// 	currentWidth += spaceWidth
		// }
		currentLine.WriteString(word)
		currentWidth += wordWidth
	}
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}
	return lines
}

// 在图片中绘制文字
func drawText(img *image.RGBA, text string, x, y int, maxWidth int, col color.Color, face font.Face) {
	lines := wordWrap(text, face, maxWidth)

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	for _, line := range lines {
		lineWidth := d.MeasureString(line).Ceil()
		// 计算该行居中对齐的 X 坐标
		centeredX := x + (maxWidth-lineWidth)/2
		// 设置当前行的 X 坐标
		d.Dot.X = fixed.I(centeredX)
		d.DrawString(line)
		d.Dot.Y += face.Metrics().Height
	}
}

// 读取并解析 TTF 字体文件
func loadFont(path string, size float64) (font.Face, error) {
	fontData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	f, err := opentype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}

	return face, nil
}
