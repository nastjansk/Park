package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"
)

// Формат вывода даты и времени
const layout  = "2006-01-02 15:04:05"

// Начало периода парковки
var minT = time.Date(2021, 7, 29, 0, 0, 0, 0, time.Local).Unix()
// Окончание периода парковки
var maxT = time.Date(2021, 7, 30, 0, 0, 0, 0, time.Local).Unix()
var delta = maxT - minT

//  Кол-во заездов/выездов за период (объем генерируемых данных)
var carCount = 2000

// Структуры для считывания и хранения данных
type parkInfo struct {
	ParkInfo []parkTime `json:"parkInfo"`
}

// Пара - въезд/выезд
type parkTime struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Генератор интервала въезда/выезда
func genParkInterval() parkTime {
	t1 := rand.Int63n(delta) + minT
	t2 := rand.Int63n(delta) + minT

	if t2 < t1 {
		t2,t1 = t1,t2
	}

	return parkTime{time.Unix(t1, 0), time.Unix(t2, 0)}
}

// Генератор файла данных
func genData(fName string) {
	data := parkInfo{}

	for i := 0; i < carCount; i++ {
		data.ParkInfo = append(data.ParkInfo, genParkInterval())
	}

	file, _ := json.MarshalIndent(data, "", " ")

	_ = ioutil.WriteFile(fName, file, 0644)
}

func main() {

	fName := "ParkingInfo"+strconv.Itoa(carCount)+".json"

	// Gen parking data
	genData(fName)

	data := parkInfo{}
	file, _ := ioutil.ReadFile(fName)

	_ = json.Unmarshal([]byte(file), &data)

	var maxAuto int = 0

	// FullScan (seconds)
	for i := minT; i < maxT; i ++ {
		countInMoment := 0
		// В каждую секунду смотрим сколько машин на парковке
		for _, k := range data.ParkInfo {
			if k.Start.Unix() <= i && k.End.Unix() >= i {
				countInMoment = countInMoment + 1
			}
		}
		if countInMoment > maxAuto {
			maxAuto = countInMoment
		}
	}

	fmt.Printf("MaxAuto = %v\n", maxAuto)

}
