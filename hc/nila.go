package hc

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/megamsys/libgo/hc"
	"github.com/megamsys/megamd/meta"
)

func init() {
	hc.AddChecker("nilavu", healthCheck)
}

func healthCheck() error {
	filename := filepath.Join(meta.MC.Home, "nilavu.yml")
	if _, err := os.Stat(filename); err == nil {
		return hc.ErrDisabledComponent
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		//match for flags riak, http_api, http_logs and store it in a bytesbuffer
		//line := scanner.Text()
		// if line.match("riak" || "http_api" || "http_logs") {
		//    append(line,buffer)
		// }
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
