
ns stddbc

assert = func(condition, errtext)
	condtype = type(condition)
	condval = case( condtype,
		'bool', 	condition,
		'function', call(condition),
		true
	)

	if(condval, true, error(errtext))
end

endns
