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
	ListenerPort = 10000
	UpstreamHost = "localhost"
	UpstreamPort = 50052
)

type Options struct {
	WindowSize uint32
}

type ManualSource struct {
	rwLock sync.RWMutex

	config.UnsafeConfigServiceServer

	hasChange    bool
	serviceNames map[string]struct{}

	options Options
}

var _ Source = (*ManualSource)(nil)

func NewManualSource() *ManualSource {
	return &ManualSource{
		hasChange:    false,
		serviceNames: make(map[string]struct{}),
	}
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
	routeName := "test-route"

	clusters := []*cluster.Cluster{}
	endpoints := []*endpoint.ClusterLoadAssignment{}

	for name := range s.serviceNames {
		cl, ep := s.makeCluster(name)
		clusters = append(clusters, cl)
		endpoints = append(endpoints, ep)
	}

	listeners := []*listener.Listener{
		makeListener(listenerName, routeName),
	}
	routes := []*route.RouteConfiguration{
		makeRoute(routeName, clusters[0].Name),
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
		s.hasChange = true
	default:
		log.Errorf("Invalid option key: %s", opt.Key)
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

func makeListener(listenerName string, routeName string) *listener.Listener {
	routerConfig, _ := anypb.New(&router.Router{})
	// HTTP filter configuration
	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource: &core.ConfigSource{
					ResourceApiVersion:    core.ApiVersion_V3,
					ConfigSourceSpecifier: &core.ConfigSource_Ads{},
				},
				RouteConfigName: routeName,
			},
		},
		HttpFilters: []*hcm.HttpFilter{{
			Name:       wellknown.Router,
			ConfigType: &hcm.HttpFilter_TypedConfig{TypedConfig: routerConfig},
		}},
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

func makeRoute(routeName string, clusterName string) *route.RouteConfiguration {
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
				InitialFetchTimeout:   durationpb.New(0 * time.Second),
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

	hashVal := hashConfig(cluster)
	newName := fmt.Sprintf("%s-%s", clusterName, hashVal)

	cluster.Name = newName
	cluster.LoadAssignment.ClusterName = newName
	cluster.EdsClusterConfig.ServiceName = newName
	endpoint.ClusterName = newName

	return cluster, endpoint
}
