package main

import (
	"fmt"
	"math"
	"os"
)

const (
	g = 9.81
)

type Solver struct {
	prbl   Problem
	theta1 []float64
	theta2 []float64
	omega1 []float64
	omega2 []float64
}

func (s *Solver) clearSolver() {
	s.theta1 = nil
	s.theta2 = nil
	s.omega1 = nil
	s.omega2 = nil
	s.theta1 = append(s.theta1, float64(s.prbl.Segments[0].Theta0))
	s.theta2 = append(s.theta2, float64(s.prbl.Segments[1].Theta0))
	s.omega1 = append(s.omega1, float64(s.prbl.Segments[0].Omega0))
	s.omega2 = append(s.omega2, float64(s.prbl.Segments[1].Omega0))
}

func (s *Solver) setData(p Problem) {
	s.prbl = p
	if s.theta1 == nil {
		s.theta1 = append(s.theta1, float64(p.Segments[0].Theta0))
		s.theta2 = append(s.theta2, float64(p.Segments[1].Theta0))
		s.omega1 = append(s.omega1, float64(p.Segments[0].Omega0))
		s.omega2 = append(s.omega2, float64(p.Segments[1].Omega0))
	}
}

func theta1Derivative(omega1 float64) float64 {
	return omega1
}

func theta2Derivative(omega2 float64) float64 {
	return omega2
}

func omega1Derivative(
	theta1 float64,
	theta2 float64,
	omega1 float64,
	omega2 float64,
	firstSegMass float64,
	secondSegMass float64,
	firstSegLength float64,
	secondSegLength float64) float64 {
	return (9*omega1*omega1*firstSegLength*secondSegMass*math.Sin(2*theta1-2*theta2) +
		+12*secondSegMass*secondSegLength*omega2*omega2*math.Sin(theta1-theta2) +
		+12*(3*math.Sin(-2*theta2+theta1)*secondSegMass/4+math.Sin(theta1)*(secondSegMass+
			+5*secondSegMass/4))*g) / (firstSegLength * (9*secondSegMass*math.Cos(2*theta1-2*theta2) -
		-8*firstSegMass - 15*secondSegMass))
}

func omega2Derivative(theta1 float64, theta2 float64, omega1 float64, omega2 float64, firstSegMass float64, secondSegMass float64, firstSegLength float64, secondSegLength float64) float64 {
	return (-9*g*math.Sin(2*theta1-theta2)*(firstSegMass+2*secondSegMass) -
		9*secondSegMass*secondSegLength*omega2*omega2*math.Sin(2*theta1-2*theta2) -
		12*firstSegLength*omega1*omega1*(firstSegMass+3*secondSegMass)*math.Sin(theta1-theta2) +
		3*g*math.Sin(theta2)*(firstSegMass+6*secondSegMass)) / (secondSegLength *
		(9*secondSegMass*math.Cos(2*theta1-2*theta2) - 8*firstSegMass - 15*secondSegMass))
}

func (s *Solver) eulerSolver() {

	file, _ := os.Create("euler.txt")

	var steps = int(float32(s.prbl.Duration) / s.prbl.Solver.SolverParameters.Dt)
	for i := 1; i < steps; i++ {
		// prev iteration values
		var cTheta1 = s.theta1[i-1]
		var cTheta2 = s.theta2[i-1]
		var cOmega1 = s.omega1[i-1]
		var cOmega2 = s.omega2[i-1]

		// current iteration values
		var theta1 = cTheta1 + float64(s.prbl.Solver.SolverParameters.Dt)*theta1Derivative(cOmega1)
		var theta2 = cTheta2 + float64(s.prbl.Solver.SolverParameters.Dt)*theta2Derivative(cOmega2)
		var omega1 = cOmega1 + float64(s.prbl.Solver.SolverParameters.Dt)*omega1Derivative(cTheta1,
			cTheta2,
			cOmega1,
			cOmega2,
			float64(s.prbl.Segments[0].Mass),
			float64(s.prbl.Segments[1].Mass),
			float64(s.prbl.Segments[0].Length),
			float64(s.prbl.Segments[1].Length))
		var omega2 = cOmega2 + float64(s.prbl.Solver.SolverParameters.Dt)*omega2Derivative(cTheta1,
			cTheta2,
			cOmega1,
			cOmega2,
			float64(s.prbl.Segments[0].Mass),
			float64(s.prbl.Segments[1].Mass),
			float64(s.prbl.Segments[0].Length),
			float64(s.prbl.Segments[1].Length))

		// append current step values to solver slices
		s.theta1 = append(s.theta1, theta1)
		s.theta2 = append(s.theta2, theta2)
		s.omega1 = append(s.omega1, omega1)
		s.omega2 = append(s.omega2, omega2)
		file.WriteString(FloatToStr(s.theta1[i]) + " ")
	}
	//fmt.Print(len(s.theta1))

}

