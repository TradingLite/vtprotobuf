// Copyright (c) 2021 PlanetScale Inc. All rights reserved.
// Copyright (c) 2013, The GoGo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package factory

import (
	"fmt"
	"github.com/planetscale/vtprotobuf/generator"
	"github.com/twmb/murmur3"
	"google.golang.org/protobuf/compiler/protogen"
	"hash"
	"hash/fnv"
)

var (
	_ = protogen.GoImportPath("tradinglite.com/core/types")
)

func init() {
	generator.RegisterFeature("factory", func(gen *generator.GeneratedFile) generator.FeatureGenerator {
		return &factory{GeneratedFile: gen}
	})
}

type factory struct {
	*generator.GeneratedFile
	once bool
}

var _ generator.FeatureGenerator = (*factory)(nil)

func (p *factory) Name() string {
	return "factory"
}

func (p *factory) GenerateFile(file *protogen.File) bool {
	p.P(`func init() {`)
	for _, message := range file.Messages {
		p.register(message)
	}
	p.P(`}`)
	p.P()
	for _, message := range file.Messages {
		p.factoryName(message)
	}
	p.P()
	for _, message := range file.Messages {
		p.factoryHash(message)
	}
	p.P()
	return p.once
}

func (p *factory) GenerateHelpers() {

}

func (p *factory) register(message *protogen.Message) {
	for _, nested := range message.Messages {
		p.register(nested)
	}

	if message.Desc.IsMapEntry() {
		return
	}

	p.once = true
	//sizeName := "MessageNameVT"
	typeName := message.GoIdent
	//typeName := fmt.Sprintf("%s.%s", f.Package, name)

	// types.RegisterHash(0xcabae491495e1ce2, "api.TickTrade", func() types.Message { return &TickTrade{} })
	//p.P(`types.RegisterHash(0x`, message.Desc.Hash(), `, "`, typeName, `", func() types.Message { return &`, typeName, `{}})`)

	typeHash := fmt.Sprintf("%#x", Hash(string(message.Desc.FullName())))
	p.P(`types.RegisterProto(`, typeHash, `, "`, message.Desc.FullName(), `", func() types.ProtoMessageVT { return &`, typeName, `{}})`)
}

func (p *factory) factoryName(message *protogen.Message) {
	for _, nested := range message.Messages {
		p.factoryName(nested)
	}
	if message.Desc.IsMapEntry() {
		return
	}
	p.once = true
	typeName := message.GoIdent
	p.P(`func (`, typeName, `) MessageNameVT() string {`, `return "`, message.Desc.FullName(), `"`, `}`)
}

func (p *factory) factoryHash(message *protogen.Message) {
	for _, nested := range message.Messages {
		p.factoryHash(nested)
	}
	if message.Desc.IsMapEntry() {
		return
	}
	p.once = true
	typeName := message.GoIdent
	typeHash := fmt.Sprintf("%#x", Hash(string(message.Desc.FullName())))
	//p.P(`/*`)
	//p.P(message.Desc.FullName())
	//p.P(`*/`)
	p.P(`func (`, typeName, `) MessageHashVT() uint32 {`, `return `, typeHash, ``, `}`)
}

func FNV64a(text string) uint64 {
	algorithm := fnv.New64()
	return uint64Hasher(algorithm, text)
}

func uint64Hasher(algorithm hash.Hash64, text string) uint64 {
	algorithm.Write([]byte(text))
	return algorithm.Sum64()
}
func Hash(text string) uint32 {
	return murmur3.StringSum32(text)
}
