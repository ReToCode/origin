package router

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/openshift/origin/pkg/cmd/util"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
	routeinternalclientset "github.com/openshift/origin/pkg/route/generated/internalclientset"
	projectinternalclientset "github.com/openshift/origin/pkg/project/generated/internalclientset"
	routeapi "github.com/openshift/origin/pkg/route/apis/route"

	"errors"
	"github.com/golang/glog"
	"github.com/openshift/origin/pkg/router/smartlbplugin"
	"github.com/openshift/origin/pkg/cmd/util/clientcmd"
	"github.com/openshift/origin/pkg/router/controller"
	"fmt"
)

type SmartLBPluginOptions struct {
	Config *clientcmd.Config

	SmartLBApiUrls string
	RouterSelection
}

// NewCommandSmartLBPlugin provides CLI handler for the smart lb plugin.
func NewCommandSmartLBPlugin(name string) *cobra.Command {
	options := &SmartLBPluginOptions{
		Config: clientcmd.NewConfig(),
	}
	options.Config.FromFile = true

	cmd := &cobra.Command{
		Use:   name,
		Short: "Start the smart lb plugin",
		Long:  "Start the plugin that synchronizes the current routes to the external smart load balancer",
		Run: func(c *cobra.Command, args []string) {
			cmdutil.CheckErr(options.Validate())
			cmdutil.CheckErr(options.Run())
		},
	}

	flag := cmd.Flags()
	options.Config.Bind(flag)
	options.Bind(flag)
	options.RouterSelection.Bind(flag)

	return cmd
}

func (p *SmartLBPluginOptions) Bind(flat *pflag.FlagSet) {
	flat.StringVar(&p.SmartLBApiUrls, "smart-lb-api-urls", util.Env("SMART_LB_API_URLS", ""), "Specify the URLs of smart load balancer API")
}

func (p *SmartLBPluginOptions) Validate() error {
	if p.SmartLBApiUrls == "" {
		return errors.New("smart load balancer APIs must be specified")
	}

	return nil
}

func (p *SmartLBPluginOptions) RouteAdmitterFunc() controller.RouteAdmissionFunc {
	return func(route *routeapi.Route) error {
		if err := p.AdmissionCheck(route); err != nil {
			return err
		}

		switch route.Spec.WildcardPolicy {
		case routeapi.WildcardPolicyNone:
			return nil

		case routeapi.WildcardPolicySubdomain:
			return fmt.Errorf("Wildcard routes are not supported by this plugin")
		}

		return fmt.Errorf("unknown wildcard policy %v", route.Spec.WildcardPolicy)
	}
}

func (p *SmartLBPluginOptions) Run() error {
	glog.Infof("Starting smart load balancer plugin for remote api: %v", p.SmartLBApiUrls)

	smartLBPlugin := smartlbplugin.NewSmartLBPlugin(p.SmartLBApiUrls)

	_, kc, err := p.Config.Clients()
	if err != nil {
		return err
	}
	routeClient, err := routeinternalclientset.NewForConfig(p.Config.OpenShiftConfig())
	if err != nil {
		return err
	}
	projectClient, err := projectinternalclientset.NewForConfig(p.Config.OpenShiftConfig())
	if err != nil {
		return err
	}

	// Handle all the routes
	statusPlugin := controller.NewStatusAdmitter(smartLBPlugin, routeClient, "smart-lb-plugin", "")
	uniqueHostPlugin := controller.NewUniqueHost(statusPlugin, p.RouteSelectionFunc(), p.RouterSelection.DisableNamespaceOwnershipCheck, statusPlugin)
	plugin := controller.NewHostAdmitter(uniqueHostPlugin, p.RouteAdmitterFunc(), false, p.RouterSelection.DisableNamespaceOwnershipCheck, statusPlugin)

	factory := p.RouterSelection.NewFactory(routeClient, projectClient.Projects(), kc)
	controller := factory.Create(plugin, false, false)
	controller.Run()

	// Handle all the router pods
	kClient, err := clientset.NewForConfig(p.Config.OpenShiftConfig())
	if err != nil {
		return err
	}
	smartlbplugin.CreateAndRunRouterInformer(kClient, smartLBPlugin)

	// Do your job now
	select {}
}
