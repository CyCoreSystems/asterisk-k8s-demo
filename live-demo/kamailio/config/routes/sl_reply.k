# Default policy is to allow stateless replies, because ideally restarting
# the proxy won't break call flow.  However, we log them in order to detect
# possible security anomalies or other problems.

onreply_route["default_stateless_reply"] {
	xlog("L_ERR", "[R-DEFAULT-STATELESS-REPLY:$ci] !> "
		      "Received stateless reply $rs ($rm) from $si:$sp\n");

	# Pass up to core for stateless forwarding.
}
