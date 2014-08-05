package pushgo

import (
	"fmt"
	"net/http"
    "encoding/json"
    "log"
    "time"
    "math"
)

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
    subtype string
    place Place
    published bool
    visible bool
    children []Location
    meta interface{}
    begin int64
    end int64
}

func (location Location) valid(reference Place) (bool) {
    //check if published or hidden
    log.Printf("Reference location: %f, %f", reference.point.latitude, reference.point.longitude);
    if ( true == location.public() && true == location.nearby() && true == location.timely() ) {
        return true;
    }
    return false;
}

//if is published and publicly visible
func (location Location) public() (bool) {
    return location.published && location.visible;
}

//if no begin/end set or ( begin/end exists && valid )
func (location Location) timely() (bool) {
    if ( 0 == location.begin && 0 == location.end ) {
        return true;
    } 
    var timestamp = time.Now().Local().Unix();
    if ( 0 != location.end ) {
        if ( location.end < timestamp ) {
            return false;
        }
    } 
    if ( 0 != location.begin ) {
        if ( location.begin > timestamp ) {
            return false;
        }
    }
    return true;
}

//if reference point is inside location boundry
func (location Location) nearby() (bool) {
    return true;
}


//TODO: pass valid() func as generic arg
func collapse( nodes []Location, reference Place, found []Location ) ( []Location )  {
    var deep []Location = []Location{}
    if len(nodes) > 0 {
        for _, child := range nodes {
            if true == child.valid(reference) {
                found = append( found, []Location{
                 child,   
                }... )
            }
            if children := child.children; len(children) > 0 {
                deep = append( deep, collapse( children, reference, []Location{} )... );
            }
        }
    }
    return append( found, deep... );
}

func init() {
    log.Printf("Started %d", time.Now().Local().Unix())

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
        var united_states = Location{
            name: "United States",
            subtype: "country",
            published: true,
            visible: false,
            begin: 0,
            end: 0,
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
                    name: "CA",
                    subtype: "state",
                    published: true,
                    visible: true,
                    begin: 0,
                    end: 0,
                    children: []Location{
                        Location{
                            name: "Davis",
                            subtype: "locale",
                            begin: 0,
                            end: 0,
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
                                    subtype: "location",
                                    published: true,
                                    visible: true,
                                    begin: 0,
                                    end: 0,
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
                                    meta: nil,
                                },
                                Location{
                                    name: "Guadajara Grill",
                                    subtype: "location",
                                    published: true,
                                    visible: true,
                                    begin: 0,
                                    end: 0,
                                    place: Place{
                                        point: Point{
                                            latitude: 38.5597532,
                                            longitude: -121.7568926,
                                        },
                                        radius: Radius{
                                            avg: 50.0,
                                        },
                                    },
                                    children: []Location{
                                        Location{
                                            name: "Guadalajara Porch",
                                            subtype: "location",
                                            published: true,
                                            visible: true,
                                            begin: 1407211227,
                                            end: 1407211740,
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
                                            meta: nil,
                                        },
                                        Location{
                                            name: "Guadajara Interior",
                                            subtype: "location",
                                            published: false,
                                            visible: true,
                                            begin: 0,
                                            end: 1407211227,
                                            place: Place{
                                                point: Point{
                                                    latitude: 38.5597532,
                                                    longitude: -121.7568926,
                                                },
                                                radius: Radius{
                                                    avg: 50.0,
                                                },
                                            },
                                            children: []Location{},
                                            meta: nil,
                                        },    
                                    },
                                    meta: nil,
                                },
                            },
                            meta: nil,
                        },
                    },
                    meta: nil,
                },
            },
            meta: nil,
        }

        var reference = Place{
            point: Point{
                latitude: 38.5445404,
                longitude: -121.7398277,
            },
            radius: Radius{
                avg: 5.0,
            },
        }
        for _, v := range collapse( []Location{ united_states }, reference, make([]Location,0) ) {
            log.Print( "Place " + v.name );
        }

        response, _ := json.Marshal(united_states.children[0].children[0].name)
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
