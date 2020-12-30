
ns stdbase64_test

import ut_fwk

ASSURE = ut_fwk.VERIFY

import common_test_util
import stdbase64
import stdbytes

test-encode-decode = func()
	data = call(stdbytes.str-to-bytes 'abcde')

	enc-ok enc-err encoded = call(stdbase64.encode data):
	_ = call(ASSURE enc-ok enc-err)

	dec-ok dec-err decoded = call(stdbase64.decode encoded):
	_ = call(ASSURE dec-ok dec-err)

	s = call(stdbytes.string decoded)
	call(ASSURE eq(s 'abcde') plus('unexpected result: ' s))
end

endns

