package ceftools

import (
	"github.com/nsf/termbox-go"
	"os"
	"strconv"
)

var CELLWIDTH = 12
var width = 0
var height = 0
var sortreverse = false

func Viewer(bycol bool) error {
	println("Loading...")

	// Read the input
	cef, err := Read(os.Stdin, bycol)
	if err != nil {
		return err
	}

	// Launch the terminal viewer
	if err := termbox.Init(); err != nil {
		return err
	}
	defer termbox.Close()

	// Set up the screen
	termbox.HideCursor()
	termbox.Clear(termbox.ColorYellow, termbox.ColorBlack)

	// Redraw everything
	width, height := termbox.Size()
	redraw(cef, width, height)

	// Main event loop
loop:
	for {
		evt := termbox.PollEvent()
		switch evt.Type {
		case termbox.EventResize:
			redraw(cef, evt.Width, evt.Height)
			width = evt.Width
			height = evt.Height
		case termbox.EventKey:
			if evt.Key == termbox.KeyCtrlC || evt.Ch == 'q' || evt.Ch == 'Q' {
				break loop
			}
			switch evt.Key {
			case termbox.KeyArrowLeft:
				left--
				if left <= 0 {
					left = 0
				}
			case termbox.KeyArrowRight:
				left++
			case termbox.KeyArrowDown:
				top++
			case termbox.KeyArrowUp:
				top--
				if top <= 0 {
					top = 0
				}
			}
			switch evt.Ch {
			case 'A':
				offsetX -= width / CELLWIDTH
				if offsetX < 0 {
					offsetX = 0
				}
			case 'D':
				offsetX += width / CELLWIDTH
			case 'S':
				offsetY += height
			case 'W':
				offsetY -= height
				if offsetY < 0 {
					offsetY = 0
				}
			case 'w':
				offsetY--
				if offsetY < 0 {
					offsetY = 0
				}
			case 's':
				offsetY++
			case 'a':
				offsetX--
				if offsetX < 0 {
					offsetX = 0
				}
			case 'd':
				offsetX++
			case 'h':
				offsetX = 0
				offsetY = 0
				top = 0
				left = 0
			case 'z':
				offsetX = cef.Columns - width/CELLWIDTH + len(cef.RowAttributes) + 1 - left
				offsetY = cef.Rows - height + len(cef.ColumnAttributes) + 3 - top
			case '+':
				CELLWIDTH++
			case '-':
				CELLWIDTH--
				if CELLWIDTH < 4 {
					CELLWIDTH = 4
				}
			case '0':
				CELLWIDTH = 12
			case 'o':
				// Sort
				sort_by := ""
				if left == len(cef.RowAttributes) {
					break
				}
				if left > len(cef.RowAttributes) {
					sort_by = "#" + strconv.Itoa(left-len(cef.RowAttributes)+offsetX)
					result, err := cef.SortNumerical(sort_by, sortreverse)
					if err == nil {
						cef = result
					}
				} else {
					result, err := cef.SortNumerical(cef.RowAttributes[left].Name, sortreverse)
					if err == nil {
						cef = result
					} else {
						result, err := cef.SortByRowAttribute(cef.RowAttributes[left].Name, sortreverse)
						if err == nil {
							cef = result
						}
					}
				}
				sortreverse = !sortreverse
			case 't':
				// Transpose
				result := new(Cef)
				result.RowAttributes = cef.ColumnAttributes
				result.ColumnAttributes = cef.RowAttributes
				result.Headers = cef.Headers
				result.Rows = cef.Columns
				result.Columns = cef.Rows
				result.Flags = cef.Flags
				result.Matrix = make([]float32, len(cef.Matrix))
				for col := 0; col < result.Columns; col++ {
					for row := 0; row < result.Rows; row++ {
						result.Set(col, row, cef.Get(row, col))
					}
				}
				cef = result
				temp := offsetX
				offsetX = offsetY
				offsetY = temp
				temp = left
				left = top
				top = temp
			}
			redraw(cef, width, height)
		}
	}
	return nil
}

