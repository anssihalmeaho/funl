
ns impure_test

import ut_fwk

ASSURE = ut_fwk.VERIFY

testBufferedChannel = proc()
	ch = chan(2)

	ret-1 = recwith(ch map('wait' false))
	tst-result-1 = call(ASSURE eq(ret-1 list(false '')) plus('Unexpected value from recwith: ' str(ret-1)))

	ret-2 = list(
		send(ch 'any value' map('wait' false))
		send(ch 'any value' map('wait' false))
		send(ch 'any value' map('wait' false))
	)
	tst-result-2 = call(ASSURE eq(ret-2 list(true true false)) plus('Unexpected value from send: ' str(ret-2)))

	and(tst-result-1 tst-result-2)
end

testOneFiberCreateAndReply = proc()
	replyCh = chan()
	testMsg = 'this is test message'
	fib = proc()
		send(replyCh, testMsg)
	end
	_ = spawn(call(fib))
	reply = recv(replyCh)

	call(ASSURE, eq(testMsg, reply), plus('Unexpected reply = ', str(reply)))
end

testOneFiberCreateAndReply2 = proc()
	replyCh = chan()
	testMsg = 'this is test message'
	_ = spawn(send(replyCh, testMsg))
	reply = recv(replyCh)

	call(ASSURE, eq(testMsg, reply), plus('Unexpected reply = ', str(reply)))
end

testTwoFiberCommAsRing = proc()
	fib = proc(ch, nextch)
		send(nextch, plus(recv(ch), '.'))
	end

	testMsg = 'this is test message'
	ch1 = chan()
	ch2 = chan()
	ch3 = chan()
	_ = spawn(call(fib, ch2, ch3))
	_ = spawn(call(fib, ch1, ch2))
	_ = send(ch1, testMsg)
	reply = recv(ch3)
	call(ASSURE, eq(plus(testMsg, '..'), reply), plus('Unexpected reply = ', str(reply)))
end

testStackedFibers = proc()
	amount = 10

	fib = proc(count, replych)
		nextCount = plus(count, 1)

		if( eq(nextCount, amount),
			send(replych, list(minus(nextCount, 1), str(nextCount))),
			call(proc(repch)
					wchan = chan()
					_ = spawn(call(fib, nextCount, wchan))
					val = recv(wchan)
					ncnt = minus(head(val), 1)
					nstr = plus(last(val), str(ncnt))
					send(repch, list(ncnt, nstr))
				end,
				replych
			)
		)
	end

	resultch = chan()
	_ = spawn(call(fib, 1, resultch))
	resultVal = recv(resultch)

	call(ASSURE, eq(resultVal, list(1, '1087654321')), plus('Unexpected resultVal = ', str(resultVal)))
end

testSelectReceive = proc()
	resultCh = chan()
	channels = list(chan(), chan(), chan())
	quitch = chan()
	controlCh = chan()

	sender = proc(num)
		send(ind(channels, num), num)
	end

	receiver = proc(count, s)
		getMsgHandler = func(char)
			proc(message)
				_ = send(controlCh, true)
				plus(char, str(message), ':')
			end
		end

		msg = select(
			ind(channels, 0), call(getMsgHandler, 'A'),
			ind(channels, 1), call(getMsgHandler, 'B'),
			ind(channels, 2), call(getMsgHandler, 'C'),
			quitch, func(m) s end
		)

		while( lt(count, 3),
			plus(count, 1),
			plus(s, msg),
			send(resultCh, msg)
		)
	end

	_ = spawn(call(receiver, 0, ''))

	_ = spawn(call(sender, 0))
	_ = spawn(call(sender, 1))
	_ = spawn(call(sender, 2))

	wait = proc(n)
		while( lt(n, 3),
			call(
				proc(nn)
					_ = recv(controlCh)
					plus(nn, 1)
				end,
				n
			),
			true
		)
	end

	_ = call(wait, 0)
	_ = send(quitch, 'quit')
	result = split(recv(resultCh), ':')

	allOk = and(
		in(result, 'A0'),
		in(result, 'B1'),
		in(result, 'C2'),
		in(result, ''),
		eq(len(result), 4)
	)
	call(ASSURE, allOk, plus('Unexpected result = ', str(result)))
end

testSelectReceiveWithListOfHandlers = proc()
	resultCh = chan()
	channels = list(chan(), chan(), chan())
	quitch = chan()
	controlCh = chan()

	sender = proc(num)
		send(ind(channels, num), num)
	end

	receiver = proc(count, s)
		getMsgHandler = func(char)
			proc(message)
				_ = send(controlCh, true)
				plus(char, str(message), ':')
			end
		end

		handlers = list(
			call(getMsgHandler, 'A'),
			call(getMsgHandler, 'B'),
			call(getMsgHandler, 'C'),
			func(m) s end
		)

		msg = select(append(channels, quitch), handlers)

		while( lt(count, 3),
			plus(count, 1),
			plus(s, msg),
			send(resultCh, msg)
		)
	end

	_ = spawn(call(receiver, 0, ''))

	_ = spawn(call(sender, 0))
	_ = spawn(call(sender, 1))
	_ = spawn(call(sender, 2))

	wait = proc(n)
		while( lt(n, 3),
			call(
				proc(nn)
					_ = recv(controlCh)
					plus(nn, 1)
				end,
				n
			),
			true
		)
	end

	_ = call(wait, 0)
	_ = send(quitch, 'quit')
	result = split(recv(resultCh), ':')

	allOk = and(
		in(result, 'A0'),
		in(result, 'B1'),
		in(result, 'C2'),
		in(result, ''),
		eq(len(result), 4)
	)
	call(ASSURE, allOk, plus('Unexpected result = ', str(result)))
end

endns
