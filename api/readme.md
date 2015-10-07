> api.chevent.xyz

### Endpoints
**get /version**

Latest version of chevent-web over latest git commit hash

**get /events**
* limit int
* page int

**post /events**
* name string
* date *UTC* string
* free bool
* image string
* url string
* description string
* quota int
* speakers.0.name
* speakers.0.subject ...
