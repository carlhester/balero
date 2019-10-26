package main

import "fmt"
import "net/http"
import "io/ioutil"
import "strconv"
import "sort"
import "strings"

import . "balero/sendalerts"
import . "balero/json2struct"

func main() {
	url := "http://api.bart.gov/api/etd.aspx?cmd=etd&orig=" + STATION + "&key=" + KEY + "&dir=" + DIR + "&json=y"
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

		if isTargetStation(train.Abbreviation) {
			targetTrains = append(targetTrains, train.Abbreviation)

			for _, est := range train.Est {
				targetMinutes = append(targetMinutes, est.Minutes)
			}
		}

	}

	fmt.Printf("\n%s ", targetTrains)
	intMinutes := convertStrMinutesToInt(targetMinutes)
	sort.Ints(intMinutes)
	fmt.Printf("%d\n", intMinutes)

	// cast output to strings to send in alert
	strTargetTrains := strings.Join(targetTrains, " ")
	//strMinutes := strconv.Itoa(intMinutes)
	alertMsg := fmt.Sprint(strTargetTrains, intMinutes)
	SendSNS(alertMsg, PHONE)

	for index, _ := range intMinutes[:len(intMinutes)-2] {
		twoTrainDelta := intMinutes[index+2] - intMinutes[index]
		if twoTrainDelta <= TIMEWIN {
			fmt.Printf("Match! %d %d %d : %d\n\n", intMinutes[index], intMinutes[index+1], intMinutes[index+2], twoTrainDelta)
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
