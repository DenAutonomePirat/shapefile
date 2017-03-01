import triangle
import triangle.plot
import matplotlib.pyplot as plt
import geojson # or import geojson
import matplotlib.pyplot as plt 
from descartes import PolygonPatch

with open("Limfjordenudsnit.geojson") as json_file:
    json_data = geojson.load(json_file) # or geojson.load(json_file)


poly = json_data.features[0].geometry
print poly

pts = poly.coordinates[0]
#print pts
import numpy as np
tri = {}
tri["vertices"] = pts


markers = triangle.triangulate(tri)
#print markers


fig = plt.figure() 
ax = fig.gca() 
BLUE = '#6699cc'

for i in markers["triangles"]:

	for j in markers["triangles"][i]:
		coordinates = markers["vertices"][j]
		
	p ={}
	p["coordinates"] = coordinates

	p["type"] = "polygon"
	print p




	ax.add_patch(PolygonPatch(p, fc=BLUE, ec=BLUE, alpha=0.5, zorder=2 ))





#ax.add_patch(PolygonPatch(poly, fc=BLUE, ec=BLUE, alpha=0.5, zorder=2 ))

ax.axis('scaled')
plt.show()