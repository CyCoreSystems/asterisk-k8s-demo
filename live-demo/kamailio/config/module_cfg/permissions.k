##
## Module: permissions
##

# Database URL, if database enabled.

#!ifdef WITH_DB
#modparam("permissions", "db_url", "DB_DRIVER://DB_USER:DB_PASSWORD@DB_HOST/DB_NAME")
#!endif

# Check all branches, instead of just the RURI of the first branch.

modparam("permissions", "check_all_branches", 1)

# Whether to cache the contents of the 'trusted' table or check it
# every time.

# modparam("permissions", "db_mode", 1)
