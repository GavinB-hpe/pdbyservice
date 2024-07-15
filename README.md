# pdbyservice

```
Usage of ./pdbyservice:
  -D int
    	How many days back to search. Cannot go further back than the data in the DB of course. (default 30)
  -F value
    	Regex to filter services with. Can be specified multiple times.
  -O	If set, only record data for services running OnPrem
  -P	If set, only record data for production services
  -R	If set, skip incidents that are resolved
  -S	If set, only record data for services running on SaaS
  -b int
    	How many days to bucket together in the graph.  (default 7)
  -d string
    	File with service data (default "services-data.json")
  -db string
    	Filename for sqlite3 or URI of DB (default "pdinfo")
  -o string
    	File used to store list of unknown services seen (default "services-seen-but-unknown.json")
  -s	If set, show the incident data
  -t string
    	Type of DB used e.g. sqlite3 (default "sqlite")
  -w int
    	Max width of column in characters (default 40)
```
