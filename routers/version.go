package router

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

// VersionInfo tells the current code base
type VersionInfo struct {
	GitCommit string
	Version   string
}

const (
	versionPrefix = "/version"
)

var versionInfo VersionInfo

func init() {
	if err := RegisterRouter(versionPrefix, versionNameSpace()); err != nil {
		logs.Error("Failed to register router: '%s'.", versionPrefix)
	} else {
		logs.Debug("Register router '%s' registered.", versionPrefix)
	}
}

// SetVersionInfo sets the gitcommit and verion from Make command
func SetVersionInfo(commit string, version string) {
	versionInfo.GitCommit = commit
	versionInfo.Version = version
}

// versionNameSpace defines the version router
func versionNameSpace() *beego.Namespace {
	ns := beego.NewNamespace(versionPrefix,
		beego.NSCond(func(ctx *context.Context) bool {
			return true
		}),
		beego.NSGet("/", func(ctx *context.Context) {
			data, _ := json.Marshal(versionInfo)
			ctx.Output.Body(data)
		}),
	)

	return ns
}
