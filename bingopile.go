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

type Games struct {
	Header string
	Games []Slides
}

type Cards struct {
	Header string
	Cards []CardFmt
}

type CardFmt struct {
	Header [len(game)]string
	Content [len(game)][len(game)]string
}

const CardTemp = `<html><head><title>{{.Header}}</title><body>
{{range .Cards}}
<table style="break-inside: avoid; border: 1px solid black; float: left; margin: 15px">
<tr>
{{range .Header}}
<th style="height: 90px; width: 90px; font-size: xx-large">{{.}}</th>
{{end}}
{{range .Content}}
<tr>
{{range .}}
<td style="border: 1px solid black; height: 90px; width: 90px; text-align: center; font-size:small">{{.}}</td>
{{end}}
</tr>
{{end}}
</table>
{{end}}
</body></html>`
const SlideTemp = `<html><head><title>{{.Header}}</title><body>
{{range .Games}}
<h1>Begin {{.Header}}</h1>
{{range .Slides}}
<div style="break-inside: avoid; break-after: always; break-before: always; width: 100%">

<div style="float: right; width: 20%">
<p>Previous:</p>
<ul>
{{range .Hist}}
<li>{{.}}</li>
{{end}}
</ul>
</div>

<div style="width: 70%"><h1>{{.Saint.Bingo}}</h1><h2 style="text-align: center">{{.Saint.Name}}</h2>
<p style="text-align: center"><img src="{{.Saint.Photo}}" alt="{{.Saint.Name}}" style="max-height: 600px; height: 100%; width: auto" /></p>
<p style="text-align: center">{{.Saint.Facts}}</p>
</div>

</div>
{{end}}
<h1 style="break-before: always">END {{.Header}}</h1>
{{end}}
</body></html>`

const frei = true

func gengameslides(bingodata []Saint, gamenum int) Slides {
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
		Header: "GAME " + strconv.Itoa(gamenum),
		Slides: slides,
	}
	return slideobj
}

func htmlout(templ *template.Template, data interface{}, outfile string) {
	fh, err := os.Create(outfile)
	if err != nil {
		panic(err)
	}
	defer fh.Close()
	if templ.Execute(fh, data) != nil {
		panic(err)
	}
}

func slidesout(games []Slides, outfile string) {
	gamesobj := Games {
		Header: game,
		Games: games,
	}
	var slidetemp = template.Must(template.New("slides").Parse(SlideTemp))
	htmlout(slidetemp, gamesobj, outfile)
}

func cardsout(cards []CardFmt, outfile string) {
	var cardtemp = template.Must(template.New("cards").Parse(CardTemp))
	cardsobj := Cards {
		Header: game,
		Cards: cards,
	}
	htmlout(cardtemp, cardsobj, outfile)
}

func refmtcard(in Card) CardFmt {
	var newtitle [len(game)]string
	var newcontent [len(game)][len(game)]string
	for i := 0; i < len(game); i++ {
		newtitle[i] = string(in.Header[i])
	}
	for i := 0; i < len(game); i++ {
		for j := 0; j < len(game); j++ {
			newcontent[i][j] = in.Content[j][i]
		}
	}
	return CardFmt {
		Header: newtitle,
		Content: newcontent,
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var filename string = os.Args[1]
	participants, err := strconv.Atoi(os.Args[2])
	var cardfile string = os.Args[3]
	numgames, err := strconv.Atoi(os.Args[4])
	var slidefile string = os.Args[5]
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
	var fmtcards []CardFmt
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
		if frei {
			cards[i].Content[len(game)/2][len(game)/2] = "FREE"
		}
		fmtcards = append(fmtcards, refmtcard(cards[i]))
	}
	fmt.Println("Slides")
	var mygames []Slides
	for i := 1; i <= numgames; i++ {
		mygames = append(mygames, gengameslides(bingodata, i))
	}
	slidesout(mygames, slidefile)
	cardsout(fmtcards, cardfile)
}
