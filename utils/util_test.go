package utils

import (
	"testing"
	"fmt"
)

func TestInitCache(t *testing.T) {


	//InitLogs()
	InitCache()

	LogDebug("1212")


	SetCache("testname","xiaozhun",1)


	name:=[]byte{}
	GetCache("testname",name)

	fmt.Println("name :",name)



}


