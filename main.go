package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
)

const DEBUG = false

const WINDOW_SIZE = 512

const (
	MAX_DISPLAY_SIZE = 10000
)

/*
JSON Structure Defining Spinner Image Names
*/
const (
	all_spinner_str = `[ "assets/images/spinner/spinner_3d0000.png",
    "assets/images/spinner/spinner_3d0001.png",
    "assets/images/spinner/spinner_3d0002.png",
    "assets/images/spinner/spinner_3d0003.png",
    "assets/images/spinner/spinner_3d0004.png",
    "assets/images/spinner/spinner_3d0005.png",
    "assets/images/spinner/spinner_3d0006.png",
    "assets/images/spinner/spinner_3d0007.png",
    "assets/images/spinner/spinner_3d0008.png",
    "assets/images/spinner/spinner_3d0009.png",
    "assets/images/spinner/spinner_3d0010.png",
    "assets/images/spinner/spinner_3d0011.png",
    "assets/images/spinner/spinner_3d0012.png",
    "assets/images/spinner/spinner_3d0013.png",
    "assets/images/spinner/spinner_3d0014.png",
    "assets/images/spinner/spinner_3d0015.png",
    "assets/images/spinner/spinner_3d0016.png",
    "assets/images/spinner/spinner_3d0017.png",
    "assets/images/spinner/spinner_3d0018.png",
    "assets/images/spinner/spinner_3d0019.png",
    "assets/images/spinner/spinner_3d0020.png",
    "assets/images/spinner/spinner_3d0021.png",
    "assets/images/spinner/spinner_3d0022.png",
    "assets/images/spinner/spinner_3d0023.png",
    "assets/images/spinner/spinner_3d0024.png",
    "assets/images/spinner/spinner_3d0025.png",
    "assets/images/spinner/spinner_3d0026.png",
    "assets/images/spinner/spinner_3d0027.png",
    "assets/images/spinner/spinner_3d0028.png",
    "assets/images/spinner/spinner_3d0029.png",
    "assets/images/spinner/spinner_3d0030.png",
    "assets/images/spinner/spinner_3d0031.png",
    "assets/images/spinner/spinner_3d0032.png",
    "assets/images/spinner/spinner_3d0033.png",
    "assets/images/spinner/spinner_3d0034.png",
    "assets/images/spinner/spinner_3d0035.png",
    "assets/images/spinner/spinner_3d0036.png",
    "assets/images/spinner/spinner_3d0037.png",
    "assets/images/spinner/spinner_3d0028.png",
    "assets/images/spinner/spinner_3d0039.png",
    "assets/images/spinner/spinner_3d0040.png",
    "assets/images/spinner/spinner_3d0041.png",
    "assets/images/spinner/spinner_3d0042.png",
    "assets/images/spinner/spinner_3d0043.png",
    "assets/images/spinner/spinner_3d0044.png",
    "assets/images/spinner/spinner_3d0045.png",
    "assets/images/spinner/spinner_3d0046.png",
    "assets/images/spinner/spinner_3d0047.png" ]`

	all_parts_str = `[
    "assets/images/parts/white_left_foot.png",
    "assets/images/parts/white_left_hand.png",
    "assets/images/parts/white_right_foot.png",
    "assets/images/parts/white_right_hand.png"
	]`
	all_banner_str = `[
    "assets/images/banner/touchtostart.png"
	]`
)

/*
JSON Structure Defining Image Sizes

	# Part Type 4 is Spinner
	 self.partTypeWidth = [250, 250, 300, 300, 350]
	 self.partTypeHeight = [300, 300, 400, 400, 250]
*/
const (
	all_image_width_str  = `[ 250, 250, 300, 300 ]`
	all_image_height_str = `[300, 300, 400, 400]`
)

/*
JSON Structure Defining Wiggle
*/
const (
	all_spinner_wiggle_str = `[
        0, 47, 46, 45, 44, 43, 42, 41, 
		40, 40, 41, 42, 43, 44, 45, 46,
		47, 0, 1, 2, 3, 4, 5, 6, 7, 8, 
		8, 7, 6, 5, 4, 3, 2, 1, 47, 46, 
		45, 44, 43, 42, 41, 40, 41, 42, 
		43, 44, 45, 46, 47, 0
	]`
)

