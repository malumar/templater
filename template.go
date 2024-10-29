package templater

import (
	"bytes"
	"fmt"
	htmltpl "html/template"
	"io/ioutil"
	"os"
	texttpl "text/template"
)

const RecoverErrorMessage = "Template Error"

type Recover struct {
	TemplateName string
	SourceError  error
}

func PrintIfError(err error) error {
	if err == nil {
		return nil
	}

	if r, ok := err.(Recover); ok {
		fmt.Printf("\n%v\n", r.Format())
	} else {
		fmt.Printf("\n%v\n", err)
	}

	return err
}

func DebugOrExit(dest *bytes.Buffer, err error) {
	if dest != nil && dest.Len() > 0 {
		fmt.Printf("\n%s\n", dest.String())
	}
	PrintIfErrorAndExit(err)
}

func PrintIfErrorAndExit(err error) {
	if PrintIfError(err) != nil {
		os.Exit(1)
	}
}

func (r Recover) Error() string {
	return RecoverErrorMessage
}

func (r Recover) Format() error {
	return fmt.Errorf("%s in `%s`: %v", RecoverErrorMessage, r.TemplateName, r.Error())
}

func HtmlFromString(name, content string) *Parser {
	return newParser(nil, true, name, content)
}

func TextFromString(name, content string) *Parser {
	return newParser(nil, false, name, content)
}

func newParser(err error, ishtml bool, name, content string) *Parser {
	return &Parser{
		err:     err,
		html:    ishtml,
		name:    name,
		content: content,
	}
}

func TextFromFile(filename string) *Parser {
	content, err := ioutil.ReadFile(filename)
	return newParser(err, false, filename, string(content))
}

func HtmlFromFile(filename string) *Parser {
	content, err := ioutil.ReadFile(filename)
	return newParser(err, true, filename, string(content))
}

type Parser struct {
	err     error
	name    string
	content string
	html    bool
}

func (p *Parser) parseHtml(data interface{}, dest *bytes.Buffer) error {

	t := htmltpl.New(p.name).Funcs(functions)
	t, err := t.Parse(p.content)
	if err != nil {
		return err
	}
	return t.Execute(dest, data)

}

func (p *Parser) parseText(data interface{}, dest *bytes.Buffer) error {

	t := texttpl.New(p.name).Funcs(functions)
	t, err := t.Parse(p.content)
	if err != nil {
		return err
	}
	return t.Execute(dest, data)

}

func (p *Parser) Parse(data interface{}, dest *bytes.Buffer) (err error) {
	if p.err != nil {
		return p.err
	}
	/*
		defer func(returnError *error) {
			recovered:=false
			if r := recover(); r != nil {
				ERROR.Println("Mieliśmy recovery po panic!")
				var ok bool
				var err error
				err, ok = r.(error)

				if !ok {
					err = fmt.Errorf("pkg: %v -- err %v", r)
				}

				*returnError = err
				recovered=true
			}
	*/
	defer func(returnError *error) {
		if r := recover(); r != nil {
			errx, ok := r.(error)
			if !ok {
				errx = fmt.Errorf("pkg: %v -- err %v", r)
			}

			*returnError = Recover{
				TemplateName: p.name,
				SourceError:  errx,
			}
		}
	}(&err)

	if p.html {
		err = p.parseHtml(data, dest)
	} else {
		err = p.parseText(data, dest)
	}

	return
}

func (p *Parser) ParseToOutput(data interface{}) (out *bytes.Buffer, err error) {
	if p.err != nil {
		return nil, p.err
	}
	out = &bytes.Buffer{}
	if err = p.Parse(data, out); err != nil {
		out.Reset()
		return nil, err
	} else {
		return out, nil
	}

}
