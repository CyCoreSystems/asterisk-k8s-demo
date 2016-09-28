route {
	if(!mf_process_maxfwd_header("10")) {
		sl_send_reply("483", "Too Many Hops (10)");
		exit;
	}

	if(!sanity_check()) {
		xlog("L_ERR", "[R-MAIN:$ci] !> "
			      "Malformed request from $si:$sp\n");

		exit;
	}


	### -- Sequential request and ACK handling -- ###

	if(has_totag()) {
	
		if(loose_route()) {
			if(is_method("BYE")) {
            rtpproxy_manage();
				xlog("L_INFO", 
				     "[R-MAIN:$ci] -> "
				     "BYE (hangup) received from $si:$sp\n");

			}
	
			if(!is_method("INVITE")) {
				xlog("L_INFO", 
				     "[R-MAIN:$ci] -> "
				     "Other sequential $rm received from "
				     "$si:$sp\n");
			}	


			if(!t_relay()) {
            rtpproxy_manage();
				sl_reply_error();
         }
	
			exit;
		} 

		else {
			# Request is not loose-routed, but has a To-tag;
			if(is_method("ACK")) {
            rtpproxy_manage();
				if(t_check_trans()) 
					t_relay();

				exit;
			}

			sl_send_reply("403", "Forbidden");
		}

		exit;
	}


	### -- Stateful CANCEL handling -- ###

	if(is_method("CANCEL")) {

      rtpproxy_manage();

		xlog("L_INFO", 
			"[R-MAIN:$ci] CANCEL received from $si (RURI=$ru)\n");

		if(!t_relay_cancel()) {

			xlog("L_ERR", "[R-MAIN:$ci] !> "
				      "Corresponding INVITE transaction for "
				      "CANCEL found, but error occurred\n");

			sl_send_reply("500", "Internal Server Error");
			exit;
		}
	
		xlog("L_ERR", "[R-MAIN:$ci] !> "
			      "Corresponding INVITE transaction for CANCEL "
			      "does not (or, no longer) exists\n");

		exit;
		
	}

	## -- All other (initial) request handling -- ##

	t_check_trans();

	if(uri == myself) {

		if(is_method("OPTIONS")) {
			options_reply();
			exit;
		}

		else if(is_method("INVITE")) {
			xlog("L_INFO", "[R-MAIN:$ci] "
				       "Initial INVITE received from $si "
				       "(RURI=$ru/From=$fu)\n");

			record_route();

			route(INVITE_RESOLVER);

			exit;
		} 

		else {

			sl_send_reply("403", "Forbidden");
			exit;
		}
	} 

	else {

		xlog("L_INFO", "[R-MAIN:$ci] !> Rejecting non-local request "
			       "$rm with RURI $ru\n");

		sl_send_reply("403", "Forbidden");
		exit;
	}
}


event_route[tm:local-request] {
	xlog("L_INFO", "[R-TM-LOCAL-REQUEST:$ci] Local request $rm to $ru\n");
}
