/*******************************************************************************
*  MIT License
*
*  Copyright (c) 2024 Ji-Young Park(jiyoung.park.dev@gmail.com)
*
*  Permission is hereby granted, free of charge, to any person obtaining a copy
*  of this software and associated documentation files (the "Software"), to deal
*  in the Software without restriction, including without limitation the rights
*  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
*  copies of the Software, and to permit persons to whom the Software is
*  furnished to do so, subject to the following conditions:
*
*      The above copyright notice and this permission notice shall be included in all
*      copies or substantial portions of the Software.
*
*      THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
*      IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
*      FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
*      AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
*      LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
*      OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
*      SOFTWARE.
*******************************************************************************/

package source

import (
	"context"
	"crypto/md5"
	"fmt"
	"sync"
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	dynForwardProxyCluster "github.com/envoyproxy/go-control-plane/envoy/extensions/clusters/dynamic_forward_proxy/v3"
	dynForwardProxy "github.com/envoyproxy/go-control-plane/envoy/extensions/common/dynamic_forward_proxy/v3"
	dynForwardProxyFilter "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/dynamic_forward_proxy/v3"
	router "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"

	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	upstream "github.com/envoyproxy/go-control-plane/envoy/extensions/upstreams/http/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jparklab/xds-test-server/pkg/generated/config"
	"github.com/jparklab/xds-test-server/pkg/utils/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	ListenerPort                   = 10000
	UpstreamHost                   = "localhost"
	UpstreamPort                   = 9901
	DynamicForwardProxyClusterName = "dynamic_forward_proxy"
)

type Options struct {
	WindowSize             uint32
	RDSInitialFetchTimeout *durationpb.Duration
	UseDynamicForwarding   bool
}

type ManualSource struct {
	rwLock sync.RWMutex

	config.UnsafeConfigServiceServer

	hasChange    bool
	serviceNames map[string]struct{}

	options Options
	index   int64
}

var _ Source = (*ManualSource)(nil)

func NewManualSource() *ManualSource {
	s := &ManualSource{
		hasChange:    false,
		serviceNames: make(map[string]struct{}),
	}

	/*
		// NOTE(jpark): this is to force making updates in routes for testing
		go func() {
			for {
				time.Sleep(1 * time.Second)
				s.rwLock.Lock()
				s.index++
				s.options.WindowSize = uint32(s.index + 128*1024)
				s.hasChange = true
				s.rwLock.Unlock()

			}
		}()
	*/

	return s
}

func (s *ManualSource) Name() string {
	return "manual"
}
func (s *ManualSource) HasChange() bool {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	return s.hasChange
}

func (s *ManualSource) Resources() InternalResources {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	s.hasChange = false
	if len(s.serviceNames) == 0 {
		return InternalResources{}
	}

	listenerName := "test-listener"
	routeName := fmt.Sprintf("test-route-%d", s.index/2)
	// routeName := fmt.Sprintf("test-route")

	clusters := []*cluster.Cluster{}
	endpoints := []*endpoint.ClusterLoadAssignment{}

	if s.options.UseDynamicForwarding {
		clusters = append(clusters, s.makeDynForwardCluster(DynamicForwardProxyClusterName))
	} else {
		for name := range s.serviceNames {
			cl, ep := s.makeCluster(name)
			clusters = append(clusters, cl)
			endpoints = append(endpoints, ep)
		}
	}

	// NOTE(jpark): we add a route for the first service only since that's all we need for testing.
	listeners := []*listener.Listener{
		s.makeListener(listenerName, routeName, s.options.UseDynamicForwarding),
	}

	routes := []*route.RouteConfiguration{}
	if s.options.UseDynamicForwarding {
		routes = append(routes,
			makeRoute(routeName, DynamicForwardProxyClusterName, s.options.UseDynamicForwarding),
		)
	} else {
		routes = append(routes,
			makeRoute(routeName, clusters[0].Name, s.options.UseDynamicForwarding),
		)
	}

	return InternalResources{
		Listeners: listeners,
		Routes:    routes,
		Clusters:  clusters,
		Endpoints: endpoints,
	}
}

func (s *ManualSource) AddService(name string) {
	log.Infof("Add a new service: %s", name)
	s.serviceNames[name] = struct{}{}
	s.hasChange = true
}

