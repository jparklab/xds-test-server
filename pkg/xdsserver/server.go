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

package xdsserver

import (
	"context"
	"fmt"
	"time"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/jparklab/xds-test-server/pkg/utils/log"
	"github.com/jparklab/xds-test-server/pkg/xdsserver/source"
	"github.com/pkg/errors"
)

const defaultNodeID = "test-id"

type XDSServer struct {
	snapshotCache cache.SnapshotCache

	nodeId string
	source source.Source

	version       int
	batchInterval time.Duration
}

type XDSServerOption func(*XDSServer)

func WithNodeId(nodeId string) XDSServerOption {
	return func(s *XDSServer) {
		s.nodeId = nodeId
	}
}
func WithBatchInterval(batchInterval time.Duration) XDSServerOption {
	return func(s *XDSServer) {
		s.batchInterval = batchInterval
	}
}

func NewXDSServer(
	snapshotCache cache.SnapshotCache,
	source source.Source,
	opts ...XDSServerOption,
) *XDSServer {
	s := &XDSServer{
		snapshotCache: snapshotCache,
		source:        source,
		version:       0,

		nodeId:        defaultNodeID,
		batchInterval: 1 * time.Second,
	}
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *XDSServer) Start(ctx context.Context, initDone chan<- struct{}) error {
	log.Infof("Starting xDS server")

	firstSnapshot := true
	ticker := time.NewTicker(s.batchInterval)

	for {
		select {
		case <-ticker.C:
			updated, err := s.refreshSnapshot(ctx)
			if err != nil {
				log.With("error", err).Error("Failed to refresh snapshot")
			}

			if firstSnapshot && updated {
				initDone <- struct{}{}
				firstSnapshot = false
			}

		case <-ctx.Done():
			log.Infof("Context has been canceled, exiting xDS server")
			return nil
		}
	}
}

// refreshSnapshot refreshes the snapshot cache if there is a change in the source.
// returns true if snapshot is updated
func (s *XDSServer) refreshSnapshot(ctx context.Context) (bool, error) {
	if !s.source.HasChange() {
		return false, nil
	}

	s.version++
	log.Infof("Updating snapshot (version: %d)", s.version)

	resources := s.source.Resources()
	snapshot, err := cache.NewSnapshot(
		fmt.Sprintf("%d", s.version),
		map[resource.Type][]types.Resource{
			resource.ClusterType:  source.ToResourceList(resources.Clusters),
			resource.RouteType:    source.ToResourceList(resources.Routes),
			resource.ListenerType: source.ToResourceList(resources.Listeners),
			resource.EndpointType: source.ToResourceList(resources.Endpoints),
		},
	)
	if err != nil {
		return false, errors.Wrap(err, "Failed to generate snapshot")
	}

	err = snapshot.Consistent()
	if err != nil {
		referencedResources := cache.GetAllResourceReferences(snapshot.Resources)
		referencedResponseType := map[types.ResponseType]struct{}{
			types.Endpoint: {},
			types.Route:    {},
		}

		for idx, items := range snapshot.Resources {
			responseType := types.ResponseType(idx)
			if _, ok := referencedResponseType[responseType]; ok {
				typeURL, err := cache.GetResponseTypeURL(responseType)
				if err != nil {
					return false, err
				}
				referenceSet := referencedResources[typeURL]

				var missing, extra []string
				for name := range items.Items {
					if _, ok := referenceSet[name]; !ok {
						extra = append(extra, name)
					}
				}

				for name := range referenceSet {
					if _, ok := items.Items[name]; !ok {
						missing = append(missing, name)
					}
				}

				if len(missing) > 0 || len(extra) > 0 {
					log.Errorf(
						"mismatched %s reference and resource lengths: %d missing: %v, %d: extra: %v",
						typeURL, len(missing), missing, len(extra), extra,
					)
				}
			}
		}

		return false, errors.Wrap(err, "Failed to check snapshot consistency")
	}

	// TODO: log diffs
	log.Infof("Update snapshot")
	err = s.snapshotCache.SetSnapshot(ctx, s.nodeId, snapshot)
	if err != nil {
		return false, errors.Wrap(err, "Failed to set snapshot")
	}

	return true, nil
}
