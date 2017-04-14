package main

import (
	"bufio"
	"fmt"
	"github.com/jonas-p/go-shp"
	"github.com/kellydunn/golang-geo"
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

			//Create .node file
			node, err := os.Create("triangle/temp.node")
			if err != nil {
				panic(err)
			}
			defer node.Close()

			nodeWriter := bufio.NewWriter(node)

			line := fmt.Sprintf("#Gennerated by importShapefile.go\n%v 2 0 0\n", p.NumPoints)
			_, err = nodeWriter.WriteString(line)
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
			nodeWriter.Flush()

			//Create .poly file
			poly, err := os.Create("triangle/temp.poly")
			if err != nil {
				panic(err)
			}
			defer poly.Close()

			polyWriter := bufio.NewWriter(poly)

			line = fmt.Sprintf("#Gennerated by importShapefile.go\n0 2 0 0\n")
			_, err = polyWriter.WriteString(line)
			if err != nil {
				panic(err)
			}

			line = fmt.Sprintf("%v 0\n", p.NumPoints-p.NumParts)
			_, err = polyWriter.WriteString(line)
			if err != nil {
				panic(err)
			}
			segmentno := 0

			for i := int32(0); i < p.NumParts; i++ {
				endOfPart := int32(0)
				if i >= p.NumParts-1 {
					endOfPart = p.NumPoints
				} else {
					endOfPart = p.Parts[i+1]
				}
				for j := p.Parts[i]; j < endOfPart-2; j++ {
					line = fmt.Sprintln(segmentno, j, j+1)
					_, err = polyWriter.WriteString(line)
					if err != nil {
						panic(err)
					}
					segmentno++

				}
				line = fmt.Sprintln(segmentno, endOfPart-2, p.Parts[i])
				_, err = polyWriter.WriteString(line)
				if err != nil {
					panic(err)
				}
				segmentno++
			}

			line = fmt.Sprintf("#Holes\n75\n")
			_, err = polyWriter.WriteString(line)
			if err != nil {
				panic(err)
			}
			//Find a point inside the hole
			for i := int32(1); i < p.NumParts; i++ {
				var found bool = false
				for j := int32(0); found == false; j++ {
					temp := p.Points[p.Parts[i]+j]
					pA := geo.NewPoint(temp.X, temp.Y)
					temp = p.Points[p.Parts[i]+j+1]
					pB := geo.NewPoint(temp.X, temp.Y)
					temp = p.Points[p.Parts[i]+j+1]
					pC := geo.NewPoint(temp.X, temp.Y)
					diff := pB.BearingTo(pA) - pB.BearingTo(pC)
					if diff > 40 && diff < 120 {
						found = true
						pH := pB.PointAtDistanceAndBearing(0.001, diff)

						line = fmt.Sprintln(i, pH.Lat(), pH.Lng())
						_, err = polyWriter.WriteString(line)
						if err != nil {
							panic(err)
						}
					}
				}

			}

			//Add holes

			polyWriter.Flush()

			fmt.Println(p.Points[p.Parts[1]].Y)

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