func (s *ManualSource) SetOption(opt *config.SetOptionCommand) {
	switch opt.Key {
	case "window_size":
		log.Infof("set window size: %d", opt.GetInt64Value())
		s.options.WindowSize = uint32(opt.GetInt64Value())
	case "use_dynamic":
		log.Infof("use dynamic proxy: %t", opt.GetBoolValue())
		s.options.UseDynamicForwarding = opt.GetBoolValue()
	case "rds_initial_fetch_timeout":
		log.Infof("set rds initial fetch timeout: %d", opt.GetInt64Value())
		s.options.RDSInitialFetchTimeout = &durationpb.Duration{
			Seconds: opt.GetInt64Value(),
		}
	default:
		log.Errorf("Invalid option key: %s", opt.Key)
		return
	}

	if len(s.serviceNames) > 0 {
		s.hasChange = true
	}

}

func (s *ManualSource) SendConfigCommand(ctx context.Context, cmd *config.Command) (*empty.Empty, error) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	switch cmd.Type {
	case config.CommandType_COMMAND_ADD_SERVICE:
		s.AddService(cmd.GetAddService().ServiceName)
	case config.CommandType_COMMAND_REMOVE_SERVICE:
		log.Fatalf("not implemented")
	case config.CommandType_COMMAND_SET_OPTION:
		s.SetOption(cmd.GetSetOption())
	}

	return &empty.Empty{}, nil
}

func (s *ManualSource) makeListener(listenerName string, routeName string, useDynamicForwarding bool) *listener.Listener {
	httpFilters := []*hcm.HttpFilter{}

	if useDynamicForwarding {
		forwardProxyConfig, _ := anypb.New(&dynForwardProxyFilter.FilterConfig{
			ImplementationSpecifier: &dynForwardProxyFilter.FilterConfig_DnsCacheConfig{
				DnsCacheConfig: &dynForwardProxy.DnsCacheConfig{
					Name: "dns_cache",
				},
			},
		})
		httpFilters = append(httpFilters, &hcm.HttpFilter{
			Name: "envoy.filters.http.dynamic_forward_proxy",
			ConfigType: &hcm.HttpFilter_TypedConfig{
				TypedConfig: forwardProxyConfig,
			},
		})
	}

	routerConfig, _ := anypb.New(&router.Router{})
	httpFilters = append(httpFilters, &hcm.HttpFilter{
		Name:       wellknown.Router,
		ConfigType: &hcm.HttpFilter_TypedConfig{TypedConfig: routerConfig},
	})

	// HTTP filter configuration
	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource: &core.ConfigSource{
					ResourceApiVersion:    core.ApiVersion_V3,
					ConfigSourceSpecifier: &core.ConfigSource_Ads{},
					InitialFetchTimeout:   s.options.RDSInitialFetchTimeout,
				},
				RouteConfigName: routeName,
			},
		},
		HttpFilters: httpFilters,
	}

	pbst, err := anypb.New(manager)
	if err != nil {
		panic(err)
	}

	return &listener.Listener{
		Name: listenerName,
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: ListenerPort,
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{{
			Filters: []*listener.Filter{{
				Name: wellknown.HTTPConnectionManager,
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: pbst,
				},
			}},
		}},
	}
}

func makeRoute(routeName string, clusterName string, useDynamicProxy bool) *route.RouteConfiguration {
	if useDynamicProxy {
		perRouteConfig, _ := anypb.New(&dynForwardProxyFilter.PerRouteConfig{
			HostRewriteSpecifier: &dynForwardProxyFilter.PerRouteConfig_HostRewriteHeader{
				HostRewriteHeader: "x-host-rewrite",
			},
		})

		return &route.RouteConfiguration{
			Name: routeName,
			VirtualHosts: []*route.VirtualHost{{
				Name:    "local_service",
				Domains: []string{"*"},
				Routes: []*route.Route{{
					Match: &route.RouteMatch{
						PathSpecifier: &route.RouteMatch_Prefix{
							Prefix: "/",
						},
					},
					Action: &route.Route_Route{
						Route: &route.RouteAction{
							ClusterSpecifier: &route.RouteAction_Cluster{
								Cluster: clusterName,
							},
						},
					},
					TypedPerFilterConfig: map[string]*anypb.Any{
						"envoy.filters.http.dynamic_forward_proxy": perRouteConfig,
					},
				}},
			}},
		}
	}

	return &route.RouteConfiguration{
		Name: routeName,
		VirtualHosts: []*route.VirtualHost{{
			Name:    "local_service",
			Domains: []string{"*"},
			Routes: []*route.Route{{
				Match: &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Prefix{
						Prefix: "/",
					},
				},
				Action: &route.Route_Route{
					Route: &route.RouteAction{
						ClusterSpecifier: &route.RouteAction_Cluster{
							Cluster: clusterName,
						},
						HostRewriteSpecifier: &route.RouteAction_HostRewriteLiteral{
							HostRewriteLiteral: UpstreamHost,
						},
					},
				},
			}},
		}},
	}
}

