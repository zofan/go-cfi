package cfi

import (
	"fmt"
	"github.com/zofan/go-fwrite"
	"github.com/zofan/go-req"
	"github.com/zofan/go-scraper"
	"github.com/zofan/go-xmlre"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	typesRe       = xmlre.Compile(`<h2>Category (.*?) - (.*?)</h2>`)
	subtypesRe    = xmlre.Compile(`<td class="left pl-20">(.*?)</td><td class="left">`)
	subtypeHeadRe = xmlre.Compile(`<h1 class="page-heading">Attributes for CFI code group : (.*?)</h1>`)
	bracketsRe    = xmlre.Compile(`\s*\[[^\[\]]+]\s*`)
	sectionRe     = xmlre.Compile(`<section class="mt-20">(.*?)</section>`)
	sectionHeadRe = xmlre.Compile(`<h2>(\d)<sup>\w+</sup>attribute of group \w+ : (.*?)</h2>`)
	sectionKVRe   = xmlre.Compile(`<div class="columns small-1 medium-1 large-1">(.*?)</div><div class="columns small-1 medium-5 large-10">(.*?)</div>`)
)

func Update() error {
	r := req.New(req.DefaultConfig)

	var types = make(map[string]string)
	var attributes = make(map[string]string)

	{
		resp := r.Get(`http://www.iotafinance.com/en/Classification-of-Financial-Instrument-codes-CFI-ISO-10962.html#cache=99999h`)
		if resp.Error() != nil {
			return resp.Error()
		}

		body := scraper.ReplaceEntities(string(resp.ReadAll()))

		{
			matches := typesRe.FindAllStringSubmatch(body, -1)
			for _, m := range matches {
				types[scraper.ClearHtml(m[1])] = scraper.ClearHtml(m[2])
			}
		}

		{
			matches := subtypesRe.FindAllStringSubmatch(body, -1)
			for _, m := range matches {
				types[scraper.ClearHtml(m[1])] = ``
			}
		}
	}

	for t := range types {
		fmt.Println(t)

		if len(t) != 2 {
			continue
		}

		time.Sleep(time.Second)

		resp := r.Get(`http://www.iotafinance.com/en/Attributes-CFI-Codes-Group-` + t + `.html#cache=99999h`)
		if resp.Error() != nil {
			return resp.Error()
		}

		body := scraper.ReplaceEntities(string(resp.ReadAll()))

		headMatch := subtypeHeadRe.FindStringSubmatch(body)
		types[t] = bracketsRe.ReplaceAllString(headMatch[1], ``)

		matches := sectionRe.FindAllStringSubmatch(body, -1)

		for _, match := range matches {
			headMatch := sectionHeadRe.FindStringSubmatch(match[1])
			kvMatches := sectionKVRe.FindAllStringSubmatch(match[1], -1)

			attributes[t+headMatch[1]] = headMatch[2]

			for _, m := range kvMatches {
				attributes[t+scraper.ClearHtml(headMatch[1])+scraper.ClearHtml(m[1])] = scraper.ClearHtml(m[2])
			}
		}
	}

	var code []string

	code = append(code, `package cfi`)
	code = append(code, ``)
	code = append(code, `var (`)
	code = append(code, `	types = `+fmt.Sprintf(`%#v`, types))
	code = append(code, ``)
	code = append(code, `	attributes = `+fmt.Sprintf(`%#v`, attributes))
	code = append(code, `)`)
	code = append(code, ``)

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	_ = fwrite.WriteRaw(dir+`/dict.go`, []byte(strings.Join(code, "\n")))

	return nil
}
