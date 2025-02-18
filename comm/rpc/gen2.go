package rpc

import (
	"os"
)

func Generate2(path, pkg, outpkg, outfile string) error {

	m, err := GetMothodInfo(path, pkg, outpkg, outfile)
	if err != nil {
		return err
	}
	w, err := os.OpenFile(outfile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	defer w.Close()
	if err != nil {
		return err
	}
	err = doTemplate(w, m, `// Code generated by gitee.com/magic-mountain/asmb/utils/rpc/gen.go.
	package {{.OutPkg}}
	
	import (
	{{range .Imports}}	{{.}}
	{{end}}
	)
	
	{{range .Infos}}
	type {{.Name}}Service struct {
	{{range .Include}}{{.}}Service
	{{end}}
		Clist *concurrent.ConcurrentList
		Lc    {{.Name}}
	}
	{{end}}
	
	{{range .Infos}}
	{{$name := .Name}}
	{{range .Methods}}
	func (s *{{$name}}Service) {{.Name}}({{.NamedParams}}) ({{.Results}}) {    
		{{if .Issupport}}	
		if s.Lc != nil {
			return s.Lc.{{.Name}}({{.ParamNames}})
		}
		ret := make(chan struct{})
		{{.Namedresults}}		
		go func() {
			if s.Clist==nil || s.Clist.Front() == nil {
				err = errors.New("list is empty")
				ret <- struct{}{}
				ret = nil
				return
			}
			ch := make(chan struct{}, 3)
			Front := s.Clist.Front()
			for element := s.Clist.Front(); element != nil; element = element.Next() {
				if ret == nil || element.Next() == Front {
					break
				}
				ch <- struct{}{}
				go func(e *list.Element) {
					rs := e.Value.(*rpc.RpcService)
					t := time.Now()
					{{.Resultsnames}} = rs.V.({{$name}}).{{.Name}}({{.ParamNames}})				
					rs.Reaction_time = time.Since(t)
					if err != nil {
						// 队尾
						s.Clist.MoveToBack(e)
						fmt.Printf("warning:%s,%v\n", rs.Url, err)
					} else {
						if ret != nil {
							// 队首
							if rs.Reaction_time < time.Millisecond*100 {
								s.Clist.MoveToFront(e)
							}
							ret <- struct{}{}
							ret = nil
						} else {
							// 向后一位
							if rs.Reaction_time > time.Second*1 {
								s.Clist.MoveAfter(e, e.Next())
							}
						}
					}
					<-ch
					if len(ch) == 0 && ret != nil {
						err = errors.New("all services are unavailable")
						ret <- struct{}{}
						ret = nil
					}
				}(element)
			}			
		}()
	
		select {
		case <-p0.Done():
			err = errors.New("ctx canceled")
			return {{.Resultsnames}}
		case  <-ret:		
			return {{.Resultsnames}}	
		}	
		{{else}}
		panic("not support")
		{{end}}
	}
	{{end}}
	{{end}}
	
	{{range .Infos}}var _ {{.Name}} = new({{.Name}}Service)
	{{end}}	
	`)

	if err != nil {
		return err
	}
	return nil
}

// func (s *IwalletService) List(p0 context.Context) ([]address.Address, error) {
// 	ret := make(chan struct{})
// 	var r0 []address.Address
// 	var err error
// 	go func() {
// 		ch := make(chan struct{}, 3)
// 		if s.Clist.Front() == nil {
// 			err = errors.New("list is empty")
// 			ret <- struct{}{}
// 			ret = nil
// 			return
// 		}
// 		for element := s.Clist.Front(); element != nil; element = element.Next() {
// 			ch <- struct{}{}
// 			if ret == nil {
// 				break
// 			}
// 			go func(e *list.Element) {
// 				rs := e.Value.(*rpc.RpcService)
// 				t := time.Now()
// 				r0, err = rs.V.(Iwallet).List(p0)
// 				rs.Reaction_time = time.Since(t)
// 				if err != nil {
// 					// 队尾
// 					s.Clist.MoveToBack(e)
// 					fmt.Printf("%s,%v\n", rs.Url, "error http invod")
// 				} else {
// 					if ret != nil {
// 						// 队首
// 						if rs.Reaction_time < time.Millisecond*100 {
// 							s.Clist.MoveToFront(e)
// 						}
// 						ret <- struct{}{}
// 						ret = nil
// 					} else {
// 						// 向后一位
// 						if rs.Reaction_time > time.Second*1 {
// 							s.Clist.MoveAfter(e, e.Next())
// 						}
// 					}
// 				}
// 				<-ch
// 				if len(ch) == 0 && e.Next() == nil && ret != nil {
// 					err = errors.New("All services are unavailable")
// 					ret <- struct{}{}
// 					ret = nil
// 				}
// 			}(element)
// 		}
// 	}()
// 	select {
// 	case <-p0.Done():
// 		err = errors.New("ctx canceled")
// 		return r0, err
// 	case <-ret:
// 		return r0, err
// 	}
// }