func hashConfig(cluster *cluster.Cluster) string {
	h := md5.New()
	blob, err := proto.Marshal(cluster)
	if err != nil {
		panic(err)
	}

	_, err = h.Write(blob)
	if err != nil {
		panic(err)
	}

	sum := h.Sum(nil)
	return fmt.Sprintf("%x", sum)
}

func (s *ManualSource) makeCluster(clusterName string) (*cluster.Cluster, *endpoint.ClusterLoadAssignment) {
	cluster := &cluster.Cluster{
		Name:                 clusterName,
		ConnectTimeout:       durationpb.New(5 * time.Second),
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_EDS},
		EdsClusterConfig: &cluster.Cluster_EdsClusterConfig{
			ServiceName: clusterName,
			EdsConfig: &core.ConfigSource{
				ResourceApiVersion:    core.ApiVersion_V3,
				ConfigSourceSpecifier: &core.ConfigSource_Ads{},
				InitialFetchTimeout:   durationpb.New(10 * time.Second),
			},
		},
		LbPolicy:        cluster.Cluster_ROUND_ROBIN,
		DnsLookupFamily: cluster.Cluster_V4_ONLY,
		LoadAssignment: &endpoint.ClusterLoadAssignment{
			ClusterName: clusterName,
		},
	}

	endpoint := &endpoint.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpoint.LocalityLbEndpoints{{
			LbEndpoints: []*endpoint.LbEndpoint{{
				HostIdentifier: &endpoint.LbEndpoint_Endpoint{
					Endpoint: &endpoint.Endpoint{
						Address: &core.Address{
							Address: &core.Address_SocketAddress{
								SocketAddress: &core.SocketAddress{
									Protocol: core.SocketAddress_TCP,
									Address:  "127.0.0.1",
									PortSpecifier: &core.SocketAddress_PortValue{
										PortValue: UpstreamPort,
									},
								},
							},
						},
					},
				},
			}},
		}},
	}

	if s.options.WindowSize > 0 {
		log.Infof("window size: %d", s.options.WindowSize)
		protoOption, err := anypb.New(&upstream.HttpProtocolOptions{
			UpstreamProtocolOptions: &upstream.HttpProtocolOptions_ExplicitHttpConfig_{
				ExplicitHttpConfig: &upstream.HttpProtocolOptions_ExplicitHttpConfig{
					ProtocolConfig: &upstream.HttpProtocolOptions_ExplicitHttpConfig_Http2ProtocolOptions{
						Http2ProtocolOptions: &core.Http2ProtocolOptions{
							InitialStreamWindowSize: wrapperspb.UInt32(s.options.WindowSize),
						},
					},
				},
			},
		})
		if err != nil {
			panic(err)
		}

		cluster.TypedExtensionProtocolOptions = map[string]*anypb.Any{
			"envoy.extensions.upstreams.http.v3.HttpProtocolOptions": protoOption,
		}
	}

	//* NOTE(jpark): comment/uncomment this to test unique names
	hashVal := hashConfig(cluster)
	newName := fmt.Sprintf("%s-%s", clusterName, hashVal)
	//*/
	//  newName := clusterName

	cluster.Name = newName
	cluster.LoadAssignment.ClusterName = newName
	cluster.EdsClusterConfig.ServiceName = newName
	endpoint.ClusterName = newName

	return cluster, endpoint
}

func (s *ManualSource) makeDynForwardCluster(clusterName string) *cluster.Cluster {
	cacheConfig, _ := anypb.New(&dynForwardProxyCluster.ClusterConfig{
		ClusterImplementationSpecifier: &dynForwardProxyCluster.ClusterConfig_DnsCacheConfig{
			DnsCacheConfig: &dynForwardProxy.DnsCacheConfig{
				Name: "dns_cache",
			},
		},
	})
	cluster := &cluster.Cluster{
		Name:     clusterName,
		LbPolicy: cluster.Cluster_CLUSTER_PROVIDED,
		ClusterDiscoveryType: &cluster.Cluster_ClusterType{
			ClusterType: &cluster.Cluster_CustomClusterType{
				Name:        "envoy.clusters.dynamic_forward_proxy",
				TypedConfig: cacheConfig,
			},
		},
		ConnectTimeout: durationpb.New(time.Duration(s.index) * time.Second),
	}

	return cluster
}
