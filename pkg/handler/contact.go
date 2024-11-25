package handler

import (
	"context"
	"fmt"
	"log"

	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/awslambda"
	"github.com/ncostamagna/go-http-utils/request"
	"github.com/ncostamagna/go-http-utils/response"
	"github.com/starry-axul/axul-contact/internal/contact"
)

func NewGetHandler(endpoints contact.Endpoints) *awslambda.Handler {
	return awslambda.NewHandler(
		endpoint.Endpoint(endpoints.Get),
		decodeGetHandler,
		EncodeResponse,
		HandlerErrorEncoder(),
		awslambda.HandlerFinalizer(HandlerFinalizer()))
}

func NewGetAllHandler(endpoints contact.Endpoints) *awslambda.Handler {
	return awslambda.NewHandler(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllHandler,
		EncodeResponse,
		HandlerErrorEncoder(),
		awslambda.HandlerFinalizer(HandlerFinalizer()))
}

func NewStoreHandler(endpoints contact.Endpoints) *awslambda.Handler {
	return awslambda.NewHandler(
		endpoint.Endpoint(endpoints.Create),
		decodeStoreHandler,
		EncodeResponse,
		HandlerErrorEncoder(),
		awslambda.HandlerFinalizer(HandlerFinalizer()))
}

func NewUpdateHandler(endpoints contact.Endpoints) *awslambda.Handler {
	return awslambda.NewHandler(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateHandler,
		EncodeResponse,
		HandlerErrorEncoder(),
		awslambda.HandlerFinalizer(HandlerFinalizer()))
}

func NewDeleteHandler(endpoints contact.Endpoints) *awslambda.Handler {
	return awslambda.NewHandler(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteHandler,
		EncodeResponse,
		HandlerErrorEncoder(),
		awslambda.HandlerFinalizer(HandlerFinalizer()))
}

func NewAlertHandler(endpoints contact.Endpoints) *awslambda.Handler {
	return awslambda.NewHandler(
		endpoint.Endpoint(endpoints.Alert),
		decodeAlertHandler,
		EncodeResponse,
		HandlerErrorEncoder(),
		awslambda.HandlerFinalizer(HandlerFinalizer()))
}

func decodeGetHandler(_ context.Context, payload []byte) (interface{}, error) {
	var event events.APIGatewayProxyRequest
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	var req contact.GetReq
	if err := request.DecodeMap(event.PathParameters, &req); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	return req, nil
}

func decodeGetAllHandler(_ context.Context, payload []byte) (interface{}, error) {

	var event events.APIGatewayProxyRequest
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, response.BadRequest(err.Error())
	}
	fmt.Println(string(payload))

	var req contact.GetAllReq
	err := request.DecodeMap(event.QueryStringParameters, &req)
	fmt.Println(req)
	if err != nil {
		return nil, response.BadRequest(err.Error())
	}

	return req, nil
}

func decodeStoreHandler(_ context.Context, payload []byte) (interface{}, error) {

	var gateway events.APIGatewayProxyRequest
	err := json.Unmarshal(payload, &gateway)
	if err != nil {
		return nil, response.BadRequest(err.Error())
	}

	var req contact.StoreReq
	if gateway.Body != "" {
		if err := json.Unmarshal([]byte(gateway.Body), &req); err != nil {
			return nil, response.BadRequest(err.Error())
		}
	}
	return req, nil
}

func decodeUpdateHandler(_ context.Context, payload []byte) (interface{}, error) {

	var req contact.UpdateReq

	return req, nil
}

func decodeDeleteHandler(_ context.Context, payload []byte) (interface{}, error) {

	var req contact.DeleteReq

	return req, nil
}

func decodeAlertHandler(_ context.Context, payload []byte) (interface{}, error) {

	var req contact.AlertReq

	return req, nil
}

func EncodeResponse(_ context.Context, resp interface{}) ([]byte, error) {
	var res response.Response
	if resp == nil {
		return APIGatewayProxyResponse(response.InternalServerError("unknown response"))
	}

	switch r := resp.(type) {
	case response.Response:
		res = r
	default:
		res = response.InternalServerError("unknown response type")
	}
	return APIGatewayProxyResponse(res)
}

// HandlerErrorEncoder
func HandlerErrorEncoder() awslambda.HandlerOption {
	return awslambda.HandlerErrorEncoder(
		awslambda.ErrorEncoder(errorEncoder()),
	)
}

// HandlerFinalizer -
func HandlerFinalizer() func(context.Context, []byte, error) {
	return func(ctx context.Context, resp []byte, err error) {
		if err != nil {
			log.Println("err", err)
		}
	}
}

