#
# This is the OpenShift Origin smart loadbalancer plugin
#
# The standard name for this image is openshift/origin-smart-lb-plugin
#
FROM openshift/origin

LABEL io.k8s.display-name="OpenShift Origin Smart LB Plugin" \
      io.k8s.description="This is a component of OpenShift Origin and contains an plugin for a smart load balancer outside multiple OpenShift clusters. This plugin tells the smart load balancer about the routes of this cluster" \
      io.openshift.tags="openshift,smart-lb-plugin"

RUN ln -sf /usr/bin/openshift /usr/bin/openshift-smart-lb-plugin

USER 1002
ENTRYPOINT ["/usr/bin/openshift-smart-lb-plugin"]
