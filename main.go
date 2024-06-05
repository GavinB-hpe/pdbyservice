package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gavinB-hpe/pdbyservice/dbtalker"
	"github.com/gavinB-hpe/pdbyservice/globals"

	"github.com/gavinB-hpe/pdbyservice/model"
	"github.com/gavinB-hpe/pdbyservice/pdanalyser"
)

// var for flags
var dbtype string
var dbdetails string
var bucketsize int

func prettifiedOutput(sc map[string]int, sn map[string]string, keys []string) {
	for _, k := range keys {
		fmt.Println(sn[k], " ", sc[k])
	}
}

func main() {
	flag.StringVar(&dbtype, "t", globals.DEFAULTDBTYPE, "Type of DB used e.g. sqlite3")
	flag.StringVar(&dbdetails, "db", globals.DEFAULTDBDETAILS, "Filename for sqlite3 or URI of DB")
	flag.IntVar(&bucketsize, "b", globals.DEFAULTBUCKETSIZE, "How many days to bucket together in the graph. ")
	flag.Parse()
	if bucketsize <= 0 {
		log.Fatalln("Invalid bucketsize value. Must be > 0")
	}
	if dbdetails == "" {
		log.Fatalln("dbdetails cannot be empty")
	}
	dbtalker := dbtalker.NewDBTalker(model.ConnectDatabase(dbtype, dbdetails))
	// gd := graphdrawer.NewGraphDrawer(dbtalker, bucketsize)
	// keys := readKeys(searchkeyfilename)
	// go timer(timerchanservice)
	// go timer(timerchanpolicy)
	// go timer(timerchanstatus)
	// go drawGraphService(&gd, keys, timerchanservice)
	// go drawGraphPolicy(&gd, keys, timerchanpolicy)
	// go drawGraphStatus(&gd, keys, timerchanstatus)
	scounts, snames, sortedkeys := pdanalyser.PDanalyse(dbtalker)
	prettifiedOutput(scounts, snames, sortedkeys)
	// webserver.ServerIt(globals.OUTPUTFILENAME, fmt.Sprintf("%s:%d", address, port), tlsmode)
}
