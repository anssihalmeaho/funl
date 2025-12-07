
ns stdcsv2

lines-to-maps = func(lines)
	import stdfu

	if(empty(lines)
		list()
		call(func()
			header = head(lines)
			call(stdfu.apply rest(lines) func(line)
				call(stdfu.zip header line)
			end)
		end)
	)
end

csv-to-maps = func(bytes)
	import stdcsv

	ok err lines = if( eq(len(argslist()) 2)
		call(stdcsv.read-all bytes last(argslist()))
		call(stdcsv.read-all bytes)
	):
	if(ok
		list(true '' call(lines-to-maps lines))
		list(false err map())
	)
end

maps-to-lines = func(header maps)
	import stdfu

	convert-to-list = func(onemap)
		call(stdfu.foreach
			header
			func(key onelist)
				value = get(onemap key)
				append(onelist value)
			end
			list()
		)
	end

	other-lines = call(stdfu.apply maps convert-to-list)
	add(other-lines header)
end

maps-to-csv = func(header maps)
	import stdcsv

	lines = call(maps-to-lines header maps)
	if( eq(len(argslist()) 3)
		call(stdcsv.write-all lines last(argslist()))
		call(stdcsv.write-all lines)
	)
end

endns

