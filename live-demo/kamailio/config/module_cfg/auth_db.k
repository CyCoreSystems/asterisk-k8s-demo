##
## Module: auth_db
##

#!ifdef WITH_DIGEST_AUTH
#!ifdef WITH_DB

# Database URL, if enabled.

#modparam("auth_db", "db_url", "DB_DRIVER://DB_USER:DB_PASSWORD@DB_HOST/DB_NAME")
modparam("auth_db", "db_url", "DB_DRIVER://DB_NA_USER:DB_NA_PASSWORD@DB_NA_HOST/DB_NA_NAME")

# The table column containing the username.

#modparam("auth_db", "user_column", "username")
modparam("auth_db", "user_column", "name")

# The table column containing the domain (if domain support is enabled).

#modparam("auth_db", "domain_column", "domain")

# The table column containing the password (plain text).

#modparam("auth_db", "password_column", "password")
modparam("auth_db", "password_column", "secret")

# Indicates whether the module should expect passwords to be stored in 
# plain text (option 0, stored in password_column) or pre-calculated HA1
# digests (option 1, stored in password_column_2, not set here).  If the
# parameter is set to 1, passwords will be calculated from the
# 'password_column'.  This should be set to 1.

modparam("auth_db", "calculate_ha1", 1)

# Use the domain attribute when looking up users in the table, and also
# in authentication challenges.

#modparam("auth_db", "use_domain", 1)
modparam("auth_db", "use_domain", 0)

# List of additional columns to fetch when performing an authorisation 
# and the AVPs in which to store their values, in this format:
# 
#    $pvar=credential
#
# This can take the form of a colon-separated list:
#
#    $pvar1=credential1;$pvar2=credential2
#
# If an AVP name is not specified, e.g.
# 
#    credential;$pvar1=credential1
#
# then the credential will be stored in $avp(s:credential) by default.

modparam("auth_db", "load_credentials", "$avp(s:auth_user_id)=id")

#!endif 
#!endif
