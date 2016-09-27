##
## Module: rtimer
##

# Route timers

# Timer to execute reporting of registration successes
modparam("rtimer", "timer", "name=nextalarm_registers;interval=1;mode=1;")
modparam("rtimer", "exec", "timer=nextalarm_registers;route=NEXTALARM_REGISTRATION_CONSUMER")
