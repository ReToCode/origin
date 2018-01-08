package smartlbplugin

import (
	"strings"
	"sync"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/watch"

	"fmt"

	"bytes"
	"encoding/json"
	"github.com/golang/glog"
	routeapi "github.com/openshift/origin/pkg/route/apis/route"
	kapi "k8s.io/kubernetes/pkg/api"
	kapi_v1 "k8s.io/kubernetes/pkg/api/v1"
	"net/http"
	"time"
)

// Should all be configurable, but this is fine for the PoC
const (
	defaultWeight = 1
	updateInterval = 5
	httpPort      = 80
	httpsPort     = 443
)

// RouterHost defines a node with a router on it
type RouterHost struct {
	Name      string `json:"name"`
	HostIP    string `json:"hostIP"`
	HTTPPort  int    `json:"httpPort"`
	HTTPSPort int    `json:"httpsPort"`
}

// Route defines a URL with a weight
type Route struct {
	URL    string `json:"url"`
	Weight int    `json:"weight"`
}

// SmartLBPlugin is the plugin
type SmartLBPlugin struct {
	apiUrls []string
	clusterKey string

	Routes      map[string]Route      `json:"routes"`
	RouterHosts map[string]RouterHost `json:"routerHosts"`
	mux         sync.Mutex
}

// NewSmartLBPlugin is the constructor
func NewSmartLBPlugin(apiUrls string, clusterKey string) *SmartLBPlugin {
	glog.Infof("Starting new plugin for %v", apiUrls)
	p := &SmartLBPlugin{
		apiUrls:     strings.Split(apiUrls, ","),
		clusterKey: clusterKey,
		Routes:      make(map[string]Route),
		RouterHosts: make(map[string]RouterHost),
	}

	updater := time.NewTicker(updateInterval * time.Second)

	go func() {
		for {
			select {
				case <-updater.C:
					p.updateSmartLoadbalancers()
			}
		}
	}()

	return p
}

// Updates the smart load balancer via its API
func (p *SmartLBPlugin) updateSmartLoadbalancers() {
	p.mux.Lock()
	data, err := json.Marshal(p)
	if err != nil {
		glog.Error("Error marshaling plugin data to JSON.", err.Error())
		return
	}
	p.mux.Unlock()

	for _, u := range p.apiUrls {
		req, err := http.NewRequest("POST", u + "/api/cluster/" + p.clusterKey, bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			glog.Errorf("Error while calling %v. Err: %v", u, err.Error())
			return
		}
		glog.V(4).Infof("Response of %v was %v", u, resp.StatusCode)
		resp.Body.Close()
	}
}

// HandlePod filters to router pods and maintains an internal state of all router hosts
func (p *SmartLBPlugin) HandlePod(eventType watch.EventType, pod *kapi_v1.Pod) error {
	if pod.Namespace == "default" && pod.Status.HostIP != "" {

		// This should be a list and configurable, but for this PoC this is good enough:
		if pod.Labels["deploymentconfig"] == "router" {
			glog.Infof("router pod changes: Type: %v, Name: %v, HostIP: %v",
				eventType, pod.Name, pod.Status.HostIP)

			p.mux.Lock()
			switch eventType {
			case watch.Modified, watch.Added:
				p.RouterHosts[pod.Name] = RouterHost{Name: pod.Name, HostIP: pod.Status.HostIP, HTTPPort: httpPort, HTTPSPort: httpsPort}
			case watch.Deleted:
				delete(p.RouterHosts, pod.Name)
			}
			p.mux.Unlock()
		}
	}

	return nil
}

// HandleRoute maintains an internal state of all routes in this Cluster
func (p *SmartLBPlugin) HandleRoute(eventType watch.EventType, route *routeapi.Route) error {
	glog.Infof("Processing route: %v", route.Spec.Host)

	p.mux.Lock()

	switch eventType {
	case watch.Modified, watch.Added:
		p.Routes[route.Spec.Host] = Route{URL: route.Spec.Host, Weight: defaultWeight}
	case watch.Deleted:
		delete(p.Routes, route.Spec.Host)
	}

	p.mux.Unlock()

	return nil
}

// HandleEndpoints Endpoints are not relevant for this plugin
func (p *SmartLBPlugin) HandleEndpoints(eventType watch.EventType, endpoints *kapi.Endpoints) error {
	return nil
}

// HandleNamespaces Namespaces are not relevant for this plugin
func (p *SmartLBPlugin) HandleNamespaces(namespaces sets.String) error {
	return fmt.Errorf("namespace handling is not implemented for this plugin")
}

// HandleNode Nodes are not relevant for this plugin
func (p *SmartLBPlugin) HandleNode(eventType watch.EventType, node *kapi.Node) error {
	return fmt.Errorf("node handling is not implemented for this plugin")
}

// Commit No-op
func (p *SmartLBPlugin) Commit() error {
	return nil
}
