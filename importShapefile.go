package main

import (
	"bufio"
	"fmt"
	"github.com/jonas-p/go-shp"

	"log"
	"os"
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

			fmt.Printf("No of vertexes in Polygon: %v\n", p.NumPoints)
			fmt.Printf("No of parts in Polygon: %v\n", p.NumParts)

			//Create node file
			node, err := os.Create("triangle/temp.node")
			if err != nil {
				panic(err)
			}
			defer node.Close()

			nodeWriter := bufio.NewWriter(node)

			poly, err := os.Create("triangle/temp.poly")
			if err != nil {
				panic(err)
			}
			defer poly.Close()

			polyWriter := bufio.NewWriter(poly)

			line := fmt.Sprintf("#Gennerated by importShapefile.go\n%v 2 0 0\n", p.NumPoints)
			_, err = nodeWriter.WriteString(line)
			if err != nil {
				panic(err)
			}

			line = fmt.Sprintf("#Gennerated by importShapefile.go\n0 2 0 0\n")
			_, err = polyWriter.WriteString(line)
			if err != nil {
				panic(err)
			}

			line = fmt.Sprintf("%v 0\n", p.NumParts)
			_, err = polyWriter.WriteString(line)
			if err != nil {
				panic(err)
			}

			for i := int32(0); i < p.NumPoints; i++ {
				//fmt.Printf("%v %v %v\n", i-p.Parts[i], p.Points[i].X, p.Points[i].Y)
				line = fmt.Sprintf("%v %v %v\n", i, p.Points[i].X, p.Points[i].Y)
				_, err := nodeWriter.WriteString(line)
				if err != nil {
					panic(err)
				}

			}
			segmentNo := 0
			for i := int32(0); i < p.NumParts; i++ {
				partStart := p.Parts[i]
				partsEnd := partStart
				if i >= p.NumParts-1 {
					partsEnd = p.NumPoints
				} else {
					partsEnd = p.Parts[i+1]
				}
				for j := partStart; j < partsEnd; j++ {
					end := j
					if j == partsEnd-1 {
						j = partStart
					}
					line = fmt.Sprintf("%v %v %v\n", segmentNo, j, end)
					_, err = polyWriter.WriteString(line)
					if err != nil {
						panic(err)
					}
					segmentNo++
				}
			}

			nodeWriter.Flush()
			polyWriter.Flush()

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
