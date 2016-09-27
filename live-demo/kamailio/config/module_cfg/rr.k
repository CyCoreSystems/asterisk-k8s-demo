## 
## Module: rr
##

# Enable full value for loose-routing parameter (;lr=on) in Route header
# fields.  Some broken UAs strip a bare ;lr parameter with no value 
# assigned.

modparam("rr", "enable_full_lr", 1)

# Append request's From tag to Record-Route;  this is useful for determining
# whether sequential in-dialog requests come from the caller or the callee
# (i.e. in the case of BYE, who hung up the call).

modparam("rr", "append_fromtag", 1)

# Insert two Record-Route headers in certain situations (TCP<->UDP, 
# two disconnected networks, etc.)

modparam("rr", "enable_double_rr", 1)

# If enabled, will add username to the Record-Route URI.  This should be 
# disabled.

modparam("rr", "add_username", 0)
