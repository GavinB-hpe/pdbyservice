package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gavinB-hpe/pdbyservice/dbtalker"
	"github.com/gavinB-hpe/pdbyservice/globals"

	"github.com/gavinB-hpe/pdbyservice/model"
	"github.com/gavinB-hpe/pdbyservice/pdanalyser"
	"github.com/gavinB-hpe/pdbyservice/utils"
)

// var for flags
var dbtype string
var dbdetails string
var bucketsize int
var unknownservicelistfilename string
var servicedatafilename string
var productionOnly bool
var saasonly bool
var onpremonly bool
var daysback int
var showincidents bool
var skipresolved bool
var maxcolwidth int

func prettifiedOutput(sc map[string]int, sn map[string]string, keys []string) {
	toto := make(map[string]int, 0)
	for _, k := range keys {
		serviceref := fmt.Sprintf("%s : %s", k, sn[k])
		toto[serviceref] = sc[k]
	}
	utils.PrintMapIntAsSortedTable("Service", "#Incidents", toto)
}

func saveToFileAsJson(fn string, sn map[string]string) {
	fl, err := os.Create(fn)
	if err != nil {
		log.Panic(err)
	}
	defer fl.Close()
	encoder := json.NewEncoder(fl)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(sn); err != nil {
		log.Fatalf("could not encode map to JSON: %v", err)
	}
}

func readServiceData(fn string) *map[string]map[string]string {
	fl, err := os.Open(fn)
	if err != nil {
		log.Panic(err)
	}
	defer fl.Close()
	decoder := json.NewDecoder(fl)
	var toto map[string]map[string]string
	err = decoder.Decode(&toto)
	if err != nil {
		log.Panic(err)
	}
	return &toto
}

func checkSeenServices(sn map[string]string, sd *map[string]map[string]string) map[string]string {
	unknown := make(map[string]string)
	for k, v := range sn {
		mp := (*sd)[k]
		if mp == nil {
			fmt.Println("WARNING : Unknown service ", k)
			unknown[k] = v
		}
	}
	return unknown
}

func getPerServiceIncidents(sr bool, key string, dbt *dbtalker.DBTalker) [][]string {
	block := make([][]string, 0)
	// fmt.Println("-----------------------------------------------------------------")
	// fmt.Println(k, "  ==>  ")
	var incidents []model.PDInfoType
	dbt.DB.Model(model.PDInfoType{}).Where("service_id = ?", key).Find(&incidents)
	for _, i := range incidents {
		if sr && i.Status != "resolved" || !sr {
			block = append(block, []string{i.ID, i.Summary, i.CreatedAt, i.Priority, i.Urgency, i.Status, i.ServiceName})
		}
	}
	return block
}

func printIncidents(sr bool, keys []string, mxc int, dbt *dbtalker.DBTalker) {
	headers := []string{"ID", "Summary", "CreatedAt", "Priority", "Urgency", "Status", "ServiceName"}
	for _, k := range keys {
		block := getPerServiceIncidents(sr, k, dbt)
		if len(block) <= 0 {
			fmt.Println("No incidents for service ", k)
			return
		}
		fmt.Println("Incidents for service", k, "=", block[0][len(headers)-1])
		utils.Print2DArrayAsTable(mxc, headers, block)
		fmt.Println()
	}
}

func main() {
	flag.StringVar(&dbtype, "t", globals.DEFAULTDBTYPE, "Type of DB used e.g. sqlite3")
	flag.StringVar(&dbdetails, "db", globals.DEFAULTDBDETAILS, "Filename for sqlite3 or URI of DB")
	flag.IntVar(&bucketsize, "b", globals.DEFAULTBUCKETSIZE, "How many days to bucket together in the graph. ")
	flag.IntVar(&daysback, "D", 30, "How many days back to search. Cannot go further back than the data in the DB of course.")
	flag.IntVar(&maxcolwidth, "c", globals.MAXCOLWIDTH, "Max width of column in characters")
	flag.StringVar(&unknownservicelistfilename, "o", globals.DEFAULTUNKNOWNSERVICELIST, "File used to store list of unknown services seen")
	flag.StringVar(&servicedatafilename, "d", globals.DEFAULTSERVICEDATAFILENAME, "File with service data")
	flag.BoolVar(&productionOnly, "P", false, "If set, only record data for production services")
	flag.BoolVar(&saasonly, "S", false, "If set, only record data for services running on SaaS")
	flag.BoolVar(&onpremonly, "O", false, "If set, only record data for services running OnPrem")
	flag.BoolVar(&showincidents, "s", false, "If set, show the incident data")
	flag.BoolVar(&skipresolved, "R", false, "If set, skip incidents that are resolved")
	flag.Parse()

	servicedata := readServiceData(servicedatafilename)
	if bucketsize <= 0 {
		log.Fatalln("Invalid bucketsize value. Must be > 0")
	}
	if dbdetails == "" {
		log.Fatalln("dbdetails cannot be empty")
	}
	dbtalker := dbtalker.NewDBTalker(model.ConnectDatabase(dbtype, dbdetails))
	// get data
	scounts, snames, sortedkeys := pdanalyser.PDanalyse(productionOnly, saasonly, onpremonly, daysback, skipresolved, servicedata, dbtalker)
	// output
	prettifiedOutput(scounts, snames, sortedkeys)
	unknown := checkSeenServices(snames, servicedata)
	saveToFileAsJson(unknownservicelistfilename, unknown)
	// print out incidents
	if showincidents {
		printIncidents(skipresolved, sortedkeys, maxcolwidth, dbtalker)
	}
}
