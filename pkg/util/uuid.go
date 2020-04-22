package util

import (
	"github.com/lxdlam/vertex/pkg/common"
	uuid "github.com/satori/go.uuid"
)

var counter int64 = 0

// GenNewUUID is a simple proxy function to github.com/satori/go.uuid
//
// The default policy is v4, i.e., random generate.
// If failed, using v1 policy as fallback.
func GenNewUUID() (ret string) {
	defer func() {
		if err := recover(); err != nil {
			common.Warn("uuid: fallback to v1 policy")
			ret = uuid.NewV1().String()
		}
	}()

	return uuid.NewV4().String()
}