// The offsets of the main matrix
var offsetX = 0
var offsetY = 0

// The top left attributes shown
var top = 0
var left = 0

func redraw(cef *Cef, w, h int) {
	termbox.Clear(termbox.ColorYellow, termbox.ColorBlack)
	drawToolbar(h-1, termbox.ColorBlack, termbox.ColorYellow)

	// Draw the column attributes
	for ix := top; ix < len(cef.ColumnAttributes); ix++ {
		// Draw the column attribute name
		drawCell(cef.ColumnAttributes[ix].Name, (len(cef.RowAttributes)-left)*CELLWIDTH, ix-top, termbox.ColorGreen, termbox.ColorBlack)
		// Draw the column attribute values
		for j := offsetX; j < cef.Columns; j++ {
			drawCell(cef.ColumnAttributes[ix].Values[j], (len(cef.RowAttributes)-left+j+1-offsetX)*CELLWIDTH, ix-top, termbox.ColorCyan, termbox.ColorBlack)
			if (len(cef.RowAttributes)-left+j+1-offsetX)*CELLWIDTH > w {
				break
			}
		}
	}

	// Draw the row attribute names
	for ix := left; ix < len(cef.RowAttributes); ix++ {
		drawCell(cef.RowAttributes[ix].Name, (ix-left)*CELLWIDTH, len(cef.ColumnAttributes)-top+1, termbox.ColorGreen, termbox.ColorBlack)
	}

	// Draw the rows
	for row := offsetY; row < cef.Rows; row++ {
		if row+len(cef.ColumnAttributes)-top+3-offsetY >= h {
			break
		}
		// Draw the row attribute values
		for ix := left; ix < len(cef.RowAttributes); ix++ {
			drawCell(cef.RowAttributes[ix].Values[row], (ix-left)*CELLWIDTH, row+len(cef.ColumnAttributes)-top+2-offsetY, termbox.ColorCyan, termbox.ColorBlack)
		}
		// Draw the row matrix values
		for col := offsetX; col < cef.Columns; col++ {
			value := float64(cef.Get(row, col))
			number := strconv.FormatFloat(value, 'f', 1, 32)
			if value > 10 {
				number = strconv.FormatFloat(value, 'f', 0, 32)
			}
			drawCell(number, (col+len(cef.RowAttributes)-left+1-offsetX)*CELLWIDTH, row+len(cef.ColumnAttributes)-top+2-offsetY, termbox.ColorWhite, termbox.ColorBlack)
			if (col+len(cef.RowAttributes)-left+1-offsetX)*CELLWIDTH > w {
				break
			}
		}

	}

	xy := " " + strconv.Itoa(offsetX) + " " + strconv.Itoa(offsetY) + " "
	drawToolbarItem(xy, h-1, w-len([]rune(xy)), termbox.ColorWhite, termbox.ColorBlue)

	termbox.Flush()
}

func drawCell(text string, x, y int, fg, bg termbox.Attribute) {
	ix := 0
	for _, char := range text {
		if ix >= CELLWIDTH-1 {
			break
		}
		termbox.SetCell(x+ix, y, char, fg, bg)
		ix++
	}
	if ix == CELLWIDTH-1 {
		termbox.SetCell(x+ix-1, y, '…', termbox.ColorRed, bg)
	}
}

func drawToolbar(y int, fg, bg termbox.Attribute) {
	x := 0
	for _, text := range []string{"[q] Quit ", "[wasd] Pan ", "[WASD] Pan quickly ", "[←↑→↓] Pan all ", "[+-0] Zoom ", "[h] Home ", "[z] End ", "[o] Sort ", "[t] Transpose "} {
		x += drawToolbarItem(text, y, x, termbox.ColorBlack|termbox.AttrBold, termbox.ColorYellow) + 1
	}
}

func drawToolbarItem(text string, y int, x int, fg, bg termbox.Attribute) int {
	ix := 0
	for _, char := range text {
		termbox.SetCell(x+ix, y, char, fg, bg)
		ix++
	}
	return ix
}
