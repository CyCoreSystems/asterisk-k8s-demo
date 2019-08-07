# Asterisk on Kubernetes Demo

This repository contains code and markup for the deployment of a highly scalable
voice application on Kubernetes using Kamailio, Asterisk, and NATS.

Supporting tools:

  - ARI core Go library [github.com/CyCoreSystems/ari](https://github.com/CyCoreSystems/ari)
  - ARI NATS proxy [github.com/CyCoreSystems/ari-proxy](https://github.com/CyCoreSystems/ari-proxy)
  - Netdiscover cloud networking discovery tool [github.com/CyCoreSystems/netdiscover](https://github.com/CyCoreSystems/netdiscover)
  - Asterisk Config kubernetes-based Asterisk templating and update engine [github.com/CyCoreSystems/asterisk-config](https://github.com/CyCoreSystems/asterisk-config)
  - Kamailio Dispatchers kubernetes-based update tool [github.com/CyCoreSystems/dispatchers](https://github.com/CyCoreSystems/dispatchers)
  - AudioSocket [github.com/CyCoreSystems/audiosocket](https://github.com/CyCoreSystems/audiosocket)

## Getting started

There are a number of kubernetes YAML files in the [k8s](live-demo/k8s/)
directory.  Some have numerical prefixes indicating that they should be deployed
in a particular order.  For the most part, getting the demo off the ground is as
easy as installing these YAML files using the usual `kubectl apply -f
<filename.yaml>` method.  However, there are a few things which must still be
done by hand.

### Asterisk config

The required configuration for Asterisk has been stripped down a lot, but there
are still a few things which need to be set up:  ARI, dialplan, and PJSIP.
Examples are included in the [asteriskconfig](live-demo/asteriskconfig)
directory.  However, you will need to update the
[inbound.conf.tmpl](live-demo/asteriskconfig/extensions.d/inbound.conf.tmpl)
file with your own DIDs (telephone numbers).

Once configured, you will need to load this configuration in to kubernetes.

  1. create a .zip file of the contents of the `asteriskconfig` directory:
    - `cd asteriskconfig`
    - `zip -r ../asterisk-config.zip *`
  2. load that .zip file in as a Secret
    - `kubectl -n voip create secret generic asterisk-config --from-file=asterisk-config.zip`

### Kamailio nodeSelector

The default kamailio DaemonSet looks for a GKE nodepool named `kamailio`.  If
this nodepool does not exist, kamailio will not be scheduled to run anywhere.
Therefore, you should either create the nodepool or modify the kamailio
DaemonSet to look for a different `nodeSelector`.

### Google Voice API key

If you intend to use the Google Speech APIs demo, you will need your own API key
loaded.  When you create an API key on Google, you are given the option to
download it as a `.JSON` file.  Do so, then load that file in as `key.json` in a
Secret named `speech-key`.

  - `kubectl -n voip create secret generic speech-key --from-file=key.json`

### Firewall rules

Depending on the environment your kubernetes is deployed to, there are any
number of ways to configure the firewall.  Fundamentally, though, UDP ports 5060
and 10000-30000 need to flow into the nodes on which the kamailio (and rtpproxy)
Pods are running.

On GCP, this is fairly easy.  You can create a special Node Pool on which the kamailio
Pods will be scheduled which have special instance tags applied.  Then, you can
tell the GCP firewall to allow the UDP ports 5060,10000-30000 into instances
with those special tags.

The kamailio deployment currently expects a nodepool to be available and named
`kamailio` in order to schedule kamailio Pods.

