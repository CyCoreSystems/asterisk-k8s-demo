failure_route[DISPATCHER_ROLLOVER] {
	xlog("L_INFO", "[R-DISPATCHER-ROLLOVER:$ci] "
		       "Relay to dispatcher-selected gateway failed\n");

	if(t_is_expired()) {
		xlog("L_INFO", "[R-DISPATCHER-ROLLOVER:$ci] !> "
			       "Transaction expired - aborting\n");
		exit;
	}

	if(t_is_canceled()) {
		xlog("L_INFO", "[R-DISPATCHER-ROLLOVER:$ci] !> "
			       "Transaction cancelled - aborting\n");
		exit;
	}

	if(!ds_next_domain()) {
		# This should happen when we are out of gateways/have tried
		# the last one in the route set.

		xlog("L_ERR", "[R-DISPATCHER-ROLLOVER:$ci] !> "
			      "No more gateways in route set\n");

		t_reply("503", "Out of gateways");
		exit;
	}

	xlog("L_INFO", "[R-DISPATCHER-ROLLOVER:$ci] -> "
		       "Attempting relay to new gateway: $rd:$rp\n");

	t_on_failure("DISPATCHER_ROLLOVER");

	t_relay();
}
