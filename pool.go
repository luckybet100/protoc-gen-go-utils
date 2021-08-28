package main

import (
	"google.golang.org/protobuf/compiler/protogen"
)

type genPoolConfig struct {
	objectName string
	structName string
	poolName   string
}

func genSyncPool(gFile *protogen.GeneratedFile, config genPoolConfig) {
	gFile.P("var ", config.poolName, " = sync.Pool{")
	gFile.P("New: func() interface{} { ")
	gFile.P("return new(", config.structName, ")")
	gFile.P("},")
	gFile.P("}")
	gFile.P()
}

func genAcquire(gFile *protogen.GeneratedFile, config genPoolConfig) {
	gFile.P("func Acquire", config.structName, "() *", config.structName, "{")
	gFile.P(config.objectName, ":=", config.poolName, ".Get()", ".(*", config.structName, ")")
	gFile.P("runtime.SetFinalizer(", config.objectName, ", func (", config.objectName, " *", config.structName, ") {")
	gFile.P("Release", config.structName, "(", config.objectName, ")")
	gFile.P("})")
	gFile.P("return ", config.poolName, ".Get().(*", config.structName, ")")
	gFile.P("}")
	gFile.P()
}

func genRelease(gFile *protogen.GeneratedFile, config genPoolConfig) {
	gFile.P("func Release", config.structName, "(", config.objectName, " *", config.structName, ") {")
	gFile.P("if ", config.objectName, " != nil {")
	gFile.P("runtime.SetFinalizer(", config.objectName, ", nil)")
	gFile.P(config.objectName, ".Reset()")
	gFile.P(config.poolName, ".Put(", config.objectName, ")")
	gFile.P("}")
	gFile.P("}")
	gFile.P()
}

func GeneratePool(gFile *protogen.GeneratedFile, message *protogen.Message) {
	config := genPoolConfig{
		structName: message.GoIdent.GoName,
		objectName: LcFirst(message.GoIdent.GoName),
		poolName:   LcFirst(message.GoIdent.GoName) + "Pool",
	}
	gFile.QualifiedGoIdent(protogen.GoIdent{
		"sync",
		"sync",
	})
	gFile.QualifiedGoIdent(protogen.GoIdent{
		"runtime",
		"runtime",
	})
	genSyncPool(gFile, config)
	genAcquire(gFile, config)
	genRelease(gFile, config)
}
