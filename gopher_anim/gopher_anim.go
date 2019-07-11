package gopher_anim

import (
	"encoding/csv"
	"encoding/json"
	// "fmt"
	"github.com/JChouCode/naruto-go/gopher"
	// "github.com/JChouCode/naruto-go/projectile"
	"github.com/faiface/pixel"
	"image"
	_ "image/png"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"
)

func loadAnimationSheet(imgPath, csvPath string, fWidth float64) (sheet pixel.Picture, anims map[string][]pixel.Rect, err error) {
	// open and load the spritesheet
	sheetFile, err := os.Open(imgPath)
	if err != nil {
		return nil, nil, err
	}
	defer sheetFile.Close()
	sheetImg, _, err := image.Decode(sheetFile)
	if err != nil {
		return nil, nil, err
	}
	sheet = pixel.PictureDataFromImage(sheetImg)

	// create a slice of frames inside the spritesheet
	var frames []pixel.Rect
	for x := 0.0; x+fWidth <= sheet.Bounds().Max.X; x += fWidth {
		frames = append(frames, pixel.R(
			x,
			0,
			x+fWidth,
			sheet.Bounds().H(),
		))
	}

	descFile, err := os.Open(csvPath)
	if err != nil {
		return nil, nil, err
	}
	defer descFile.Close()

	anims = make(map[string][]pixel.Rect)

	// load the animation information, name and interval inside the spritesheet
	desc := csv.NewReader(descFile)
	for {
		anim, err := desc.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		name := anim[0]
		start, _ := strconv.Atoi(anim[1])
		end, _ := strconv.Atoi(anim[2])
		anims[name] = frames[start : end+1]
	}
	return sheet, anims, nil
}

func addAnim(imgPath string, name string, fWidth float64, anims map[string][]pixel.Rect) (sheet pixel.Picture, newAnim map[string][]pixel.Rect, err error) {
	// open and load the spritesheet
	sheetFile, err := os.Open(imgPath)
	if err != nil {
		return nil, nil, err
	}
	defer sheetFile.Close()
	sheetImg, _, err := image.Decode(sheetFile)
	if err != nil {
		return nil, nil, err
	}
	sheet = pixel.PictureDataFromImage(sheetImg)
	// create a slice of frames inside the spritesheet
	var frames []pixel.Rect
	for x := 0.0; x+fWidth <= sheet.Bounds().Max.X; x += fWidth {
		frames = append(frames, pixel.R(
			x,
			0,
			x+fWidth,
			sheet.Bounds().H(),
		))
	}
	anims[name] = frames
	return sheet, anims, nil
}

func LoadAnimationJson(imgPath string, jsonPath string) (sheet pixel.Picture, anims map[string][]pixel.Rect, err error) {
	// open and load the spritesheet
	sheetFile, err := os.Open(imgPath)
	if err != nil {
		return nil, nil, err
	}
	defer sheetFile.Close()
	sheetImg, _, err := image.Decode(sheetFile)
	if err != nil {
		return nil, nil, err
	}
	sheet = pixel.PictureDataFromImage(sheetImg)
	height := sheet.Bounds().H()
	// width := sheet.Bounds().W()

	temp, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, nil, err
	}

	type Pos struct {
		X float64
		Y float64
		W float64
		H float64
	}
	var data map[string][]Pos
	err2 := json.Unmarshal(temp, &data)
	if err2 != nil {
		return nil, nil, err
	}
	// fmt.Println(data)
	anims = make(map[string][]pixel.Rect)
	// load the animation information, name and interval inside the spritesheet

	for k, v := range data {
		for _, i := range v {
			if _, ok := anims[k]; ok {
				anims[k] = append(anims[k], pixel.R(i.X, height+i.Y-i.H, i.X+i.W, height+i.Y))
			} else {
				anims[k] = []pixel.Rect{pixel.R(i.X, height+i.Y-i.H, i.X+i.W, height+i.Y)}
			}
		}
	}
	// fmt.Println(anims)
	return sheet, anims, nil
}

type move_state int

const (
	idle move_state = iota
	run
	jump
	crouch
	cwalk
	throw
)

var rate = 1.0 / 10
var counter = 0.0

type GopherAnim struct {
	sprite *pixel.Sprite
	frame  pixel.Rect
	state  move_state

	sheet pixel.Picture
	anims map[string][]pixel.Rect
	throw bool
	dir   float64
}

// Initialize GopherAnim
func New(imgPath string, jsonPath string) GopherAnim {
	sheet, anims, err := LoadAnimationJson(imgPath, jsonPath)
	if err != nil {
		panic(err)
	}
	return GopherAnim{pixel.NewSprite(nil, pixel.Rect{}), pixel.Rect{}, idle, sheet, anims, false, 0}
}

func (ga *GopherAnim) Update(g gopher.Gopher, dt float64) {
	counter += dt

	var tempState move_state
	//Update state
	switch {
	case ga.IsThrow():
		tempState = throw
	case g.IsJump():
		tempState = jump
	case g.IsCrouch():
		if g.GetVel().Len() > 0 {
			tempState = cwalk
		} else {
			tempState = crouch
		}
	case g.GetVel().Len() > 0:
		tempState = run
	case g.GetVel().Len() == 0:
		tempState = idle
	}

	if tempState != ga.state {
		ga.state = tempState
		counter = 0
	}

	switch ga.state {
	case idle:
		i := int(math.Floor(counter / rate))
		ga.frame = ga.anims["side"][i%len(ga.anims["side"])]
	case run:
		// fmt.Print("run")
		i := int(math.Floor(counter / rate))
		ga.frame = ga.anims["run"][i%len(ga.anims["run"])]
	case jump:
		// fmt.Print("jump")
		i := 0
		switch {
		case g.GetVel().Y < gopher.GetJumpY()*1/3:
			i++
		case g.GetVel().Y < gopher.GetJumpY()*2/3:
			i++
		}
		ga.frame = ga.anims["jump"][i]
	case crouch:
		ga.frame = ga.anims["crouch"][2]
	case cwalk:
		i := int(math.Floor(counter / rate))
		ga.frame = ga.anims["cwalk"][i%len(ga.anims["cwalk"])]
	case throw:
		i := int(math.Floor(counter / rate))
		ga.frame = ga.anims["throw"][i%len(ga.anims["throw"])]
		if i == 2 {
			ga.throw = false
		}
	}
	ga.dir = g.GetDir()
}

func (ga *GopherAnim) Draw(t pixel.Target, g gopher.Gopher) {
	// fmt.Print(ga.frame)
	ga.sprite.Set(ga.sheet, ga.frame)
	ga.sprite.Draw(t, pixel.IM.
		// ScaledXY(pixel.ZV, pixel.V(
		// 	g.GetBody().W()/ga.sprite.Frame().W(),
		// 	g.GetBody().H()/ga.sprite.Frame().H(),
		// )).
		ScaledXY(pixel.ZV, pixel.V(ga.dir, 1)).
		Moved(g.GetBody().Center()),
	)
}

func (ga *GopherAnim) GetSheet() pixel.Picture {
	return ga.sheet
}

func (ga *GopherAnim) GetFrame() pixel.Rect {
	return ga.frame
}

func (ga *GopherAnim) Throw() {
	ga.throw = true
}

func (ga *GopherAnim) IsThrow() bool {
	return ga.throw
}
