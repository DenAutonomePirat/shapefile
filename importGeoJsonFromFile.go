package main

import (
	//"fmt"
	"fmt"
	"github.com/paulmach/go.geojson"
	"io/ioutil"
	"log"
	"reflect"
)

func main() {
	raw, err := ioutil.ReadFile("Limfjorden.geojson") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	fc1, err := geojson.UnmarshalFeatureCollection(raw)

	f2 := fc1.Features[0].Geometry

	fmt.Println(reflect.TypeOf(f2))

	fmt.Println(f2)
}
