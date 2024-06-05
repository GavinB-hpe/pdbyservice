package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gavinB-hpe/pdbyservice/dbtalker"
	"github.com/gavinB-hpe/pdbyservice/globals"

	// "github.com/gavinB-hpe/pdbyservice/graphdrawer"
	"github.com/gavinB-hpe/pdbyservice/model"
	"github.com/gavinB-hpe/pdbyservice/pdanalyser"
)

// var for flags
var dbtype string
var dbdetails string
var refreshdelay int
var searchkeyfilename string
var bucketsize int
var tlsmode bool
var address string
var port int

// Other var
var timerchanservice = make(chan bool, 1)
var timerchanpolicy = make(chan bool, 1)
var timerchanstatus = make(chan bool, 1)

func timer(ch chan bool) {
	ch <- true
	for {
		time.Sleep(time.Duration(refreshdelay) * time.Second)
		ch <- true
	}
}

func readKeys(skfname string) []string {
	kf, err := os.Open(skfname)
	if err != nil {
		log.Fatalf("Could not open %s for reading\n", skfname)
	}
	defer kf.Close()
	fileScanner := bufio.NewScanner(kf)
	fileScanner.Split(bufio.ScanLines)
	var searchkeys []string
	for fileScanner.Scan() {
		searchkeys = append(searchkeys, fileScanner.Text())
	}
	return searchkeys
}

/*

// service, title, policy, status, type
func drawGraph(gd *graphdrawer.GraphDrawer, keys []string, timerchan chan bool) {
	drawGraphService(gd, keys, timerchan)
}

func drawGraphService(gd *graphdrawer.GraphDrawer, keys []string, timerchan chan bool) {
	for {
		<-timerchan
		gd.Draw(graphdrawer.KnownChoices[0], keys)
	}
}

func drawGraphPolicy(gd *graphdrawer.GraphDrawer, keys []string, timerchan chan bool) {
	for {
		<-timerchan
		gd.Draw(graphdrawer.KnownChoices[2], keys)
	}
}

func drawGraphStatus(gd *graphdrawer.GraphDrawer, keys []string, timerchan chan bool) {
	for {
		<-timerchan
		gd.Draw(graphdrawer.KnownChoices[3], keys)
	}
}
*/

func prettifiedOutput(sc map[string]int, sn map[string]string, keys []string) {
	for _, k := range keys {
		fmt.Println(sn[k], " ", sc[k])
	}
}

func main() {
	flag.StringVar(&dbtype, "t", globals.DEFAULTDBTYPE, "Type of DB used e.g. sqlite3")
	flag.StringVar(&dbdetails, "db", globals.DEFAULTDBDETAILS, "Filename for sqlite3 or URI of DB")
	flag.IntVar(&refreshdelay, "r", globals.DEFAULTREFRESHDELAY, "How long between graph regeneration")
	flag.IntVar(&bucketsize, "b", globals.DEFAULTBUCKETSIZE, "How many days to bucket together in the graph. ")
	flag.StringVar(&searchkeyfilename, "k", globals.DEFAULTSEARCHKEYFILENAME, "File with list of keys to search for in the description")
	flag.StringVar(&address, "a", globals.DEFAULTADDRESS, "Address to listen on.")
	flag.IntVar(&port, "p", globals.DEFAULTPORT, "Port to listen on.")
	flag.BoolVar(&tlsmode, "tls", false, "Enable TLS mode. Requires setting "+globals.TLSCERTPATHKEY+" and "+globals.TLSKEYPATHKEY+" in the environment")
	flag.Parse()
	if bucketsize <= 0 {
		log.Fatalln("Invalid bucketsize value. Must be > 0")
	}
	if refreshdelay <= 1 {
		log.Fatalln("Invalid refresh delay value. Must be > 0")
	}
	if dbdetails == "" {
		log.Fatalln("dbdetails cannot be empty")
	}
	if port < 1000 {
		log.Println("Remember to sudo setcap CAP_NET_BIND_SERVICE=+eip pdbyservice")
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
