package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var text1 = "[yellow]Leverage agile [red]frameworks to provide a robust synopsis for high level overviews. Iterative approaches to text1 strategy foster collaborative thinking to further the overall value proposition. Organically grow the holistic world view of disruptive innovation via workplace diversity and empowerment. Bring to the table win-win survival strategies to ensure proactive domination. At the end of the day, going forward, a new normal that has evolved from generation X is on the runway heading towards a streamlined cloud solution. User generated content in real-time will have multiple touchpoints for offshoring. Capitalize on low hanging fruit to identify a ballpark value added activity to beta test. Override the digital divide with additional clickthroughs from DevOps. Nanotechnology immersion along the information highway will close the loop on focusing solely on the bottom line. [[yellow]] press Enter, then Tab/Backtab for word selections"
var text2 = "Le Lorem Ipsum est simplement du faux texte employé dans la composition et la mise en page avant impression. Le Lorem Ipsum est le faux texte standard de l'imprimerie depuis les années 1500, quand un imprimeur anonyme assembla ensemble des morceaux de texte pour réaliser un livre spécimen de polices de texte. Il n'a pas fait que survivre cinq siècles, mais s'est aussi adapté à la bureautique informatique, sans que son contenu n'en soit modifié. Il a été popularisé dans les années 1960 grâce à la vente de feuilles Letraset contenant des passages du Lorem Ipsum, et, plus récemment, par son inclusion dans des applications de mise en page de texte, comme Aldus PageMaker."

const corporate = `Leverage agile frameworks to provide a robust synopsis for high level overviews. Iterative approaches to corporate strategy foster collaborative thinking to further the overall value proposition. Organically grow the holistic world view of disruptive innovation via workplace diversity and empowerment.

Bring to the table win-win survival strategies to ensure proactive domination. At the end of the day, going forward, a new normal that has evolved from generation X is on the runway heading towards a streamlined cloud solution. User generated content in real-time will have multiple touchpoints for offshoring.

Capitalize on low hanging fruit to identify a ballpark value added activity to beta test. Override the digital divide with additional clickthroughs from DevOps. Nanotechnology immersion along the information highway will close the loop on focusing solely on the bottom line.

[yellow]Press Enter, then Tab/Backtab for word selections`

type Display[T any] struct {
	App      *tview.Application
	TextView *tview.TextView
	Channels []chan T
}

func New[T any](channelNumber int) *Display[T] {
	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	channels := make([]chan T, channelNumber)
	for i := 0; i < channelNumber; i++ {
		channels = append(channels, make(chan T))
	}
	return &Display[T]{
		App:      app,
		TextView: textView,
		Channels: channels,
	}
}

func main() {
	display := New[string](1)

	numSelections := 0
	go func() {
		for i, word := range strings.Split(corporate, " ") {
			if i%5 == 0 {
				display.TextView.Clear()
			}
			if word == "the" {
				word = "[red]the[white]"
			}
			if word == "to" {
				word = fmt.Sprintf(`["%d"]to[""]`, numSelections)
				numSelections++
			}
			fmt.Fprintf(display.TextView, "%s ", word)
			time.Sleep(200 * time.Millisecond)
		}
	}()

	display.TextView.SetDoneFunc(func(key tcell.Key) {
		currentSelection := display.TextView.GetHighlights()
		if key == tcell.KeyEnter {
			if len(currentSelection) > 0 {
				display.TextView.Highlight()
			} else {
				display.TextView.Highlight("0").ScrollToHighlight()
			}
		} else if len(currentSelection) > 0 {
			index, _ := strconv.Atoi(currentSelection[0])
			if key == tcell.KeyTab {
				index = (index + 1) % numSelections
			} else if key == tcell.KeyBacktab {
				index = (index - 1 + numSelections) % numSelections
			} else {
				return
			}
			display.TextView.Highlight(strconv.Itoa(index)).ScrollToHighlight()
		}
	})
	display.TextView.SetBorder(true)
	if err := display.App.SetRoot(display.TextView, true).SetFocus(display.TextView).Run(); err != nil {
		panic(err)
	}
}

// func main() {

// 	app := tview.NewApplication()
// 	runFlex := tview.NewFlex().SetDirection(tview.FlexRow)
// 	runView := tview.NewTextView().SetDynamicColors(true).SetRegions(true).
// 		SetWordWrap(true).SetChangedFunc(func() {
// 	})
// 	errView := tview.NewTextView().SetDynamicColors(true).SetRegions(true).
// 		SetWordWrap(true)
// 	runView.SetTitle(" TEXT 1 ").SetBorder(true)
// 	errView.SetTitle(" TEXT 2 ").SetBorder(true)

// 	runFlex.AddItem(runView, 0, 1, true).AddItem(errView, 0, 1, true)

// 	go func() {
// 		for _, word := range strings.Split(text1, " ") {
// 			runView.Clear()
// 			fmt.Fprintf(runView, "%s ", word)
// 			fmt.Println("1")
// 			time.Sleep(200 * time.Millisecond)
// 		}
// 	}()

// 	go func() {
// 		for _, word := range strings.Split(text2, " ") {
// 			errView.Clear()
// 			fmt.Fprintf(errView, "%s ", word)
// 			fmt.Println("2")
// 			time.Sleep(250 * time.Millisecond)
// 		}
// 	}()

// 	if err := app.SetRoot(runFlex, true).Run(); err != nil {
// 		panic(err)
// 	}

// }
