/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package scanner

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/mandelsoft/mdgen/labels"
	"github.com/mandelsoft/mdgen/labels/format"
)

func init() {
	Tokens.Register("numberrange", ParseNumberRange)
}

// TODO move functionality from parsing to resolution

func ParseNumberRange(p Parser, e Element) (Element, error) {
	var err error

	if len(e.Tags()) == 0 {
		return nil, e.Errorf("at least tag for range type required")
	}

	var limit int64 = -1
	var master string
	var abbrev string

	parts := e.Tags()

	for i, p := range parts[1:] {
		off := strings.Index(p, "=")
		if off < 0 {
			return nil, e.Errorf("argument %d [%s] requires assignment", i+1, p)
		}
		f := p[:off]
		v := p[off+1:]
		if f == "" {
			return nil, e.Errorf("argument %d [%s] required nin-empty field name", i+1, p)
		}
		switch f {
		case "master":
			comps := strings.Split(v, ":")
			master = comps[0]
			switch len(comps) {
			case 1:
			case 2:
				if !strings.HasPrefix(comps[1], "#") {
					return nil, e.Errorf("master limit must start with #")
				}
				limit, err = strconv.ParseInt(comps[1][1:], 10, 8)
				if err != nil {
					return nil, e.Errorf("master limit must be number, but found %s: %s", comps[1], err)
				}
				if limit < 0 {
					return nil, e.Errorf("invalid master limit %d", limit)
				}
				limit--
			default:
				return nil, e.Errorf("expected master spec <name>[:<level limit>]")
			}
		case "abbrev":
			abbrev = v
		default:
			return nil, e.Errorf("argument %d [%s] uses unknown field %s (use master or abbrev)", i+1, p, f)
		}
	}

	comps := strings.Split(parts[0], ":")

	lvl := -1
	lvl_index := 0
	typ_index := 0
	switch len(comps) {
	case 1:
	case 2:
		if strings.HasPrefix(comps[lvl_index], "#") {
			lvl_index = 1
		} else {
			typ_index = 1
		}
	case 3:
		lvl_index = 2
	default:
		return nil, e.Errorf("expected label spec <name>[:<type>][:<level>]")
	}

	if lvl_index > 0 {
		if !strings.HasPrefix(comps[lvl_index], "#") {
			return nil, e.Errorf("section level must start with #")
		}
		i, err := strconv.ParseInt(comps[lvl_index][1:], 10, 32)
		if err != nil {
			return nil, e.Errorf("number required for section level, but found %q: %s", comps[lvl_index], err.Error())
		}
		lvl = int(i) - 1
	}

	var rule labels.Rule
	name := comps[0]
	sep := ""

	if typ_index > 0 {
		typ := comps[typ_index]
		r, s := utf8.DecodeRuneInString(comps[typ_index])
		if strings.Contains(format.Separators, string(r)) {
			typ = typ[s:]
			sep = string(r)
		}

		switch typ {
		case "numbered":
			rule = labels.NewNumbered(name, lvl)
		case "void":
			rule = labels.NewVoid(name, lvl)
		default:
			f, err := format.FormatFor(typ)
			if err != nil {
				return nil, e.Errorf("unknown label type %q: %s", typ, err)
			}
			rule = labels.NewFreeForm(name, f, lvl)
		}
	}
	err = p.State.Container.SetLabelRule(&e.location, name, abbrev, sep, rule, lvl)
	if err != nil {
		return nil, e.Errorf("%s", err)
	}
	if master != "" {
		err = p.State.Container.SetLabelMaster(name, master, sep, int(limit))
		if err != nil {
			return nil, e.Errorf("%s", err)
		}
	}
	return p.tokenizer.NextElement()
}
