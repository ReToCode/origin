package router

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/openshift/origin/pkg/cmd/util"

	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"

	"errors"
	"github.com/golang/glog"
)

type SmartLBPlugin struct {
	SmartLBApiUrl string
}

// NewCommandSmartLBPlugin provides CLI handler for the smart lb plugin.
func NewCommandSmartLBPlugin(name string) *cobra.Command {
	plugin := &SmartLBPlugin {
	}

	cmd := &cobra.Command{
		Use: name,
		Short: "Start the smart lb plugin",
		Long: "Start the plugin that synchronizes the current routes to the external smart load balancer",
		Run: func(c *cobra.Command, args []string) {
			cmdutil.CheckErr(plugin.Validate())
			cmdutil.CheckErr(plugin.Run())
		},
	}

	plugin.Bind(cmd.Flags())

	return cmd
}

func (p *SmartLBPlugin) Bind(flat *pflag.FlagSet) {
	flat.StringVar(&p.SmartLBApiUrl, "smart-lb-api-url", util.Env("SMART_LB_API_URL", ""), "Specify the URL of smart load balancer API")
}

func (p *SmartLBPlugin) Validate() error {
	if p.SmartLBApiUrl == "" {
		return errors.New("smart load balancer API must be specified")
	}

	return nil
}

func (p *SmartLBPlugin) Run() error {
	glog.Infof("Starting smart load balancer plugin for remote api: %v", p.SmartLBApiUrl)



	// Do your job now
	select {}
}