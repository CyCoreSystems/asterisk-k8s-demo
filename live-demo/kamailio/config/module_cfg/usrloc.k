##########################################################################
# Kamailio 3.x configuration basis - usrloc module parameters.           #
# 									 #
# For use by Evariste Systems LLC and authorised third parties only.     #
#									 #
# Boilerplate parameters, code and comment annotations by:               #
#									 #
# Alex Balashov <abalashov@evaristesys.com>                              #
##########################################################################

#!ifdef WITH_REGISTRAR

##
## Module: usrloc
##

# The branch flag to be used as an indicator of NAT (if the Contact URI
# of the binding is NAT'd).  It will be imported and used by all other
# modules that draw on usrloc.

modparam("usrloc", "nat_bflag", 1)

# Database-related options, if WITH_DB is defined.

#!ifdef WITH_DB

# Table column containing usernames.

modparam("usrloc", "user_column", "username")

# Table column containing domains.

modparam("usrloc", "domain_column", "domain")

# The database URL, if enabled.

modparam("usrloc", "db_url", "DB_DRIVER://DB_USER:DB_PASSWORD@DB_HOST/DB_NAME")

# Database backing mode to use.  Options are:
#
#  0 - Disable database entirely, use in-memory storage only.
#  1 - Write-Through: All changes in memory are immediately 
#      reflected in database table.
#  2 - Write-Back: Changes made in memory and synchronised
#      to database on timer intervals.  Timer also deletes
#      expired contacts from database.
#  3 - Database only: No memory cache is kept, all operations
#      performed directly on the database.  Timer deletes expired
#      contacts from the database 

modparam("usrloc", "db_mode", 0)

#!endif

# Order multiple contact bindings (if allowed) for a given AOR by
# timestamp;  otherwise, they will be ordered based on q-value.

modparam("usrloc", "desc_time_order", 0)

# Interval of timer used to delete expired contacts, synchronise with
# the database, and other periodic tasks.

modparam("usrloc", "timer_interval", 60)

# Match AOR lookups on Contact only (1) or (2) on Contact & also
# Call-ID (non-3261 compliant).

modparam("usrloc", "matching_mode", 0)

# Maximum retransmission dampening time (in seconds).  Retransmissions
# within this interval that have the same Call-ID and CSeq will be accepted
# and replied to in the same manner as the original request was replied to,
# though no update will be done to the bindings.  If this dampening time
# lapses, an error is reported.

modparam("usrloc", "cseq_delay", 10)

# The number of rows read from the database in bulk when loading the
# location table records at startup.

modparam("usrloc", "fetch_rows", 1000)

# The number of buckets in the hash table used to store location records.
# The bucket size of this table is 2^(hash_size).  

modparam("usrloc", "hash_size", 12)

#!endif
