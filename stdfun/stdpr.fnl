
ns stdpr

get-pr = func(do-print)
	if( do-print
		func(text value)
			_ = print(text value)
			value
		end

		func(_ value) value end
	)
end

get-pp-pr = func(do-print)
	if( do-print
		func(text value)
			import stdpp

			_ = print(text call(stdpp.form value))
			value
		end

		func(_ value) value end
	)
end

endns
