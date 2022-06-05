package internal

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/matryer/is"
)

type (
	// for each test case
	testCase struct {
		caseName string
		url      string
		query    requestQuery
		expected response
	}
	// request query parameters
	requestQuery struct {
		startTime string
		endTime   string
	}
	// response - for now we check just status code
	response struct {
		code int
	}
)

func TestAPI(t *testing.T) {
	// creating mock sql db
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// starting up our api and connection to mock db
	server, err := Setup(db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//creating test cases
	var testCases = []testCase{
		{
			caseName: "timeseries request",
			url:      "/timeseries",
			query:    requestQuery{"2020-01-01T23:01", "2023-01-01T12:03"},
			expected: response{http.StatusOK},
		},
		{
			caseName: "cpu statistic request",
			url:      "/statistic/cpu",
			query:    requestQuery{"2020-01-01T23:01", "2023-01-01T12:03"},
			expected: response{http.StatusOK},
		},
		{
			caseName: "concurrency statistic request",
			url:      "/statistic/concurrency",
			query:    requestQuery{"2020-01-01T23:01", "2023-01-01T12:03"},
			expected: response{http.StatusOK},
		},
	}
	// used for asserting conditions
	i := is.New(t)
	// to ensure rows with three columns are returned for all 200 api calls
	columns := []string{"col1", "col2", "col3"}
	// testing each case
	for _, v := range testCases {
		t.Run(v.caseName, func(t *testing.T) {
			// mocking sql db depending on type of request

			switch {
			case strings.Contains(v.url, ts):
				mock.ExpectPrepare("select . from").ExpectQuery().
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(time.Now(), 2.33, 54352))
			case strings.Contains(v.url, cpu):
				mock.ExpectPrepare("select AVG").ExpectQuery().
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(54.43, 94.3, 12.3))
			case strings.Contains(v.url, concurrency):
				mock.ExpectPrepare("select AVG").ExpectQuery().
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow(26423.5, 32432, 12343))
			}
			// creating request
			r := httptest.NewRequest(http.MethodGet, v.url, nil)
			// setting request parameters
			q := r.URL.Query()
			q.Add("startTime", v.query.startTime)
			q.Add("endTime", v.query.endTime)
			r.URL.RawQuery = q.Encode()
			w := httptest.NewRecorder()
			// calling api
			server.Api.ServeHTTP(w, r)
			resp := w.Result()
			//printing response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error("could not read body")
			}
			fmt.Println(string(body))
			// we make sure that all expectations were met
			i.Equal(v.expected.code, resp.StatusCode)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}

}
