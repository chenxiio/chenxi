package rpc

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"golang.org/x/xerrors"
)

type methodMeta struct {
	node  ast.Node
	ftype *ast.FuncType
}

type Visitor struct {
	Methods map[string]map[string]*methodMeta
	Include map[string][]string
}
type methodInfo struct {
	Name                                     string
	node                                     ast.Node
	Tags                                     map[string][]string
	Issupport                                bool
	Resultsnames, Namedresults               string
	NamedParams, ParamNames, Results, DefRes string
}

type strinfo struct {
	Name    string
	Methods map[string]*methodInfo
	Include []string
}

type meta struct {
	Infos   map[string]*strinfo
	Imports map[string]string
	OutPkg  string
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	st, ok := node.(*ast.TypeSpec)
	if !ok {
		return v
	}

	iface, ok := st.Type.(*ast.InterfaceType)
	if !ok {
		return v
	}
	if v.Methods[st.Name.Name] == nil {
		v.Methods[st.Name.Name] = map[string]*methodMeta{}
	}
	for _, m := range iface.Methods.List {
		switch ft := m.Type.(type) {
		case *ast.Ident:
			v.Include[st.Name.Name] = append(v.Include[st.Name.Name], ft.Name)
		case *ast.FuncType:
			v.Methods[st.Name.Name][m.Names[0].Name] = &methodMeta{
				node:  m,
				ftype: ft,
			}
		}
	}

	return v
}

// func main() {
// 	// latest (v1)
// 	if err := generate("../acontract/api", "api", "api", "../acontract/api/proxy_gen.go"); err != nil {
// 		fmt.Println("error: ", err)
// 	}

// 	// v0

// }

func typeName(e ast.Expr, pkg string) (string, error) {
	switch t := e.(type) {
	case *ast.SelectorExpr:
		return t.X.(*ast.Ident).Name + "." + t.Sel.Name, nil
	case *ast.Ident:
		pstr := t.Name
		if !unicode.IsLower(rune(pstr[0])) && pkg != "api" {
			//pstr = "" + pstr // todo src pkg name
		}
		return pstr, nil
	case *ast.ArrayType:
		subt, err := typeName(t.Elt, pkg)
		if err != nil {
			return "", err
		}
		return "[]" + subt, nil
	case *ast.StarExpr:
		subt, err := typeName(t.X, pkg)
		if err != nil {
			return "", err
		}
		return "*" + subt, nil
	case *ast.MapType:
		k, err := typeName(t.Key, pkg)
		if err != nil {
			return "", err
		}
		v, err := typeName(t.Value, pkg)
		if err != nil {
			return "", err
		}
		return "map[" + k + "]" + v, nil
	case *ast.StructType:
		if len(t.Fields.List) != 0 {
			return "", xerrors.Errorf("can't struct")
		}
		return "struct{}", nil
	case *ast.InterfaceType:
		if len(t.Methods.List) != 0 {
			return "", xerrors.Errorf("can't interface")
		}
		return "interface{}", nil
	case *ast.ChanType:
		subt, err := typeName(t.Value, pkg)
		if err != nil {
			return "", err
		}
		if t.Dir == ast.SEND {
			subt = "->chan " + subt
		} else {
			subt = "<-chan " + subt
		}
		return subt, nil
	default:
		return "", xerrors.Errorf("unknown type")
	}
}

