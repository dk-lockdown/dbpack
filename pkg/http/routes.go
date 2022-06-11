/*
 * Copyright 2022 CECTC, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package http

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/cectc/dbpack/pkg/config"
	"github.com/cectc/dbpack/pkg/log"
)

func InitHttpServer(ctx context.Context, conf config.HttpConf) {
	lis, err := net.Listen("tcp4", fmt.Sprintf(":%d", conf.Port))
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		<-ctx.Done()
		if err := lis.Close(); err != nil {
			log.Error(err)
		}
	}()
	handler := RegisterRoutes()
	server := &http.Server{
		Handler: handler,
	}
	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("unable create status server: %+v", err)
	}
	grpcServer := RegisterGrpc(conf)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("unable create status server: %+v", err)
	}
	log.Infof("start api server :  %s", lis.Addr())
}

func RegisterGrpc(conf config.HttpConf) *grpc.Server {
	s := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             conf.Grpc.EnforcementPolicy.MinTime,
		PermitWithoutStream: conf.Grpc.EnforcementPolicy.PermitWithoutStream,
	}), grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     conf.Grpc.ServerParameters.MaxConnectionIdle,
		MaxConnectionAge:      conf.Grpc.ServerParameters.MaxConnectionAge,
		MaxConnectionAgeGrace: conf.Grpc.ServerParameters.MaxConnectionAgeGrace,
		Time:                  conf.Grpc.ServerParameters.Time,
		Timeout:               conf.Grpc.ServerParameters.Timeout,
	}))
	RegisterProxyServiceServer(s, svc)
	return s
}

func RegisterRoutes() http.Handler {
	router := mux.NewRouter().SkipClean(true).UseEncodedPath()
	// Add healthcheck router
	registerHealthCheckRouter(router)

	// Add server metrics router
	registerMetricsRouter(router)

	// Add status router
	registerStatusRouter(router)

	return router
}
