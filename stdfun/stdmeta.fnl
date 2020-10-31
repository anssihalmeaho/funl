
ns stdmeta

/*
todo

- list('list' list-schema)
    -> validates with list schema

- list('exact-keys' 'f1' 'f2')
    -> should have exactly following keys (not less, not more)

- list('len' 5)
     -> length of list or map (key-values) checked
	 -> ranges too ?

- list('all-in-list' list('type' 'int'))
     -> requirement(s) for all items in list
	 -> could it be also for all key-values in map ?

- validation with customized function(s)...?

*/

# gathers documentation from schema
get-doc = func(schema)
	import stdfu

	visit-map = func(map-schema indent result-str)
		map-fields-visit = func(kvitem output)
			get-field-visitor = func(mkey)
				func(check out-str)
					fout = case( head(check)
						'required'
							'required, '
						'type'
							plus('type: ' head(rest(check)) ', ')
						'map'
							call(visit-map head(rest(check)) plus(indent '    ') 'map: ')
						'doc'
							call(func()
								doc-lines = call(stdfu.apply rest(check) func(item) plus('\n' indent '  -> ' item) end)
								plus(doc-lines: ', ')
							end)
						'in'
							plus('allowed: [ ' call(stdfu.apply rest(check) func(x) plus(str(x) ' ') end): '], ')
						error('illegal tag: ' str(head(check)))
					)
					plus(out-str fout)
				end
			end

			keyv checklist = kvitem:
			field-output = call(stdfu.loop call(get-field-visitor keyv) checklist '')
			plus(output '\n' indent str(keyv) ' : ' field-output)
		end

		if( eq(type(map-schema) 'map')
			call(func()
				map-output = call(stdfu.loop map-fields-visit keyvals(map-schema) '')
				plus('\n' indent 'map: ' map-output)
			end)
			'requires map'
		)
	end

	if( and( eq(type(schema) 'list') not(empty(schema)))
		case( head(schema)
			'map' call(visit-map head(rest(schema)) '' '')
			'none'
		)
		'requires non-empty list'
	)
end

# validates data against schema
validate = func(schema srcdata)
	import stdfu

	validate-map = func(map-schema map-data field-path)
		map-fields-checker = func(kvitem results)
			get-field-checker = func(mkey)
				func(check result-list)
					_ msg-list = result-list:
					passed message = case( head(check)
						'required'
							if( in(map-data mkey)
								result-list
								list(false append(msg-list sprintf('required field %v not found (%s)' mkey field-path)))
							)

						'in'
							call(func()
								key-found val = getl(map-data mkey):
								allowed = rest(check)
								cond(
									not(key-found)
										result-list
									not(in(allowed val))
										list(false append(msg-list sprintf('field %v is not in allowed set (%v not in: %v)(%s)' mkey val allowed field-path)))
									result-list
								)
							end)

						'type'
							call(func()
								key-found val = getl(map-data mkey):
								required-type = head(rest(check))
								cond(
									not(key-found)
										result-list
									not(eq(type(val) required-type))
										list(false append(msg-list sprintf('field %v is not required type (got: %v, expected: %v)(%s)' mkey type(val) required-type field-path)))
									result-list
								)
							end)

						'map'
							call(func()
								submap-found submap-data = getl(map-data mkey):
								cond(
									not(submap-found)
										list(false append(msg-list sprintf('field %v (sub map) not found (%s)' mkey field-path)))
									not(eq(type(submap-data) 'map'))
										list(false append(msg-list sprintf('field %v is not map (%s)' mkey field-path)))
									call(validate-map head(rest(check)) submap-data plus(field-path ' -> ' str(mkey)))
								)
							end)

						'doc'
							result-list

						list(false append(msg-list sprintf('unknown validator: %s' str(head(check)))))
					):
					list(passed message)
				end
			end

			keyv checklist = kvitem:
			call(stdfu.loop call(get-field-checker keyv) checklist results)
		end

		if( eq(type(map-schema) 'map')
			call(stdfu.loop map-fields-checker keyvals(map-schema) list(true list()))
			list(false list('requires map'))
		)
	end

	if( and( eq(type(schema) 'list') not(empty(schema)))
		case( head(schema)
			'map' call(validate-map head(rest(schema)) srcdata '')
			list(false list(sprintf('unknown validator: %s' str(head(schema)) )))
		)
		list(false list('requires non-empty list'))
	)
end

endns

