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
	gFile.P("Release", config.structName, "(", config.objectName, ")")
	gFile.P("return nil, err")
	gFile.P("}")
	gFile.P("return ", config.objectName, ", nil")
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

func GenFdbMethods(gFile *protogen.GeneratedFile, message *protogen.Message) {
	config := genFDBConfig{
		structName: message.GoIdent.GoName,
		objectName: LcFirst(message.GoIdent.GoName),
	}
	gFile.QualifiedGoIdent(protogen.GoIdent{
		"proto",
		"google.golang.org/protobuf/proto",
	})
	genLoadFromFDB(gFile, config)
	genStoreFDB(gFile, config)
}
