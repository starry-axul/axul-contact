package contact

import (
	"context"
	"errors"
	"fmt"
	"github.com/ncostamagna/go-http-utils/meta"
	"github.com/ncostamagna/go-http-utils/response"
	"strconv"
	"time"
)

const (
	layoutISO = "2006-01-02 15:04:05"
)

// Endpoints struct
type (
	StoreReq struct {
		Auth      Authentication
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Nickname  string `json:"nickname"`
		Gender    string `json:"gender"`
		Phone     string `json:"phone"`
		Birthday  string `json:"birthday"`
	}

	GetReq struct {
		Auth Authentication
		ID   string
	}

	GetAllReq struct {
		Auth      Authentication
		Days      int64
		Birthday  string
		Firstname string
		Lastname  string
		Month     int16
		Limit     int
		Page      int
	}

	AlertReq struct {
		Birthday string
	}

	Authentication struct {
		ID    string
		Token string
	}

	UpdateReq struct {
		ID        string  `json:"id"`
		Firstname *string `json:"firstname"`
		Lastname  *string `json:"lastname"`
		Nickname  *string `json:"nickname"`
		Gender    *string `json:"gender"`
		Phone     *string `json:"phone"`
		Birthday  *string `json:"birthday"`
	}

	DeleteReq struct {
		ID string
	}

	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoints struct {
		Create Controller
		Update Controller
		Get    Controller
		GetAll Controller
		Delete Controller
		Alert  Controller
	}
)

// MakeEndpoints handler endpoints
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Update: makeUpdateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Delete: makeDeleteEndpoint(s),
		Alert:  makeAlertEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(StoreReq)

		if req.Firstname == "" {
			return nil, response.BadRequest("first name is required")
		}

		if req.Lastname == "" {
			return nil, response.BadRequest("last name is required")
		}

		if req.Nickname == "" {
			return nil, response.BadRequest("nick name is required")
		}

		birthday, err := time.Parse(layoutISO, fmt.Sprintf("%s 17:00:00", req.Birthday))

		if err != nil {
			return nil, response.BadRequest(err.Error())
		}

		c, err := s.Create(ctx, req.Firstname, req.Lastname, req.Nickname, req.Gender, req.Phone, birthday)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", c, nil), nil

	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetAllReq)

		//if err := s.authorization(ctx, req.Auth.ID, req.Auth.Token); err != nil {
		//	return nil, response.Unauthorized(err.Error())
		//}

		fd := time.Now().UTC()
		f := Filter{
			Firstname: req.Firstname,
			Lastname:  req.Lastname,
			Month:     req.Month,
			firstDate: time.Date(fd.Year(), fd.Month(), fd.Day(), 0, 0, 0, 0, time.UTC),
		}

		if req.Birthday != "" {
			days, err := strconv.Atoi(req.Birthday)
			if err != nil {
				return nil, response.BadRequest("Invalid birthday format in Query String")
			}

			f.Birthday = &days
		}

		if req.Days > 0 {
			f.RangeDays = &req.Days
		}

		count, err := s.Count(ctx, f)
		fmt.Println(count)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, "25")
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		cs, err := s.GetAll(ctx, f, meta.Offset(), meta.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("Success", cs, meta), nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(UpdateReq)
		var birthday *time.Time

		if req.Firstname != nil && *req.Firstname == "" {
			return nil, response.BadRequest("first name is required")
		}

		if req.Lastname != nil && *req.Lastname == "" {
			return nil, response.BadRequest("last name is required")
		}

		if req.Nickname != nil && *req.Nickname == "" {
			return nil, response.BadRequest("nick name is required")
		}

		if req.Birthday != nil {
			b, err := time.Parse(layoutISO, fmt.Sprintf("%s 17:00:00", *req.Birthday))
			if err != nil {
				return nil, response.BadRequest(err.Error())
			}
			birthday = &b
		}

		if err := s.Update(ctx, req.ID, req.Firstname, req.Lastname, req.Nickname, req.Gender, req.Phone, birthday); err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil, nil), nil
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetReq)

		contact, err := s.Get(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		return response.OK("Success", contact, nil), nil
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(DeleteReq)

		if err := s.Delete(ctx, req.ID); err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil, nil), nil
	}
}

func makeAlertEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AlertReq)

		cs, err := s.Alert(ctx, req.Birthday)
		if err != nil {
			return nil, err
		}

		return response.OK("Success", cs, nil), nil
	}
}
