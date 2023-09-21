package main

import (
	"image"
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

const (
	frameTime = 10
)

type pendulumDrawingWidget struct {
	widget.BaseWidget

	theta1 float64
	theta2 float64
	pr     *Problem
	raster *canvas.Raster
}

func (pdw *pendulumDrawingWidget) setNewAngles(newTheta1 float64, newTheta2 float64) {
	pdw.theta1 = newTheta1
	pdw.theta2 = newTheta2
}

func newPendulumDrawingWidget(cTheta1 float64, cTheta2 float64, p *Problem) *pendulumDrawingWidget {
	pdw := &pendulumDrawingWidget{theta1: cTheta1, theta2: cTheta2, pr: p}
	pdw.raster = canvas.NewRaster(pdw.drawFrame)
	pdw.ExtendBaseWidget(pdw)
	return pdw
}

func (pdw *pendulumDrawingWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(pdw.raster)
}

func (pdw *pendulumDrawingWidget) animate(s *Solver) {
	var amountOfFrames = len(s.theta1)
	for i := 0; i < amountOfFrames; i++ {
		start := time.Now().UnixMilli()
		pdw.setNewAngles(s.theta1[i], s.theta2[i])
		pdw.Refresh()
		end := time.Now().UnixMilli()
		diff := end - start
		time.Sleep(time.Millisecond * time.Duration(frameTime-diff))
	}
}

func (pwd *pendulumDrawingWidget) drawFrame(w, h int) image.Image {
	var prbl = pwd.pr
	var theta1 = pwd.theta1
	var theta2 = pwd.theta2

	dest := image.NewRGBA(image.Rect(0, 0, 2000, 1800))
	gc := draw2dimg.NewGraphicContext(dest)
	if prbl == nil {
		return dest
	}
	var fixedPointX float64 = 1000.
	var fixedPointY float64 = 900.
	var firstSegmentLenght = float64(prbl.Segments[0].Length) * 80
	var secondSegmentLenght = float64(prbl.Segments[1].Length) * 80
	var firstSegmentEndX = float64(fixedPointX + firstSegmentLenght*math.Sin(theta1))
	var firstSegmentEndY = float64(fixedPointY + firstSegmentLenght*math.Cos(theta1))
	var secondSegmentEndX = float64(firstSegmentEndX + secondSegmentLenght*math.Sin(theta2))
	var secondSegmentEndY = float64(firstSegmentEndY + secondSegmentLenght*math.Cos(theta2))
	// Set some properties
	gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
	gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0xff})
	gc.SetLineWidth(10)
	gc.BeginPath()
	gc.MoveTo(fixedPointX, fixedPointY)
	draw2dkit.RoundedRectangle(gc, fixedPointX-5, fixedPointY-5, fixedPointX+5, fixedPointY+5, 2, 2)
	gc.FillStroke()
	gc.MoveTo(fixedPointX, fixedPointY)
	gc.LineTo(firstSegmentEndX, firstSegmentEndY)
	gc.FillStroke()
	draw2dkit.Circle(gc, firstSegmentEndX, firstSegmentEndY, 2)
	gc.FillStroke()
	gc.MoveTo(firstSegmentEndX, firstSegmentEndY)
	gc.LineTo(secondSegmentEndX, secondSegmentEndY)
	gc.FillStroke()
	return dest
}
