package platform

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"image/color"
	_ "image/png"
	"math"
	"math/rand"
)

type Platform struct {
	rect  pixel.Rect
	color color.Color
}

func New(coord pixel.Vec, height float64, width float64) Platform {
	rect := pixel.Rect{coord, coord.Add(pixel.V(width, height))}
	color := randomNiceColor()
	return Platform{rect, color}
}

func (p *Platform) Draw(imd *imdraw.IMDraw) {
	imd.Color = p.color
	imd.Push(p.rect.Min, p.rect.Max)
	imd.Rectangle(0)
}

func (p Platform) GetRect() pixel.Rect {
	return p.rect
}

func randomNiceColor() pixel.RGBA {
again:
	r := rand.Float64()
	g := rand.Float64()
	b := rand.Float64()
	len := math.Sqrt(r*r + g*g + b*b)
	if len == 0 {
		goto again
	}
	return pixel.RGB(r/len, g/len, b/len)
}
