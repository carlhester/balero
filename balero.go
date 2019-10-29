package main

import "flag"
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
	station := config.Station
	KEY := config.KEY
	dir := config.Dir
	targetLine := config.TargetLine
	timeWindow := config.TimeWindow
	sendSMS := config.SendSMS

	flag.StringVar(&station, "station", "MONT", "Starting station abbreviation")
	flag.StringVar(&dir, "dir", "n", "Train direction")
	flag.StringVar(&targetLine, "line", "YELLOW", "Target Line")
	flag.BoolVar(&sendSMS, "SMS", false, "Send SMS alerts (if configured)")
	flag.IntVar(&timeWindow, "time", 15, "Time window threshold for a match")
	flag.Parse()

	url := "http://api.bart.gov/api/etd.aspx?cmd=etd&orig=" + station + "&key=" + KEY + "&dir=" + dir + "&json=y"
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
			if strings.EqualFold(est.Color, targetLine) {
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
	timeStamp := fmt.Sprintf("%s", currTime.Format("Jan _2 15:04:05"))

	intMin := convertStrMinutesToInt(targetMinutes)
	sort.Ints(intMin)

	fmt.Printf("Time: %s\n", timeStamp)
	fmt.Printf("targetTrains: %s \n", targetTrains)
	fmt.Printf("intMin: %d \n", intMin)

	alertMsg := timeStamp
	numResults := 0

	if len(intMin) > 2 {
		for i, _ := range intMin[:len(intMin)-2] {
			twoTrainDelta := intMin[i+2] - intMin[i]
			if twoTrainDelta <= timeWindow {
				fmt.Printf("Match! %d %d %d : %d\n\n", intMin[i], intMin[i+1], intMin[i+2], twoTrainDelta)
				partAlertMsg := fmt.Sprintf("%s %d %d %d : %d", targetTrains, intMin[i], intMin[i+1], intMin[i+2], twoTrainDelta)
				alertMsg = fmt.Sprintf("%s\n%s\n", alertMsg, partAlertMsg)
				numResults += 1
			}
		}
	}
	if numResults > 0 {
		fmt.Printf("Alert:\n%s", alertMsg)
		if sendSMS == true {
			SendSNS(alertMsg)
		}
	}
}

func convertStrMinutesToInt(minutes []string) []int {
	var intMin []int
	for _, strMin := range minutes {
		if strMin == "Leaving" {
			strMin = "0"
		}
		i, err := strconv.Atoi(strMin)
		if err != nil {
			panic(err.Error())
		}
		intMin = append(intMin, i)
	}
	return intMin
}
