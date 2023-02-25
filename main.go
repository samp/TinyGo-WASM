// Set environment and build target
// $Env:GOOS = "js"; $Env:GOARCH = "wasm"
// Go: Build canvas_go.wasm from main.go
// go build -o canvas_go.wasm main.go
// TinyGo: Build canvas_tiny.wasm from main.go
// tinygo build -o canvas_tiny.wasm -target wasm ./main.go

package main

import (
	// This library bridges the gap between Go and the browser environment
	// https://golang.org/pkg/syscall/js
	"syscall/js"
)

var (
	// js.Value can be any JS object/type/constructor
	window, doc, body, canvas, ctx, buttonWrapper, buttonClear js.Value
	windowSize                                                 struct{ w, h float64 }
	// gs is at the highest scope, all others can access it
	gs = state{brushSize: 5, brushX: 0, brushY: 0, mouseDown: false}
)

func main() {
	// Create empty channel
	// https://stackoverflow.com/a/47262117
	forever := make(chan bool)

	setup()

	// Render every frame
	var renderer js.Func
	renderer = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		window.Call("requestAnimationFrame", renderer)
		return nil
	})
	window.Call("requestAnimationFrame", renderer)

	// Handle the pointerdown event
	var mouseDownEventHandler js.Func = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		mouseX := args[0].Get("clientX").Float()
		mouseY := args[0].Get("clientY").Float()
		if mouseX < 100 && mouseY < 50 {
			go clear()
		} else {
			gs.mouseDown = true
			gs.brushX = mouseX
			gs.brushY = mouseY
			go draw(mouseX, mouseY)
		}
		return nil
	})
	window.Call("addEventListener", "pointerdown", mouseDownEventHandler)

	// Handle the pointerup event
	var mouseUpEventHandler js.Func = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		mouseX := args[0].Get("clientX").Float()
		mouseY := args[0].Get("clientY").Float()
		if mouseX < 100 && mouseY < 50 {
			go clear()
		} else {
			gs.mouseDown = false
			gs.brushX = mouseX
			gs.brushY = mouseY
			go draw(mouseX, mouseY)
		}
		return nil
	})
	window.Call("addEventListener", "pointerup", mouseUpEventHandler)

	// Handle the pointermove event
	var mouseMoveEventHandler js.Func = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		mouseX := args[0].Get("clientX").Float()
		mouseY := args[0].Get("clientY").Float()
		if gs.mouseDown == true {
			go drawWithArc(gs.brushX, gs.brushY, mouseX, mouseY)
			gs.brushX = mouseX
			gs.brushY = mouseY
		}
		return nil
	})
	window.Call("addEventListener", "pointermove", mouseMoveEventHandler)

	// Attempt to receive from empty channel
	// Forces program to stay "awake" forever
	<-forever
}

func setup() {
	// Get all the document info from js
	window = js.Global()
	doc = window.Get("document")
	body = doc.Get("body")
	windowSize.h = window.Get("innerHeight").Float()
	// Enable if using HTML buttons
	//windowSize.h = window.Get("innerHeight").Float() - 100
	windowSize.w = window.Get("innerWidth").Float()

	// Create canvas element
	canvas = doc.Call("createElement", "canvas")
	canvas.Set("height", windowSize.h)
	canvas.Set("width", windowSize.w)
	body.Call("appendChild", canvas)

	// Create HTML buttons
	// Disabled for now - it's not any faster than canvas buttons
	/*buttonWrapper = doc.Call("createElement", "div")
	buttonWrapper.Set("style", "display:flex")
	buttonClear = doc.Call("createElement", "button")
	buttonClear.Set("innerText", "Clear")
	body.Call("appendChild", buttonWrapper)
	buttonWrapper.Call("appendChild", buttonClear)*/

	// Bring canvas to code
	ctx = canvas.Call("getContext", "2d")

	// Draw the clear button
	go drawButtonClear()
}

// Draw a single point
func draw(x float64, y float64) {
	ctx.Set("fillStyle", "black")
	ctx.Call("beginPath")
	// Arc between the previous point and the current position.
	// Stops the broken line effect when the mouse moves faster than the browser can draw
	ctx.Call("arc", x, y, gs.brushSize, 0, 3.14159*2, false)
	ctx.Call("fill")
	ctx.Call("closePath")

	// Every frame, redraw the clear button - keeps it on top of the brush
	go drawButtonClear()
}

// Draw a line
func drawWithArc(oldX float64, oldY float64, newX float64, newY float64) {
	ctx.Call("beginPath")
	ctx.Call("moveTo", oldX, oldY)
	ctx.Call("lineTo", newX, newY)
	ctx.Set("lineCap", "round")
	ctx.Set("lineWidth", gs.brushSize*2)
	ctx.Call("stroke")
	ctx.Call("closePath")

	// Every frame, redraw the clear button - keeps it on top of the brush
	go drawButtonClear()
}

func clear() {
	// Clear the screen
	ctx.Call("clearRect", 0, 0, windowSize.w, windowSize.h)

	// Draw the clear button
	go drawButtonClear()
	gs.mouseDown = false
}

func drawButtonClear() {
	// Draw the clear button
	// Box
	ctx.Call("beginPath")
	ctx.Call("rect", 3, 3, 100, 50)
	ctx.Set("fillStyle", "white")
	ctx.Call("fill")
	ctx.Set("lineWidth", 2)
	ctx.Set("strokeStyle", "black")
	ctx.Call("stroke")
	ctx.Call("closePath")

	// Text
	ctx.Call("beginPath")
	ctx.Set("font", "20px sans-serif")
	ctx.Set("fillStyle", "black")
	ctx.Call("fillText", "Clear", 25, 35)
	ctx.Call("closePath")
}

// This is the global state storage
type state struct {
	brushX, brushY, brushSize float64
	mouseDown                 bool
}