/*
JSON Structure Defining Colors

	self.curColorRGB = ["rgb(248,9,43)", "rgb(9,19,248)", "rgb(14,248,9)", "rgb(228,248,9)"]
*/
const (
	all_colors_str = `[
		{red:248, green: 9, blue:43},
		{red:9, green:19, 248, blue: 248},
		{red:14, green: 248, blue:9},
		{red:228 green:248, blue: 9}
	]`
)

const (
	all_speed_ctl_str = `[ 1, 2, 5, 10]`
)

func DbgPrint(str ...interface{}) {
	if DEBUG {
		fmt.Println(str...)
		return
	}
}

/*
type tappableRaster struct {
	fyne.CanvasObject
	OnTapped func()
}

func NewTappableRaster(raster fyne.CanvasObject, onTapped func()) *tappableRaster {
	return &tappableRaster{CanvasObject: raster, OnTapped: onTapped}
}
*/

type Spinner struct {
	mode int // 0 intro screen
	// 1 playing
	spinner_mode int // 0-4
	tick         int // 0-48
	// Spinner
	spinner_wiggle      []int
	cur_spinner_image   int
	spinner_image_x     int
	spinner_image_y     int
	spinner_image_names []string
	spinner_images      []*canvas.Image
	// Banner
	banner_image_x int
	banner_image_y int
	banner_image   *canvas.Image
	// Parts
	parts_image_x     []int
	parts_image_y     []int
	parts_image_names []string
	parts_images      []*canvas.Image
	// Colors
	all_colors []Color
	// Speed
	speed_setting        int
	speed_seconds        int
	speed_ticks_per_play int
	// Window
	cur_w, cur_h       int
	size               int
	black_out_left     int
	centering_left_adj int
	centering_top_adj  int
	black_out_top      int
	// Color
	cur_color_num int
	cur_part      int
}

type Color struct {
	Red   uint8 `json:"red"`
	Green uint8 `json:"green"`
	Blue  uint8 `json:"blue"`
}

type Point struct {
	x float64
	y float64
}

/*
 * General Functions
 */

func NewColor(r, g, b uint8) Color {
	c := Color{
		Red:   r,
		Green: g,
		Blue:  b,
	}
	return c
}

func NewPoint(set_x, set_y float64) Point {
	p := Point{
		x: set_x,
		y: set_y,
	}
	return p
}

/*
 * Mandel functions
 */

/*
func (m *Spinner) CalcOneDot() {
	var p Point

	realx := m.min_x + float64(m.cur_x)*m.span_one_dot
	realy := m.min_y + m.span - float64(m.cur_y)*m.span_one_dot

	p = NewPoint(realx, realy)

	color := m.CalcOnePointColor(p)

	m.tiles[m.cur_x][m.cur_y].red = color.red
	m.tiles[m.cur_x][m.cur_y].green = color.green
	m.tiles[m.cur_x][m.cur_y].blue = color.blue
}
*/

func (m *Spinner) ResetWindow(w, h int) {
	// Check
	if (w > MAX_DISPLAY_SIZE) || (h > MAX_DISPLAY_SIZE) {
		fmt.Println("Monitor is too big")
		panic(1)
	}
	// New Window Size
	m.cur_w = w
	m.cur_h = h
	// New Mandelbrot Size
	max_val := 0
	min_val := 0
	// Choose the smallest so it looks okay on Mobile platform
	if w > h {
		max_val = w
		min_val = h
	} else {
		max_val = h
		min_val = w
	}
	//max_mult64 := (max_val / 64) * 64
	min_mult64 := (min_val / 64) * 64
	// Blackout and Center
	if max_val == w {
		// wider than tall
		m.black_out_left = (w - m.size) >> 1
		m.centering_left_adj = 0
		m.centering_top_adj = (m.size - min_mult64) >> 1
		m.black_out_top = (h - min_mult64) >> 1
	} else {
		// taller than wide
		m.black_out_left = (w - min_mult64) >> 1
		m.centering_left_adj = (m.size - min_mult64) >> 1
		m.centering_top_adj = 0
		m.black_out_top = (h - m.size) >> 1
	}
}

