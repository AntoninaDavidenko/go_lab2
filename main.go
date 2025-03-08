package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type CalculationResult struct {
	KCoal    float64
	ECoal    float64
	KFuelOil float64
	EFuelOil float64
}

func calculate(coal, fuelOil float64) CalculationResult {
	const (
		lowerHeatOfCombustionCoal     = 20.47
		ashCarriedOutCoal             = 0.80
		efficiencyOfCleaning          = 0.985
		coalA                         = 25.20
		combustibleSubstancesInFlyAsh = 1.5
		lowerHeatOfCombustionFuelOil  = 40.40
		ashCarriedOutFuelOil          = 1.0
		ashContentOfDryMass           = 0.15
	)

	kCoal := (math.Pow(10, 6) / lowerHeatOfCombustionCoal) *
		ashCarriedOutCoal *
		(coalA / (100 - combustibleSubstancesInFlyAsh)) *
		(1 - efficiencyOfCleaning)
	eCoal := math.Pow(10, -6) * kCoal * lowerHeatOfCombustionCoal * coal

	kFuelOil := (math.Pow(10, 6) / lowerHeatOfCombustionFuelOil) *
		ashCarriedOutFuelOil *
		(ashContentOfDryMass / 100) *
		(1 - efficiencyOfCleaning)
	eFuelOil := math.Pow(10, -6) * kFuelOil * lowerHeatOfCombustionFuelOil * fuelOil

	return CalculationResult{kCoal, eCoal, kFuelOil, eFuelOil}
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	if r.Method == http.MethodPost {
		r.ParseForm()
		coal, _ := strconv.ParseFloat(r.FormValue("coal"), 64)
		fuelOil, _ := strconv.ParseFloat(r.FormValue("fuelOil"), 64)
		result := calculate(coal, fuelOil)
		tmpl.Execute(w, result)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
