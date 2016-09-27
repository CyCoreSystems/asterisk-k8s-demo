##
## Module: ctl 
##

# The transport used for the binary RPC protocol.  This is specified in 
# the form:  [protocol:]address_port|path
# 
# The following transport protocols are supported:
# 
#   - unixd: UNIX datagram socket.
#   - unix[s]: UNIX stream sockets.
#   - tcp - TCP network socket (tcp:listen_addr:port)
#   - udp - UDP network socket (udp:listen_addr:port)
# 
# Multiple transport interfaces can be specified by setting this parameter
# multiple times.

modparam("ctl", "binrpc", "unix:/var/run/kamailio/kamailio_ctl")
modparam("ctl", "binrpc", "udp:127.0.0.1:9998")

# Permissions 'mode' used on the UNIX domain socket or FIFO.
# Default: 0600.

modparam("ctl", "mode", 0600)

# The EUID under which the UNIX domain socket or FIFO will be created.
# Default: not set (the EUID of the Kamailio process will be used).

#modparam("ctl", "user", "root")

# The EGID under which the UNIX domain socket or FIFO will be created.
# Default: not set (the EGID of the Kamailio process will be used).

#modparam("ctl", "group", "root")

# FIFO used for the obsolete FIFO protocol.  This should be disabled.

#modparam("ctl", "fifo", "fifo:/tmp/kamailio_fifo")

# Enable or disable automatic type conversion for the method parameters.
# This should be disabled.

modparam("ctl", "autoconversion", 0)

# The size of the buffer for the RPC reply, in kilobytes.  
# Default: 4 (kB)

modparam("ctl", "binrpc_max_body_size", 4)

# The size of the structure buffer for the RPC reply, in kilobytes.
# Default: 1 (kB)

modparam("ctl", "binrpc_struct_max_body_size", 2)
