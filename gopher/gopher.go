package gopher

import (
	"github.com/JChouCode/naruto-go/platform"
	"github.com/faiface/pixel"
	"math"
)

// Gravity constant
const gravity = 700

// Velocity of gopher
var (
	runX    float64 = 230
	crouchX float64 = 150
	jumpY   float64 = 500
)

type Gopher struct {
	body   pixel.Rect
	vel    pixel.Vec
	jump   bool
	crouch bool
}

// Initialize Gopher
func New() Gopher {
	// return Gopher{pixel.R(-6, -7, 6, 7), pixel.V(0, 0), false}
	return Gopher{pixel.R(1, 1, 41, 61), pixel.V(0, 0), false, false}
}

// Gopher Jump
func (g *Gopher) Jump() {
	g.vel.Y += jumpY
	g.body = g.body.Moved(pixel.V(0, g.vel.Y))
	g.jump = true
}

func (g *Gopher) Update(ctrl pixel.Vec, dt float64, platforms []platform.Platform) {
	if ctrl.Y < 0 {
		g.crouch = true
	} else {
		g.crouch = false
	}
	switch {
	// Running forward
	case ctrl.X > 0:
		if ctrl.Y < 0 {
			g.vel.X = crouchX
		} else {
			g.vel.X = runX
		}
		// Running backward
	case ctrl.X < 0:
		if ctrl.Y < 0 {
			g.vel.X = -crouchX
		} else {
			g.vel.X = -runX
		}
		// Not moving
	case ctrl.X == 0:
		g.vel.X = 0
	}
	// Apply gravity
	if g.jump {
		g.vel.Y -= gravity * dt
	}

	g.body = g.body.Moved(g.vel.Scaled(dt))

	//Check if hit platform
	if g.jump && g.vel.Y < 0 {
		for _, p := range platforms {
			if g.IsCollide(p) {
				g.vel.Y = 0
				g.body = g.body.Moved(pixel.V(0, p.GetRect().Max.Y-g.body.Min.Y))
				g.jump = false
				continue
			}
		}
	} else {
		temp := true
		for _, p := range platforms {
			if g.IsCollide(p) {
				temp = false
			}
		}
		if temp {
			g.jump = true
		}
	}

	// Check if hit ground
	if g.body.Min.Y <= 1 {
		g.vel.Y = 0
		g.body = g.body.Moved(pixel.V(0, 1-g.body.Min.Y))
		g.jump = false
	}

	if !g.jump && ctrl.Y > 0 {
		g.vel.Y = jumpY
		g.jump = true
	}
}

func (g Gopher) IsCollide(p platform.Platform) bool {
	return g.body.Max.X-g.body.W()/2 <= p.GetRect().Max.X && g.body.Min.X+g.body.W()/2 >= p.GetRect().Min.X && math.Abs(g.body.Min.Y-p.GetRect().Max.Y) <= 7
}

func (g *Gopher) IsJump() bool {
	return g.jump
}

func (g *Gopher) IsCrouch() bool {
	return g.crouch
}

func (g *Gopher) GetBody() pixel.Rect {
	return g.body
}

func (g *Gopher) GetVel() pixel.Vec {
	return g.vel
}

func (g *Gopher) GetDir() float64 {
	switch {
	case g.vel.X > 0:
		return 1
	case g.vel.X < 0:
		return -1
	case g.vel.X == 0:
		return 1
	}
	return 1
}

func GetJumpY() float64 {
	return jumpY
}
