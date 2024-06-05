package pdanalyser

import (
	"sort"

	"github.com/gavinB-hpe/pdbyservice/dbtalker"
	"github.com/gavinB-hpe/pdbyservice/globals"
)

func PDanalyse(dbt *dbtalker.DBTalker) (map[string]int, map[string]string, []string) {
	urgency := globals.DEFAULTURGENCYVALUES
	status := globals.DEFAULTSTATUSVALUES
	scount := make(map[string]int, 0)
	snames := make(map[string]string, 0)
	for _, pdi := range dbt.GetIncidents(urgency, status) {
		scount[pdi.ServiceID] += 1
		snames[pdi.ServiceID] = pdi.ServiceName
	}
	keys := make([]string, 0)
	for k := range scount {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return scount[keys[i]] > scount[keys[j]] })
	return scount, snames, keys
}
