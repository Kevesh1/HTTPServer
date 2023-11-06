package router
 
import (
    "context"
    "fmt"
    "net/http"
    "regexp"
    "strings"
)
 
type route struct {
    method       string           //HTTP method
    pattern      *regexp.Regexp   //Match the request path
    innerHandler http.HandlerFunc //Handler
    paramKeys    []string         //Parameters in order in the path
}
 
//router will store all routes defined for our server, essentially doing all the routing
 
type router struct {
    routes []route
}
 
func newRouter() *router {
    return &router{routes: []route{}}
}
 
// handler for logging for the route
func (r *route) handler(res http.ResponseWriter, req *http.Request) {
    request := fmt.Sprint(req.Method, " ", req.URL)
    fmt.Println("received ", request)
	start := time.Now()
    r.innerHandler(utils.NewResponseWriter(res), req)
	res.Time = time.Since(start).Milliseconds(
	fmt.Printf("%s resolved with %s\n", request, res)
	)
}
 
// Function for adding a route to the routes struct.
// Extract parameters from the request and replace them with regex pattern
// ^ matches the start of a line or string
// $ matches the end of a line or string
// () Groups characters together to apply metacharacters to the entire group
// [] Defines a character class, allowing you to match any one of the characters inside the brackets
// + Matches one or more occurrences of the preceding character or group
func (r *router) addRoute(method, endpoint string, handler http.HandlerFunc) {
    //handle path parameters
    pathParamPattern := regexp.MustCompile(":([a-z]+)")
    matches := pathParamPattern.FindAllStringSubmatch(endpoint, -1)
    paramKeys := []string{}
    if len(matches) > 0 {
        // replace path parameter definition with regex pattern to capture any string
        endpoint = pathParamPattern.ReplaceAllLiteralString(endpoint, "([^/]+)")
        // store the names of path parameters
        for i := 0; i < len(matches); i++ {
            paramKeys = append(paramKeys, matches[i][1])
        }
    }
 
    route := route{method, regexp.MustCompile("^" + endpoint + "$"), handler, paramKeys}
    r.routes = append(r.routes, route)
}
 
func (r *router) GET(pattern string, handler http.HandlerFunc){
    r.addRoute(http.MethodGet, pattern, handler)
}
 
// EXAMPLE:
// route{"GET", "/chat/([^/]+)/user/([^/]+)", someHandler, ["chatid", "userid"]}
router.GET("/chat/:chatid/user/:userid", someHandler)
 
 
// Document this one
// Loop through the routes and if a matching method (GET or POST) is found
// send the response and updated context to the route's handler. If no match
// is found but paths matches we return a list of allowed methods
 
func (r *router) ServeHTTP(res http.ResponseWriter, req *http.Request){
    allow := []string{}
    for _, route := range r.routes {
        matches := route.pattern.FindStringSubmatch(req.URL.Path)
        if len(matches) > 0 {
            if req.Method == route.method{
                route.handler(res, buildContext(req, route.paramKeys, matches[1:]))
                return
            }
            allow = append(allow, route.method)
        }
    }
    if len(allow) > 0 {
        res.Header().Set("Allow", strings.Join(allow, ", "))
        res.WriteHeader(http.StatusMethodNotAllowed) //405 error
        return
    }
    http.NotFound(res, req) // 404 error
}
 
type ContextKey string
 
//Return a copy of the request with added context, including path parameters
 
func buildContext(req *http.Request, paramKeys, paramValues []string) *http.Request {
    ctx := req.Context()
    for i := 0; i < len(paramKeys); i++ {
        ctx = context.WithValue(ctx, ContexKey(paramKeys[i]), paramValues[i])
    }
    return req.WithContext(ctx)
}