## 
## Module: db_postgres
##

# The number of times to retry connection to the Postgres database server
# in the event that it is lost.  

modparam("db_postgres", "retries", 3)
