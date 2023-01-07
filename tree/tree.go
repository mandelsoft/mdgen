/*
 * SPDX-FileCopyrightText: 2023 Mandelsoft.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package tree

import (
	"fmt"
	"path"
	"strings"

	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/vfs"

	"github.com/mandelsoft/mdgen/labels"
	"github.com/mandelsoft/mdgen/scanner"
	"github.com/mandelsoft/mdgen/statements/section"
	utils "github.com/mandelsoft/mdgen/utils"
)

type Tree = *tree

type tree struct {
	documents  map[string]scanner.Document
	resolution *Resolution
}

func NewTree() Tree {
	return &tree{
		documents: map[string]scanner.Document{},
	}
}

func (t *tree) Print(gap string) {
	ngap := gap + "  "
	fmt.Printf("%sTREE:\n", gap)
	fmt.Printf("%s  documents:\n", gap)
	for _, n := range utils.StringMapKeys(t.documents) {
		d := t.documents[n]
		d.Print(ngap)
	}
}

func (t *tree) GetDocument(refpath string) *DocumentInfo {
	return t.resolution.documents[refpath]
}

func ForFolder(path string, fss ...vfs.FileSystem) (Tree, error) {
	fs := utils.OptionalDefaulted(osfs.New(), fss...)
	tr := NewTree()

	var err error

	if ok, nerr := vfs.IsFile(fs, path); nerr == nil && ok {
		err = scanFile(tr, path, fs, "/")
	} else {
		err = scanDir(tr, path, fs, "/")

	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return tr, nil
}

func scanDir(tr Tree, p string, fs vfs.FileSystem, refpath string) error {
	fmt.Printf("%s: scanning %s\n", refpath, p)
	list, err := vfs.ReadDir(fs, p)
	if err != nil {
		return err
	}
	for _, f := range list {
		if f.IsDir() {
			err = scanDir(tr, path.Join(p, f.Name()), fs, path.Join(refpath, f.Name()))
		} else {
			if strings.HasSuffix(f.Name(), ".mdg") {
				err = scanFile(tr, path.Join(p, f.Name()), fs, path.Join(refpath, f.Name()[:len(f.Name())-4]))
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func scanFile(tr Tree, p string, fs vfs.FileSystem, refpath string) error {
	fmt.Printf("%s: reading %s\n", refpath, p)
	file, err := fs.Open(p)
	if err != nil {
		return err
	}

	pa := scanner.NewParser(p, refpath, file)

	doc, err := pa.Parse()
	if err != nil {
		return err
	}

	tr.documents[refpath] = doc
	return nil
}

func (t *tree) Resolve() error {
	res := NewResolution(t.documents)
	t.resolution = res

	fmt.Printf("resolve blocks...\n")
	err := t.ResolveBlocks(res)
	if err != nil {
		return err
	}

	fmt.Printf("resolve structure...\n")
	err = t.ResolveStructural(res)
	if err != nil {
		return err
	}

	fmt.Printf("resolve number ranges...\n")
	err = t.ResolveNumberRanges(res)
	if err != nil {
		return err
	}

	fmt.Printf("reference list:\n")
	for l, r := range res.refindex {
		fmt.Printf("  %s: %s#%s\n", l, r.GetRefPath(), r.Anchor())
	}

	fmt.Printf("resolve values...\n")
	err = t.ResolveValues(res)
	if err != nil {
		return err
	}

	return nil
}

func (t *tree) ResolveBlocks(res *Resolution) error {
	for _, di := range res.documents {
		di.document.RequestNumberRanges(di.context)
		err := di.Walk(scanner.Resolve[scanner.Register](di.context))
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *tree) ResolveStructural(res *Resolution) error {
	found := map[scanner.Document]*docref{}
	for _, di := range res.documents {
		for l, ref := range di.context.docrefs.links {
			var target scanner.Document
			if l.IsTag() {
				ri := res.refindex[l]
				if ri == nil {
					return fmt.Errorf("%s: structural tag reference %q cannot be resolved", ref.location, di.Source())
				}
				target = ri.Context().GetDocument()
			} else {
				target = t.documents[l.Path()]
				if target == nil {
					return fmt.Errorf("%s: structural document reference %q cannot be resolved", ref.location, di.Source())
				}
			}
			ti := res.documents[target.GetRefPath()]
			if ti.structinfo != nil {
				return fmt.Errorf("%s: duplicate structural usage of document %s: %s", ref.location, target.Source(), found[target].location)
			}
			fmt.Printf("%s: found structural usage in %s\n", target.Source(), di.Source())
			sect := false
			for _, n := range target.GetNodes() {
				if _, ok := n.(section.SectionNode); ok {
					if sect {
						return fmt.Errorf("%s: structural document %s may contain only one top level section\n", ref.location, target.Source())
					}
					sect = true
				}
			}
			ti.structinfo = NewStructInfo(di)
			ref.docinfo = ti
			found[target] = ref
		}
	}
	return nil
}

func (t *tree) ResolveNumberRanges(res *Resolution) error {
	for _, di := range res.documents {
		err := t.resolveDocumentOrder(res, di, utils.History{})
		if err != nil {
			return err
		}
	}
	for _, di := range res.documents {
		if di.IsRoot() {
			err := t.resolveNumberRanges(di)
			if err != nil {
				return err
			}
		}
	}

	for _, di := range res.documents {
		if di.IsRoot() {
			resolved := map[string]bool{}
			for _, nr := range di.rootinfo.ranges {
				t.resolveLabels(di, nr, nil, resolved)
			}
		}
	}
	return nil
}

func (t *tree) resolveDocumentOrder(res *Resolution, di *DocumentInfo, hist utils.History) error {
	if di.rootinfo != nil {
		return nil
	}

	hist, cycle := hist.Add(di.Source())
	if cycle != nil {
		return fmt.Errorf("structural cycle %s", hist)
	}

	var root *RootInfo
	master := di.structinfo
	if master != nil {
		err := t.resolveDocumentOrder(res, res.documents[master.docinfo.GetRefPath()], hist)
		if err != nil {
			return err
		}
		root = res.documents[master.docinfo.GetRefPath()].rootinfo
		fmt.Printf("%s: sub structure for %s\n", di.Source(), root.docinfo.Source())
	} else {
		fmt.Printf("%s: root document\n", di.Source())
		root = NewRootInfo(di)
	}
	di.rootinfo = root
	return nil
}

func (t *tree) resolveLabels(di *DocumentInfo, nr *NumberRangeInfo, hist utils.History, resolved map[string]bool) error {
	if resolved[nr.Type()] {
		return nil
	}
	resolved[nr.Type()] = true
	hist, cycle := hist.Add(nr.Type())
	if len(cycle) > 0 {
		return fmt.Errorf("number range dependency cycle: %s", cycle)
	}
	if nr.master != "" {
		master := di.rootinfo.ranges[nr.master]
		if master == nil {
			return fmt.Errorf("number range %s: dependency to non-existing number range %s", nr.Type(), nr.master)
		}
		err := t.resolveLabels(di, master, hist, resolved)
		if err != nil {
			return err
		}
	}
	fmt.Printf("%s: generate %s labels\n", di.Source(), nr.Type())
	nr.CreateLabels(nil)
	return nil
}

func (t *tree) resolveNumberRanges(di *DocumentInfo) error {
	root := di.rootinfo
	fmt.Printf("%s: found numberranges: %s\n", di.Source(), strings.Join(utils.SortedMapKeys(di.context.numberranges), ", "))
	rules := di.document.GetLabelRules()
	for typ := range di.context.numberranges {
		outer := root.ranges[typ]
		l := rules[typ]
		if outer == nil {
			var r labels.Rule
			master := ""
			limit := -1
			sep := ""
			lvl := 0
			abbrev := ""
			if l != nil {
				if l.Rule != nil {
					r = l.Rule
				}
				sep = l.Separator
				if l.Level >= 0 {
					lvl = l.Level
				}
				if l.Abbrev != "" {
					abbrev = l.Abbrev
				}
				master = l.Master
				limit = l.Limit
			}
			if r == nil {
				r = labels.NewNumbered(typ, lvl)
			} else {
				if r.Level() < 0 {
					r = r.WithLevel(lvl)
				}
			}
			var provider func() scanner.HierarchyLabel
			if master != "" {
				provider = func() scanner.HierarchyLabel {
					nr := t.resolution.current.context.GetNumberRange(master)
					l := nr.Actual()
					if limit >= 0 {
						for l.Level() > limit {
							l = l.Parent()
						}
					}
					return l
				}
				fmt.Printf("%s: initialize number range %s with master %s: %s\n", di.Source(), typ, master, r.Format())
			} else {
				fmt.Printf("%s: initialize number range %s: %s\n", di.Source(), typ, r.Format())
			}

			var loc *scanner.Location
			if l != nil {
				loc = l.Location
			}
			nr := &NumberRangeInfo{NumberRange: scanner.NewNumberRange(typ, abbrev, provider), master: master, location: loc}
			nr.SetRule(sep, r)
			root.ranges[typ] = nr
		} else {
			if l != nil {
				return fmt.Errorf("%s: document level %s numberrange not possible for sub document", di.Source(), typ)
			}
			fmt.Printf("%s: found reused number range %s from root document\n", di.Source(), typ)
		}
	}

	for _, nr := range root.ranges {
		err := checkCycle(di, nr, nil)
		if err != nil {
			return err
		}
	}
	fmt.Printf("  resolve %s\n", di.GetRefPath())
	err := di.Walk(scanner.Resolve[scanner.LabelResolver](di.context))
	if err != nil {
		return err
	}

	for _, ref := range di.context.docrefs.order {
		err := t.resolveNumberRanges(ref.docinfo)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkCycle(di *DocumentInfo, nr *NumberRangeInfo, hist utils.History) error {
	loc := di.Source()

	if nr.location != nil {
		loc = nr.location.String()
	}
	hist, cycle := hist.Add(nr.Type())
	if len(cycle) > 0 {
		return fmt.Errorf("%s: dependency cycle for number ranges: %s", loc, cycle)
	}
	if nr.master == "" {
		return nil
	}
	master := di.rootinfo.ranges[nr.master]
	if master == nil {
		return fmt.Errorf("%s: unkown master number range %s for %s", loc, nr.master, nr.Type())
	}
	return checkCycle(di, master, hist)
}

func (t *tree) ResolveValues(res *Resolution) error {
	var found map[string][]unresolved
	last := -1

	for last != 0 {
		found = map[string][]unresolved{}
		cnt := 0
		for rp, di := range t.resolution.documents {
			if di.document.IsTemplate() {
				continue
			}
			di.context.unresolved = nil
			err := di.Walk(scanner.Resolve[scanner.ValueResolver](di.context))
			if err != nil {
				return err
			}
			found[rp] = di.context.unresolved
			cnt += len(di.context.unresolved)
			if len(di.context.unresolved) > 0 {
				fmt.Printf("%s: found %d unresolved nodes:\n", di.Source(), len(di.context.unresolved))
				for _, u := range di.context.unresolved {
					fmt.Printf("   %s: %s\n", u.Location(), u.err)
				}
			}
		}
		if last > 0 {
			if cnt > last {
				panic("oops: growing number of problems ")
			}
			if cnt == last {
				break
			}
		}
		last = cnt
	}
	if last > 0 {
		msg := ""
		for n, l := range found {
			msg += fmt.Sprintf("%s: found %d unresolved nodes:\n", res.documents[n].Source(), len(l))
			for _, u := range l {
				msg += fmt.Sprintf("   %s: %s\n", u.Location(), u.err)
			}
		}
		return fmt.Errorf("%s", msg)
	}
	return nil
}

func (t *tree) Emit(tw TreeWriter) error {
	for _, di := range t.resolution.documents {
		if di.document.IsTemplate() {
			continue
		}
		w, target, err := tw.Document(di.document.GetRefPath())
		if err != nil {
			return err
		}
		fmt.Printf("writing %s\n", di.document.GetRefPath())
		/*
			toc := scanner.DocTOCIds(di.context, scanner.SECTION_TYPE)
				for _, c := range toc {
				fmt.Printf("%*s-%s\n", c.Level(), "", c.Id())
			}
		*/
		err = di.Emit(scanner.NewWriter(w), target)
		err2 := w.Close()
		if err2 != nil {
			return err2
		}
		if err != nil {
			return err
		}
	}
	return nil
}
