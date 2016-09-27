package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"k8s.io/client-go/1.4/pkg/watch"

	"github.com/CyCoreSystems/ari/client/nc"
	"github.com/CyCoreSystems/dispatchers/pkg/endpoints"
	"github.com/pkg/errors"
)

var outputFilename string
var proxies []string
var namespace string
var endpointName string

func init() {
	proxies = []string{}

	flag.StringVar(&outputFilename, "o", "/etc/asterisk/pjsip.d/proxies.conf", "Output file for proxies config")
	flag.StringVar(&endpointName, "name", "kamailio", "Name of the Service to source proxy list")
	flag.StringVar(&namespace, "ns", "default", "Namespace of the Service to source proxy list")
}

func main() {
	var failureCount int
	flag.Parse()

	_, err := update()
	if err != nil {
		fmt.Println("Failed to update proxies list", err)
	}
	err = export()
	if err != nil {
		panic("Failed to export proxies.conf: " + err.Error())
	}
	err = notify()
	if err != nil {
		panic("Failed to notify asterisk of update: " + err.Error())
	}

	for failureCount < 10 {
		err := maintain()
		if err != nil {
			fmt.Println("Error: ", err)
			failureCount++
		}
	}
}

func maintain() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	changes := make(chan error, 10)

	// Listen to each of the namespaces
	go watchNamespace(ctx, changes, namespace)

	for {
		err := <-changes
		if err != nil {
			return errors.Wrap(err, "error maintaining watch")
		}

		// Update and check if the changes were significant
		changed, err := update()
		if err != nil {
			return errors.Wrap(err, "error updating proxies list")
		}

		if changed {
			err = export()
			if err != nil {
				return errors.Wrap(err, "failed to export proxies.conf")
			}

			err = notify()
			if err != nil {
				return errors.Wrap(err, "failed to notify asterisk of update")
			}
		}
	}
}

func watchNamespace(ctx context.Context, changes chan error, ns string) {
	w, err := endpoints.Watch(ns)
	if err != nil {
		changes <- errors.Wrap(err, "failed to watch namespace")
		return
	}
	defer w.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-w.ResultChan():
			if ev.Type == watch.Error {
				changes <- errors.Wrap(err, "watch error")
				return
			}
			changes <- nil
		}
	}
}

// export pushes the current dispatcher sets to file
func export() error {
	f, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	return proxiesTemplate.Execute(f, proxies)
}

// notify signals to kamailio to reload its dispatcher list
func notify() error {
	n, err := nc.New(nc.Options{
		URL: "nats://nats:4222",
	})
	if err != nil {
		return errors.Wrap(err, "failed to connect to NATS")
	}

	return errors.Wrap(n.Asterisk.ReloadModule("res_pjsip.so"), "failed to reload PJSIP")
}

// update updates the list of proxies
func update() (changed bool, err error) {
	list, err := endpoints.Get(namespace, endpointName)
	if err != nil {
		return
	}

	if differ(proxies, list) {
		changed = true
	}
	proxies = list

	return
}
