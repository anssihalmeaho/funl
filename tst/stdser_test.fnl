
ns stdser_test

import ut_fwk

ASSURE = ut_fwk.VERIFY

import common_test_util
import stdbytes
import stdser

test-bytearray-encode-decode = func()
	data = call(stdbytes.str-to-bytes 'abcde')

	enc-ok enc-err encoded = call(stdser.encode data):
	_ = call(ASSURE enc-ok enc-err)

	dec-ok dec-err decoded = call(stdser.decode encoded):
	_ = call(ASSURE dec-ok dec-err)

	s = call(stdbytes.string decoded)
	call(ASSURE eq(s 'abcde') plus('unexpected result: ' s))
end

endns

