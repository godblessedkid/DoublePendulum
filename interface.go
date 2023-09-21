package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

func FloatToStr(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 2, 64)
}

// Problem struct
type Problem struct {
	Title    string `json:"title"`
	Segments []struct {
		Length float32 `json:"length"`
		Mass   float32 `json:"mass"`
		Theta0 float32 `json:"theta0"`
		Omega0 float32 `json:"omega0"`
	} `json:"segments"`
	Duration int `json:"duration"`
	Solver   struct {
		SolverType       string `json:"type"`
		SolverParameters struct {
			Dt float32 `json:"dt"`
		} `json:"parameters"`
	} `json:"solver"`
}

// load Json by its path to a Problem object
func loadJson(pathToJson string, prbl *Problem) {

	jsonFile, err := os.Open(pathToJson)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened " + pathToJson)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &prbl)
}

func makeStringListFromProblem(prbl *Problem) []string {

	var strList []string
	if prbl.Duration == 0 {
		return strList
	}

	strList = append(
		strList,
		"Problem Parameters :",
		"-------------------",
		"Segment 1: ",
		"Length: "+FloatToStr(float64(prbl.Segments[0].Length))+" m",
		"Mass: "+FloatToStr(float64(prbl.Segments[0].Mass))+" kg",
		"Omega0: "+FloatToStr(float64(prbl.Segments[0].Omega0))+" rad",
		"Theta0: "+FloatToStr(float64(prbl.Segments[0].Theta0))+" rad",
		"-------------------",
		"Segment 2: ",
		"Length: "+FloatToStr(float64(prbl.Segments[1].Length))+" m",
		"Mass: "+FloatToStr(float64(prbl.Segments[1].Mass))+" kg",
		"Omega0: "+FloatToStr(float64(prbl.Segments[1].Omega0))+" rad",
		"Theta0: "+FloatToStr(float64(prbl.Segments[1].Theta0))+" rad",
		"-------------------",
		"Duration: "+strconv.Itoa(prbl.Duration)+" sec",
		"Solver method: "+prbl.Solver.SolverType,
		"dt: "+FloatToStr(float64(prbl.Solver.SolverParameters.Dt)))

	return strList

}
