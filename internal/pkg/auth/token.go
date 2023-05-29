package auth

import (
	"fmt"
	"github.com/zbysir/writeflow/internal/pkg/util"
	"time"
)

var t = time.Now()

func CreateToken(key string) string {
	return util.MD5(fmt.Sprintf("hollow%vhollow", key))
}

func CheckToken(key string, token string) bool {
	return CreateToken(key) == token
}
