package main

import (
	"os"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
	"text/template"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Value struct {
	Status Status `json:"status"`
}

func main() {
	go autoReload()

	http.HandleFunc("/", rendering)
	// menggunakan port :80 karena localhost:8080 error di laptop saya
	http.ListenAndServe(":80", nil)
}

func rendering(w http.ResponseWriter, r *http.Request) {
	inputJSON, err := os.ReadFile("input.json")
	if err != nil {
		log.Fatal("Error read input.json")
	}

	var input Value
	err = json.Unmarshal(inputJSON, &input)
	if err != nil {
		log.Fatal("Error uncompiling input.json")
	}

	t, err := template.ParseFiles("template.html")
	if err != nil {
		log.Fatal("Error parsing index.html")
	}

	var statusWater string 
	var statusWind string

	var water = input.Status.Water
	var wind = input.Status.Wind

	if water < 6 {
		statusWater = "Aman"
	} else if water >= 6 && water <= 8 {
		statusWater = "Siaga"
	} else {
		statusWater = "Bahaya"
	}

	if wind < 7 {
		statusWind = "Aman"
	} else if wind >= 7 && wind <= 15 {
		statusWind = "Siaga"
	} else {
		statusWind = "Bahaya"
	}

	output := map[string]interface{}{
		"water":       water,
		"wind":        wind,
		"statusWater": statusWater,
		"statusWind":  statusWind,
	}

	t.Execute(w, output)
}

func autoReload(){
	
	var	min = 1
	var max = 100


	for {
		input := Value{
			Status: Status{
				// untuk memberi nomor acak
				Water	:  rand.Intn(max-min) + min,
				Wind	:  rand.Intn(max-min) + min,
			},
		}

		inputJSON, err := json.MarshalIndent(&input, "", "  ")
		if err != nil {
			log.Fatal("Error compile input json")
		}

		err = os.WriteFile("input.json", inputJSON, 0644)
		if err != nil {
			log.Fatal("Error writing input.json")
		}

		time.Sleep(15 * time.Second)
	}
}
