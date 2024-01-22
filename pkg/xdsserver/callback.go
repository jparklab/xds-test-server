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
	"strings"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	xds "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/jparklab/xds-test-server/pkg/utils/log"
)

var XDSCallBacks = xds.CallbackFuncs{
	DeltaStreamOpenFunc: func(ctx context.Context, streamID int64, typeURL string) error {
		log.With(
			"stream_id", streamID,
			"event", XDSEventConnected,
			"type", extractTypeName(typeURL),
		).Infof("XDS Event")
		return nil
	},
	DeltaStreamClosedFunc: func(streamID int64, node *core.Node) {
		log.With(
			"stream_id", streamID,
			"event", XDSEventClosed,
		).Infof("XDS Event")
	},
	StreamDeltaRequestFunc: func(streamID int64, req *discovery.DeltaDiscoveryRequest) error {
		typeName := extractTypeName(req.TypeUrl)
		nonce := req.ResponseNonce
		metadata := getNodeMetadata(req.Node)

		eventLogger := log.With(
			"stream_id", streamID,
			"resource_type", typeName,
		).With(metadata...)

		if nonce != "" {
			if req.ErrorDetail != nil {
				eventLogger.With(
					"event", XDSEventNACK,
					"nonce", nonce,
				).Infof("XDS Event")
			} else {
				eventLogger.With(
					"event", XDSEventACK,
					"nonce", nonce,
				).Infof("XDS Event")
			}
		}
		for _, name := range req.ResourceNamesSubscribe {
			eventLogger.With(
				"event", XDSEventSubscribe,
				"resource_name", name,
			).Infof("XDS Event")
		}
		for _, name := range req.ResourceNamesUnsubscribe {
			eventLogger.With(
				"event", XDSEventUnsubscribe,
				"resource_name", name,
			).Infof("XDS Event")
		}
		if len(req.ResourceNamesSubscribe) == 0 && len(req.ResourceNamesUnsubscribe) == 0 {
			eventLogger.With(
				"event", XDSEventWildcardRequest,
				"nonce", nonce,
			).Infof("XDS Event")
		}

		return nil
	},
	StreamDeltaResponseFunc: func(streamID int64, req *discovery.DeltaDiscoveryRequest, res *discovery.DeltaDiscoveryResponse) {
		typeName := extractTypeName(req.TypeUrl)
		metadata := getNodeMetadata(req.Node)

		eventLogger := log.With(
			"stream_id", streamID,
			"resource_type", typeName,
		).With(metadata...)

		eventLogger.Info("response")

		for _, name := range res.RemovedResources {
			eventLogger.With(
				"event", XDSEventRemove,
				"resource_name", name,
			).Infof("XDS Event")
		}

		for _, name := range res.Resources {
			eventLogger.With(
				"event", XDSEventResponse,
				"resource_name", name,
				"nonce", res.Nonce,
			).Infof("XDS Event")
		}
	},
}

func extractTypeName(typeURL string) string {
	tokens := strings.Split(typeURL, ".")
	return tokens[len(tokens)-1]
}

func getNodeMetadata(node *core.Node) []interface{} {
	metadata := node.GetMetadata().AsMap()

	keyValueList := make([]interface{}, 0, len(metadata)*2)
	for key, value := range metadata {
		keyValueList = append(keyValueList, key, value)
	}

	return keyValueList
}
