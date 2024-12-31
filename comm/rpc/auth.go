package rpc

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"

	logging "github.com/ipfs/go-log/v2"
)

var AuthNewCmd = &cli.Command{
	Name:  "authnew",
	Usage: "创建一个新token",

	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "perm, p",                                                             // 配置名称
			Value: "read",                                                                // 缺省配置值
			Usage: "permission to assign to the token, one of: read, write, sign, admin", // 配置描述
		},
		cli.StringFlag{
			Name:  "secret, s", // 配置名称
			Value: "123456",    // 缺省配置值
			Usage: "jwt 加密密码,请修改默认密码",
			//Destination: &Cfg.APISecret,
		},
	},
	Action: func(cctx *cli.Context) error {

		perm := cctx.String("perm")
		var ap []auth.Permission = nil
		if perm != "" {
			idx := 0
			for i, p := range AllPermissions {
				if auth.Permission(perm) == p {
					idx = i + 1
				}
			}
			ap = AllPermissions[:idx]
		}

		//	Cfg.APISecret = cctx.String("secret")
		// if idx == 0 {
		// 	return fmt.Errorf("--perm flag has to be one of: %s", AllPermissions)
		// }

		// slice on [:idx] so for example: 'sign' gives you [read, write, sign]

		token, err := AuthNew(context.Background(), ap, cctx.String("s"))
		if err != nil {
			return err
		}

		// TODO: Log in audit log when it is implemented

		fmt.Println(string(token))
		return nil

	},
}

type Config struct {
	Token     string
	APISecret string
}

var Cfg = &Config{APISecret: "123456", Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.ae22SYWjZ_RJRRhfWDpVFzWThu_6EQ-iBgAn8vdrR-w"}

var log = logging.Logger("rpc")

type jwtPayload struct {
	Allow []auth.Permission
}

const (
	// When changing these, update docs/md too
	PermNone  auth.Permission = "none"
	PermRead  auth.Permission = "read" // default
	PermWrite auth.Permission = "write"
	PermSign  auth.Permission = "sign"  // Use wallet keys for signing
	PermAdmin auth.Permission = "admin" // Manage permissions
)

var AllPermissions = []auth.Permission{PermNone, PermRead, PermWrite, PermSign, PermAdmin}
var DefaultPerms = []auth.Permission{PermNone}

func AuthVerify(ctx context.Context, token string) ([]auth.Permission, error) {

	var payload jwtPayload
	if _, err := jwt.Verify([]byte(token), jwt.NewHS256([]byte(Cfg.APISecret)), &payload); err != nil {
		return nil, xerrors.Errorf("JWT Verification failed: %w", err)
	}

	return payload.Allow, nil
}

func AuthNew(ctx context.Context, perms []auth.Permission, pwd string) ([]byte, error) {
	p := jwtPayload{
		Allow: perms, // TODO: consider checking validity
	}

	return jwt.Sign(&p, jwt.NewHS256([]byte(pwd)))
}

func AuthHeader() http.Header {
	if len(Cfg.Token) != 0 {
		headers := http.Header{}
		headers.Add("Authorization", "Bearer "+Cfg.Token)
		return headers
	}
	log.Warn("API Token not set and requested, capabilities might be limited.")
	return nil
}
func AuthHeaderWtoken(t string) http.Header {
	headers := http.Header{}
	if len(t) > 0 {
		headers.Add("Authorization", "Bearer "+t)
	}

	return headers
}

func GetRpcClient(t, url, name string, out interface{}, timeout time.Duration) (func(), error) {
	closer, err := jsonrpc.NewMergeClient(context.Background(), url+"/"+name+"/v0", name, []interface{}{out}, AuthHeaderWtoken(t), jsonrpc.WithTimeout(timeout))
	if err != nil {
		return nil, err
	}
	return closer, nil
}

func RegisterRpc(v, outInternal interface{}, out interface{}, name string) *jsonrpc.RPCServer {

	rpcServer := jsonrpc.NewServer()

	// //c = &outChain
	auth.PermissionedProxy(AllPermissions, DefaultPerms, v, outInternal)
	rpcServer.Register(name, out)
	//rpcServer.Register("asmb", acontext.Chain)

	return rpcServer
}

type CorsHandler struct {
	Origin string

	http.HandlerFunc
}

func (h *CorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.Method)

	//fmt.Println(r.Method, r.Header["Origin"])
	w.Header().Add("Access-Control-Request-Method", "POST, GET, PUT, DELETE, OPTIONS")
	if len(r.Header["Origin"]) > 0 {
		if strings.Contains(h.Origin, r.Header["Origin"][0]) || h.Origin == "*" {
			w.Header().Add("Access-Control-Allow-Origin", r.Header["Origin"][0])
		}
	}

	// for k, _ := range r.Header {
	// 	w.Header().Add("Access-Control-Allow-Headers", k)
	// }
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")

	//.Headers("Access-Control-Request-Method", "*", "Access-Control-Allow-Origin", "*").al
	if r.Method == http.MethodOptions {

		//		Access-Control-Allow-Origin
		w.WriteHeader(http.StatusOK)
		return
	}

	h.HandlerFunc(w, r)
}
