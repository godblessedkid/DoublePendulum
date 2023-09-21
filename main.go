package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// current Problem
	var prbl Problem
	paramList := binding.NewStringList()
	paramList.Append("Upload some parameters to solver!")

	// App initialization
	myApp := app.New()
	w := myApp.NewWindow("Double Pendulum")
	pendulumParameters := myApp.NewWindow("Parameters")
	pendulumParameters.SetCloseIntercept(func() {
		pendulumParameters.Hide()
	})
	pendulumParameters.Resize(fyne.NewSize(300, 670))
	w.Resize(fyne.NewSize(1000, 900))
	var solver Solver
	// Main Menu
	loadJsonMenuItem := fyne.NewMenuItem("Load Json", func() {
		fileDialog := dialog.NewFileOpen(
			func(uc fyne.URIReadCloser, _ error) {
				if uc != nil {
					jsonFilepath := uc.URI().Path()
					fmt.Println(jsonFilepath)
					loadJson(jsonFilepath, &prbl)
					paramListToSet := makeStringListFromProblem(&prbl)
					paramList.Set(paramListToSet)
					pendulumParameters.Show()
					//Solver
					solver.setData(prbl)
					pdw := newPendulumDrawingWidget(solver.theta1[0], solver.theta2[0], &prbl)
					w.SetContent(pdw)
				}
			}, w)
		fileDialog.SetFilter(
			storage.NewExtensionFileFilter([]string{".json"}))
		fileDialog.Show()
	})

	CheckItem := fyne.NewMenuItem("Show loaded parameters", func() {
		pendulumParameters.Show()
	})

	SolveItem := fyne.NewMenuItem("Start Simulation", func() {
		if solver.theta1 != nil {
			solver.clearSolver()
			var solverType string = solver.prbl.Solver.SolverType
			switch solverType {
			case "euler":
				solver.eulerSolver()
			case "Runge-Kutta":
				solver.RungeKutta()
			default:
				solver.eulerSolver()
			}
			pdw := newPendulumDrawingWidget(solver.theta1[0], solver.theta2[0], &prbl)
			pdw.Resize(fyne.NewSize(495, 450))
			w.SetContent(pdw)
			pdw.animate(&solver)

		}
	})

	//Parameters list item
	list := widget.NewListWithData(paramList,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	list.Resize(fyne.NewSize(500, 500))
	list.UnselectAll()
	content := container.New(layout.NewAdaptiveGridLayout(1), list)
	pendulumParameters.SetContent(content)

	menu := fyne.NewMenu("Options", loadJsonMenuItem, CheckItem /*StartItem,*/, SolveItem)
	mainMenu := fyne.NewMainMenu(menu)
	w.SetMainMenu(mainMenu)

	w.Show()
	w.SetMaster()
	myApp.Run()
}
