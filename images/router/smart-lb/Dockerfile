#
# This is the OpenShift Origin smart loadbalancer plugin
#
# The standard name for this image is openshift/origin-smart-lb-plugin
#
FROM openshift/origin

#RUN INSTALL_PKGS="haproxy" && \
#    yum install -y $INSTALL_PKGS && \
#    rpm -V $INSTALL_PKGS && \
#    yum clean all && \
#    mkdir -p /var/lib/haproxy/router/{certs,cacerts} && \
#    mkdir -p /var/lib/haproxy/{conf,run,bin,log} && \
#    touch /var/lib/haproxy/conf/{{os_http_be,os_edge_http_be,os_tcp_be,os_sni_passthrough,os_reencrypt,os_route_http_expose,os_route_http_redirect,cert_config,os_wildcard_domain}.map,haproxy.config} && \
#    setcap 'cap_net_bind_service=ep' /usr/sbin/haproxy && \
#    chown -R :0 /var/lib/haproxy && \
#    chmod -R g+w /var/lib/haproxy


LABEL io.k8s.display-name="OpenShift Origin Smart LB Plugin" \
      io.k8s.description="This is a component of OpenShift Origin and contains an plugin for a smart load balancer outside multiple OpenShift clusters. This plugin tells the smart load balancer about the routes of this cluster" \
      io.openshift.tags="openshift,smart-lb-plugin"
USER 1002
ENV SMART_LB_API_URL=localhost:9999/api
ENTRYPOINT ["/usr/bin/openshift-smart-lb-plugin"]
