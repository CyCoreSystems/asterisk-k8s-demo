#!ifdef WITH_DB

##
## Module: uri_db
##

# Database URL, if enabled.

modparam("uri_db", "db_url", "DB_DRIVER://DB_USER:DB_PASSWORD@DB_HOST/DB_NAME")

# Use the "uri" table instead of the "subscriber" table for checking
# values in check_to()/check_from().

modparam("uri_db", "use_uri_table", 0)

# Table against which checks should be done;  this should be set to 
# 'subscriber' and will only work if 'use_uri_table' is set to 0.

modparam("uri_db", "db_table", "subscriber")

#!endif
