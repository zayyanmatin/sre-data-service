package internal

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

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
	dateFormat  = "2006-01-02T15:04"
)

func Setup(db *sql.DB) (*Server, error) {
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
	startString := ctx.QueryParam("startTime")
	endString := ctx.QueryParam("endTime")
	startTime, err := time.Parse(dateFormat, startString)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, models.Message{Error: fmt.Errorf("could not parse start time: %w", err).Error()})
	}
	endTime, err := time.Parse(dateFormat, endString)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, models.Message{Error: fmt.Errorf("could not parse end time: %w", err).Error()})
	}
	return w.queryTimeSeries(ctx, startTime, endTime)
}

func (w *Server) getCpu(ctx echo.Context) error {
	//retrieving query parameters
	startString := ctx.QueryParam("startTime")
	endString := ctx.QueryParam("endTime")
	startTime, err := time.Parse(dateFormat, startString)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, models.Message{Error: fmt.Errorf("could not parse start time: %w", err).Error()})
	}
	endTime, err := time.Parse(dateFormat, endString)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, models.Message{Error: fmt.Errorf("could not parse end time: %w", err).Error()})
	}
	return w.queryStatistic(ctx, cpu, startTime, endTime)
}

func (w *Server) getConcurrency(ctx echo.Context) error {
	//retrieving query parameters
	startString := ctx.QueryParam("startTime")
	endString := ctx.QueryParam("endTime")
	startTime, err := time.Parse(dateFormat, startString)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, models.Message{Error: fmt.Errorf("could not parse start time: %w", err).Error()})
	}
	endTime, err := time.Parse(dateFormat, endString)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, models.Message{Error: fmt.Errorf("could not parse end time: %w", err).Error()})
	}
	return w.queryStatistic(ctx, concurrency, startTime, endTime)
}

func (w *Server) queryTimeSeries(ctx echo.Context, start, end time.Time) error {
	//preparing sql statement before arguments are populated
	stmt, err := w.Db.Prepare("select *  from sre.timeseriesv2 where  ?<=ts and ts<=? order by ts")
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
		var (
			ts          sql.NullTime
			cpu         sql.NullFloat64
			concurrency sql.NullInt32
		)
		if err := rows.Scan(&ts, &cpu, &concurrency); err != nil {
			return ctx.JSON(http.StatusInternalServerError, models.Message{Error: fmt.Errorf("could not scan row: %w", err).Error()})
		}
		timeseries := models.Timeseries{
			Timestamp: ts.Time, Cpu: float32(cpu.Float64), Concurrency: uint32(concurrency.Int32)}
		all = append(all, timeseries)
	}

	return ctx.JSON(http.StatusOK, all)
}
func (w *Server) queryStatistic(ctx echo.Context, variable string, start, end time.Time) error {
	//preparing sql statement before arguments are populated
	stmt, err := w.Db.Prepare(fmt.Sprintf(`select AVG(%s), MAX(%s), MIN(%s) from sre.timeseriesv2 where  ?<=ts and ts<=?`, variable, variable, variable))
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
