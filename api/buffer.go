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

package api

type IoBuffer interface {
	// Len returns the number of bytes of the unread portion of the buffer;
	// b.Len() == len(b.Bytes()).
	Len() int

	// Bytes returns all bytes from buffer, without draining any buffered data.
	// It can be used to get fixed-length content, such as headers, body.
	// Note: do not change content in return bytes, use write instead
	Bytes() []byte

	// Read reads the next len(p) bytes from the buffer or until the buffer
	// is drained. The return value n is the number of bytes read. If the
	// buffer has no data to return, err is io.EOF (unless len(p) is zero);
	// otherwise it is nil.
	Read(p []byte) (n int, err error)

	// Write appends the contents of p to the buffer, growing the buffer as
	// needed. The return value n is the length of p; err is always nil. If the
	// buffer becomes too large, Write will panic with ErrTooLarge.
	Write(p []byte) (n int, err error)

	// Truncate discards all but the first n unread bytes from the buffer
	// but continues to use the same allocated storage.
	// It panics if n is negative or greater than the length of the buffer.
	Truncate(n int)

	// Reset resets the buffer to be empty,
	// but it retains the underlying storage for use by future writes.
	// Reset is the same as Truncate(0).
	Reset()
}
