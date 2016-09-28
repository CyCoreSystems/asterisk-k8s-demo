## 
## Module: dispatcher
##

# Force overwriting of destination set even when it is already populated.

modparam("dispatcher", "force_dst", 1)

# Enable failover/rollover on gateways that don't respnd.

modparam("dispatcher", "flags", 2)

# Set AVPs which are used internally (as arrays) by dispatcher to
# implement serial-forking based failover.

modparam("dispatcher", "dst_avp", "$avp(s:d_dst_set)")
modparam("dispatcher", "grp_avp", "$avp(s:d_grp_set)")
modparam("dispatcher", "cnt_avp", "$avp(s:d_cnt_avp)")
modparam("dispatcher", "sock_avp", "$avp(dssocket)")

# Probe dispatchers
# Mode 0: only probe gateways in the probing state
# Mode 1: probe all gateways
# Mode 2: only proble gateways which are both inactive AND probing

#modparam("dispatcher", "ds_probing_mode", 0)
#modparam("dispatcher", "ds_probing_threshhold", 5)
#modparam("dispatcher", "ds_ping_interval", 5)

# Where to find the list of dispatchers

modparam("dispatcher", "list_file", "/data/kamailio/dispatcher.list")
