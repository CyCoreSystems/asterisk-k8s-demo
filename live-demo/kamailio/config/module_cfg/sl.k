##
## Module: sl 
##

# Default status code of a stateless reply, if unspecified.

modparam("sl", "default_code", 505)

# The default reason phrase of a stateless reply, if unspecified.

modparam("sl", "default_reason", "Server Error")

# Whether to bind to the TM module in order to send a stateful reply
# if a transaction has been created, provided that send_reply() is used,
# rather than sl_send_reply().

modparam("sl", "bind_tm", 1)
