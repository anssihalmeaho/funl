
ns main

main = proc()
	import stdfu

	fizzbuzz = func(int-list)
		is-divisible-by = func(dividend divisor)
			eq(mod(dividend divisor) 0)
		end

		choose-string = func(i result-str)
			is-divisible-by-5 = call(is-divisible-by i 5)
			is-divisible-by-3 = call(is-divisible-by i 3)
			plus(result-str '\n' cond(
				and(is-divisible-by-5 is-divisible-by-3) 'FizzBuzz'
				is-divisible-by-5 'Buzz'
				is-divisible-by-3 'Fizz'
				str(i)
			))
		end

		call(stdfu.loop choose-string int-list '')
	end

	call(fizzbuzz call(stdfu.generate 1 100 func(n) n end))
end

endns
