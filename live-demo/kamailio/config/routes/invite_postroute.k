route[INVITE_POSTROUTE] {

   rtpengine_manage();

	t_on_reply("MAIN_REPLY");

	if(!t_relay()) 
		sl_reply_error();

}
