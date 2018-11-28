// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package descriptor provides functions for obtaining protocol buffer
// descriptors for generated Go types.
//
// Deprecated: Do not use. The new v2 Message interface provides direct support
// for programmatically interacting with the descriptor information.
package descriptor

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// extractFile extracts a FileDescriptorProto from a gzip'd buffer.
func extractFile(gz []byte) (*descriptorpb.FileDescriptorProto, error) {
	r, err := gzip.NewReader(bytes.NewReader(gz))
	if err != nil {
		return nil, fmt.Errorf("failed to open gzip reader: %v", err)
	}
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to uncompress descriptor: %v", err)
	}

	fd := new(descriptorpb.FileDescriptorProto)
	if err := proto.Unmarshal(b, fd); err != nil {
		return nil, fmt.Errorf("malformed FileDescriptorProto: %v", err)
	}

	return fd, nil
}

// Message is a proto.Message with a method to return its descriptor.
//
// Message types generated by the protocol compiler always satisfy
// the Message interface.
type Message interface {
	proto.Message
	Descriptor() ([]byte, []int)
}

// ForMessage returns a FileDescriptorProto and a DescriptorProto from within it
// describing the given message.
func ForMessage(msg Message) (fd *descriptorpb.FileDescriptorProto, md *descriptorpb.DescriptorProto) {
	gz, path := msg.Descriptor()
	fd, err := extractFile(gz)
	if err != nil {
		panic(fmt.Sprintf("invalid FileDescriptorProto for %T: %v", msg, err))
	}

	md = fd.MessageType[path[0]]
	for _, i := range path[1:] {
		md = md.NestedType[i]
	}
	return fd, md
}
