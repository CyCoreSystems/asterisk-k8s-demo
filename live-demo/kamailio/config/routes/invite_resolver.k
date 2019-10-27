
route[INVITE_RESOLVER] {
	xlog("L_DBG", "[R-INVITE-RESOLVER:$ci] Entering INVITE resolver\n");

   if(!($rU =~ "^1[0-9]{10}$" ) ) {
      xlog("L_INFO", "[R-INVITE-RESOLVER:$ci] Disallowing non-NANPA call\n");
      exit;
   }

   $var(disp_set) = 1;

   $fs = "udp:PRIVATE_IP:5080";

   xlog("L_INFO", "[R-INVITE-RESOLVER:$ci] Processing dispatcher set $var(disp_set)\n");

   if(!ds_select_domain("$var(disp_set)", "4")) {
      # This should only happen if the route set is empty.
      sl_send_reply("503", "Out of Gateways");
      xlog("L_ERR", "[R-INVITE-RESOLVER:$ci] !> "
               "No gateways available!\n");
      exit;
   }

   xlog("L_INFO", "[R-INVITE-RESOLVER:$ci] -> "
             "Selected gateway: $rd:$rp\n");

   t_on_failure("DISPATCHER_ROLLOVER");

   route(INVITE_POSTROUTE);
}

