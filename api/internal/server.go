package internal

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	"github.com/zayyanmatin/sre-data-service/models"
)

type Server struct {
	Api *echo.Echo
	Db  *sql.DB
}

const (
	cpu         = "cpu"
	ts          = "timeseries"
	concurrency = "concurrency"
)

func Start(db *sql.DB) (*Server, error) {
	var server Server

	// creating echo instance
	e := echo.New()

	// Routes
	e.GET("/timeseries", server.getTimeseries)
	e.GET("/statistic/cpu", server.getCpu)
	e.GET("/statistic/concurrency", server.getConcurrency)

	server.Api = e
	server.Db = db

	return &server, nil

}

func (w *Server) Close() error {
	// closing echo instance
	if err := w.Api.Close(); err != nil {
		return fmt.Errorf("unable to close api: %w", err)
	}

	// closing db instance
	if err := w.Db.Close(); err != nil {
		return fmt.Errorf("unable to close db: %w", err)
	}
	return nil
}

func (w *Server) getTimeseries(ctx echo.Context) error {
	//retrieving query parameters
	startTime := ctx.QueryParam("startTime")
	endTime := ctx.QueryParam("endTime")
	return w.queryTimeSeries(ctx, startTime, endTime)
}

func (w *Server) getCpu(ctx echo.Context) error {
	//retrieving query parameters
	startTime := ctx.QueryParam("startTime")
	endTime := ctx.QueryParam("endTime")
	return w.queryStatistic(ctx, cpu, startTime, endTime)
}

func (w *Server) getConcurrency(ctx echo.Context) error {
	//retrieving query parameters
	startTime := ctx.QueryParam("startTime")
	endTime := ctx.QueryParam("endTime")
	return w.queryStatistic(ctx, concurrency, startTime, endTime)
}

func (w *Server) queryTimeSeries(ctx echo.Context, start, end string) error {
	//preparing sql statement before arguments are populated
	stmt, err := w.Db.Prepare("select *  from sre.timeseries where  ?<=ts and ts<=?")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.Message{Error: fmt.Errorf("could not prepare statement: %w", err).Error()})
	}
	defer stmt.Close()
	// querying sql statement by supplying arguments
	rows, err := stmt.Query(start, end)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.Message{Error: fmt.Errorf("could not query statement: %w", err).Error()})
	}
	defer rows.Close()
	// populating  array struct by iterating each row
	var all []models.Timeseries
	for rows.Next() {
		var t models.Timeseries
		if err := rows.Scan(&t.Timestamp, &t.Cpu, &t.Concurrency); err != nil {
			return ctx.JSON(http.StatusInternalServerError, fmt.Errorf("could not scan row: %w", err))
		}
		all = append(all, t)
	}

	return ctx.JSON(http.StatusOK, all)
}
func (w *Server) queryStatistic(ctx echo.Context, variable, start, end string) error {
	//preparing sql statement before arguments are populated
	stmt, err := w.Db.Prepare(fmt.Sprintf(`select AVG(%s), MAX(%s), MIN(%s) from sre.timeseries where  ?<=ts and ts<=?`, variable, variable, variable))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.Message{Error: fmt.Errorf("could not prepare statement: %w", err).Error()})
	}
	defer stmt.Close()
	// querying sql statement by supplying arguments
	rows, err := stmt.Query(start, end)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, models.Message{Error: fmt.Errorf("could not query statement: %w", err).Error()})
	}
	defer rows.Close()
	// populating  array struct by iterating each row
	var all []models.Statistic
	for rows.Next() {
		var (
			avg, max, min sql.NullFloat64
		)

		if err := rows.Scan(&avg, &max, &min); err != nil {
			return ctx.JSON(http.StatusInternalServerError, models.Message{Error: fmt.Errorf("could not scan row: %w", err).Error()})

		}
		statistic := models.Statistic{
			Avg: avg.Float64, Max: max.Float64, Min: min.Float64,
		}
		all = append(all, statistic)
	}
	return ctx.JSON(http.StatusOK, all)
}
