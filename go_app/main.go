package store

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


func Country() (Location) {
        var united_states = Location{
            name: "United States",
            subtype: "country",
            published: true,
            visible: false,
            begin: 0,
            end: 0,
            place: Place{
                point: Point{
                    latitude: 39.8282,
                    longitude: -98.5795,
                },
                radius: Radius{ 
                    avg: 2156500.0, //.5(horizontal width of US) in meters
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
                    place: Place{
                        point: Point{
                            latitude: 36.1700,
                            longitude: -119.7462,
                        },
                        radius: Radius{
                            avg: 676000.0,
                        },
                    },
                    children: []Location{
                        Location{
                            name: "Davis",
                            subtype: "locale",
                            begin: 0,
                            published: true,
                            visible: true,
                            end: 0,
                            place: Place{
                                point: Point{
                                    latitude: 38.5539,
                                    longitude: -121.7381,
                                },
                                radius: Radius{ 
                                    avg: 25000.00,
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
                                            latitude: 38.5492101,
                                            longitude: -121.6961637,
                                        },
                                        radius: Radius{
                                            avg: 100.0,
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

    return united_states;
}



func (location Location) valid(reference Place) (bool) {
    //checks if 1) public 2) nearby and 3) timely
    if ( true == location.public() && true == location.nearby(location.place, reference) && true == location.timely() ) {
        return true;
    }
    return false;
}

//if is published and publicly visible
func (location Location) public() (bool) {
    return location.published;
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

func toRadians(degrees float64) (float64) {
    return (degrees * (math.Pi/180))
}

func distance(target Place, reference Place) (float64) {
    var deltaLat float64 = toRadians(target.point.latitude - reference.point.latitude);
    var deltaLon float64 = toRadians(target.point.longitude - reference.point.longitude);
    var a float64 = math.Sin(deltaLat/float64(2)) * math.Sin(deltaLat/float64(2)) + math.Cos(toRadians(reference.point.latitude)) * math.Cos(toRadians(target.point.latitude)) * math.Sin(deltaLon/float64(2)) * math.Sin(deltaLon/float64(2));
    //6371 = Earth's radius in km 
    return (float64(6371000) * float64(2) * math.Atan2( math.Sqrt(a), math.Sqrt(1-a) ) );
}
//if reference point is inside location boundry
func (location Location) nearby(target Place, reference Place) (bool) {
    var radius_meters float64 = target.radius.avg + reference.radius.avg;
    //log.Printf("Reference location: %f, %f", reference.point.latitude, reference.point.longitude);
    log.Printf("Target location: %f, %f", target.point.latitude, target.point.longitude);
    //log.Printf("Radius: %d meters", radius_meters );
    var distance float64 = distance(target, reference)
    log.Printf("Distance: %f meters", distance);
    return distance < float64(radius_meters);
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
                if children := child.children; len(children) > 0 {
                    deep = append( deep, collapse( children, reference, []Location{} )... );
                }
            } else {
                log.Print("Invalid node " + child.name);
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
        var original Location = Country()
        var reference = Place{
            point: Point{
                latitude: 38.5501232, //38.5445404,
                longitude: -121.695873, //-121.7398277,
            },
            radius: Radius{
                avg: 20.0,
            },
        }
        for _, v := range collapse( []Location{ original }, reference, make([]Location,0) ) {
            log.Print( "Place " + v.name );
        }
        response, _ := json.Marshal(original.children[0].children[0].name)
        fmt.Fprint(w, string(response))
    } )
}

