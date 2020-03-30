
ns ut_fwk

VERIFY = func(condition, errText)
	if(condition,
		true,
		call(func()
			_ = print(plus(name(VERIFY), ' fail: '), errText)
			false
		end)
	)
end

endns
