
ns common_test_util

isAllTrueInList = func(l)
	checker = func(li)
		if(empty(li),
			true,
			if(head(li), 
				call(checker, rest(li)),
				false
			)
		)
	end
	
	call(checker, l)
end

endns

