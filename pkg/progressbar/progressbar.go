package progressbar

import "fmt"

type Bar struct {
	style            string
	speed            string
	length           float64
	ArrayLength      int
	current_progress int
}

func NewProgress(ArrayLength int, style string) *Bar {
	if style == "" {
		style = "█"
	}
	return &Bar{ArrayLength: ArrayLength, style: style, length: float64(ArrayLength) / 50}
}

func (bar *Bar) AddProgressCount(current_count int) {
	for i := int32(float64(bar.current_progress)/bar.length) + 1; i <= int32(float64(current_count)/bar.length); i++ {
		bar.speed += bar.style
	}
	bar.current_progress = current_count // current_count += 1
	fmt.Printf("\r[%-50s]%0.2f%%  	%8d/%d", bar.speed, bar.percentage(), bar.current_progress, bar.ArrayLength)

}

func (bar *Bar) ProgressEnd() {
	for i := int32(float64(bar.current_progress)/bar.length) + 1; i <= 50; i++ {
		bar.speed += bar.style
	}
	fmt.Printf("\r[%-50s]%0.2f%%  	%8d/%d", bar.speed, bar.percentage(), bar.current_progress, bar.ArrayLength)

}

func (bar *Bar) percentage() float64 {
	return float64(bar.current_progress) / float64(bar.ArrayLength) * 100
}

func (bar *Bar) progress_stop() {
	fmt.Println() // print a new line after progress bar
}
