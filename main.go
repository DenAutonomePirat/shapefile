package main

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"github.com/denautonomepirat/golang-geo"
	"github.com/denautonomepirat/shapefile/triangle"
	gj "github.com/kpawlik/geojson"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math"
	"time"
)

type MapStore struct {
	session     *mgo.Session
	points      *mgo.Collection
	polygons    *mgo.Collection
	pointsCache map[int]*Point
}

func NewMapStore() *MapStore {

	session, err := mgo.Dial("localhost")
	//defer session.Close() //hmm why not?
	if err != nil {
		panic(err)
	}

	points := session.DB("map").C("points")
	polygons := session.DB("map").C("polygons")

	session.DB("map")
	return &MapStore{
		session:     session,
		points:      points,
		polygons:    polygons,
		pointsCache: make(map[int]*Point),
	}
}

type Point struct {
	ID         bson.ObjectId   `bson:"_id,omitempty" json:"_id"`
	Type       string          `bson:"type" json:"type"`
	Properties PointProperties `bson:"properties" json:"properties"`
	Geometry   Geometry        `bson:"geometry" json:"geometry"`
	mapStore   *MapStore
	timeStamp  time.Time
}

type PointProperties struct {
	POINTA    int        `bson:"POINTA,omitempty" json:"POINTA"`
	POINTB    int        `bson:"POINTB,omitempty" json:"POINTB"`
	POINTC    int        `bson:"POINTC,omitempty" json:"POINTC"`
	Node      int        `bson:"node,omitempty" json:"node"`
	Naighbors []Naighbor `bson:"neighbornodes,omitempty" json:"neighbornodes"`
}
type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}
type Naighbor struct {
	Node     int     `bson:"node,omitempty" json:"node"`
	Distance float64 `bson:"distance,omitempty" json:"distance"`
	Bearing  float64 `bson:"bearing,omitempty" json:"bearing"`
}

func main() {
	points := triangle.ImportShapefile("limfjorden/Limfjorden.shp")
	fmt.Println(len(points))
	go run()
	select {}

}
func run() {
	MapStore := NewMapStore()

	count, err := MapStore.points.Count()
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Number of points", count)

	// count, err = MapStore.polygons.Count()
	// if err != nil {
	// 	log.Println(err.Error())
	// }
	// log.Println("Number of polygons", count)

	// query the database
	// db.points.find({geometry:{$nearSphere:{$geometry:{type:"Point",coordinates:[9.5704684,56.9268073]}}}})
	// err = points.Find(bson.M{"geometry": bson.M{"$nearSphere": bson.M{"$geometry": results.Geometry}}}).One(&results)
	//node := findNearestNode(points, geo.NewPoint(57.01189192290, 9.52040228844))
	//fmt.Println(node)

	// result := MapStore.findNode(12)
	// response, err := json.MarshalIndent(result, " ", "   ")

	// if err != nil {
	// 	panic(err)
	// }
	//Uncomment to update neighbors on mongo
	//for i := 0; i < count; i++ {
	//	MapStore.UpdateNeighbours(i)
	//}
	log.Println("done Mingling")
	// t1 and t2 are *Tile objects from inside the world.
	start := time.Now()

	fmt.Println(time.Now().Sub(start))

	start = time.Now()
	path, distance, found := MapStore.Path(MapStore.findNode(3990), MapStore.findNode(6523))
	if !found {
		log.Println("Could not find path")
	}
	fmt.Println(time.Now().Sub(start))
	fmt.Println("total distance", distance)
	fmt.Println("no of nodes in path: ", len(path))
	fmt.Println("no of nodes checked: ", len(MapStore.pointsCache))

	start = time.Now()

	fmt.Println(path)

	fc := gj.NewFeatureCollection([]*gj.Feature{})

	// feature

	for i := 0; i < len(path); i++ {
		point := MapStore.findNode(path[i])

		//lat := gj.Coord(point.Geometry.Coordinates[0])
		//p := gj.NewPoint(gj.Coordinate{12, 3.123})
		p := gj.NewPoint(gj.Coordinate{gj.Coord(point.Geometry.Coordinates[0]), gj.Coord(point.Geometry.Coordinates[1])})
		f1 := gj.NewFeature(p, nil, nil)
		fc.AddFeatures(f1)

	}
	if gjstr, err := gj.Marshal(fc); err != nil {
		panic(err)
	} else {
		fmt.Println(gjstr)
	}

	// path is a slice of Pather objects which you can cast back to *Tile.

}

func (m *MapStore) findNode(node int) *Point {

	p := m.pointsCache[node]
	if p == nil {

		results := Point{mapStore: m, timeStamp: time.Now().Add(5 * time.Minute)}
		err := m.points.Find(bson.M{"properties.node": node}).One(&results)

		if err != nil {
			log.Println(err.Error())
		}
		m.pointsCache[node] = &results
		return &results
	} else {
		return p
	}
}

func (m *MapStore) findNearestNode(g *geo.Point) int {
	var basePoint Point
	err := m.points.Find(bson.M{"geometry": bson.M{"$nearSphere": bson.M{"$geometry": bson.M{
		"type": "Point",
		"coordinates": []float64{g.Lng(),
			g.Lat()},
	},
	//"$maxDistance": 300,
	},
	}}).One(&basePoint)

	if err != nil {
		panic(err.Error())
	} else {
		return basePoint.Properties.Node
	}
}