func (s *Solver) RungeKutta() {
	fmt.Println("USING RUNGE-KUTTA SOLVER")

	file, _ := os.Create("runge.txt")

	var steps = int(float32(s.prbl.Duration) / s.prbl.Solver.SolverParameters.Dt)
	for i := 1; i < steps; i++ {
		// prev iteration values
		var cTheta1 = s.theta1[i-1]
		var cTheta2 = s.theta2[i-1]
		var cOmega1 = s.omega1[i-1]
		var cOmega2 = s.omega2[i-1]

		var theta1K1 = theta1Derivative(cOmega1)
		var theta2K1 = theta2Derivative(cOmega2)
		var omega1K1 = omega1Derivative(cTheta1,
			cTheta2,
			cOmega1,
			cOmega2,
			float64(s.prbl.Segments[0].Mass),
			float64(s.prbl.Segments[1].Mass),
			float64(s.prbl.Segments[0].Length),
			float64(s.prbl.Segments[1].Length))
		var omega2K1 = omega2Derivative(cTheta1,
			cTheta2,
			cOmega1,
			cOmega2,
			float64(s.prbl.Segments[0].Mass),
			float64(s.prbl.Segments[1].Mass),
			float64(s.prbl.Segments[0].Length),
			float64(s.prbl.Segments[1].Length))

		var theta1K2 = theta1Derivative(cOmega1 + float64(s.prbl.Solver.SolverParameters.Dt)*omega1K1/2.)
		var theta2K2 = theta2Derivative(cOmega2 + float64(s.prbl.Solver.SolverParameters.Dt)*omega2K1/2.)
		var omega1K2 = omega1Derivative(cTheta1+float64(s.prbl.Solver.SolverParameters.Dt)*theta1K1/2.,
			cTheta2+float64(s.prbl.Solver.SolverParameters.Dt)*theta2K1/2.,
			cOmega1+float64(s.prbl.Solver.SolverParameters.Dt)*omega1K1/2.,
			cOmega2+float64(s.prbl.Solver.SolverParameters.Dt)*omega2K1/2.,
			float64(s.prbl.Segments[0].Mass),
			float64(s.prbl.Segments[1].Mass),
			float64(s.prbl.Segments[0].Length),
			float64(s.prbl.Segments[1].Length))
		var omega2K2 = omega2Derivative(cTheta1+float64(s.prbl.Solver.SolverParameters.Dt)*theta1K1/2.,
			cTheta2+float64(s.prbl.Solver.SolverParameters.Dt)*theta2K1/2.,
			cOmega1+float64(s.prbl.Solver.SolverParameters.Dt)*omega1K1/2.,
			cOmega2+float64(s.prbl.Solver.SolverParameters.Dt)*omega2K1/2.,
			float64(s.prbl.Segments[0].Mass),
			float64(s.prbl.Segments[1].Mass),
			float64(s.prbl.Segments[0].Length),
			float64(s.prbl.Segments[1].Length))

		var theta1K3 = theta1Derivative(cOmega1 + float64(s.prbl.Solver.SolverParameters.Dt)*omega1K2/2.)
		var theta2K3 = theta2Derivative(cOmega2 + float64(s.prbl.Solver.SolverParameters.Dt)*omega2K2/2.)
		var omega1K3 = omega1Derivative(cTheta1+float64(s.prbl.Solver.SolverParameters.Dt)*theta1K2/2.,
			cTheta2+float64(s.prbl.Solver.SolverParameters.Dt)*theta2K2/2.,
			cOmega1+float64(s.prbl.Solver.SolverParameters.Dt)*omega1K2/2.,
			cOmega2+float64(s.prbl.Solver.SolverParameters.Dt)*omega2K2/2.,
			float64(s.prbl.Segments[0].Mass),
			float64(s.prbl.Segments[1].Mass),
			float64(s.prbl.Segments[0].Length),
			float64(s.prbl.Segments[1].Length))
		var omega2K3 = omega2Derivative(cTheta1+float64(s.prbl.Solver.SolverParameters.Dt)*theta1K2/2.,
			cTheta2+float64(s.prbl.Solver.SolverParameters.Dt)*theta2K2/2.,
			cOmega1+float64(s.prbl.Solver.SolverParameters.Dt)*omega1K2/2.,
			cOmega2+float64(s.prbl.Solver.SolverParameters.Dt)*omega2K2/2.,
			float64(s.prbl.Segments[0].Mass),
			float64(s.prbl.Segments[1].Mass),
			float64(s.prbl.Segments[0].Length),
			float64(s.prbl.Segments[1].Length))

		var theta1K4 = theta1Derivative(cOmega1 + float64(s.prbl.Solver.SolverParameters.Dt)*omega1K3)
		var theta2K4 = theta2Derivative(cOmega2 + float64(s.prbl.Solver.SolverParameters.Dt)*omega2K3)
		var omega1K4 = omega1Derivative(cTheta1+float64(s.prbl.Solver.SolverParameters.Dt)*theta1K3,
			cTheta2+float64(s.prbl.Solver.SolverParameters.Dt)*theta2K3,
			cOmega1+float64(s.prbl.Solver.SolverParameters.Dt)*omega1K3,
			cOmega2+float64(s.prbl.Solver.SolverParameters.Dt)*omega2K3,
			float64(s.prbl.Segments[0].Mass),
			float64(s.prbl.Segments[1].Mass),
			float64(s.prbl.Segments[0].Length),
			float64(s.prbl.Segments[1].Length))
		var omega2K4 = omega2Derivative(cTheta1+float64(s.prbl.Solver.SolverParameters.Dt)*theta1K3,
			cTheta2+float64(s.prbl.Solver.SolverParameters.Dt)*theta2K3,
			cOmega1+float64(s.prbl.Solver.SolverParameters.Dt)*omega1K3,
			cOmega2+float64(s.prbl.Solver.SolverParameters.Dt)*omega2K3,
			float64(s.prbl.Segments[0].Mass),
			float64(s.prbl.Segments[1].Mass),
			float64(s.prbl.Segments[0].Length),
			float64(s.prbl.Segments[1].Length))

		var theta1 = cTheta1 + float64(s.prbl.Solver.SolverParameters.Dt)*(theta1K1+2.*theta1K2+2.*theta1K3+theta1K4)/6.
		var theta2 = cTheta2 + float64(s.prbl.Solver.SolverParameters.Dt)*(theta2K1+2.*theta2K2+2.*theta2K3+theta2K4)/6.
		var omega1 = cOmega1 + float64(s.prbl.Solver.SolverParameters.Dt)*(omega1K1+2.*omega1K2+2*omega1K3+omega1K4)/6.
		var omega2 = cOmega2 + float64(s.prbl.Solver.SolverParameters.Dt)*(omega2K1+2.*omega2K2+2*omega2K3+omega2K4)/6.

		// append current step values to solver slices
		s.theta1 = append(s.theta1, theta1)
		s.theta2 = append(s.theta2, theta2)
		s.omega1 = append(s.omega1, omega1)
		s.omega2 = append(s.omega2, omega2)

		file.WriteString(FloatToStr(s.theta1[i]) + " ")
	}
	fmt.Print(len(s.theta1))
}