func GetMothodInfo(path, pkg, outpkg, outfile string) (*meta, error) {
	fset := token.NewFileSet()
	apiDir, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	outfile, err = filepath.Abs(outfile)
	if err != nil {
		return nil, err
	}
	pkgs, err := parser.ParseDir(fset, apiDir, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return nil, err
	}

	ap := pkgs[pkg]

	v := &Visitor{make(map[string]map[string]*methodMeta), map[string][]string{}}
	ast.Walk(v, ap)

	m := &meta{
		OutPkg:  outpkg,
		Infos:   map[string]*strinfo{},
		Imports: map[string]string{},
	}

	for fn, f := range ap.Files {
		if strings.HasSuffix(fn, "gen.go") {
			continue
		}

		//fmt.Println("F:", fn)
		cmap := ast.NewCommentMap(fset, f, f.Comments)

		for _, im := range f.Imports {
			m.Imports[im.Path.Value] = im.Path.Value
			if im.Name != nil {
				m.Imports[im.Path.Value] = im.Name.Name + " " + m.Imports[im.Path.Value]
			}
		}

		for ifname, methods := range v.Methods {
			if _, ok := m.Infos[ifname]; !ok {
				m.Infos[ifname] = &strinfo{
					Name:    ifname,
					Methods: map[string]*methodInfo{},
					Include: v.Include[ifname],
				}
			}
			info := m.Infos[ifname]
			for mname, node := range methods {
				filteredComments := cmap.Filter(node.node).Comments()

				if _, ok := info.Methods[mname]; !ok {
					var params, pnames []string
					//var issp bool = false
					var issppams bool = false
					var isspresult bool = false
					for _, param := range node.ftype.Params.List {
						pstr, err := typeName(param.Type, outpkg)
						if err != nil {
							return nil, err
						}
						if strings.Contains(pstr, "context") {
							issppams = true
						}
						c := len(param.Names)
						if c == 0 {
							c = 1
						}

						for i := 0; i < c; i++ {
							pname := fmt.Sprintf("p%d", len(params))
							pnames = append(pnames, pname)
							params = append(params, pname+" "+pstr)
						}
					}

					results := []string{}
					resultsnames := []string{}
					namedresults := []string{}
					if node.ftype.Results == nil {
						goto next
					}
					for index, result := range node.ftype.Results.List {
						rs, err := typeName(result.Type, outpkg)
						if err != nil {
							return nil, err
						}
						results = append(results, rs)
						if rs != "error" {
							resultsnames = append(resultsnames, fmt.Sprintf("r%d", index))
							namedresults = append(namedresults, fmt.Sprintf("var r%d %s", index, rs))
						} else {
							isspresult = true
							resultsnames = append(resultsnames, "err")
							namedresults = append(namedresults, "var err error\n")
						}
					}
				next:

					defRes := ""
					if len(results) > 1 {
						defRes = results[0]
						switch {
						case defRes[0] == '*' || defRes[0] == '<', defRes == "interface{}":
							defRes = "nil"
						case defRes == "bool":
							defRes = "false"
						case defRes == "string":
							defRes = `""`
						case defRes == "int", defRes == "int64", defRes == "uint64", defRes == "uint":
							defRes = "0"
						default:
							defRes = "*new(" + defRes + ")"
						}
						defRes += ", "
					}

					info.Methods[mname] = &methodInfo{
						Name:         mname,
						node:         node.node,
						Tags:         map[string][]string{},
						Issupport:    issppams && isspresult,
						Resultsnames: strings.Join(resultsnames, ", "),
						Namedresults: strings.Join(namedresults, "\n"),
						NamedParams:  strings.Join(params, ", "),
						ParamNames:   strings.Join(pnames, ", "),
						Results:      strings.Join(results, ", "),
						DefRes:       defRes,
					}
				}

				// try to parse tag info
				if len(filteredComments) > 0 {
					tagstr := filteredComments[len(filteredComments)-1].List[0].Text
					tagstr = strings.TrimPrefix(tagstr, "//")
					tl := strings.Split(strings.TrimSpace(tagstr), " ")
					for _, ts := range tl {
						tf := strings.Split(ts, ":")
						if len(tf) != 2 {
							continue
						}
						if tf[0] != "perm" { // todo: allow more tag types
							continue
						}
						info.Methods[mname].Tags[tf[0]] = tf
					}
				}
			}
		}
	}
	return m, nil
}

func Generate(path, pkg, outpkg, outfile string) error {

	/*jb, err := json.MarshalIndent(Infos, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jb))*/
	m, err := GetMothodInfo(path, pkg, outpkg, outfile)
	if err != nil {
		return err
	}

	w, err := os.OpenFile(outfile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	err = doTemplate(w, m, `// Code generated by gitee.com/magic-mountain/asmb/utils/rpc/gen.go. 
package {{.OutPkg}}

import (
{{range .Imports}}	{{.}}
{{end}}
)
`)
	if err != nil {
		return err
	}

	err = doTemplate(w, m, `
{{range .Infos}}
type {{.Name}}Struct struct {
{{range .Include}}{{.}}Struct
{{end}}
	Internal struct {
{{range .Methods}}
		{{.Name}} func({{.NamedParams}}) ({{.Results}}) `+"`"+`{{range .Tags}}{{index . 0}}:"{{index . 1}}"{{end}}`+"`"+`{{end}}
	}
}

type {{.Name}}Stub struct {
{{range .Include}}
	{{.}}Stub
{{end}}
}
{{end}}

{{range .Infos}}
{{$name := .Name}}
{{range .Methods}}
func (s *{{$name}}Struct) {{.Name}}({{.NamedParams}}) ({{.Results}}) {
	return s.Internal.{{.Name}}({{.ParamNames}})
}

{{end}}
{{end}}

{{range .Infos}}var _ {{.Name}} = new({{.Name}}Struct)
{{end}}

`)
	return err
}

func doTemplate(w io.Writer, info interface{}, templ string) error {
	t := template.Must(template.New("").
		Funcs(template.FuncMap{}).Parse(templ))

	return t.Execute(w, info)
}
