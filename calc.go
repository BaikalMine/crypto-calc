package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {

		userHash := r.URL.Query().Get("hashrate")

		resp, err := http.Get("https://etherchain.org/api/basic_stats")
		if err != nil {
			log.Fatalln(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		var result map[string]interface{}
		jsonErr := json.Unmarshal(body, &result)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		stats := result["currentStats"].(map[string]interface{})
		difficulty := stats["difficulty"]
		price_usd := stats["price_usd"]
		block_time := stats["block_time"]
		block_reward := stats["block_reward"]

		blockTime := block_time.(float64)
		diff := int(difficulty.(float64))
		networkhashrate := (float64(diff) / blockTime) / 1e9

		userhashrate, err := strconv.ParseInt(userHash, 10, 64)
		if err != nil {
			panic(err)
		}

		userhash := float64(userhashrate * 1e6)
		netHash := float64(networkhashrate * 1e9)

		userRatio := userhash / netHash
		blocksPerMin := 60.0 / blockTime
		ethPerMin := blocksPerMin * block_reward.(float64)
		earningsmin := userRatio * ethPerMin
		earningshour := earningsmin * 60
		earningsday := earningshour * 24
		earningsweek := earningsday * 7
		earningsmonth := earningsday * 30
		earningsyear := earningsday * 365

		priceUSD := price_usd.(float64)
		mUSD := earningsmin * priceUSD
		hUSD := earningshour * priceUSD
		dUSD := earningsday * priceUSD
		wUSD := earningsweek * priceUSD
		MUSD := earningsmonth * priceUSD
		yUSD := earningsyear * priceUSD

		fmt.Fprintf(w, " Минута %f ETH  %f USD", earningsmin, mUSD)
		fmt.Fprintf(w, " Час %f ETH  %f USD", earningshour, hUSD)
		fmt.Fprintf(w, " День %f ETH  %f USD", earningsday, dUSD)
		fmt.Fprintf(w, " Неделя %f ETH  %f USD", earningsweek, wUSD)
		fmt.Fprintf(w, " Месяц %f ETH  %f USD", earningsmonth, MUSD)
		fmt.Fprintf(w, " Год %f ETH  %f USD", earningsyear, yUSD)

		message := make(map[string]interface{})
		messageETH := make(map[string]interface{})
		messageUSD := make(map[string]interface{})

		messageETH["Minute"] = earningsmin
		messageETH["Hour"] = earningshour
		messageETH["Day"] = earningsday
		messageETH["Mounth"] = earningsmonth
		messageETH["Year"] = earningsyear

		messageUSD["Minute"] = mUSD
		messageUSD["Hour"] = hUSD
		messageUSD["Day"] = dUSD
		messageUSD["Mounth"] = MUSD
		messageUSD["Year"] = yUSD

		message["ETH"] = messageETH
		message["USD"] = messageUSD

		bytesRepresentation, err := json.Marshal(message)
		if err != nil {
			log.Fatalln(err)
		}

		resp1, err := http.Post("http://localhost:8182/post", "application/json", bytes.NewBuffer(bytesRepresentation))
		if err != nil {
			panic(err)
		}

		var result1 map[string]interface{}

		json.NewDecoder(resp1.Body).Decode(&result1)

	})
	fmt.Println("Enter this in your browser --->>> http://localhost:8181/user?hashrate=@YourHashrate(MH)@")
	http.ListenAndServe(":8181", nil)

}
