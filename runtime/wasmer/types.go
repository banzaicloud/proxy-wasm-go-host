// Copyright (c) 2023 Cisco and/or its affiliates. All rights reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//       https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package wasmer

import (
	"log"
	"reflect"

	wasmerGo "github.com/wasmerio/wasmer-go/wasmer"
)

func convertFromGoType(t reflect.Type) *wasmerGo.ValueType {
	switch t.Kind() { //nolint:exhaustive
	case reflect.Int32:
		return wasmerGo.NewValueType(wasmerGo.I32)
	case reflect.Int64:
		return wasmerGo.NewValueType(wasmerGo.I64)
	case reflect.Float32:
		return wasmerGo.NewValueType(wasmerGo.F32)
	case reflect.Float64:
		return wasmerGo.NewValueType(wasmerGo.F64)
	default:
		log.Printf("[wasmer][type] convertFromGoType unsupported type: %v", t.Kind().String())
	}

	return nil
}

func convertToGoTypes(in wasmerGo.Value) reflect.Value {
	switch in.Kind() { //nolint:exhaustive
	case wasmerGo.I32:
		return reflect.ValueOf(in.I32())
	case wasmerGo.I64:
		return reflect.ValueOf(in.I64())
	case wasmerGo.F32:
		return reflect.ValueOf(in.F32())
	case wasmerGo.F64:
		return reflect.ValueOf(in.F64())
	}

	return reflect.Value{}
}

func convertFromGoValue(val reflect.Value) wasmerGo.Value {
	switch val.Kind() { //nolint:exhaustive
	case reflect.Int32:
		return wasmerGo.NewI32(int32(val.Int()))
	case reflect.Int64:
		return wasmerGo.NewI64(val.Int())
	case reflect.Float32:
		return wasmerGo.NewF32(float32(val.Float()))
	case reflect.Float64:
		return wasmerGo.NewF64(val.Float())
	default:
		log.Printf("[wasmer][type] convertFromGoValue unsupported val type: %v", val.Kind().String())
	}

	return wasmerGo.Value{}
}
