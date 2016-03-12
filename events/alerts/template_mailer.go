package alerts

import (
	"bytes"
	"errors"
	"github.com/megamsys/vertice/meta"
	"html/template"
	"os"
	"path/filepath"
)

func body(name string, mp map[string]string) (string, error) {
	if meta.MC == nil {
		return "", errors.New(`[meta] Is it there in vertice.conf ?`)
	}

	f := filepath.Join(meta.MC.Dir, "mailer", name+".html")
	if _, err := os.Stat(f); err != nil {
		return "", err
	}
	var w bytes.Buffer
	t, err := template.ParseFiles(f)
	if err != nil {
		return "", err
	}

	if err = t.Execute(&w, mp); err != nil {
		return "", err
	}
	return w.String(), nil
}
