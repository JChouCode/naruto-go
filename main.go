package main

import (
	"fmt"
	"github.com/JChouCode/naruto-go/hero"
	"github.com/JChouCode/naruto-go/hero_anim"
	"github.com/JChouCode/naruto-go/platform"
	"github.com/JChouCode/naruto-go/projectile"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	// "image"
	_ "image/png"
	"math/rand"
	"os"
	"time"
)

//Initialize window
func initWindow(t string, w float64, h float64) *pixelgl.Window {
	cfg := pixelgl.WindowConfig{
		Title:  t,
		Bounds: pixel.R(0, 0, w, h),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	// win.SetSmooth(true)
	return win
}

func run() {

	win := initWindow("Gopher-Run-GO", 1024, 768)

	rand.Seed(time.Now().UnixNano())

	hero := hero.New()
	anim := hero_anim.New("neji.png", "neji.json")

	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(200, 400), atlas)
	basicTxt.Color = colornames.Black

	platforms := []platform.Platform{
		platform.New(pixel.V(200, 100), 30, 200),
		platform.New(pixel.V(350, 250), 30, 200),
		platform.New(pixel.V(500, 400), 30, 200),
		platform.New(pixel.V(650, 550), 30, 200),
	}
	var projectiles []projectile.Projectile

	imd := imdraw.New(anim.GetSheet())
	imd2 := imdraw.New(projectile.GetSheet())

	last := time.Now()

	for !win.Closed() {
		basicTxt.Clear()
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colornames.Blanchedalmond)

		ctrl := pixel.ZV
		if win.Pressed(pixelgl.KeyLeft) {
			ctrl.X--
		}
		if win.Pressed(pixelgl.KeyRight) {
			ctrl.X++
		}
		if win.Pressed(pixelgl.KeyDown) {
			ctrl.Y--
		}
		if win.JustPressed(pixelgl.KeyUp) {
			ctrl.Y++
		}
		if win.JustPressed(pixelgl.KeyA) {
			anim.Throw()
			projectiles = append(projectiles, projectile.New(hero, 21, 9))
		}
		if win.JustPressed(pixelgl.KeyQ) {
			os.Exit(1)
		}

		hero.Update(ctrl, dt, platforms)
		anim.Update(hero, dt)
		for index := range projectiles {
			(&projectiles[index]).Update(dt)
		}

		k := 0
		for _, proj := range projectiles {
			if !proj.Offscreen() {
				projectiles[k] = proj
				k++
			}
		}
		projectiles = projectiles[:k]

		imd.Clear()
		imd2.Clear()

		for _, p := range platforms {
			p.Draw(imd)
		}

		anim.Draw(imd, hero)
		for _, p := range projectiles {
			p.Draw(imd2)
		}

		imd.Draw(win)
		imd2.Draw(win)

		// Debug
		fmt.Fprintln(basicTxt, anim.GetFrame())
		fmt.Fprintln(basicTxt, len(projectiles))
		if len(projectiles) > 0 {
			fmt.Fprintln(basicTxt, projectiles[0].GetBody())
			fmt.Fprintln(basicTxt, projectiles[0].GetFrame())
		}
		basicTxt.Draw(win, pixel.IM)

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
