// MIT License
//
// Copyright (c) 2020 Pedro Rodrigues
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"fmt"
	"math"
	"time"

	t "github.com/csixteen/simulated-evolution/pkg/types"
	u "github.com/csixteen/simulated-evolution/pkg/utils"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	TreesPNG   = "trees.png"
	AnimalsPNG = "animals2.png"

	ScreenWidth  = 1024
	ScreenHeight = 768
)

var (
	// Camera control
	camPos       = pixel.ZV
	camSpeed     = 500.0
	camZoom      = 1.0
	camZoomSpeed = 1.2

	// These are needed just to display the FPS in the Window title
	frames = 0
	second = time.Tick(time.Second)
)

func detectKeyPress(w *pixelgl.Window, dt float64) {
	if w.Pressed(pixelgl.KeyLeft) {
		camPos.X -= camSpeed * dt
	}
	if w.Pressed(pixelgl.KeyRight) {
		camPos.X += camSpeed * dt
	}
	if w.Pressed(pixelgl.KeyDown) {
		camPos.Y -= camSpeed * dt
	}
	if w.Pressed(pixelgl.KeyUp) {
		camPos.Y += camSpeed * dt
	}
	camZoom *= math.Pow(camZoomSpeed, w.MouseScroll().Y)
}

func addToBatch(
	cam pixel.Matrix,
	b *pixel.Batch,
	sprites pixel.Picture,
	r pixel.Rect,
	pos u.Point,
) {
	entity := pixel.NewSprite(sprites, r)
	entity.Draw(b, pixel.IM.Moved(cam.Unproject(pixel.V(pos.X, pos.Y))))
}

func drawWorld(
	cam pixel.Matrix,
	win *pixelgl.Window,
	tb, ab *pixel.Batch,
	tSheet, aSheet *SpritesSheet,
	world *t.World,
) {
	tb.Clear()
	ab.Clear()

	for pos, e := range world.Entities {
		switch e.EntityType() {
		case "tree":
			addToBatch(cam, tb, tSheet.sprites, tSheet.frames[e.Id()], pos)
		case "animal":
			l := len(aSheet.frames)
			addToBatch(cam, ab, aSheet.sprites, aSheet.frames[e.Id()%l], pos)
		}
	}

	tb.Draw(win)
	ab.Draw(win)
}

func run() {
	world := t.NewWorld(ScreenWidth, ScreenHeight)
	treesSheet := NewSpritesSheet(TreesPNG, 32)
	animalsSheet := NewSpritesSheet(AnimalsPNG, 48)

	treesBatch := pixel.NewBatch(&pixel.TrianglesData{}, treesSheet.sprites)
	animalsBatch := pixel.NewBatch(&pixel.TrianglesData{}, animalsSheet.sprites)

	cfg := pixelgl.WindowConfig{
		Title:  "Simulated evolution",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		detectKeyPress(win, dt)

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		//-----------------------------------------
		//       World updates and rendering

		win.Clear(colornames.Whitesmoke)

		world.Update()

		drawWorld(cam, win, treesBatch, animalsBatch, treesSheet, animalsSheet, world)

		win.Update()
		//------------------------------------------

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