func (m *MapStore) UpdateNeighbours(n int) {

	var basePoint Point
	var point []Point
	//Find basePoint
	err := m.points.Find(bson.M{"properties.node": n}).One(&basePoint)
	if err != nil {
		log.Println(err.Error())
	} else {

		err = m.points.Find(bson.M{
			"$or": []bson.M{
				{
					"properties.POINTA": bson.M{"$in": []int{basePoint.Properties.POINTA, basePoint.Properties.POINTB, basePoint.Properties.POINTC}}, "properties.POINTB": bson.M{"$in": []int{basePoint.Properties.POINTA, basePoint.Properties.POINTB, basePoint.Properties.POINTC}},
				},
				{
					"properties.POINTB": bson.M{"$in": []int{basePoint.Properties.POINTA, basePoint.Properties.POINTB, basePoint.Properties.POINTC}}, "properties.POINTC": bson.M{"$in": []int{basePoint.Properties.POINTA, basePoint.Properties.POINTB, basePoint.Properties.POINTC}},
				},
				{
					"properties.POINTC": bson.M{"$in": []int{basePoint.Properties.POINTA, basePoint.Properties.POINTB, basePoint.Properties.POINTC}}, "properties.POINTA": bson.M{"$in": []int{basePoint.Properties.POINTA, basePoint.Properties.POINTB, basePoint.Properties.POINTC}},
				},
			},
		}).All(&point)

		if err != nil {
			log.Println(err.Error())
		}
		p1 := geo.NewPoint(basePoint.Geometry.Coordinates[0], basePoint.Geometry.Coordinates[1])

		basePoint.Properties.Naighbors = make([]Naighbor, 0, 3)

		for i := 0; i < len(point); i++ {
			if point[i].Properties.Node != n {
				var naighbor Naighbor
				p2 := geo.NewPoint(point[i].Geometry.Coordinates[0], point[i].Geometry.Coordinates[1])
				naighbor.Distance = p1.GreatCircleDistance(p2)
				naighbor.Bearing = math.Mod(p1.BearingTo(p2)+2*math.Pi, 2*math.Pi)
				naighbor.Node = point[i].Properties.Node
				basePoint.Properties.Naighbors = append(basePoint.Properties.Naighbors, naighbor)

			}
		}
		m.points.Update(bson.M{"properties.node": n}, basePoint)

		response, err := json.MarshalIndent(basePoint, " ", "   ")

		if err != nil {
			panic(err)
		}
		fmt.Println(string(response))

	}
}

//From here on, all credit goes to https://github.com/beefsack/go-astar

// Path calculates a short path and the distance between the two Point nodes.
//
// If no path is found, found will be false.
func (m *MapStore) Path(from, to *Point) (path []int, distance float64, found bool) {
	nm := nodeMap{}
	nq := &priorityQueue{}
	counter := 0
	heap.Init(nq)
	fromNode := nm.get(from.Properties.Node)
	fromNode.open = true
	heap.Push(nq, fromNode)
	for {
		if nq.Len() == 0 {
			// There's no path, return found false.
			return
		}
		current := heap.Pop(nq).(*node)
		current.open = false
		current.closed = true

		if current == nm.get(to.Properties.Node) {
			// Found a path to the goal.
			p := []int{}
			curr := current
			for curr != nil {
				p = append(p, curr.point)

				curr = curr.parent
			}
			return p, current.cost, true
		}

		counter++

		for i := 0; i < len(m.findNode(current.point).Properties.Naighbors); i++ {
			point := m.findNode(current.point).Properties
			if point.Naighbors[i].Node != current.point {

				//cost := current.cost + point.Naighbors[i].Distance
				cost := current.cost + 1

				neighborNode := nm.get(point.Naighbors[i].Node)
				if cost < neighborNode.cost {
					if neighborNode.open {
						heap.Remove(nq, neighborNode.index)
					}
					neighborNode.open = false
					neighborNode.closed = false
				}
				if !neighborNode.open && !neighborNode.closed {
					neighborNode.cost = cost
					neighborNode.open = true
					p := m.findNode(neighborNode.point)
					p1 := geo.NewPoint(p.Geometry.Coordinates[0], p.Geometry.Coordinates[1])
					p2 := geo.NewPoint(to.Geometry.Coordinates[0], to.Geometry.Coordinates[1])

					neighborNode.rank = cost + math.Abs(p1.Lat()-p2.Lat()) + math.Abs(p1.Lng()-p2.Lng())
					neighborNode.rank = cost + p1.GreatCircleDistance(p2)
					neighborNode.parent = current
					heap.Push(nq, neighborNode)
				}
			}
		}
	}
}

// A priorityQueue implements heap.Interface and holds Nodes.  The
// priorityQueue is used to track open nodes by rank.
type priorityQueue []*node

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].rank < pq[j].rank
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	no := x.(*node)
	no.index = n
	*pq = append(*pq, no)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	no := old[n-1]
	no.index = -1
	*pq = old[0 : n-1]
	return no
}

// astar is an A* pathfinding implementation.

// node is a wrapper to store A* data for a Point node.
type node struct {
	point  int
	cost   float64
	rank   float64
	parent *node
	open   bool
	closed bool
	index  int
}

// nodeMap is a collection of nodes keyed by Point nodes for quick reference.
type nodeMap map[int]*node

// get gets the Point object wrapped in a node, instantiating if required.
func (nm nodeMap) get(p int) *node {
	n, ok := nm[p]
	if !ok {
		n = &node{
			point: p,
		}
		nm[p] = n
		fmt.Println(p)
	}
	return n
}
