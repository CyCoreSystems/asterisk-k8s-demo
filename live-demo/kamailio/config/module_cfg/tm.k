##
## Module: tm
##

# Timer (in milliseconds) that fires if no final reply for a request 
# or ACK for a negative INVITE reply arrives. 
#
# Despite what the documentation implies, this timer is actually for
# _any_ reply received to a request, including provisional (1xx) 
# replies.  If a provisional reply is provided in a time under this
# timer, it is discharged;  otherwise, when it fires, the call will 
# be considered to have failed and will go to a failure route.  This
# option can be tuned for fast failover via serial forking.

modparam("tm", "fr_timer", 3000)	# 3 seconds

# Timer (in milliseconds) that fires if no final reply for an INVITE
# arrives after a provisional message was received.  
#
# Note that this timer can be restarted when a provisional message is
# received by the restart_fr_on_each_reply message.
#
# In practical terms, this timer is the "final call pickup" timer, e.g.
# the maximum amount of time the call is allowed to "ring" (provided
# restart_fr_on_each_reply is not enabled) before being considered
# failed.

modparam("tm", "fr_inv_timer", 90000)	# 1.5 minutes.

# Initial retransmission period (in milliseconds).  This is equivalent to 
# timer T1 in RFC 3261, and is normally set at 500 milliseconds.  It can, 
# for some purposes, be reasonable to make this more aggressive.

modparam("tm", "retr_timer1", 500)

# Maximum retransmission period.  The retransmission interval starts at
# timer T1 (retr_timer1 above) and backs off until it reaches this value 
# (timer T2 in RFC 3261), after which it stays constant until maximum
# transaction lifetime.  It can, for some purposes, be reasonable to make
# this more aggressive.

modparam("tm", "retr_timer2", 4000)

# If set, timed-out INVITE transactions will be replied to the sender
# with a locally generated 408 time out.  If not set, this branch will
# be silently dropped with no feedback.

modparam("tm", "noisy_ctimer", 1)

# Restart the fr_inv_timer every time a provisional reply is received.
# Sometimes multiple identical or lower provisional replies are 
# received over the course of call setup (e.g. 180, then 183, then 180
# again, etc.).  
#
# If this is not set, the fr_inv_timer count will be restarted only 
# on the initial provisional reply or sequentially increasing provisional
# replies (e.g. 100, then 180, then 183).  
#
# Disabling this can be useful with certain UAs that continuously 
# retransmit provisional replies without providing a timely conclusion 
# to the call that makes impossible the implementation of certain 
# failover techniques, such as automatic rollover to another destination
# after X seconds.

modparam("tm", "restart_fr_on_each_reply", 0)

# Automatically send backward 100 Trying provisional messages when
# processing INVITEs.  
#
# Note that 100 replies are hop-by-hop;  if this option is turned off,
# all 100 replies must be sent manually (using t_reply(), say).  
# However, this behaviour can also be toggled manually on a per-transaction
# basis using t_set_auto_inv_100().

modparam("tm", "auto_inv_100", 1)

# The reason string in 100 replies, e.g. "Trying".

modparam("tm", "auto_inv_100_reason", "Trying")

# If this option is enabled, and the final reply to more than one parallel
# branch is a 401 or 407 challenge, then all the WWW-Authenticate and 
# Proxy-Authenticate headers from all the 401 and 407 replies will be 
# aggregated into a new final reply.  This option does nothing if there
# is only one branch on which such a reply arrives.
#
# If disabled, only the first 401/407 received from any branch will be 
# forwarded.
# 
# Note that this option should be enabled to be compliant with RFC 3261
# proxy behaviour, but it may be useful to turn it off in certain situations.

modparam("tm", "aggregate_challenges", 1)

# Generate CANCEL and negative ACK requests from the INVITE message that
# was sent out instead of building them from the received request.  
# This requires that the outgoing INVITE be partially re-parsed.  However,
# this behaviour is required for RFC 3261 compliance.  

modparam("tm", "reparse_invite", 1)

# Header fields prefixed by this parameter are included in the CANCEL
# and negative ACK messages if they were present in the outgoing INVITE.
# This only affects those headers that are not covered by RFC 3261, 
# which would be neither mandated nor prohibited.  This parameter
# requires reparse_invite in order to work.
#
# Under most circumstances, this value should be blank.

modparam("tm", "ac_extra_hdrs", "")

# If set and blacklisting is enabled, every source of 503 replies is 
# automatically added to the blacklist.  Note that this could be very
# easly exploited in DoS attacks where SIP is open to the public.
#
# This option should be disabled under all imaginable circumstances.

modparam("tm", "blst_503", 0)

# Bitmap of methods that trigger blacklisting on transaction timeouts.
# This option has no effect if blacklists are disabled.
#
# Bit values for various methods:
#
#   INVITE    : 1
#   CANCEL    : 2
#   ACK       : 4
#   BYE       : 8
#   INFO      : 16
#   REGISTER  : 32
#   SUBSCRIBE : 64
#   NOTIFY    : 128
#   OTHER     : 256 
#   (see parser/msg_parser.h for further details)

modparam("tm", "blst_methods_add", 0)

# Bitmap of methods (from table in description for 'blst_methods_add') 
# that are looked up in the blacklist before stateful forwarding.
# 
# This option has no effect if blacklisting is disabled.

modparam("tm", "blst_methods_lookup", 0)

