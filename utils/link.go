/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"fmt"
	"path"
	"strings"
)

type Link struct {
	anchor string
	path   string
	tag    string
}

func (l Link) IsValid() bool {
	return l.tag != "" || l.anchor != "" || l.path != ""
}

func NewLink(path string, anchor string) Link {
	return Link{anchor: anchor, path: path}
}

func NewTagLink(tag string) Link {
	return Link{tag: tag}
}

func (l Link) Anchor() string {
	if l.tag != "" {
		return l.tag
	}
	return l.anchor
}

func (l Link) Path() string {
	return l.path
}

func (l Link) IsTag() bool {
	return l.tag != ""
}

func (l Link) Tag() string {
	return l.tag
}

func (l Link) String() string {
	if l.tag != "" {
		return "#" + l.tag
	}
	if l.anchor == "" {
		return l.path
	}
	return l.path + "#" + l.anchor
}

func (l Link) Abs(base string, global bool) (Link, error) {
	var r Link

	if l.tag != "" {
		return l, nil
	}
	if path.IsAbs(l.path) {
		return l, nil
	}

	if global {
		r.path = base
	}
	if l.path != "" {
		b := base
		if strings.HasPrefix(base, "/") {
			b = base[1:]
		}
		if strings.HasPrefix(path.Join(path.Dir(b), l.path), ".") {
			return r, fmt.Errorf("link %q outside of tree %q", l.path, base)
		}
		r.path = path.Join(path.Dir(base), l.path)
	}

	r.anchor = l.anchor
	return r, nil
}

func ExtendAnchor(a string, ns string) (string, error) {
	if a == "" || ns == "" {
		return a, nil
	}
	b := ns
	if strings.HasPrefix(ns, "/") {
		b = ns[1:]
		if strings.HasPrefix(path.Join(b, a), ".") {
			return a, fmt.Errorf("anchor %q outside of tree %q", a, ns)
		}
	}
	return path.Join(ns, a), nil
}

func ParseAbsoluteLink(link string, base string, global bool) (Link, error) {
	l, err := ParseLink(link)
	if err != nil {
		return l, err
	}
	return l.Abs(base, global)
}

func ParseLink(link string, asAnchor ...bool) (Link, error) {
	var r Link

	comps := strings.Split(link, "#")
	if len(comps) > 2 {
		return r, fmt.Errorf("invalid link target %q", link)
	}
	p := comps[0]
	a := ""
	if len(comps) == 1 {
		if len(asAnchor) == 0 || !asAnchor[0] {
			if p == "" {
				return r, fmt.Errorf("invalid link target %q", link)
			}
			r.path = p
			return r, nil
		}
		p = ""
		a = comps[0]
	} else {
		a = comps[1]
	}

	if p == "" {
		if path.IsAbs(a) {
			r.tag = a
		} else {
			r.anchor = a
		}
	} else {
		if path.IsAbs(a) {
			return r, fmt.Errorf("invalid abolute anchor %q (path not possible)", comps[1])
		} else {
			r.path = p
			r.anchor = a
		}
	}
	return r, nil
}
