package micro_service

import (
	"catface/app/global/variable"
	"context"
	"fmt"
	"strings"

	"github.com/carlmjohnson/requests"
)

func TestLinkPythonService() bool {
	err := requests.URL(FetchPythonServiceUrl("link_test")).Fetch(context.Background())
	return err == nil
}

func FetchPythonServiceUrl(url string) string {
	// 检查 url 是否以 / 开头，如果是则去掉开头的 /
	if strings.HasPrefix(url, "/") {
		url = url[1:]
	}

	return fmt.Sprintf(`http://%s:%v/%s/%s`,
		variable.ConfigYml.GetString("PythonService.Host"),
		variable.ConfigYml.GetString("PythonService.Port"),
		variable.ConfigYml.GetString("PythonService.TopUrl"),
		url)
}
