package main

import "fmt"
import "net/http"
import "io/ioutil"
import "strconv"
import "strings"
import "sort"
import "time"

import "balero/config"
import . "balero/sendalerts"
import . "balero/json2struct"

func main() {
	url := "http://api.bart.gov/api/etd.aspx?cmd=etd&orig=" + config.STATION + "&key=" + config.KEY + "&dir=" + config.DIR + "&json=y"
	resp, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		panic(err.Error())
	}

	usableData := RawDataIntoDataStruct(data)

	var targetTrains []string
	var targetMinutes []string

	for _, train := range usableData.Root.Station[0].Etd {

		for _, est := range train.Est {

			if isTargetLine(est.Color) {
				targetTrains = append(targetTrains, train.Abbreviation)
				targetMinutes = append(targetMinutes, est.Minutes)
			}
		}
	}

	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err.Error())
	}
	currTime := time.Now()
	currTime = currTime.In(loc)
	fmt.Printf(currTime.Format("Jan _2 15:04:05"))

	fmt.Printf("\ntargetTrains: %s ", targetTrains)
	intMinutes := convertStrMinutesToInt(targetMinutes)
	sort.Ints(intMinutes)
	fmt.Printf("intMinutes: %d\n", intMinutes)

	if len(intMinutes) > 2 {
		for i, _ := range intMinutes[:len(intMinutes)-2] {
			twoTrainDelta := intMinutes[i+2] - intMinutes[i]
			if twoTrainDelta <= config.TIMEWIN {
				fmt.Printf("Match! %d %d %d : %d\n\n", intMinutes[i], intMinutes[i+1], intMinutes[i+2], twoTrainDelta)
				alertMsg := fmt.Sprintf("%s %d %d %d : %d", targetTrains, intMinutes[i], intMinutes[i+1], intMinutes[i+2], twoTrainDelta)
				SendSNS(alertMsg)
			}
		}
	}
}

func convertStrMinutesToInt(minutes []string) []int {
	var intMinutes []int
	for _, strMin := range minutes {
		if strMin == "Leaving" {
			strMin = "0"
		}
		i, err := strconv.Atoi(strMin)
		if err != nil {
			panic(err.Error())
		}
		intMinutes = append(intMinutes, i)
	}
	return intMinutes
}

func isTargetStation(station string) bool {
	switch station {
	case
		"ANTC", "CONC", "NCON", "PITT", "PHIL", "WCRK":
		return true
	}
	return false
}

func isTargetLine(line string) bool {
	for _, lineColor := range config.TargetStations {
		if strings.EqualFold(lineColor, line) {
			return true
		}
	}
	return false
}
