##
# General-purpose reply route for logging/diagnostic purposes, as well
# as a landing spot for hooks of features that require modifying or
# dropping replies, such as far-end NAT traversal.

onreply_route[MAIN_REPLY] {
	xlog("L_INFO", "[R-MAIN-REPLY:$ci] Reply '$rs $rr' from $si:$sp\n");

   # Fix the Contact, if we are talking to a NATed endpoint
   if($avp(nated) == 1) {
      fix_nated_contact();
   }

   rtpproxy_manage();
}
