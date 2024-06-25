package model

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gavinB-hpe/pdbyservice/globals"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
)

const (
	DBTYPE = "sqlite3" // for now
)

type PDInfoType struct {
	ID                  string    `json:"id" binding:"required"`
	Summary             string    `json:"summary" binding:"required"`
	Description         string    `json:"description" binding:"required"`
	CreatedAt           string    `json:"createdat" binding:"required"`
	CreatedAtT          time.Time `json:"createdatt"`
	LastStatusChangeAt  string    `json:"laststatuschangeat" binding:"required"`
	LastStatusChangeAtT time.Time `json:"laststatuschangeatt" binding:"required"`
	HTMLURL             string    `json:"htmlurl" binding:"required"`
	Priority            string    `json:"priority" binding:"required"`
	Urgency             string    `json:"urgency" binding:"required"`
	Status              string    `json:"status" binding:"required"`
	Policy              string    `json:"policy" binding:"required"`
	ServiceName         string    `json:"servicename" binding:"required"`
	ServiceID           string    `json:"serviceid" binding:"required"`
}

func (pdi *PDInfoType) ToString() string {
	return fmt.Sprintf("ID=%s, Summary=%s, CreatedAt=%s, Priority=%s, Urgency=%s, Status=%s, ServiceName=%s", pdi.ID, pdi.Summary, pdi.CreatedAt, pdi.Priority, pdi.Urgency, pdi.Status, pdi.ServiceName)
}

var DB *gorm.DB
var ExportMetrics = false

func setupSQLiteDatabase(name string, values interface{}) (*gorm.DB, error) {
	filename := "./" + name + ".db"
	filename2 := strings.Replace(filename, ".db.db", ".db", 1)
	if filename != filename2 {
		log.Println("WARNING - changed filename from ", filename, " to ", filename2)
		filename = filename2
	}
	log.Println("Connecting to SQLite3 database at ", filename)
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		log.Println("Failed to connect to database "+name+" using filename "+filename, err)
		return nil, err
	}
	if ExportMetrics {
		if err := db.Use(otelgorm.NewPlugin()); err != nil {
			log.Fatal(err)
		}
	}
	db.AutoMigrate(values)
	return db, nil
}

// get the databases created
func ConnectDatabase(tp string, name string) *gorm.DB {
	if tp != globals.DEFAULTDBTYPE {
		log.Fatalln("Unsupport DB type")
	}
	db, err := setupSQLiteDatabase(name, &PDInfoType{})
	if err != nil {
		panic(err)
	}
	DB = db
	return DB
}
