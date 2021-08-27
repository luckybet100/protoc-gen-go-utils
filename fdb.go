package main

import "google.golang.org/protobuf/compiler/protogen"

type genFDBConfig struct {
	objectName string
	structName string
}

func genLoadFromFDB(gFile *protogen.GeneratedFile, config genFDBConfig) {
	gFile.P("func Load", config.structName, "FDB(tr fdb.ReadTransaction, key subspace.Subspace) (*", config.structName, ", error) {")
	gFile.P("value, err := tr.Get(key).Get()")
	gFile.P("if err != nil {")
	gFile.P("return nil, err")
	gFile.P("} else if value == nil {")
	gFile.P("return nil, nil")
	gFile.P("}")
	gFile.P(config.objectName, ":=", "Acquire", config.structName, "()")
	gFile.P("if err = proto.Unmarshal(value,", config.objectName, "); err != nil {")
	gFile.P("return nil, err")
	gFile.P("}")
	gFile.P("return ", config.objectName, ", nil")
	gFile.P("}")
	gFile.P()
}

func genMustLoadFromFDB(gFile *protogen.GeneratedFile, config genFDBConfig) {
	gFile.P("func MustLoad", config.structName, "FDB(tr fdb.ReadTransaction, key subspace.Subspace) *", config.structName, " {")
	gFile.P("value := tr.Get(key).MustGet()")
	gFile.P("if value == nil {")
	gFile.P("return nil")
	gFile.P("}")
	gFile.P(config.objectName, ":=", "Acquire", config.structName, "()")
	gFile.P("if err := proto.Unmarshal(value,", config.objectName, "); err != nil {")
	gFile.P("panic(err)")
	gFile.P("}")
	gFile.P("return ", config.objectName, "")
	gFile.P("}")
	gFile.P()
}

func genStoreFDB(gFile *protogen.GeneratedFile, config genFDBConfig) {
	gFile.P("func Store", config.structName, "FDB(tr fdb.Transaction, key subspace.Subspace, ", config.objectName, " *", config.structName, ") error {")
	gFile.P("value, err := proto.Marshal(", config.objectName, ")")
	gFile.P("if err != nil {")
	gFile.P("return err")
	gFile.P("}")
	gFile.P("tr.Set(key, value)")
	gFile.P("return nil")
	gFile.P("}")
	gFile.P()
}

func genMustStoreFDB(gFile *protogen.GeneratedFile, config genFDBConfig) {
	gFile.P("func MustStore", config.structName, "FDB(tr fdb.Transaction, key subspace.Subspace, ", config.objectName, " *", config.structName, ") {")
	gFile.P("value, err := proto.Marshal(", config.objectName, ")")
	gFile.P("if err != nil {")
	gFile.P("panic(err)")
	gFile.P("}")
	gFile.P("tr.Set(key, value)")
	gFile.P("}")
	gFile.P()
}

func genRangeFDB(gFile *protogen.GeneratedFile, config genFDBConfig) {
	gFile.P("func GetRange", config.structName, "FDB(tr fdb.ReadTransaction, rng fdb.Range, opts fdb.RangeOptions) ([]*", config.structName, ", error) {")
	gFile.P("iter := tr.GetRange(rng, opts).Iterator()")
	gFile.P("result := make([]*", config.structName, ", 0, opts.Limit)")
	gFile.P("for iter.Advance() {")
	gFile.P("kv, err := iter.Get()")
	gFile.P("if err != nil {")
	gFile.P("return nil, err")
	gFile.P("}")
	gFile.P(config.objectName, ":=Acquire", config.structName, "()")
	gFile.P("err = proto.Unmarshal(kv.Value, ", config.objectName, ")")
	gFile.P("if err != nil {")
	gFile.P("return nil, err")
	gFile.P("}")
	gFile.P("result = append(result, ", config.objectName, ")")
	gFile.P("}")
	gFile.P("return result, nil")
	gFile.P("}")
	gFile.P()
}

func genMustRangeFDB(gFile *protogen.GeneratedFile, config genFDBConfig) {
	gFile.P("func MustGetRange", config.structName, "FDB(tr fdb.ReadTransaction, rng fdb.Range, opts fdb.RangeOptions) ([]*", config.structName, ", error) {")
	gFile.P("iter := tr.GetRange(rng, opts).Iterator()")
	gFile.P("result := make([]*", config.structName, ", 0, opts.Limit)")
	gFile.P("for iter.Advance() {")
	gFile.P("kv := iter.MustGet()")
	gFile.P(config.objectName, ":=Acquire", config.structName, "()")
	gFile.P("err := proto.Unmarshal(kv.Value, ", config.objectName, ")")
	gFile.P("if err != nil {")
	gFile.P("panic(err)")
	gFile.P("}")
	gFile.P("result = append(result, ", config.objectName, ")")
	gFile.P("}")
	gFile.P("return result, nil")
	gFile.P("}")
	gFile.P()
}

func GenFdbMethods(gFile *protogen.GeneratedFile, message *protogen.Message) {
	config := genFDBConfig{
		structName: message.GoIdent.GoName,
		objectName: LcFirst(message.GoIdent.GoName),
	}
	gFile.QualifiedGoIdent(protogen.GoIdent{
		"proto",
		"google.golang.org/protobuf/proto",
	})
	gFile.QualifiedGoIdent(protogen.GoIdent{
		"fdb",
		"github.com/apple/foundationdb/bindings/go/src/fdb",
	})
	gFile.QualifiedGoIdent(protogen.GoIdent{
		"subspace",
		"github.com/apple/foundationdb/bindings/go/src/fdb/subspace",
	})
	genLoadFromFDB(gFile, config)
	genMustLoadFromFDB(gFile, config)
	genStoreFDB(gFile, config)
	genMustStoreFDB(gFile, config)
	genRangeFDB(gFile, config)
	genMustRangeFDB(gFile, config)
}
