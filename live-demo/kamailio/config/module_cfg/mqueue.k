##
## Module: mqueue
##

# Messaging queue definitions

# Queue for notifications to NextAlarm about successful registrations
modparam("mqueue", "mqueue", "name=nextalarm_registration;size=50")

