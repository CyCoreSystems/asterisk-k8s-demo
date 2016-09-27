##
## Module: sqlops
##

# Generic database query module.

#!ifdef WITH_DB

# The 'sqlops' module allows one to set up any number of connection 
# "handles" which are referred to by name when passed to sql_query() 
# or sql_result_free().  They can utilise any supported database 
# in the standard "database URL" format.
#
# Examples:
# 
#   modparam("sqlops", "sqlcon", "db1=>mysql://user:password@host/dbname")
#   modparam("sqlops", "sqlcon", "our_db=>postgres://user:password@host/dbname")

modparam("sqlops", "sqlcon", 
	 "canonical=>DB_DRIVER://DB_USER:DB_PASSWORD@DB_HOST/DB_NAME")
modparam("sqlops", "sqlcon", 
	 "nextalarm=>DB_DRIVER://DB_NA_USER:DB_NA_PASSWORD@DB_NA_HOST/DB_NA_NAME")

#!endif
