package pushgo

import (
	"fmt"
	"net/http"
    "encoding/json"
    "log"
    "time"
    "math"
)

type Country struct {
    name string
    children []State
}
type State struct {
    name string
    children []Locale
}
type Radius struct {
    avg float64
}
type Point struct {
    latitude float64
    longitude float64
}
type Place struct {
    point Point
    radius Radius
}
type Location struct {
    name string
    place Place
    children []Location
}
type Locale struct {
    name string
    place Place
    children []Location
}

//  
func collapse( nodes []Location ) ( []Location )  {
    var found = make([]Location,0)
    if len(nodes) > 0 {
        for _, child := range nodes {
            if children := child.children; len(children) > 0 {
                found = append( found, collapse( children )... );
            }
        }
    }
    return found;
}

func init() {
    log.Printf("Started %d", time.Now().Local().Unix())

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
        var country = Country{
            name: "United States",
            children: []State{
                State{
                    name: "CA",
                    children: []Locale{
                        Locale{
                            name: "Davis",
                            place: Place{
                                point: Point{
                                    latitude: 38.5539,
                                    longitude: -121.7381,
                                },
                                radius: Radius{ 
                                    avg: 10000.0,
                                },
                            },
                            children: []Location{
                                Location{
                                    name: "Bistro 33",
                                    place: Place{
                                        point: Point{
                                            latitude: 38.5444038,
                                            longitude: -121.7397349,
                                        },
                                        radius: Radius{
                                            avg: 50.0,
                                        },
                                    },
                                    children: []Location{},
                                },
                            },
                        },
                    },
                },
            },
        }

        //var lat float64 = 38.5445404
        //var lon float64 = -121.7398277


        response, _ := json.Marshal(country.children[0].children[0].name)
	    fmt.Fprint(w, string(response))
    } )
}

func toRadians(degrees float64) (float64) {
    return (degrees * (math.Pi/180))
}

func nearby(latTarget float64, lonTarget float64, latRef float64, lonRef float64, radius int) (bool) {
    return (distance(latTarget, lonTarget, latRef, lonRef) < float64(radius));
}

func distance(lat1 float64, lon1 float64, lat2 float64, lon2 float64) (float64) {
    var deltaLat float64 = toRadians(lat2-lat1)
    var deltaLon float64 = toRadians(lon2-lon1)
    var a float64 = math.Sin(deltaLat/float64(2)) * math.Sin(deltaLat/float64(2)) + math.Cos(toRadians(lat1)) * math.Cos(toRadians(lat2)) * math.Sin(deltaLon/float64(2)) * math.Sin(deltaLon/float64(2));
    //6371 = Earth's radius in km 
    return (float64(6371) * float64(2) * math.Atan2( math.Sqrt(a), math.Sqrt(1-a) ) );
}
