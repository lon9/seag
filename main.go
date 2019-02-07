package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"image"
	"os"

	"github.com/gobuffalo/packr/v2"
	"github.com/kyroy/kdtree"
	"github.com/kyroy/kdtree/points"

	_ "image/jpeg"
	_ "image/png"

	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/nfnt/resize"
)

func getEmotes(box *packr.Box, n int) ([]Emote, error) {
	b, err := box.Find("emote.json")
	if err != nil {
		return nil, err
	}
	var emotes []Emote
	if err := json.Unmarshal(b, &emotes); err != nil {
		return nil, err
	}
	if n == -1 {
		return emotes, nil
	}
	return emotes[:n], nil
}

func saveHTML(box *packr.Box, art [][]Emote) (err error) {
	templateString, err := box.FindString("template.html")
	if err != nil {
		return err
	}
	t, err := template.New("t").Parse(templateString)
	if err != nil {
		return err
	}
	f, err := os.Create("index.html")
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, art)
}

func printStrings(art [][]Emote) {
	for y := 0; y < len(art); y++ {
		for x := 0; x < len(art[y]); x++ {
			fmt.Print(art[y][x].Name)
		}
		fmt.Println()
	}
}

func main() {

	var (
		input     string
		width     int
		height    int
		qFactor   int
		variation int
		isHTML    bool
	)

	flag.StringVar(&input, "i", "sample.jpg", "Input image")
	flag.IntVar(&width, "width", 24, "Output width")
	flag.IntVar(&height, "height", 0, "Output height")
	flag.IntVar(&qFactor, "q", 10, "Quality")
	flag.IntVar(&variation, "v", -1, "Variation -1 = all")
	flag.BoolVar(&isHTML, "html", false, "Whether output html or not")
	flag.Parse()

	box := packr.New("myBox", "./data")

	emotes, err := getEmotes(box, variation)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	m := resize.Resize(uint(width*qFactor), uint(height*qFactor), img, resize.NearestNeighbor)
	if err != nil {
		panic(err)
	}

	rgbTree := new(kdtree.KDTree)
	for _, emote := range emotes {
		c := colorful.Hsl(emote.HLS[0], emote.HLS[2]/255.0, emote.HLS[1]/255.0)
		rgbTree.Insert(points.NewPoint([]float64{c.R, c.G, c.B}, emote))
	}

	result := make([][]Emote, m.Bounds().Max.Y/qFactor)
	for y := 0; y < m.Bounds().Max.Y; y += qFactor {
		result[y/qFactor] = make([]Emote, m.Bounds().Max.X/qFactor)
		for x := 0; x < m.Bounds().Max.X; x += qFactor {
			c, ok := colorful.MakeColor(m.At(x, y))
			if !ok {
				c = colorful.Color{R: 1.0, G: 1.0, B: 1.0}
			}
			res := rgbTree.KNN(&points.Point{Coordinates: []float64{c.R, c.G, c.B}}, 1)
			result[y/qFactor][x/qFactor] = res[0].(*points.Point).Data.(Emote)
		}
	}
	printStrings(result)
	if isHTML {
		if err := saveHTML(box, result); err != nil {
			panic(err)
		}
	}
}
