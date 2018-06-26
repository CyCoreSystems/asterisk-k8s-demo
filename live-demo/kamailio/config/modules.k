# The directories (colon-delimeted list) to search when loading a module.
# The ordering of this list reflects the order in which they will be 
# searched, as with a shell $PATH.

loadpath "/usr/local/lib64/kamailio/modules:/usr/local/lib64/kamailio/modules_k:/usr/lib/x86_64-linux-gnu/kamailio/modules/"

# Load the modules.

loadmodule "ctl"
loadmodule "pv"
loadmodule "tm"
loadmodule "rr"
loadmodule "sl"
loadmodule "avpops"
loadmodule "tmx"
loadmodule "kex"
loadmodule "xlog"
loadmodule "maxfwd"
loadmodule "json"
loadmodule "textops"
loadmodule "textopsx"
loadmodule "siputils"
loadmodule "sanity"
loadmodule "permissions"
loadmodule "dispatcher"
loadmodule "nathelper"
loadmodule "sdpops"

#!ifdef WITH_DIALOG_TRACKING
loadmodule "dialog"
#!endif

loadmodule "rtpproxy"
loadmodule "rtimer"
loadmodule "utils"

# -- INCLUDE - Configuration parameters for those modules for which -- #
# --           there are overrides to them.                         -- #

include_file "module_cfg/xlog.k"
include_file "module_cfg/maxfwd.k"
include_file "module_cfg/ctl.k"
include_file "module_cfg/pv.k"
include_file "module_cfg/rr.k"
include_file "module_cfg/sl.k"
include_file "module_cfg/tm.k"
include_file "module_cfg/sanity.k"
include_file "module_cfg/permissions.k"
include_file "module_cfg/nathelper.k"
include_file "module_cfg/dispatcher.k"

include_file "module_cfg/rtpproxy.k"
include_file "module_cfg/rtimer.k"
include_file "module_cfg/utils.k"