# Behaviour when CANCELing a transaction for which no >= 100 reply
# was received.
#
# Option 0 - Immmediately stop request (INVITE) retransmission on the
# branch as if the branch received an immediate 487 reply (fake internal
# 487 reply).  The benefit of this is that unreplied branches will be
# terminated immediately. It also, however, introduces the risk of a 
# race condition due to a belated 2xx reply, perhaps due to network 
# delay or far-end UA congestion/load.  
#
# Option 1 - Continue retransmitting the request on unreplied branches.
# If a provisional reply is later received, a CANCEL will be immediately
# sent in an attempt to trigger an immediate 487.  This approach avoids
# the race condition liability above, but it's more resource-intensive 
# because the CANCEL attempt will keep the transaction alive for the
# whole fr_timer interval.  
# 
# Option 2 - Send and retransmit CANCEL even if the branch did not 
# receive a reply.  This has the same advantages are option 1 and 
# also avoids an extra round-trip in case of a provisional reply, but
# is not RFC 3261 compliant because 3261 only allows sending CANCELs
# on pending branches (e.g. those on which a provisional reply of some 
# description was positively received).
#
# The most reasonable approach is Option 1, unless resources are
# severely constrained due to RAM limitations and extremely high 
# call volume with a high probability of unreplied requests/branches.

modparam("tm", "cancel_b_method", 1)

# If enabled, this option causes the SIP message after a DNS failover
# to be constructed from the outgoing message buffer of the failed
# branch instead of from the received request.
# 
# It must be enabled if multiple branches are used, the SIP message 
# is different in them (i.e. modified somehow in the serial forking/rollover
# process), and at least one of those branches can result in DNS failover.
#
# This option breaks detection of outgoing sockets, though, so it is
# dangerous on multi-homed hosts where a SIP request on a new branch
# may need to exit a different interface after DNS failover.

modparam("tm", "reparse_on_dns_failover", 0)

# The reply route block that is invoked when a reply is received
# that has no associated transaction tracked statefully.  The reply
# is passed to the core for stateless forwarding after the execution
# of this route block unless that block returns 0, which would
# disallow stateless forwarding of replies that are not associated
# with a transaction.

modparam("tm", "on_sl_reply", "default_stateless_reply")

# This option selects how to treat CANCELs that do not match any 
# statefully tracked transaction.  
#
# Option 0 - Forward them statefully.
# 
# Option 1 - Forward them statelessly.
#
# Option 2 - Drop them.
#
# While the stateful approach has the hidden advantage of detecting INVITEs
# that are received after the CANCEL (due to network conditions/packet loss),
# it exposes the proxy to a DoS attack if it authenticates all request types,
# including CANCEL, by flooding with CANCELs to random destinations;  since
# the CANCEL is not successfully authenticated, each bogus CANCEL will create
# a new transaction that will live in memory for 30 seconds by default).
#
# The most reasonable behaviour is option 1.

modparam("tm", "unmatched_cancel", 1)

# Try match the Request URI when doing pre-RFC 3261 transaction matching
# (where the Via branch parameter does not contain the cookie indicating use 
# of a 3261 SIP stack).  The only reason to have this option on is for 
# interoperability with old implementations, otherwise it should absolutely
# be disabled.

modparam("tm", "ruri_matching", 0)

# If set, will try to match the topmost Via when doing pre-RFC3261 
# transaction matching (where the Via branch parameter does not contain 
# the cookie indicating use of a 3261 SIP stack).  The only reason to have
# this option on is for interoperability with old implementations, otherwise
# it should absolutely be disabled.

modparam("tm", "via1_matching", 0)

# If set, callbacks registered with TM for the event type 
# TMCB_LOCAL_RESPONSE_OUT will also be called for provisional replies.
# 
# Unless there is a specific reason for having this on (e.g. custom 
# modules that require it), it should be off.

modparam("tm", "pass_provisional_replies", 0)

# Default response code if t_reply() cannot retrieve its parameters
# (e.g. they contain a nonexistent/undefined AVP).

modparam("tm", "default_code", 500)

# Default reason string if t_reply() cannot retrieve its parameters
# (e.g. they contain a nonexistent/undefined AVP).

modparam("tm", "default_reason", "Server Internal Error")

# If set, TM will treat all 6xx replies like normal failure replies and not
# prevent serial forking from taking place on them.  Otherwise, the
# normal response to a 6xx reply is to cancel all parallel branches, 
# stop DNS failover, and any other forking.  However, serial forking
# invoked manually in a failure_route via append_branch() would still work.
# 
# Note that this behaviour is not RFC 3261 compliant;  per 3261, 6xx-class
# replies indicate that the sender of the reply has authoritative knowledge
# that the destination cannot be reached by any other means or any other
# endpoints.  However, lots of implementors misuse 6xx (especially the
# 603 Declined reply) to signal some sort of ambiguous failure that cannot
# be otherwise meaningfully categorised and does not necessarily indicate
# that.  Consequently, it may be necessary to enable this option in certain
# situations.
#
# This behaviour can be overridden on a per-transaction basis with 
# t_set_disable_6xx().

modparam("tm", "disable_6xx_block", 0) 	# We choose RFC compliance

# Treatment of locally generated ACKs for 2xx replies to local 
# transactions -- i.e. transactions created via t_uac*() or through
# the TM API or via RPC/MI).
#
# Option 0 - The ACK destination is chosen according to the RFC 
# behaviour;  the next hop is found using the Contact and the Route 
# set, and then DNS resolution is used on it, if applicable.
#
# Option 1- ACK is sent to the same address as the corresponding INVITE
# branch.
#
# Option 2 - ACK is sent to the source of the 2xx reply.
#
# The RFC 3261-prescribed behaviour is option 0.

modparam("tm", "local_ack_mode", 0)

