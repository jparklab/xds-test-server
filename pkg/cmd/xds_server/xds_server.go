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
package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	xds "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/jparklab/xds-test-server/pkg/generated/config"
	"github.com/jparklab/xds-test-server/pkg/utils/log"
	"github.com/jparklab/xds-test-server/pkg/xdsserver"
	"github.com/jparklab/xds-test-server/pkg/xdsserver/source"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	xdsPort = flag.Int("xds-port", 5000, "the port to serve xDS service requests on.")
	nodeID  = flag.String("node-id", "test-id", "Node ID")
)

func main() {
	flag.Parse()

	xdsListener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *xdsPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(128*1024*1024),
		grpc.MaxSendMsgSize(128*1024*1024),
	)
	reflection.Register(grpcServer)

	ctx, cancel := context.WithCancel(context.Background())

	snapshotCache := cache.NewSnapshotCache(
		true,
		cache.IDHash{}, nil)
	discovery.RegisterAggregatedDiscoveryServiceServer(
		grpcServer,
		xds.NewServer(ctx, snapshotCache, xdsserver.XDSCallBacks),
	)

	configSource := source.NewManualSource()
	config.RegisterConfigServiceServer(grpcServer, configSource)

	initDone := make(chan struct{}, 1)
	xdsServer := xdsserver.NewXDSServer(
		snapshotCache,
		configSource,
		xdsserver.WithNodeId(*nodeID),
	)

	// start xds server to start processing snapshots
	go xdsServer.Start(ctx, initDone)

	/* do not wait for init since we need to listen to grpc for commands
	select {
	case <-initDone:
		log.Infof("Initialization done")
		break
	case <-time.After(1 * time.Minute):
		log.Fatalf("Timed out while waiting for xDS server to initialize")
	}
	*/

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		s := <-sigCh
		log.Infof("Received a signal %v, attempting graceful shutdown", s)
		cancel()

		grpcServer.GracefulStop()
		wg.Done()
	}()

	log.Infof("Starting grpc server at port %d", *xdsPort)
	err = grpcServer.Serve(xdsListener)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	shutDownDone := make(chan struct{}, 1)
	go func() {
		defer close(shutDownDone)
		wg.Wait()
	}()

	select {
	case <-shutDownDone:
		log.Infof("Shutdown complete")
	case <-time.After(5 * time.Second):
		log.Warnf("Timed out while waiting for shutdown to complete, force shutting down")
		grpcServer.Stop()
	}

}
