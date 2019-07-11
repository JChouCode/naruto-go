package projectile

import (
	"fmt"
	"github.com/JChouCode/naruto-go/gopher"
	"github.com/JChouCode/naruto-go/gopher_anim"
	"github.com/faiface/pixel"
	_ "image/png"
	"math"
)

var velX = 500.0

var airRes = 0.999
var rate = 3.0 / 10
var counter = 0.0
var sheet pixel.Picture
var anims map[string][]pixel.Rect
var err error

type Projectile struct {
	sprite *pixel.Sprite
	frame  pixel.Rect
	sheet  pixel.Picture
	anims  map[string][]pixel.Rect
	dir    float64
	vel    pixel.Vec
	body   pixel.Rect
}

func init() {
	sheet, anims, err = gopher_anim.LoadAnimationJson("projectile.png", "projectile.json")
	fmt.Print(anims)
	if err != nil {
		panic(err)
	}
}

func New(g gopher.Gopher, height float64, width float64) Projectile {
	gBody := g.GetBody()
	startPos := gBody.Center().Add(pixel.V(gBody.W()/2, -height/2))
	return Projectile{pixel.NewSprite(nil, pixel.Rect{}), pixel.Rect{}, sheet, anims, g.GetDir(), pixel.V(velX, 0), pixel.R(startPos.X, startPos.Y, startPos.X+width, startPos.Y+height)}
}

func (p *Projectile) Update(dt float64) {
	counter += dt
	// p.vel.X = p.vel.X * airRes
	p.body = p.body.Moved(p.vel.Scaled(dt))
	i := int(math.Floor(counter / rate))
	p.frame = p.anims["kunai"][i%len(p.anims["kunai"])]
}

func (p *Projectile) Draw(t pixel.Target) {
	p.sprite.Set(p.sheet, p.frame)
	p.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(p.dir, 1)).
		Moved(p.body.Center()),
	)
}

func GetSheet() pixel.Picture {
	return sheet
}

func (p Projectile) GetVel() pixel.Vec {
	return p.vel
}

func (p Projectile) GetBody() pixel.Rect {
	return p.body
}

func (p Projectile) GetFrame() pixel.Rect {
	return p.frame
}

func (p Projectile) Offscreen() bool {
	return p.body.Min.X > 1024
}
