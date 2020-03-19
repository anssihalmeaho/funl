
ns main

import stdio
import stdfu

None = 1
X = 2
O = 3

machine = 1
player = 2

#--- gets empty grid as basis
get-initial-grid = func()
	map()
end

symbol-map = map(
	None ' '
	X    'X'
	O    'O'
)

#--- prints grid to screen
print-grid = proc(grid)
	print-line = proc(start-ind)
		v1 = if(in(grid start-ind) get(grid start-ind) None)
		v2 = if(in(grid plus(start-ind 1)) get(grid plus(start-ind 1)) None)
		v3 = if(in(grid plus(start-ind 2)) get(grid plus(start-ind 2)) None)
		
		_ = call(stdio.printf '  %s | %s | %s\n' 
			get(symbol-map v1)
			get(symbol-map v2)
			get(symbol-map v3)
		)
		_ = if(not(eq(start-ind 7)) call(stdio.printf '  --+---+---\n') '')
		true
	end
	
	print-it = proc(line-count)
		_ = call(stdio.printf '\n')
		_ = call(print-line 1)
		_ = call(print-line 4)
		_ = call(print-line 7)
		_ = call(stdio.printf '\n')
		true
	end

	call(print-it 0)
end

#--- has-game-ended
has-game-ended = call(func()
	import stdset
	import stdfu
	
	winning-list = list(
		list(1 2 3)
		list(4 5 6)
		list(7 8 9)
		
		list(1 4 7)
		list(2 5 8)
		list(3 6 9)

		list(1 5 9)
		list(3 5 7)		
	)

	check-one-line = func(line marked)
		set1 = call(stdset.list-to-set call(stdset.newset) line)
		set2 = call(stdset.list-to-set call(stdset.newset) marked)
		#_ = print('marked: ' marked)
		call(stdset.is-subset set2 set1)
	end
	
	func(grid)
		is-winner = func(mark)
			marked = keys(call(stdfu.filter grid func(_ v) eq(v mark) end))

			checkloop = func(winlist found)
				while( and( not(found) not(empty(winlist)) )
					rest(winlist)
					call(check-one-line head(winlist) marked)
					found
				)
			end

			call(checkloop winning-list false)
		end
		
		cond(
			call(is-winner X) list(true true player)
			call(is-winner O) list(true true machine)
			gt(len(grid) 8)   list(true false 0)
			list(false false 0)
		)
	end
end)

minmax = func(grid turn depth)
	swap = func(whos-turn)
		if( eq(whos-turn machine)
			player
			machine
		)
	end

	min = func(idx1 idx2)
		if( gt(idx1 idx2)
			idx2
			idx1
		)
	end

	max = func(idx1 idx2)
		if( gt(idx1 idx2)
			idx1
			idx2
		)
	end
	
	get-mark = func(whos-turn)
		if( eq(whos-turn machine)
			O
			X
		)
	end
	
	get-one-score = func(idx current-score)
		next-grid = put(grid idx call(get-mark turn))
		is-end has-winner winner = call(has-game-ended next-grid):
		score = if( is-end
		if( has-winner
					case( turn
						player  minus(0 10)
						machine 10
					)
					0 # no winner, even game
			)
			call(minmax next-grid call(swap turn) plus(depth 1))
		)
		if( eq(current-score 'no-score')
			score
			if( eq(turn machine)
				call(max score current-score)
				call(min score current-score)
			)
		)
	end
	
	looper = func(avail current-score)
		while( not(empty(avail))
			rest(avail)
			call(get-one-score head(avail) current-score 1)
			current-score
		)
	end
	
	all-idxs = list(1 2 3 4 5 6 7 8 9)
	used-idxs = keys(grid)
	available-spots = call(stdfu.filter all-idxs func(idx) not(in(used-idxs idx)) end)

	score = call(looper available-spots 'no-score')
	score
end

scores-for-availables = func(grid)
	proxy = func(item)
		ngrid = put(grid item O)
		is-end has-winner winner = call(has-game-ended ngrid):
		if( is-end
			if( has-winner 
				if( eq(winner machine)
					10
					minus(0 10)
				)
				0
			)
			call(minmax ngrid player 1)
		)
	end
	
	looper = func(avail result)
		while( not(empty(avail))
			rest(avail)
			put(result head(avail) call(proxy head(avail)))
			result
		)
	end
	
	all-idxs = list(1 2 3 4 5 6 7 8 9)
	used-idxs = keys(grid)
	available-spots = call(stdfu.filter all-idxs func(idx) not(in(used-idxs idx)) end)
	call(looper available-spots map())
end

#--- players-move
players-move = proc(grid)
	mark-move = func(index-value)
		is-valid-index = and(
			lt(index-value 10)
			gt(index-value 0)
		)
		cond(
			not(is-valid-index)  list(grid false)
			in(grid index-value) list(grid false)
			list(put(grid index-value X) true)
		)
	end

	ask-move = proc()
		_ = call(print-grid grid)
		_ = call(stdio.printline ' What is your move ? ')
		index = conv(call(stdio.readinput) 'int')
		if( eq(type(index) 'string')
			list(grid false)
			call(mark-move index)
		)
	end

	if( eq(len(grid) 8) # there's only one possibility
		call(proc()
			used-spots = keys(grid)
			missing-spot = call(stdfu.filter list(1 2 3 4 5 6 7 8 9) func(idx) not(in(used-spots idx)) end):
			call(mark-move missing-spot)
		end)

		call(proc()		
			ngrid is-valid-move = call(ask-move):
			_ = if(not(is-valid-move) call(stdio.printf '\n...invalid input, please retry...\n') true)
			while( not(is-valid-move) list(ngrid true) )
		end)
	)
end

#--- machines-move
machines-move = proc(source-grid machine-mark)	
	decide-move = func()
		score-map = call(scores-for-availables source-grid)
		grouped = call(stdfu.group-by score-map func(k v) list(v k) end)
		cond(
			in(grouped 10) head(get(grouped 10))
			in(grouped 0) head(get(grouped 0))
			in(grouped minus(0 10)) head(get(grouped minus(0 10)))
			error('something wrong...' grouped)
		)
	end
	
	decide-first = func()
		cond(
			not(in(source-grid 5)) 5
			1
		)
	end

	machines-index = if( eq(len(source-grid) 1)
		call(decide-first)
		call(decide-move)
	)
	next-grid = put(source-grid machines-index machine-mark)
	list(next-grid true)
end

#--- print-error
print-error = proc()
	call(stdio.printline 'Error, game end.')
end

#--- play 
play = proc(grid next-to-play)

	#--- get-result-printout
	get-result-printout = func(has-winner winner)
		winner-name = case( winner
			player  'You'
			machine 'Me'
			'No winner'
		)
		if( has-winner
			sprintf('\n Winner is %s ! \n' winner-name)
			sprintf('\n Game even. \n')
		)
	end

	next-grid ok = case( next-to-play
		player  call(players-move grid)
		machine call(machines-move grid O)
	):
	other-in-turn = case( next-to-play
		player  machine
		machine player
	)
	if( ok
		call(proc()
			is-end has-winner winner = call(has-game-ended next-grid):
			if( is-end
				call(proc()
					_ = call(print-grid next-grid)
					call(get-result-printout has-winner winner)
				end)
				call(play next-grid other-in-turn)
			)
		end)
		call(print-error)
	)
end

#--- main
main = proc()
	_ = call(stdio.printf ' Moves are given with numbers in range 1-9\n')
	grid = call(get-initial-grid)	
	call(play grid player)
end

endns
