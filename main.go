package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gavinB-hpe/pdbyservice/dbtalker"
	"github.com/gavinB-hpe/pdbyservice/globals"

	"github.com/gavinB-hpe/pdbyservice/model"
	"github.com/gavinB-hpe/pdbyservice/pdanalyser"
	"github.com/gavinB-hpe/pdbyservice/utils"
)

type ArrayFlags []string

// var for flags
var dbtype string
var dbdetails string
var unknownservicelistfilename string
var servicedatafilename string
var productionOnly bool
var saasonly bool
var onpremonly bool
var daysback int
var showincidents bool
var skipresolved bool
var maxcolwidth int
var servicefilters ArrayFlags

func (i *ArrayFlags) String() string {
	if *i == nil {
		return ""
	}
	toto := ""
	for _, s := range *i {
		toto += s + " "
	}
	return strings.TrimSpace(toto)
}

func (i *ArrayFlags) Set(v string) error {
	*i = append(*i, v)
	return nil
}

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
		if sr {
			if i.Status != "resolved" {
				block = append(block, []string{i.ID, i.Summary, i.CreatedAt, i.LastStatusChangeAt, i.Priority, i.Urgency, i.Status, i.AssignedName, i.ServiceName})
			}
		} else {
			block = append(block, []string{i.ID, i.Summary, i.CreatedAt, i.LastStatusChangeAt, i.Priority, i.Urgency, i.Status, i.AssignedName, i.ResolvedBy, i.ServiceName})
		}
	}
	return block
}

func printIncidents(sr bool, keys []string, mxc int, dbt *dbtalker.DBTalker) {
	var headers []string
	if sr {
		headers = []string{"ID", "Summary", "CreatedAt", "LastStatusChangeAt", "Priority", "Urgency", "Status", "AssignedName", "ServiceName"}
	} else {
		headers = []string{"ID", "Summary", "CreatedAt", "LastStatusChangeAt", "Priority", "Urgency", "Status", "AssignedName", "ResolvedBy", "ServiceName"}
	}
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
	flag.IntVar(&daysback, "D", 30, "How many days back to search. Cannot go further back than the data in the DB of course.")
	flag.IntVar(&maxcolwidth, "w", globals.MAXCOLWIDTH, "Max width of column in characters")
	flag.StringVar(&unknownservicelistfilename, "o", globals.DEFAULTUNKNOWNSERVICELIST, "File used to store list of unknown services seen")
	flag.StringVar(&servicedatafilename, "d", globals.DEFAULTSERVICEDATAFILENAME, "File with service data")
	flag.BoolVar(&productionOnly, "P", false, "If set, only record data for production services")
	flag.BoolVar(&saasonly, "S", false, "If set, only record data for services running on SaaS")
	flag.BoolVar(&onpremonly, "O", false, "If set, only record data for services running OnPrem")
	flag.BoolVar(&showincidents, "s", false, "If set, show the incident data")
	flag.BoolVar(&skipresolved, "R", false, "If set, skip incidents that are resolved")
	flag.Var(&servicefilters, "F", "Regex to filter services with. Can be specified multiple times.")
	flag.Parse()

	servicedata := readServiceData(servicedatafilename)
	if dbdetails == "" {
		log.Fatalln("dbdetails cannot be empty")
	}
	dbtalker := dbtalker.NewDBTalker(model.ConnectDatabase(dbtype, dbdetails))
	// get data
	scounts, snames, sortedkeys := pdanalyser.PDanalyse(servicefilters, productionOnly, saasonly, onpremonly, daysback, skipresolved, servicedata, dbtalker)
	// output
	prettifiedOutput(scounts, snames, sortedkeys)
	unknown := checkSeenServices(snames, servicedata)
	saveToFileAsJson(unknownservicelistfilename, unknown)
	// print out incidents
	if showincidents {
		printIncidents(skipresolved, sortedkeys, maxcolwidth, dbtalker)
	}
}