func errorEncoder() func(context.Context, error) ([]byte, error) {
	return func(_ context.Context, err error) ([]byte, error) {
		res := buildResponse(err)
		return APIGatewayProxyResponse(res)
	}
}

// buildResponse builds an error response from an error.
func buildResponse(err error) response.Response {
	if err == nil {
		return response.InternalServerError("unknown error")
	}

	switch e := err.(type) {
	case response.Response:
		return e
	}

	return response.InternalServerError("")
}

// APIGatewayProxyResponse
func APIGatewayProxyResponse(res response.Response) ([]byte, error) {
	bytes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	awsResponse := events.APIGatewayProxyResponse{
		Body:       string(bytes),
		StatusCode: res.StatusCode(),
	}
	return json.Marshal(awsResponse)
}


/*package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/digitalhouse-tech/go-lib-kit/response"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	authentication "github.com/ncostamagna/axul_auth/auth"
	"github.com/ncostamagna/axul_contact/internal/contact"
	"net/http"
	"strconv"
)

// NewHTTPServer is a server handler
func NewHTTPServer(ctx context.Context, auth authentication.Auth, endpoints contact.Endpoints) http.Handler {
	r := gin.Default()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Use(ginDecode(), authDecode(auth))

	r.POST("/contacts", gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateContact,
		encodeResponse,
		opts...,
	)))

	r.GET("/contacts", gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAll,
		encodeResponse,
		opts...,
	)))

	r.GET("/contacts/:id", gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetContact,
		encodeResponse,
		opts...,
	)))

	r.PATCH("/contacts/:id", gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateCourse,
		encodeResponse,
		opts...,
	)))

	r.DELETE("/contacts/:id", gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteContact,
		encodeResponse,
		opts...,
	)))

	r.POST("/contacts/alert", gin.WrapH(httptransport.NewServer(
		endpoint.Endpoint(endpoints.Alert),
		decodeAlert,
		encodeResponse,
		opts...,
	)))

	return r

}

func ginDecode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "params", c.Params)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func authDecode(auth authentication.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "auth", auth)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func encodeResponse(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func decodeCreateContact(ctx context.Context, r *http.Request) (interface{}, error) {
	var req contact.StoreReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	v := r.URL.Query()
	if err := authContact(ctx, v.Get("userid"), r.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	return req, nil
}

func decodeGetContact(ctx context.Context, r *http.Request) (interface{}, error) {
	pp := ctx.Value("params").(gin.Params)
	req := contact.GetReq{
		ID: pp.ByName("id"),
	}

	qs := r.URL.Query()

	if err := authContact(ctx, qs.Get("userid"), r.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	return req, nil
}

func decodeGetAll(ctx context.Context, r *http.Request) (interface{}, error) {

	v := r.URL.Query()

	d, _ := strconv.ParseInt(v.Get("days"), 0, 64)
	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	m, _ := strconv.ParseInt(v.Get("month"), 0, 64)
	req := contact.GetAllReq{
		Birthday:  v.Get("birthday"),
		Days:      d,
		Month:     int16(m),
		Firstname: v.Get("firstname"),
		Lastname:  v.Get("lastname"),
		Limit:     limit,
		Page:      page,
	}

	if err := authContact(ctx, v.Get("userid"), r.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	//req.Auth.ID = v.Get("userid")
	//req.Auth.Token = r.Header.Get("Authorization")

	return req, nil
}

func decodeDeleteContact(ctx context.Context, r *http.Request) (interface{}, error) {
	pp := ctx.Value("params").(gin.Params)
	req := contact.DeleteReq{
		ID: pp.ByName("id"),
	}

	qs := r.URL.Query()
	if err := authContact(ctx, qs.Get("userid"), r.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}
	return req, nil
}

func decodeUpdateCourse(ctx context.Context, r *http.Request) (interface{}, error) {

	var req contact.UpdateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	params := ctx.Value("params").(gin.Params)
	req.ID = params.ByName("id")

	qs := r.URL.Query()
	if err := authContact(ctx, qs.Get("userid"), r.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	return req, nil
}

func decodeAlert(ctx context.Context, r *http.Request) (interface{}, error) {

	v := r.URL.Query()

	req := contact.AlertReq{
		Birthday: v.Get("birthday"),
	}

	return req, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}

func authContact(ctx context.Context, userID, token string) error {
	a := ctx.Value("auth").(authentication.Auth)
	return a.Access(userID, token)
}
*/