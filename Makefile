# created on Sat Feb  7 22:13:04 EST 2015 by Mark bucciarelli
# Create new SVG plot.

cqrsprof.svg: gnuplot.script data.*
	gnuplot < gnuplot.script
