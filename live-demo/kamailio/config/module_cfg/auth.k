##
## Module: auth
##

#!ifdef WITH_DIGEST_AUTH

# There are three parameters that govern additional SIP message integrity
# checks against three different types of SIP messages in connection with
# which digest authentication may be used by the script writer:
#
#   - auth_checks_register - Checks on REGISTER messages.
#   - auth_checks_no_dlg - Checks on requests outside of a dialog
#                          (initial requests).
#   - auth_checks_in_dlg - Checks on sequential requests within a dialog
#			   such as a BYE or re-INVITE.
#
# Possible check options for all three of these are a bitwise AND-ed bitmask 
# of the following flags:
# 
#   - 1: Check if the request/target URI changed.
#   - 2: Check the Call-ID.
#   - 4: Check the From tag.
#   - 8: Check the source IP.

modparam("auth", "auth_checks_register", 11)	# RURI, Call-ID and src_ip.
modparam("auth", "auth_checks_no_dlg", 9)	# RURI, src_ip, Call-ID.
modparam("auth", "auth_checks_in_dlg", 15)	# Everything.

# Enable 'qop' for challenges.  This option may confuse non-RFC 3261 UAs, 
# but is otherwise relatively safe.  The 'qop' feature, together with
# 'nonce_count' (see below), provides some extra protection against replay
# attacks while still allowing credential caching on the UA side.
# 
# Valid values are "auth", "auth-int", and "" (unset).

modparam("auth", "qop", "auth")

# If enabled, the received nonce is remembered and checked against the older
# value;  it protects against replay attacks.
#
# Note: This protection only works if authentication is done in stateful
# mode, meaning that a transaction should be created prior to the 
# authentication challenge, and all replies should be sent statefully.

modparam("auth", "nonce_count", 1)

# If enabled, allows each nonce from a UA to be used only once;  all messages
# will be challenged, even retransmissions.  Stateful mode should be used
# to catch retransmissions before authentication checks, and all replies
# should be sent with t_reply().

modparam("auth", "one_time_nonce", 1)

# Number of partitions in the array used to hold both nonce_count and 
# one_time_nonce values.  Increasing this value can have a detrimental
# effect in highly concurrent situations because it can decrease the
# maximum supported in-flight nonces, although it does, in theory, 
# increase performance.  
# 
# This value should be a power of 2, up to 64.

modparam("auth", "nid_pool_no", 4)

# Maximum number of in-flight nonces for 'nonce_count'; that is, it is the
# maximum number of nonces for which state should be kept.  
# This should be a power of 2.

modparam("auth", "nc_array_size", 8388608)  # 8 MB.

# Maximum number of in-flight nonces for 'one_time_nonce'.  When this quantity
# is exceeded, older nonces are discarded and overwritten.

modparam("auth", "otn_in_flight_no", 8388608)  # 8 MB

# Shared secret, mainly used to pass credentials around multiple proxies.
# If unset, it defaults to a randomly generated string.

#modparam("auth", "secret", "")

# The expiration value of a nonce in seconds (nonces have a limited lifetime).
# This is to protect against replay attacks.  Any user agent with an expired
# nonce will be challenged again for a new nonce.  The challenge in this
# case will contain a "stale" parameter, indicating that the user does not
# need to be disturbed in an interactive fashion (e.g. re-prompted for a 
# username and password), but instead the existing credentials should be
# recalculated.

modparam("auth", "nonce_expire", 600)  # 10 minutes.

# Maximum difference (in seconds) between a nonce creation time and the
# current time, if nonce creation time appears in the future (due to 
# clock drift).

modparam("auth", "nonce_auth_max_drift", 5)  # 5 seconds

# If set to 1, www_challenge() and proxy_challenge() will send a 
# reply statelessly regardless of whether a transaction exists.  If set to 0,
# the default value, the reply will be sent statelessly or statefully
# depending on whether a transaction exists.

modparam("auth", "force_stateless_reply", 0)

#!endif 
