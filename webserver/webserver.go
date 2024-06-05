package webserver

import (
	"log"
	"net/http"
	"os"

	"github.com/gavinB-hpe/pdbyservice/globals"
	// "github.com/gavinB-hpe/pdbyservice/graphdrawer"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()
var filename = globals.OUTPUTFILENAME

func GetIndex(c *gin.Context) {
	router.LoadHTMLFiles("templates/index.html")
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

/*

func GetGraphService(c *gin.Context) {
	log.Println("GetGraphService")
	thing := graphdrawer.KnownChoices[0]
	router.LoadHTMLFiles(fmt.Sprintf("./%s/%s-%s", globals.GENERATED_ASSETS_DIRECTORY, thing, filename))
	log.Println("loaded ", thing, filename)
	c.HTML(http.StatusOK, fmt.Sprintf("%s-%s", thing, filename), gin.H{})
}

func GetGraphPolicy(c *gin.Context) {
	log.Println("GetGraphService")
	thing := graphdrawer.KnownChoices[2]
	router.LoadHTMLFiles(fmt.Sprintf("./%s/%s-%s", globals.GENERATED_ASSETS_DIRECTORY, thing, filename))
	log.Println("loaded ", thing, filename)
	c.HTML(http.StatusOK, fmt.Sprintf("%s-%s", thing, filename), gin.H{})
}

func GetGraphStatus(c *gin.Context) {
	log.Println("GetGraphService")
	thing := graphdrawer.KnownChoices[3]
	router.LoadHTMLFiles(fmt.Sprintf("./%s/%s-%s", globals.GENERATED_ASSETS_DIRECTORY, thing, filename))
	log.Println("loaded ", thing, filename)
	c.HTML(http.StatusOK, fmt.Sprintf("%s-%s", thing, filename), gin.H{})
}
*/

func ServerIt(fn string, address string, tlsmode bool) {
	filename = fn
	// LoadHTMLGlob is overwritten by later LoadHTMLFiles
	// router.LoadHTMLGlob("templates/*")
	// enable serving of fixed assets
	router.Static("/assets", "./assets")
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")
	// routes we want to react to
	router.GET("/", GetIndex)
	router.GET("/index.html", GetIndex)
	// router.GET("/graphservice", GetGraphService)
	// router.GET("/graphpolicy", GetGraphPolicy)
	// router.GET("/graphstatus", GetGraphStatus)
	certpath := os.Getenv(globals.TLSCERTPATHKEY)
	keypath := os.Getenv(globals.TLSKEYPATHKEY)

	var re error
	if tlsmode && certpath != "" && keypath != "" {
		log.Println("Running in TLS mode")
		re = router.RunTLS(address, certpath, keypath)
	} else {
		log.Println("Non-TLS mode")
		re = router.Run(address)
	}
	if re != nil {
		log.Fatalf("process failed with error: %v\n", re)
	}
}
