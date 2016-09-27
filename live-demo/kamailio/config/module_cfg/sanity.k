## 
## Module: sanity
##

# Bitmap of checks to use if sanity check is executed with no parameter.
# By default, all implemented checks are included in the sanity_check()
# function.  
#
# Bitmap values:
#
#   SIP version in Request URI         : 1 
#   Request URI scheme                 : 2
#   Presence of RFC-mandated headers   : 4
#   CSeq method consistency            : 32
#   CSeq value consistency             : 64
#   Content length actual vs. declared : 128
#   Expires value                      : 256
#   Supported extensions               : 512

modparam("sanity", "default_checks", 103) 	# 1 & 2 & 4 & 32 & 64

# Bitmap of URIs to be checked if URI check is set to be used in 
# sanity_check().

