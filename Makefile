# created on Sat Feb  7 22:13:04 EST 2015 by Mark bucciarelli
# Plot performance data.

cqrsprof.svg: gnuplot.script sqlite_data gob_data
	gnuplot < gnuplot.script

gob_data: data.0772fa0
sqlite_data: data.b8f55aa

# Make data rebuilds cqrs and cqrsprof.
data.0772fa0: makedata
	./setstore FileSystem
	./makedata $@

# Make data rebuilds cqrs and cqrsprof.
data.b8f55aa: makedata
	./setstore Sqlite
	./makedata $@
