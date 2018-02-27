package utils

import (
	"testing"
	"fmt"
	"runtime"
	"path/filepath"
	"github.com/astaxie/beego"
	"strings"
)

// 测试环境初始化
func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, "../../"+string(filepath.Separator))))
	parentPath := beego.Substr(apppath, 0, strings.LastIndex(apppath, "/"))
	beego.TestBeegoInit(parentPath)
}

func TestInitCache(t *testing.T) {

	//InitLogs()
	InitCache()
	//LogDebug("1212")
	//SetCache("testname","xiaozhun",0)
	var name string
	GetCache("testname", &name)
	fmt.Println("name :", name)

}

func TestUniqueId(t *testing.T) {
	fmt.Println(UniqueId())
	fmt.Println(UniqueId())

	fmt.Println(UniqueId())
	fmt.Println(UniqueId())

}
