package views

import (
	"blog/context"
	"blog/models"
	"errors"
	"fmt"
	"github.com/gorilla/csrf"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
)

type public interface {
	Public() string
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, pattern ...string) (Template, error) {
	t := template.New(filepath.Base(pattern[0]))
	t = t.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrf field not implemented")
			},
			"currentUser": func() (*models.User, error) {
				return nil, fmt.Errorf("currentUser not implemented")
			},
			"errors": func() []string{
				return nil
			},
		})

	t, err := t.ParseFS(fs, pattern...)
	if err != nil {
		return Template{}, fmt.Errorf("Error parsing template: %w", err)
	}

	return Template{
		t,
	}, nil
}

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request ,data interface{}, errs ...error) {
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Println(err)
		return
	}
	errMsgs := errMessages(errs...)
	tpl = tpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
		"currentUser": func() *models.User {
			return context.User(r.Context())
		},
		"errors": func() []string{
			return errMsgs
		},
	})
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "Error executing the template", http.StatusInternalServerError)
		return
	}
}

func errMessages(errs ...error) []string {
	var msgs []string
	for _, err := range errs {
		var pubErr public
		if errors.As(err, &pubErr) {
			msgs = append(msgs, pubErr.Public())
		} else {
			log.Println(err)
			msgs = append(msgs, "Something went wrong")
		}
	}
	return msgs
}