/*
func (m *Spinner) DrawOneDot(px, py, w, h int) color.Color {
	use_px := 0
	use_py := 0
	use_px = px
	use_py = py
	if (w != m.cur_w) || (h != m.cur_h) {
		m.ResetWindow(w, h)
	}

	// color_px
	color_px := use_px - m.black_out_left + m.centering_left_adj
	color_py := use_py - m.black_out_top + m.centering_top_adj

	//fmt.Printf("DrawOne px:%d,py:%d,w:%d,h:%d", px, py, w, h)
	// Black out or color
	black_color := color.RGBA{0, 0, 0, 0xff}
	if use_py < m.black_out_top {
		// Top
		return (black_color)
	} else if use_px < m.black_out_left {
		// Left
		return (black_color)
	} else if use_py >= (m.black_out_left + m.size) {
		// Right
		return (black_color)
	} else {
		return (m.DrawOneDotNotBlack(color_px, color_py))
	}
}
func (m *Spinner) DrawOneDotNotBlack(use_px, use_py int) color.Color {
	//fmt.Println("px:",px,"py:",py,"w:",w,"h:",h)
	idx_x := 0
	idx_y := 0
	gran := 64
	ret_red := uint8(m.tiles[idx_x][idx_y].red)
	ret_green := uint8(m.tiles[idx_x][idx_y].green)
	ret_blue := uint8(m.tiles[idx_x][idx_y].blue)
	ret_color := color.RGBA{ret_red, ret_green, ret_blue, 0xff}
	return (ret_color)
}
*/
/*
func (m *Spinner) DrawSpinner(image_num int) {
	m.spinner_images[m.cur_spinner_image].Hide()
	m.spinner_images[image_num].Show()
}

func (m *Spinner) UpdateSpinner() {
	if m.spinner_mode == 0 {
		// Still
		m.DrawSpinner(0)
	} else if m.spinner_mode == 1 {
		// Spin
		m.DrawSpinner(m.tick)
	} else if m.spinner_mode == 2 {
		// Reverse Spin
		m.DrawSpinner(47 - m.tick)
	} else if m.spinner_mode == 3 {
		// x3 Spin
		m.DrawSpinner((m.tick * 3) % 48)
	} else if m.spinner_mode == 4 {
		m.DrawSpinner(m.spinner_wiggle[m.tick])
	} else {
		panic(1)
	}
}
*/

func (m *Spinner) UpdateSome() {
	if m.mode == 0 {
		//m.UpdateSpinner()
		m.tick = (m.tick + 1) % 48
		if m.tick == 0 {
			//m.spinner_mode = math.rand.Int() % 4
			m.spinner_mode = (m.spinner_mode + 1) % 4
		}
	} else if m.mode == 1 {
		//m.UpdatePlay()
		// Tick Update
		m.tick = (m.tick + 1) % m.speed_ticks_per_play
	}
}

