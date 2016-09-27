##########################################################################
# Kamailio 3.x configuration basis - registrar module parameters.        #
# 									 #
# For use by Evariste Systems LLC and authorised third parties only.     #
#									 #
# Boilerplate parameters, code and comment annotations by:               #
#									 #
# Alex Balashov <abalashov@evaristesys.com>                              #
##########################################################################

#!ifdef WITH_REGISTRAR

##
## Module: registrar
##

# Default binding expiration value if the processed REGISTER message has no
# Expires headers nor expiration-related parameters in the Contact.

modparam("registrar", "default_expires", 90)

# The minimum expiration value of a contact binding (seconds);  values 
# requested by the UAC that are lower than this minimum will automatically 
# be floored to this minimum.

modparam("registrar", "min_expires", 90)

# The maximum expiration value of a contact binding (seconds);  values
# requested by the UAC that are higher than this value will automatically
# set to this ceiling.

modparam("registrar", "max_expires", 180)

# Default q-value for contact bindings.

modparam("registrar", "default_q", 0)

# Realm prefix to strip from domain prior to binding.

#modparam("registrar", "realm_prefix", "sip.")

# If enabled, will cause multiple contacts to automatically be parallel
# forked by TM after lookup().

modparam("registrar", "append_branches", 1)

# If enabled, the AOR comparison on lookup() will be case-sensitive.

modparam("registrar", "case_sensitive", 0)

# Maximum concurrent contacts for a given AOR binding.  This value controls
# the number of parallel registrations possible for one "user name."

modparam("registrar", "max_contacts", 5)

# Value of Retry-After header appended to 5xx replies to registration attempts 
# in certain situations.  If set to 0, no Retry-After header will be added.

modparam("registrar", "retry_after", 0)

# Filter to supported methods on lookup()?

modparam("registrar", "method_filtering", 0)

# Far-end NAT traversal stuff.

modparam("registrar", "received_avp", "$avp(s:recv_nat)")
modparam("registrar", "received_param", "recv_nat")

#!endif
