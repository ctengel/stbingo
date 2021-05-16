package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"math/rand"
	"time"
	"html/template"
)

type Saint struct {
	Name string
	Facts string
	Photo string
	Callable bool
	Bingo string
}

const game = "BINGO"

type Card struct {
	Header string
	Content [len(game)][len(game)]string
}

type Slide struct {
	Saint Saint
	Hist []string
	Header string
}

type Slides struct {
	Header string
	Slides []Slide
}

type Cards struct {
	Header string
	Cards []Card
}

const CardTemp = `<html><head><title>{{.Header}}</title><body>
<table>
</table>
</body></html>`
const SlideTemp = `<html><head><title>{{.Header}}</title><body>
{{range .Slides}}
<h1>{{.Saint.Name}}</h1>
<p>Previous:</p>
<ul>
{{range .Hist}}
<li>{{.}}</li>
{{end}}
</ul>
{{end}}
</body></html>`

func main() {
	rand.Seed(time.Now().UnixNano())
	var filename string = os.Args[1]
	participants, err := strconv.Atoi(os.Args[2])
	//var cardfile string = os.Args[3]
	//var slidefile string = os.Args[4]
	if err != nil {
		os.Exit(1)
	}
	var bingodata []Saint
	var bingofields [len(game)][]string
	fmt.Println("Participants", participants)
	fmt.Println("Reading CSV", filename, "...")
	f, err := os.Open(filename)
	if err != nil {
		os.Exit(1)
	}
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		os.Exit(1)
	}
	for _, line := range lines {
		randnum := rand.Intn(len(game))
		randlet := string(game[randnum])
		data := Saint {
			Name: line[0],
			Facts: line[2],
			Photo: line[1],
			Callable: len(line[2]) + len(line[1]) != 0,
			Bingo: randlet,
		}
		bingodata = append(bingodata, data)
		bingofields[randnum] = append(bingofields[randnum], line[0])
	}
	fmt.Println("Saints:")
	for _, mydat := range bingodata {
		fmt.Println(mydat.Bingo, mydat.Callable, mydat.Name, mydat.Facts, mydat.Photo)
	}
	fmt.Println("Cards:")
	var cards []Card
	for i := 0; i < participants; i++ {
		fmt.Println(i)
		newcard := Card {
			Header: game,
		}
		cards = append(cards, newcard)
		for j, bl := range game {
			var dest [len(game)]string
			perm := rand.Perm(len(bingofields[j]))
			for k := 0; k < len(game); k++ {
				dest[k] = bingofields[j][perm[k]]
			}
			fmt.Println(string(bl), dest)
			cards[i].Content[j] = dest
		}
	}
	fmt.Println("Slides")
	var slides []Slide
	var past []string
	callouts := rand.Perm(len(bingodata))
	for _, i := range callouts {
		if bingodata[i].Callable {
			newslide := Slide {
				Saint: bingodata[i],
				Hist: past,
				Header: game,
			}
			slides = append(slides, newslide)
			past = append(past, string(bingodata[i].Bingo) + " " + bingodata[i].Name)

		}
	}
	for _, slide := range slides {
		fmt.Println(slide)
	}
	slideobj := Slides {
		Header: game,
		Slides: slides,
	}
	var slidetemp = template.Must(template.New("slides").Parse(SlideTemp))
	if slidetemp.Execute(os.Stdout, slideobj) != nil {
		os.Exit(1)
	}
}
