route[INVITE_POSTROUTE] {

	t_on_reply("MAIN_REPLY");

	if(!t_relay()) 
		sl_reply_error();

}
