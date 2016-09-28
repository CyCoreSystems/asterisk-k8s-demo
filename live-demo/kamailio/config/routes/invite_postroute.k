route[INVITE_POSTROUTE] {

   rtpproxy_manage();

	t_on_reply("MAIN_REPLY");

	if(!t_relay()) 
		sl_reply_error();

}
