package dbtalker

import (
	"strings"

	"github.com/gavinB-hpe/pdbyservice/model"
)

func (dbt *DBTalker) GetIncidentsByLoop(urgency string, status string) []model.PDInfoType {
	var items []model.PDInfoType
	// TODO - search in SQL to get matching items
	dbt.DB.Find(&items)
	tmp := make([]model.PDInfoType, 0)
	urgencies := strings.Split(urgency, ",")
	for _, u := range urgencies {
		statuses := strings.Split(status, ",")
		for _, s := range statuses {
			for _, it := range items {
				if it.Urgency == u && it.Status == s {
					tmp = append(tmp, it)
				}
			}
		}

	}
	return tmp
}

// This version is 3ms faster - 20ms vs 23ms for the loop.
func (dbt *DBTalker) GetIncidents(urgency string, status string) []model.PDInfoType {
	var items []model.PDInfoType
	dbt.DB.Where("status IN ? AND urgency IN ?", strings.Split(status, ","), strings.Split(urgency, ",")).Find(&items)
	return items
}
