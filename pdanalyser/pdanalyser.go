package pdanalyser

import (
	"log"
	"regexp"
	"sort"
	"time"

	"github.com/gavinB-hpe/pdbyservice/dbtalker"
	"github.com/gavinB-hpe/pdbyservice/globals"
	"github.com/gavinB-hpe/pdbyservice/model"
)

/*
Helper function for deciding when to include the incident.
ProductionOnly = filters production only if set, *else all*
SaaSOnly and OnPrem cannot both be set.
SaaSOnly and OnPrem both unset = take all
*/
func yesok(po, so, onp bool, dt map[string]string) bool {
	if so && onp {
		log.Panic("Cannot choose both SaaS and OnPrem")
	}
	if !po && !so && !onp {
		return true // if no filters set, all matches
	}
	pocond := dt[globals.INPRODUCTION] == globals.TRUE
	socond := dt[globals.LOCATION] == globals.SAAS
	onpcond := dt[globals.LOCATION] == globals.ONPREM
	if po && pocond && !so && !onp {
		return true
	}
	if po && pocond && so && socond && !onp {
		return true
	}
	if po && pocond && onp && onpcond && !so {
		return true
	}
	if !po && !onp && so && socond {
		return true
	}
	if !po && !so && onp && onpcond {
		return true
	}
	return false

}

func dateOK(pdi model.PDInfoType, days int) bool {
	now := time.Now()
	then := now.AddDate(0, 0, -days)
	return pdi.CreatedAtT.After(then)
}

func statusOK(sr bool, i model.PDInfoType) bool {
	if sr {
		if i.Status == "resolved" {
			return false
		}
	}
	return true
}

func matchFilter(sf []string, pdi model.PDInfoType) bool {
	if len(sf) == 0 {
		return true
	}
	for _, s := range sf {
		if m, err := regexp.MatchString(s, pdi.ServiceID); err == nil && m {
			return true
		}
	}
	return false
}

// PDanalyse is a function that performs analysis on incidents based on various parameters.
// It takes in the following parameters:
// - sf []string: a slice of strings representing filter criteria
// - po bool: a boolean indicating whether to include incidents with Priority Outage status
// - so bool: a boolean indicating whether to include incidents with Service Outage status
// - onp bool: a boolean indicating whether to include incidents with On-Call No Pager status
// - days int: an integer representing the number of days to consider for analysis
// - skipr bool: a boolean indicating whether to skip incidents with Resolved status
// - sd *map[string]map[string]string: a pointer to a map containing service details
// - dbt *dbtalker.DBTalker: a pointer to a DBTalker object for interacting with the database
//
// It returns the following:
// - map[string]int: a map with service IDs as keys and the count of incidents as values
// - map[string]string: a map with service IDs as keys and the corresponding service names as values
// - []string: a slice of strings representing the service IDs in descending order of incident count
func PDanalyse(
	sf []string,
	po bool,
	so bool,
	onp bool,
	days int,
	skipr bool,
	sd *map[string]map[string]string,
	dbt *dbtalker.DBTalker) (map[string]int, map[string]string, []string) {
	urgency := globals.DEFAULTURGENCYVALUES
	status := globals.DEFAULTSTATUSVALUES
	scount := make(map[string]int, 0)
	snames := make(map[string]string, 0)
	for _, pdi := range dbt.GetIncidents(urgency, status) {
		if dateOK(pdi, days) && statusOK(skipr, pdi) && matchFilter(sf, pdi) {
			dt := (*sd)[pdi.ServiceID]
			if yesok(po, so, onp, dt) {
				scount[pdi.ServiceID] += 1
				snames[pdi.ServiceID] = pdi.ServiceName
			}
		}
	}
	keys := make([]string, 0)
	for k := range scount {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return scount[keys[i]] > scount[keys[j]] })
	return scount, snames, keys
}
