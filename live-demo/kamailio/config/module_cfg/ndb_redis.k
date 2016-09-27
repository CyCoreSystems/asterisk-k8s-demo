##
## Module: ndb_redis
##

# Define database servers to be available
modparam("ndb_redis", "server", "name=abn;addr=REDIS_HOST;port=6379")
