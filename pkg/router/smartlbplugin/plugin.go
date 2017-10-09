package smartlbplugin

import (
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apimachinery/pkg/util/sets"

	kapi "k8s.io/kubernetes/pkg/api"
	routeapi "github.com/openshift/origin/pkg/route/apis/route"
	"github.com/golang/glog"
	"fmt"
)

type SmartLBPlugin struct {
}

func NewSmartLBPlugin(apiUrls string) (*SmartLBPlugin, error) {
	glog.Infof("Starting new plugin for %v", apiUrls)

	return &SmartLBPlugin{}, nil
}

func (p *SmartLBPlugin) HandleRoute(eventType watch.EventType, route *routeapi.Route) error {
	glog.V(4).Infof("Processing route for service: %v (%v)", route.Spec.To.Name, route.Spec.Host)

	routeName := fmt.Sprintf("host: %v, port: %v, path: %v", route.Spec.Host, route.Spec.Port, route.Spec.Path)

	switch eventType {
	case watch.Modified:
		glog.V(4).Infof("Updating route %s", routeName)
	case watch.Deleted:
		glog.V(4).Infof("Deleting route %s", routeName)
	case watch.Added:
		glog.V(4).Infof("Adding new route %s", routeName)
	}

	return nil
}

// Endpoints are not relevant for this plugin
func (p *SmartLBPlugin) HandleEndpoints(eventType watch.EventType, endpoints *kapi.Endpoints) error {
	return nil
}

// Namespaces are not relevant for this plugin
func (p *SmartLBPlugin) HandleNamespaces(namespaces sets.String) error {
	return fmt.Errorf("namespace handling is not implemented for this plugin")
}

// Nodes are not relevant for this plugin
func (p *SmartLBPlugin) HandleNode(eventType watch.EventType, node *kapi.Node) error {
	return fmt.Errorf("node handling is not implemented for this plugin")
}

// No-op
func (p *SmartLBPlugin) Commit() error {
	return nil
}
