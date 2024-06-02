package listing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/tamboto2000/99-backend-exercise/internal/apps/public/config"
	"github.com/tamboto2000/99-backend-exercise/internal/common/errors"
	"github.com/tamboto2000/99-backend-exercise/pkg/gocb"
	"github.com/tamboto2000/99-backend-exercise/pkg/httpx"
)

type ListingService interface {
	CreateListing(ctx context.Context, userId int64, lt string, price int64) (Listing, error)
	GetAllListing(ctx context.Context, pageNum, pageSize int) ([]Listing, error)
}

const (
	listingSvcCfgKey     = "listing_svc"
	createListingEndpKey = "create_listing"
	getAllListingEndpKey = "get_all_listings"
)

type listingSvc struct {
	cl       httpx.Client
	host     string
	endpsCfg config.Endpoints
}

func NewListingSvc(cfg config.Config) ListingService {
	svcCfg := cfg.Services.Get(listingSvcCfgKey)
	endpsCfg := svcCfg.Endpoints

	// TODO: configuration for cb and http client
	// can be put into config
	cb := gocb.NewCircuitBreaker("listing-svc", gocb.Settings{
		ErrThreshold: 10,
		ErrInterval:  1 * time.Minute,
		Timeout:      1 * time.Minute,
		Retry:        3,
	})

	cl := http.DefaultClient
	cl.Timeout = 10 * time.Second

	return &listingSvc{
		cl:       httpx.NewClientWithCB(cl, cb),
		host:     svcCfg.Host,
		endpsCfg: endpsCfg,
	}
}

func (l *listingSvc) CreateListing(ctx context.Context, userId int64, lt string, price int64) (Listing, error) {
	endp := fmt.Sprintf("%s%s", l.host, l.endpsCfg.Get(createListingEndpKey))

	data := make(url.Values)
	data.Set("user_id", strconv.FormatInt(userId, 10))
	data.Set("listing_type", lt)
	data.Set("price", strconv.FormatInt(price, 10))
	reqBody := bytes.NewReader([]byte(data.Encode()))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endp, reqBody)
	if err != nil {
		return Listing{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := l.cl.Do(req)
	if err != nil {
		return Listing{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			var respErr createListingResponse
			if err := json.NewDecoder(resp.Body).Decode(&respErr); err != nil {
				return Listing{}, err
			}

			fields := make(errors.Fields)
			for _, e := range respErr.Errors {
				if e == "invalid user_id" {
					fields.Add("user_id", e)
				}

				if e == "invalid listing_type. Supported values: 'rent', 'sale'" {
					fields.Add("listing_type", e)
				}

				if e == "price must be greater than 0" {
					fields.Add("price", e)
				}
			}

			errVld := errors.NewErrValidation("invalid input", fields)
			return Listing{}, errVld
		}

		respBody, _ := io.ReadAll(resp.Body)
		return Listing{}, errors.New(string(respBody), errors.CodeInternal)
	}

	var respData createListingResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return Listing{}, err
	}

	return respData.Listing, nil
}

func (l *listingSvc) GetAllListing(ctx context.Context, pageNum, pageSize int) ([]Listing, error) {
	if pageNum <= 0 {
		pageNum = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	endp := fmt.Sprintf("%s%s", l.host, l.endpsCfg.Get(getAllListingEndpKey))

	uri, err := url.ParseRequestURI(endp)
	if err != nil {
		return nil, err
	}

	vals := uri.Query()
	vals.Set("page_num", strconv.Itoa(pageNum))
	vals.Set("page_size", strconv.Itoa(pageSize))

	uri.RawQuery = vals.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := l.cl.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var respData getAllListingResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}

	return respData.Listings, nil
}