func NewSpinner() Spinner {
	m := Spinner{
		mode:         0,
		spinner_mode: 0,
		tick:         0,
		// Spinner
		//spinner_wiggle	[]int
		cur_spinner_image: 0,
		//spinner_image_x : 0,
		//spinner_image_y : 0,
		//spinner_image_names	[]string
		//spinner_images	[]image.Image
		// Banner
		//banner_image_x:0
		//banner_image_y:0
		//banner_image	image.Image
		// Parts
		//parts_image_x	[]int
		//parts_image_y	[]int
		//parts_image_names	[]string
		//parts_images	[]image.Image
		// Colors
		//all_colors	    []Color
		// Speed
		speed_setting:        1,
		speed_seconds:        1,   // FIXME: Ignored
		speed_ticks_per_play: 100, // FIXME: Make configurable
		// Window
		//cur_w, cur_h    int
		// Color
		cur_color_num: 0,
		cur_part:      0,
	}
	/*
		err := json.Unmarshal([]byte(all_colors_str), &m.all_colors)
		if err != nil {
			fmt.Printf("Unable to marshal JSON due to %s", err)
			panic(1)
		}
	*/
	// Banner
	image_holder := canvas.NewImageFromFile("assets/images/banner/touchtostart.png")
	m.banner_image = image_holder
	//m.banner_image = canvas.NewImageFromFile("assets/images/banner/touchtostart.png")
	m.banner_image_x = 0 // FIXME
	m.banner_image_y = 0 // FIXME
	// Spinner
	var image_names []string
	err := json.Unmarshal([]byte(all_spinner_str), &image_names)
	if err != nil {
		fmt.Printf("Unable to marshal JSON due to %s", err)
		panic(1)
	}
	//DBGfmt.Print(image_names)
	for _, img_name := range image_names {
		//DBGfmt.Print(idx)
		//DBGfmt.Print(img_name)
		image_holder = canvas.NewImageFromFile(img_name)
		m.spinner_images = append(m.spinner_images, image_holder)
	}
	m.spinner_image_x = 350
	m.spinner_image_y = 250
	err = json.Unmarshal([]byte(all_spinner_wiggle_str), &m.spinner_wiggle)
	if err != nil {
		fmt.Printf("Unable to marshal JSON due to %s", err)
		panic(1)
	}
	// Parts
	err = json.Unmarshal([]byte(all_parts_str), &m.parts_image_names)
	if err != nil {
		fmt.Printf("Unable to marshal JSON due to %s", err)
		panic(1)
	}
	for _, img_name := range m.parts_image_names {
		m.parts_images = append(m.parts_images, canvas.NewImageFromFile(img_name))
	}
	err = json.Unmarshal([]byte(all_image_width_str), &m.parts_image_x)
	if err != nil {
		fmt.Printf("Unable to marshal JSON due to %s", err)
		panic(1)
	}
	err = json.Unmarshal([]byte(all_image_height_str), &m.parts_image_y)
	if err != nil {
		fmt.Printf("Unable to marshal JSON due to %s", err)
	}
	return m
}

/*
 * Main
 */

func main() {

	colOneContent := container.New(layout.NewVBoxLayout())

	myApp := app.New()
	myWindow := myApp.NewWindow("Spinner")
	myWindow.SetPadded(false)

	// Resize ignored by Mobile Platforms
	// - Mobile platforms are always full screen
	// - 27 is a hack determined by Ubuntu/Gnome
	//myWindow.Resize(fyne.NewSize(256, (256 + 27)))
	myWindow.Resize(fyne.NewSize(WINDOW_SIZE, (WINDOW_SIZE + 27)))

	// Control Menu Set up
	//	menuItemGenerate := fyne.NewMenuItem("Generate Background", func() {
	//		fmt.Println("In Generate Background")
	//	})
	menuItemQuit := fyne.NewMenuItem("Quit", func() {
		//fmt.Println("In DoQuit:")
		os.Exit(0)
	})
	//	menuControl:= fyne.NewMenu("Control", menuItemColor, menuItemZoom, menuItemQuit);
	//menuControl := fyne.NewMenu("Control", menuItemGenerate, menuItemQuit)
	menuControl := fyne.NewMenu("Control", menuItemQuit)
	// About Menu Set up
	menuItemAbout := fyne.NewMenuItem("About...", func() {
		dialog.ShowInformation("About Spinner v1.0.0", "Author: Craig Warner \n\ngithub.com/craig-warner/spinner-fyne", myWindow)
	})
	menuHelp := fyne.NewMenu("Help ", menuItemAbout)
	mainMenu := fyne.NewMainMenu(menuControl, menuHelp)
	myWindow.SetMainMenu(mainMenu)

	// Mandelbrot
	mySpinner := NewSpinner()
	// Raster
	//myRaster := canvas.NewRasterWithPixels(mySpinner.DrawOneDot)

	topContent := container.New(layout.NewHBoxLayout())
	topContent.Add(layout.NewSpacer())
	topContent.Add(colOneContent)
	topContent.Add(layout.NewSpacer())

	wholeContent := container.New(layout.NewVBoxLayout())
	wholeContent.Add(layout.NewSpacer())
	wholeContent.Add(topContent)
	wholeContent.Add(layout.NewSpacer())
	//wholeContent.Add(bottomContent)

	myWindow.SetContent(wholeContent)

	go func() {
		for {
			mySpinner.UpdateSome()
			//myRaster.Refresh()
			time.Sleep(time.Nanosecond * 100000000)
		}
	}()

	myWindow.ShowAndRun()
}
