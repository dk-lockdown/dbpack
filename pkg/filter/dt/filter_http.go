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

package dt

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/cectc/dbpack/pkg/filter"
	dbpackHttp "github.com/cectc/dbpack/pkg/http"
	"github.com/cectc/dbpack/pkg/log"
	"github.com/cectc/dbpack/pkg/proto"
)

const httpFilter = "HttpDistributedTransaction"

func init() {
	filter.RegistryFilterFactory(httpFilter, &httpFactory{})
}

type httpFactory struct {
}

func (factory *httpFactory) NewFilter(config map[string]interface{}) (proto.Filter, error) {
	var (
		err          error
		content      []byte
		filterConfig *HttpFilterConfig
	)

	if content, err = json.Marshal(config); err != nil {
		return nil, errors.Wrap(err, "marshal http distributed transaction filter config failed.")
	}
	if err = json.Unmarshal(content, &filterConfig); err != nil {
		log.Errorf("unmarshal http distributed transaction filter failed, %s", err)
		return nil, err
	}

	f := &_httpFilter{
		conf: filterConfig,
	}

	for _, ti := range filterConfig.TransactionInfos {
		dbpackHttp.AddTransactionInfo(&dbpackHttp.TransactionInfo{
			RequestPath: ti.RequestPath,
			Timeout:     ti.Timeout,
		})
		log.Debugf("proxy %s, will create global transaction, put xid into request header", ti.RequestPath)
	}

	for _, r := range filterConfig.TccResourceInfos {
		dbpackHttp.AddTccResourceInfo(&dbpackHttp.TccResourceInfo{
			PrepareRequestPath:  r.PrepareRequestPath,
			CommitRequestPath:   r.CommitRequestPath,
			RollbackRequestPath: r.RollbackRequestPath,
		})
		log.Debugf("proxy %s, will register branch transaction", r.PrepareRequestPath)
	}
	return f, nil
}

// TransactionInfo transaction info config
type TransactionInfo struct {
	RequestPath string `yaml:"request_path" json:"request_path"`
	Timeout     int32  `yaml:"timeout" json:"timeout"`
}

// TccResourceInfo tcc resource config
type TccResourceInfo struct {
	PrepareRequestPath  string `yaml:"prepare_request_path" json:"prepare_request_path"`
	CommitRequestPath   string `yaml:"commit_request_path" json:"commit_request_path"`
	RollbackRequestPath string `yaml:"rollback_request_path" json:"rollback_request_path"`
}

// HttpFilterConfig http filter config
type HttpFilterConfig struct {
	ApplicationID string `yaml:"appid" json:"appid"`
	BackendHost   string `yaml:"backend_host" json:"backend_host"`

	TransactionInfos []*TransactionInfo `yaml:"transaction_infos" json:"transaction_infos"`
	TccResourceInfos []*TccResourceInfo `yaml:"tcc_resource_infos" json:"tcc_resource_infos"`
}

type _httpFilter struct {
	conf *HttpFilterConfig
}

func (f *_httpFilter) GetKind() string {
	return httpFilter
}

func (f _httpFilter) PreHandle(ctx *fasthttp.RequestCtx) error {
	path := ctx.Request.RequestURI()
	method := ctx.Method()

	if !strings.EqualFold(string(method), fasthttp.MethodPost) {
		return nil
	}

	transactionInfos := dbpackHttp.GetTransactionInfos()
	transactionInfo, found := transactionInfos[strings.ToLower(string(path))]
	if found {
		result, err := f.handleHttp1GlobalBegin(ctx, transactionInfo)
		if !result {
			if err := f.handleHttp1GlobalEnd(ctx); err != nil {
				log.Error(err)
			}
		}
		return err
	}

	tccResourceInfos := dbpackHttp.GetTccResourceInfo()
	tccResourceInfo, exists := tccResourceInfos[strings.ToLower(string(path))]
	if exists {
		result, err := f.handleHttp1BranchRegister(ctx, tccResourceInfo)
		if !result {
			if err := f.handleHttp1BranchEnd(ctx); err != nil {
				log.Error(err)
			}
		}
		return err
	}
	return nil
}

func (f _httpFilter) PostHandle(ctx *fasthttp.RequestCtx) error {
	path := ctx.Request.RequestURI()
	method := ctx.Method()

	if !strings.EqualFold(string(method), fasthttp.MethodPost) {
		return nil
	}

	transactionInfos := dbpackHttp.GetTransactionInfos()
	_, found := transactionInfos[strings.ToLower(string(path))]
	if found {
		if err := f.handleHttp1GlobalEnd(ctx); err != nil {
			return err
		}
	}

	tccResourceInfos := dbpackHttp.GetTccResourceInfo()
	_, exists := tccResourceInfos[strings.ToLower(string(path))]
	if exists {
		if err := f.handleHttp1BranchEnd(ctx); err != nil {
			return err
		}
	}
	return nil
}
