// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ctxhttp

import (
	"net/http"

	"errors"
	"net/url"
	"strings"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/utils/log"
	"golang.org/x/net/context"
)

// ErrBaseUrlDoNotMatch will be returned if the request URL does not match the configured URL.
var ErrBaseUrlDoNotMatch = errors.New("The Base URLs do not match")

// WithValidateBaseUrl is a middleware which checks if the request base URL
// is equal to the one store in the configuration, if not
// i.e. redirect from http://example.com/store/ to http://www.example.com/store/
// @see app/code/Magento/Store/App/FrontController/Plugin/RequestPreprocessor.php
func WithValidateBaseUrl(cr config.ReaderPubSuber) Middleware {

	// Having the GetBool command here, means you must restart the app to take
	// changes in effect. @todo refactor and use pub/sub to automatically change
	// the isRedirectToBase value.
	checkBaseURL, err := cr.GetBool(config.Path(store.PathRedirectToBase)) // scope default
	if config.NotKeyNotFoundError(err) {
		log.Error("ctxhttp.WithValidateBaseUrl.GetBool", "err", err, "path", store.PathRedirectToBase)
	}

	redirectCode := http.StatusMovedPermanently
	if rc, err := cr.GetInt(config.Path(store.PathRedirectToBase)); rc != redirectCode && false == config.NotKeyNotFoundError(err) {
		redirectCode = http.StatusFound
	}

	return func(h Handler) Handler {
		return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			storeManager, ok := store.FromContextManagerReader(ctx)
			if !ok {
				return log.Error("ctxhttp.WithValidateBaseUrl.FromContextManagerReader", "err", errors.New("Cannot extract config.Reader from context"), "ctx", ctx)
			}

			if checkBaseURL && r.Method != "POST" {
				store, err := storeManager.Store()
				if err != nil {
					return log.Error("ctxhttp.WithValidateBaseUrl.storeManager.Store", "err", err, "ctx", ctx)
				}

				baseURL := store.BaseURL(config.URLTypeWeb, store.IsCurrentlySecure(r))
				if nil == isBaseUrlCorrect(r, baseURL) {
					redirectURL := baseURL + r.URL.Path
					http.Redirect(w, r, redirectURL, redirectCode)
					return nil
				}
			}

			return h.ServeHTTPContext(ctx, w, r)
		})
	}
}

// isBaseUrlCorrect checks if the requested host, scheme are same as the servers and
// if the path of the baseURL is included in the request URI.
func isBaseUrlCorrect(r *http.Request, baseURL string) error {
	uri, err := url.Parse(baseURL)
	if err != nil {
		return log.Error("store.isBaseUrlCorrect.url.Parse", "err", err)
	}

	if r.Host == uri.Host && r.URL.Host == uri.Host && r.URL.Scheme == uri.Scheme && strings.Contains(r.URL.RequestURI(), uri.Path) {
		return nil
	}
	return log.Error("store.isBaseUrlCorrect.compare", "err", ErrBaseUrlDoNotMatch, "r.Host", r.Host, "baseURL", uri.String(), "requestURL", r.URL.String(), "strings.Contains", []string{r.URL.RequestURI(), uri.Path})
}
