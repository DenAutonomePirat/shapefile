package main

import (
	"fmt"
	"github.com/jonas-p/go-shp"
	"log"
)

func main() {
	// open a shapefile for reading
	shape, err := shp.Open("limfjorden/Limfjorden.shp")
	if err != nil {
		log.Fatal(err)
	}
	defer shape.Close()

	// fields from the attribute table (DBF)
	fields := shape.Fields()

	// loop through all features in the shapefile
	for shape.Next() {

		n, p := shape.Shape()
		switch p := p.(type) {
		case *shp.Polygon:
			fmt.Printf("hej: %v\n", p.NumParts)
			var i int32
			for i = 0; i < p.NumParts; i++ {
				fmt.Printf("Part: %v starts at %v \n", i, p.Parts[i])
			}

		default:
			fmt.Printf("unexpected type %T\n", p)

		}

		// print feature
		// print attributes
		for k, f := range fields {
			val := shape.ReadAttribute(n, k)
			fmt.Printf("\t%v: %v\n", f, val)
		}
		fmt.Println()
	}
}
