package writeflowui

import (
	"bytes"
	"fmt"
	writeflowui "github.com/zbysir/writeflow-ui"
	"github.com/zbysir/writeflow/internal/pkg/fshook"
	"io/fs"
)

type UIConfig struct {
	ApiHost string
	WsHost  string
}

func UIFs(c UIConfig) fs.FS {
	d, _ := fs.Sub(writeflowui.Dist, "dist")

	d = fshook.NewFsHook(d, map[string]func(body []byte) []byte{
		"index.html": func(body []byte) []byte {
			return bytes.ReplaceAll(body, []byte("<head>"), append([]byte("<head>"), []byte(fmt.Sprintf(`<script>window.__service_host__ = {"api": %q, "ws": %q} </script>`, c.ApiHost, c.WsHost))...))
		},
	})

	return d
}
