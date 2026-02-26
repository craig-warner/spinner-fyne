package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const DEBUG = false

const WINDOW_SIZE = 512

const (
	MAX_DISPLAY_SIZE = 10000
)

type CalledColor int

const (
	CCRed    CalledColor = 0
	CCGreen  CalledColor = 1
	CCBlue   CalledColor = 2
	CCYellow CalledColor = 3
	CCBlack  CalledColor = 4
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
    "assets/images/parts/white_left_hand.png",
    "assets/images/parts/white_left_foot.png",
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
	all_colors_str = `[ {red:248, green:9, blue:43}, {red:9, green:19, blue:248}, {red:14, green:248, blue:9}, {red:228 green:248, blue:9} ]`
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
	// Dots
	play_dots []*canvas.Circle
	// Colors
	all_colors []Color
	// Speed
	//speed_setting        int
	//speed_seconds        int
	speed                float64 // 1.0 to 60.0 seconds between spins
	speed_ticks_per_play int
	// Window
	cur_w, cur_h       int
	size               int
	black_out_left     int
	centering_left_adj int
	centering_top_adj  int
	black_out_top      int
	// Color
	cur_color      CalledColor
	cur_part       int
	next_cur_color int // 0-3
	// Canvas
	canvas_size fyne.Size
	// Called Colors
	called_images []*canvas.Image
	// Record Parts
	left_hand_color  CalledColor
	left_foot_color  CalledColor
	right_hand_color CalledColor
	right_foot_color CalledColor
	// Record Dots
	lh_dots []*canvas.Circle
	lf_dots []*canvas.Circle
	rh_dots []*canvas.Circle
	rf_dots []*canvas.Circle
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

/*
func NewColor() Color {
	c := Color{
		Red:   0,
		Green: 0,
		Blue:  0,
	}
	return c
}
*/

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

func (m *Spinner) DrawSpinner(image_num int) {
	//fmt.Printf("Hide %d\n", image_num)
	m.spinner_images[m.cur_spinner_image].Hide()
	//fmt.Printf("Show %d\n", image_num)
	m.spinner_images[image_num].Show()
}

func (m *Spinner) UpdateSpinner() {
	var new_spinner_image_num int
	if m.tick == 48 {
		new_spinner_image_num = 0
	} else {
		if m.spinner_mode == 0 {
			// Still
			new_spinner_image_num = 0
		} else if m.spinner_mode == 1 {
			// Spin
			new_spinner_image_num = m.tick / 10
		} else if m.spinner_mode == 2 {
			// Reverse Spin
			new_spinner_image_num = 47 - (m.tick / 10)
		} else if m.spinner_mode == 3 {
			// x3 Spin
			new_spinner_image_num = ((m.tick / 10 * 3) % 48)
		} else if m.spinner_mode == 4 {
			new_spinner_image_num = (m.spinner_wiggle[(m.tick/10)%48])
		} else {
			panic(1)
		}
	}
	m.DrawSpinner(new_spinner_image_num)
	m.cur_spinner_image = new_spinner_image_num
}

func (m *Spinner) BuildDot(cc CalledColor) *canvas.Circle {
	var nc color.NRGBA
	if cc == CCRed {
		nc = color.NRGBA{m.all_colors[0].Red, m.all_colors[0].Green, m.all_colors[0].Blue, 0xff}
	} else if cc == CCGreen {
		nc = color.NRGBA{m.all_colors[1].Red, m.all_colors[1].Green, m.all_colors[1].Blue, 0xff}
	} else if cc == CCBlue {
		nc = color.NRGBA{m.all_colors[2].Red, m.all_colors[2].Green, m.all_colors[2].Blue, 0xff}
	} else if cc == CCYellow {
		nc = color.NRGBA{m.all_colors[3].Red, m.all_colors[3].Green, m.all_colors[3].Blue, 0xff}
	} else {
		panic(1)
	}
	c := canvas.NewCircle(nc)
	c.Resize(fyne.NewSize(100, 100))
	c.Hide()
	return c
}

func (m *Spinner) ShowDot(play bool, part_num int, color_num int) {
	if play {
		m.play_dots[m.next_cur_color].Show()
	} else {
		if part_num == 0 {
			m.lh_dots[color_num].Show()
		} else if part_num == 1 {
			m.lf_dots[color_num].Show()
		} else if part_num == 2 {
			m.rf_dots[color_num].Show()
		} else if part_num == 3 {
			m.rh_dots[color_num].Show()
		}
	}
}
func (m *Spinner) UpdatePlay() {
	if m.tick == 0 {
		// Hide All
		// Hide Parts
		for part_num := range 4 {
			m.parts_images[part_num].Hide()
		}
		// Hide play Dots
		for dot_num := range 4 {
			// Hide Dot for current color
			m.play_dots[dot_num].Hide()
		}
		m.cur_part = rand.Intn(4)
		m.next_cur_color = rand.Intn(4)
	} else if m.tick == 5*int(m.speed) {
		fmt.Print("Show Part")
		// Show part
		m.parts_images[m.cur_part].Show()
		//m.parts_images[m.cur_part].Refresh()
		//m.parts_images[m.cur_part].Show()
	} else if m.tick == 25*int(m.speed) || (m.tick == 30*int(m.speed)) {
		//fmt.Printf("Show Dot : %d : %d \n", m.cur_part, m.next_cur_color)
		// Draw dot
		fyne.Do(func() { m.ShowDot(true, 0, m.next_cur_color) })
		if m.cur_part == 0 {
			for cdot := range 4 {
				m.lh_dots[cdot].Hide()
			}
			//m.lh_dots[m.next_cur_color].Show()
			fyne.Do(func() { m.ShowDot(false, 0, m.next_cur_color) })
		} else if m.cur_part == 1 {
			for cdot := range 4 {
				m.lf_dots[cdot].Hide()
			}
			//m.lf_dots[m.next_cur_color].Show()
			fyne.Do(func() { m.ShowDot(false, 1, m.next_cur_color) })
		} else if m.cur_part == 2 {
			for cdot := range 4 {
				m.rf_dots[cdot].Hide()
			}
			//m.rf_dots[m.next_cur_color].Show()
			fyne.Do(func() { m.ShowDot(false, 2, m.next_cur_color) })
		} else if m.cur_part == 3 {
			for cdot := range 4 {
				m.rh_dots[cdot].Hide()
			}
			//m.rh_dots[m.next_cur_color].Show()
			fyne.Do(func() { m.ShowDot(false, 3, m.next_cur_color) })
		}
	}
}

func (m *Spinner) ResizeCanvas(s fyne.Size) {
	fmt.Print("Resize")
	m.cur_h = int(s.Height)
	m.cur_w = int(s.Width)
}

func (m *Spinner) UpdateSome(s fyne.Size) {
	if (s.Height != float32(m.cur_h)) || (s.Width != float32(m.cur_w)) {
		m.ResizeCanvas(s)
	}
	//fmt.Print(m.mode)
	m.canvas_size = s
	//fmt.Print(s, "\n")
	if m.mode == 0 {
		//fmt.Print(m.spinner_mode)
		//fmt.Print(m.tick)
		m.UpdateSpinner()
		m.tick = (m.tick + 1) % 480
		if m.tick == 0 {
			//m.spinner_mode = math.rand.Int() % 4
			m.spinner_mode = (m.spinner_mode + 1) % 5
		}
	} else if m.mode == 1 {
		//fmt.Print(m.tick)
		m.UpdatePlay()
		// Tick Update
		m.tick = (m.tick + 1) % m.speed_ticks_per_play
		//if m.tick == 0 {
		//	m.cur_part = (m.cur_part + 1) % 4
		//}
	}
}

func (m *Spinner) GetSpinnerImage(image_num int) *canvas.Image {
	return (m.spinner_images[image_num])
}

func (m *Spinner) DoPlay() {
	//

	m.next_cur_color = rand.Intn(4)
	// Set mode
	m.mode = 1
	// reset timer
	m.tick = 0
	// Hide Banner
	m.banner_image.Hide()
	// Hide Spinner
	for spinner_image_number := range 48 {
		m.spinner_images[spinner_image_number].Hide()
	}
	// Show Called Parts
	for called_part := range 4 {
		m.called_images[called_part].Show()
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
		//		speed_setting:        1,
		//		speed_seconds:        1,   // FIXME: Ignored
		speed_ticks_per_play: 200,
		// Window
		//cur_w, cur_h    int
		// Color
		cur_color: CCBlack,
		cur_part:  0,
		// Record Colors
		left_hand_color:  CCBlack,
		left_foot_color:  CCBlack,
		right_hand_color: CCBlack,
		right_foot_color: CCBlack,
	}
	m.all_colors = append(m.all_colors, NewColor(248, 9, 43))  // Red
	m.all_colors = append(m.all_colors, NewColor(9, 19, 248))  // Blue
	m.all_colors = append(m.all_colors, NewColor(14, 248, 9))  // Green
	m.all_colors = append(m.all_colors, NewColor(228, 248, 9)) // Yellow
	/*
		 Doesn't work for some reason
		var all_colors []Color
		err := json.Unmarshal([]byte(all_colors_str), &all_colors)
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
	for part_num, img_name := range m.parts_image_names {
		m.parts_images = append(m.parts_images, canvas.NewImageFromFile(img_name))
		m.parts_images[part_num].Hide()
		//m.parts_images[part_num].SetMinSize(fyne.NewSize(100, 200))
		m.parts_images[part_num].SetMinSize(fyne.NewSize(200, 200))
		// Called Images
		m.called_images = append(m.called_images, canvas.NewImageFromFile(img_name))
		m.called_images[part_num].Hide()
		//m.called_images[part_num].SetMinSize(fyne.NewSize(25, 50))
		m.called_images[part_num].SetMinSize(fyne.NewSize(50, 50))
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
	// Dots
	for play_dot_num := range 4 {
		if play_dot_num == int(CCRed) {
			m.play_dots = append(m.play_dots, m.BuildDot(CCRed))
			m.lh_dots = append(m.lh_dots, m.BuildDot(CCRed))
			m.lf_dots = append(m.lf_dots, m.BuildDot(CCRed))
			m.rh_dots = append(m.rh_dots, m.BuildDot(CCRed))
			m.rf_dots = append(m.rf_dots, m.BuildDot(CCRed))
		} else if play_dot_num == int(CCGreen) {
			m.play_dots = append(m.play_dots, m.BuildDot(CCGreen))
			m.lh_dots = append(m.lh_dots, m.BuildDot(CCGreen))
			m.lf_dots = append(m.lf_dots, m.BuildDot(CCGreen))
			m.rh_dots = append(m.rh_dots, m.BuildDot(CCGreen))
			m.rf_dots = append(m.rf_dots, m.BuildDot(CCGreen))
		} else if play_dot_num == int(CCBlue) {
			m.play_dots = append(m.play_dots, m.BuildDot(CCBlue))
			m.lh_dots = append(m.lh_dots, m.BuildDot(CCBlue))
			m.lf_dots = append(m.lf_dots, m.BuildDot(CCBlue))
			m.rh_dots = append(m.rh_dots, m.BuildDot(CCBlue))
			m.rf_dots = append(m.rf_dots, m.BuildDot(CCBlue))
		} else if play_dot_num == int(CCYellow) {
			m.play_dots = append(m.play_dots, m.BuildDot(CCYellow))
			m.lh_dots = append(m.lh_dots, m.BuildDot(CCYellow))
			m.lf_dots = append(m.lf_dots, m.BuildDot(CCYellow))
			m.rh_dots = append(m.rh_dots, m.BuildDot(CCYellow))
			m.rf_dots = append(m.rf_dots, m.BuildDot(CCYellow))
		} else {
			panic(1)
		}
	}
	// speed
	m.speed = 5.0 // 5 seconds per spin
	m.speed_ticks_per_play = int(m.speed) * 100
	return m
}

/*
 * Main
 */

func main() {

	spinnerContent := container.New(layout.NewHBoxLayout())
	bannerContent := container.New(layout.NewHBoxLayout())
	colOneContent := container.New(layout.NewVBoxLayout())
	play_stack := container.New(layout.NewStackLayout())
	calledHBox := container.New(layout.NewHBoxLayout())
	lhstack := container.New(layout.NewStackLayout())
	lfstack := container.New(layout.NewStackLayout())
	rhstack := container.New(layout.NewStackLayout())
	rfstack := container.New(layout.NewStackLayout())
	play_rect := canvas.NewRectangle(color.Black)
	play_vbox := container.New(layout.NewVBoxLayout())

	// Spinner
	mySpinner := NewSpinner()

	// Parts
	/*
		for part_num := range 4 {
			play_stack.Add(mySpinner.parts_images[part_num])
		}
	*/

	myApp := app.New()
	myWindow := myApp.NewWindow("Spinner")
	myWindow.SetPadded(false)
	// Center
	myWindow.CenterOnScreen()

	/*
		myPlayWindow := myApp.NewWindow("Playing Twister")
		myPlayWindow.CenterOnScreen()
		myPlayWindow.Hide()
	*/

	// Resize ignored by Mobile Platforms
	// - Mobile platforms are always full screen
	// - 27 is a hack determined by Ubuntu/Gnome
	//myWindow.Resize(fyne.NewSize(256, (256 + 27)))
	myWindow.Resize(fyne.NewSize(WINDOW_SIZE, (WINDOW_SIZE + 27)))
	//myPlayWindow.Resize(fyne.NewSize(WINDOW_SIZE, (WINDOW_SIZE + 27)))

	// Control Menu Set up
	//	menuItemGenerate := fyne.NewMenuItem("Generate Background", func() {
	//		fmt.Println("In Generate Background")
	//	})
	menuItemQuit := fyne.NewMenuItem("Quit", func() {
		//fmt.Println("In DoQuit:")
		os.Exit(0)
	})
	menuItemPlay := fyne.NewMenuItem("Play", func() {
		/*
			myPlayWindow := myApp.NewWindow("Playing Twister")
			myPlayWindow.CenterOnScreen()
			myPlayWindow.Show()
			myPlayWindow.Resize(fyne.NewSize(WINDOW_SIZE, (WINDOW_SIZE + 27)))
		*/
		// Play
		//play_rect := canvas.NewRectangle(color.Black)
		//play_rect.SetMinSize(fyne.NewSize(WINDOW_SIZE, WINDOW_SIZE))
		//play_rect.Show()
		/*
			play_stack := container.New(layout.NewStackLayout())
			//play_stack.Add(play_rect)
			part_stack := container.New(layout.NewStackLayout())
			part_vbox := container.New(layout.NewVBoxLayout())
			part_hbox := container.New(layout.NewHBoxLayout())
			// Dots
			for dot_num := range 10 {
				play_stack.Add(mySpinner.dots[dot_num])
			}
			// Parts
			for part_num := range 4 {
				play_stack.Add(mySpinner.parts_images[part_num])
			}
			part_hbox.Add(layout.NewSpacer())
			part_hbox.Add(part_stack)
			part_hbox.Add(layout.NewSpacer())
			part_vbox.Add(layout.NewSpacer())
			part_vbox.Add(part_hbox)
			part_vbox.Add(layout.NewSpacer())
			play_stack.Add(part_vbox)
		*/
		play_rect.Show()
		mySpinner.DoPlay()
		/*
			go func() {
				for {
					mySpinner.UpdateSome()
					time.Sleep(time.Nanosecond * 100000000)
				}
			}()
			myPlayWindow.ShowAndRun()
		*/
	})
	menuItemSpeedSettings := fyne.NewMenuItem("Speed Settings", func() {
		var popup *widget.PopUp
		new_speed := float64(mySpinner.speed)
		speedSettingsText := widget.NewLabel("Speed Settings (1 to 60 seconds per spin)")
		speedSettingsSlider := widget.NewSlider(1.0, 60.0)
		speedSettingsSlider.SetValue(float64(mySpinner.speed))
		speedSettingsSlider.OnChanged = func(s float64) {
			//cp.DbgPrint("Speed Settings Callback:", s)
			new_speed = s
		}
		popUpContent := container.NewVBox(
			speedSettingsText,
			speedSettingsSlider,
			container.NewHBox(
				layout.NewSpacer(),
				widget.NewButton("Ok", func() {
					mySpinner.speed = new_speed
					mySpinner.speed_ticks_per_play = 100 * int(new_speed)
					// fine grain pan
					popup.Hide()
				}),
				widget.NewButton("Cancel", func() {
					popup.Hide()
				}),
				layout.NewSpacer(),
			),
		)
		popup = widget.NewModalPopUp(popUpContent, myWindow.Canvas())
		popup.Show()
	})
	//	menuControl:= fyne.NewMenu("Control", menuItemColor, menuItemZoom, menuItemQuit);
	//menuControl := fyne.NewMenu("Control", menuItemGenerate, menuItemQuit)
	menuControl := fyne.NewMenu("Control", menuItemPlay, menuItemSpeedSettings, menuItemQuit)
	// About Menu Set up
	menuItemAbout := fyne.NewMenuItem("About...", func() {
		dialog.ShowInformation("About Spinner v1.0.0", "Author: Craig Warner \n\ngithub.com/craig-warner/spinner-fyne", myWindow)
	})
	menuHelp := fyne.NewMenu("Help ", menuItemAbout)
	mainMenu := fyne.NewMainMenu(menuControl, menuHelp)
	myWindow.SetMainMenu(mainMenu)

	// Raster
	//myRaster := canvas.NewRasterWithPixels(mySpinner.DrawOneDot)
	//colOneContent.Add(myRaster)
	bigger_stack := container.New(layout.NewStackLayout())
	big_stack := container.New(layout.NewStackLayout())
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(WINDOW_SIZE, WINDOW_SIZE))
	rect.Show()
	big_stack.Add(rect)
	bigger_stack.Add(rect)

	var image_holder *canvas.Image
	//var image_holders []*canvas.Image
	stack := container.New(layout.NewStackLayout())
	for image_num := range 48 {
		image_holder = mySpinner.GetSpinnerImage(image_num)
		image_holder.SetMinSize(fyne.NewSize(350, 100))
		//image_holder.Move(fyne.NewPos(0.0, 0.0))
		if image_num == 0 {
			image_holder.Show()
		} else {
			image_holder.Hide()
		}
		//image_holders = append(image_holders, image_holder)
		stack.Add(image_holder)
	}
	spinnerContent.Add(layout.NewSpacer())
	spinnerContent.Add(stack)
	spinnerContent.Add(layout.NewSpacer())

	image_holder = canvas.NewImageFromFile("assets/images/banner/touchtostart.png")
	//image_holder.FillMode = canvas.ImageFillContain // Does not work
	image_holder.SetMinSize(fyne.NewSize(270, 230))
	image_holder.Show()
	mySpinner.banner_image = image_holder
	// Play
	play_rect.SetMinSize(fyne.NewSize(WINDOW_SIZE, WINDOW_SIZE))
	play_rect.Hide()
	play_stack.Add(play_rect)
	part_stack := container.New(layout.NewStackLayout())
	part_vbox := container.New(layout.NewVBoxLayout())
	part_hbox := container.New(layout.NewHBoxLayout())
	// Phay Dots
	for dot_num := range 4 {
		play_stack.Add(mySpinner.play_dots[dot_num])
	}
	// Parts
	for part_num := range 4 {
		part_stack.Add(mySpinner.parts_images[part_num])
	}
	// Called HBox
	// Called Stacks
	// lh
	for lh_dot_num := range 4 {
		lhstack.Add(mySpinner.lh_dots[lh_dot_num])
	}
	lhstack.Add(mySpinner.called_images[0])
	// lf
	for lf_dot_num := range 4 {
		lfstack.Add(mySpinner.lf_dots[lf_dot_num])
	}
	lfstack.Add(mySpinner.called_images[1])
	// rf
	for rf_dot_num := range 4 {
		rfstack.Add(mySpinner.rf_dots[rf_dot_num])
	}
	rfstack.Add(mySpinner.called_images[2])
	// rh
	for rh_dot_num := range 4 {
		rhstack.Add(mySpinner.rh_dots[rh_dot_num])
	}
	rhstack.Add(mySpinner.called_images[3])
	// place record stacks
	calledHBox.Add(lhstack)
	calledHBox.Add(lfstack)
	calledHBox.Add(layout.NewSpacer())
	calledHBox.Add(rfstack)
	calledHBox.Add(rhstack)
	//
	part_hbox.Add(layout.NewSpacer())
	part_hbox.Add(part_stack)
	part_hbox.Add(layout.NewSpacer())
	part_vbox.Add(layout.NewSpacer())
	part_vbox.Add(part_hbox)
	part_vbox.Add(layout.NewSpacer())
	play_stack.Add(part_vbox)
	play_vbox.Add(play_stack)
	play_vbox.Add(calledHBox)
	//fyne_size := fyne.NewSize(10.0, 10.0)
	//image_holder.Resize(fyne_size)
	bannerContent.Add(layout.NewSpacer())
	bannerContent.Add(image_holder)
	bannerContent.Add(layout.NewSpacer())

	//SomeText := widget.NewLabel("Some Text")
	//colOneContent.Add(SomeText)
	colOneContent.Add(layout.NewSpacer())
	colOneContent.Add(spinnerContent)
	colOneContent.Add(bannerContent)
	colOneContent.Add(layout.NewSpacer())

	topContent := container.New(layout.NewHBoxLayout())
	topContent.Add(layout.NewSpacer())
	topContent.Add(colOneContent)
	topContent.Add(layout.NewSpacer())

	big_stack.Add(topContent)
	// Parts
	for part_num := range 4 {
		big_stack.Add(mySpinner.parts_images[part_num])
	}
	big_stack.Add(play_vbox)
	bigger_stack.Add(big_stack)

	wholeContent := container.New(layout.NewVBoxLayout())
	//wholeContent.Add(layout.NewSpacer())
	wholeContent.Add(bigger_stack)
	//wholeContent.Add(layout.NewSpacer())
	//wholeContent.Add(bottomContent)

	myWindow.SetContent(wholeContent)
	//myPlayWindow.SetContent(play_stack)

	var canvas_size fyne.Size
	go func() {
		for {
			canvas_size = myWindow.Canvas().Size()
			mySpinner.UpdateSome(canvas_size)
			time.Sleep(time.Nanosecond * 10000000)
			fyne.Do(func() {
				big_stack.Refresh()
			})
		}
	}()
	myWindow.ShowAndRun()
}
