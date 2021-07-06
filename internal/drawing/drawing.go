package drawing

import (
	"fmt"
	"image/color"
	"path/filepath"

	"github.com/fogleman/gg"
	"github.com/gxravel/bus-routes-visualizer/internal/busroutesapi/v1"
)

// TODO: config
const (
	defaultWidth   = 1000
	defaultHeight  = 300
	fontSize       = 15
	lineWidth      = 5
	xStart, yStart = 150, 150
	rPoint         = 30
	stringXOffset  = -2 * rPoint
	stringYOffset  = -1.5 * rPoint
	fontPath       = "internal/drawing/Roboto-Black.ttf"
	dataPath       = "internal/drawing/data"

	xOffset = 200
)

// DrawRoutes draws routes as a linear graph with addresses as vertexes and saves it in PNG.
// Returns path to the image.
func DrawRoutes(name string, routes []*busroutesapi.RouteDetailed) (string, error) {
	dc := gg.NewContext(defaultWidth, defaultHeight)

	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetColor(color.White)
	dc.SetLineWidth(lineWidth)

	err := dc.LoadFontFace(fontPath, fontSize)
	if err != nil {
		return "", err
	}

	var x, y, r float64 = xStart, yStart, rPoint

	for _, route := range routes {
		dc.DrawString(route.City+" "+route.Bus, 50, 50)

		for j, point := range route.Points {
			dc.DrawPoint(x, y, r)
			dc.FillPreserve()

			dc.DrawString(fmt.Sprintf("%d) %s", point.Step, point.Address), x+stringXOffset, y+stringYOffset)

			x2 := x + xOffset

			dc.LineTo(x2, y)

			if j != len(route.Points)-1 {
				dc.Stroke()
			}

			x = x2
		}
	}

	path := filepath.Join(dataPath, name) + ".png"

	return path, dc.SavePNG(path)
}
