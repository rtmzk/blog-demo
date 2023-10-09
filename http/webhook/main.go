// Copyright 2023 rtmzk
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

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	AccessApiVersion = "authorization.k8s.io/v1beta1"
	AccessKind       = "SubjectAccessReview"
)

func resp(allowed bool, reason string) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetAPIVersion(AccessApiVersion)
	obj.SetKind(AccessKind)
	obj.Object["status"] = map[string]interface{}{
		"allowed": allowed,
		"reason":  reason,
	}

	return obj
}

func main() {
	http.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp(false, "bad request"))
			return
		}
		var obj = &unstructured.Unstructured{}
		err = json.Unmarshal(b, obj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp(false, "bad request"))
			return
		}

		fmt.Println(string(b))

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp(true, ""))
	})

	http.ListenAndServeTLS(":9090", "certs/server.crt", "certs/server.key", nil)
}
