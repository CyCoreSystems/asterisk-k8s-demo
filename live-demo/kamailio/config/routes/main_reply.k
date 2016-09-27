##
# General-purpose reply route for logging/diagnostic purposes, as well
# as a landing spot for hooks of features that require modifying or
# dropping replies, such as far-end NAT traversal.

onreply_route[MAIN_REPLY] {
	xlog("L_INFO", "[R-MAIN-REPLY:$ci] Reply '$rs $rr' from $si:$sp\n");

   # a=maxptime breaks current Innomedia firmware (4.0.4.0)
   # so we remove it from the reply's SDP
   if(sdp_content()) {
      xlog("L_DBG","[R-MAIN-REPLY:$ci] Removing maxptime from SDP\n");
      sdp_remove_line_by_prefix("a=maxptime");
   }

   # Fix the Contact, if we are talking to a NATed endpoint
   if($avp(nated) == 1) {
      fix_nated_contact();
   }
}
