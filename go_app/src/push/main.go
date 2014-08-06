package push

import (
	"fmt"
	"net/http"
    "encoding/json"
    "log"
    "time"
    "math"
    "sort"
    "strings"
)

type Radius struct {
    avg float64
    min float64
    max float64
}
type Point struct {
    latitude float64
    longitude float64
}
type Place struct {
    point Point
    radius Radius
    distance float64
}
type Location struct {
    id string
    subtype string
    place Place
    published bool
    visible bool
    displayed bool
    children []Location
    meta interface{}
    begin int64
    end int64
}


func Data() (Location) {
        var united_states = Location{
            id: "United States",
            subtype: "country",
            published: true,
            visible: true,
            displayed: false,
            begin: 0,
            end: 0,
            place: Place{
                point: Point{
                    latitude: 39.8282,
                    longitude: -98.5795,
                },
                radius: Radius{ 
                    avg: 2156500.0, //.5(horizontal width of US) in meters
                    max: 2156500.0,
                    min: 0.0,
                },
                distance: 0.0,
            },
            children: []Location{
                Location{
                    id: "CA",
                    subtype: "state",
                    published: true,
                    visible: true,
                    displayed: true,
                    begin: 0,
                    end: 0,
                    place: Place{
                        point: Point{
                            latitude: 36.1700,
                            longitude: -119.7462,
                        },
                        radius: Radius{
                            avg: 676000.0,
                            max: 676000.0,
                            min: 0.0,
                        },
                        distance: 0.0,
                    },
                    children: []Location{
                        Location{
                            id: "Davis",
                            subtype: "locale",
                            begin: 0,
                            published: true,
                            visible: true,
                            displayed: true,
                            end: 0,
                            place: Place{
                                point: Point{
                                    latitude: 38.5539,
                                    longitude: -121.7381,
                                },
                                radius: Radius{ 
                                    avg: 25000.00,
                                    max: 25000.00,
                                    min: 0,
                                },
                                distance: 0.0,
                            },
                            children: []Location{
                                Location{
                                    id: "Bistro 33",
                                    subtype: "location",
                                    published: true,
                                    visible: true,
                                    displayed: true,
                                    begin: 0,
                                    end: 0,
                                    place: Place{
                                        point: Point{
                                            latitude: 38.5444038,
                                            longitude: -121.7397349,
                                        },
                                        radius: Radius{
                                            avg: 50.0,
                                            max: 50.0,
                                            min: 0,
                                        },
                                        distance: 0.0,
                                    },
                                    children: []Location{},
                                    meta: nil,
                                },
                                Location{
                                    id: "Guadajara Grill",
                                    subtype: "location",
                                    published: true,
                                    visible: true,
                                    displayed: true,
                                    begin: 0,
                                    end: 0,
                                    place: Place{
                                        point: Point{
                                            latitude: 38.5492101,
                                            longitude: -121.6961637,
                                        },
                                        radius: Radius{
                                            avg: 10.0,
                                            max: 20.0,
                                            min: 0.0,
                                        },
                                        distance: 0.0,
                                    },
                                    children: []Location{
                                        Location{
                                            id: "Guadalajara Porch",
                                            subtype: "location",
                                            published: true,
                                            visible: true,
                                            displayed: true,
                                            begin: 1407211227,
                                            end: 1407211740,
                                            place: Place{
                                                point: Point{
                                                    latitude: 38.5492101,
                                                    longitude: -121.6961637,
                                                },
                                                radius: Radius{
                                                    avg: 40.0,
                                                    max: 50.0,
                                                    min: 20.0,
                                                },
                                                distance: 0.0,
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
    //checks if 1) public and 2) timely
    if ( true == location.public() && true == location.timely() ) {
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

//TODO: pass valid() func as generic arg
func collapse( nodes []Location, reference Place, found []Location ) ( []Location )  {
    var deep []Location = []Location{}
    var keys []float64


    if len(nodes) > 0 {
        for _, child := range nodes {
            if true == child.valid(reference) {
                var target Place = child.place
                var max_radius float64 = target.radius.max + reference.radius.max;
                var meters_away float64 = distance(target, reference)
                if meters_away < float64(max_radius) {
                    child.place.distance = meters_away
                    found = append( found, []Location{
                     child,   
                    }... )
                    if children := child.children; len(children) > 0 {
                        deep = append( deep, collapse( children, reference, []Location{} )... );
                    }
                }
            }
        }

    }
    found = append(found, deep...);
    for _, child := range found {
        keys = append(keys, child.place.distance);
    }
    sort.Float64s(keys)
    var tmp []Location;
    if "ASC" == SORT { //TODO: try a switch
        n := len(keys) 
        for i := n; i > 0; i-- { 
            var value float64 = keys[i - 1]
            for _, child := range found {
                if value == child.place.distance {
                    log.Print("found VALUE", value, child.place.distance)
                    tmp = append(tmp, child);
                } else {
                    log.Print("mismatch", value, child.place.distance)
                }
            }
        }
    } else {
        for i, _ := range keys { 
            var value float64 = keys[i]
            for _, child := range found {
                if value == child.place.distance {
                    log.Print("found VALUE", value, child.place.distance)
                    tmp = append(tmp, child);
                } else {
                    log.Print("mismatch", value, child.place.distance)
                }
            }
        }
    }
    return tmp;
}

//TODO: Url parsing http://localhost:8080/presence#id%3D123%26lat%3D38.549210%26lng%3D-121.696164%26time%3D1407302303
const MAX_RADIUS float64 = 500.0
const AVG_RADIUS float64 = 20.0
const SORT string = "DESC" //ASC else DESC
const LIMIT int = 0

func init() {
    log.Printf("Started %d", time.Now().Local().Unix())
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" && r.URL.Path != "/presence" {
            http.NotFound(w, r)
            return
        }
        var matching []map[string]interface{}
        for _, v := range collapse( []Location{ Data() }, Place{
            point: Point{
                latitude: 38.5501232, //38.5445404,
                longitude: -121.695873, //-121.7398277,
            },
            radius: Radius{
                min: 0.0,
                avg: AVG_RADIUS,
                max: MAX_RADIUS,
            },
        }, make([]Location,0) ) {
            if ( 0 != LIMIT && len(matching) >= LIMIT ) {
                break;
            }
            if true == v.visible && v.place.distance < MAX_RADIUS {
                matching = append(matching, []map[string]interface{}{
                    map[string]interface{}{ //TODO: Pass inbound anonymous id
                    "uri": "/" + v.subtype + "/" + strings.Replace( v.id, " ", "_", -1 ) + "/#lat%3D" + fmt.Sprintf("%f", v.place.point.latitude) + "%26lng%3D" + fmt.Sprintf("%f", v.place.point.longitude) + "%26time%3D" + fmt.Sprintf("%d", time.Now().Local().Unix() ),
                    "begin": v.begin,
                    "end": v.end,
                    "displayed": v.displayed,
                    "latitude": v.place.point.latitude,
                    "longitude": v.place.point.longitude,
                    "radius": v.place.radius.avg,
                    "distance": v.place.distance,
                    },
                }...);
            }
                    
        }
        response, _ := json.Marshal(matching)
        log.Printf("RESPONSE: %s", string(response))
        fmt.Fprint(w, string(response))


    } )
}

