package fix_accessible

import (
	"fmt"
	"github.com/clcert/osr/databases"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	chilean_dns "github.com/clcert/osr/tasks/processes/importer/clcert/chilean-dns"
	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
	"time"
)

var RRToDerived = map[models.RRType]models.RRType{
	models.A:  models.NORR,
	models.MX: models.MX,
	models.NS: models.NS,
}

func Execute(ctx *tasks.Context) (err error) {
	minDate := time.Time{} // 0001-01-01 is default min date
	maxDate := time.Now().AddDate(0, 0, 1) // Tomorrow is default max date
	minDateStr, ok := ctx.Params["minDate"]
	if ok {
		minDate, err = time.Parse(chilean_dns.DayFormat, minDateStr)
		if err != nil {
			return fmt.Errorf("error parsing minDate: %s", err)
		}
	}
	maxDateStr, ok := ctx.Params["maxDate"]
	if ok {
		maxDate, err = time.Parse(chilean_dns.DayFormat, maxDateStr)
		if err != nil {
			return fmt.Errorf("error parsing minDate: %s", err)
		}
	}
	accessibleDate, err := chilean_dns.GetAccessibleIPs(ctx.Sources[0], ctx)
	if err != nil {
		return err
	}
	writer, err := databases.GetPostgresWriter()
	if err != nil {
		return err
	}
	if len(accessibleDate) == 0 {
		return fmt.Errorf("empty accessible date list")
	}
	ctx.Log.WithFields(logrus.Fields{
		"minDate": minDate,
		"maxDate": maxDate,
	}).Info("Finding all tasks IDs from DNS RRs")
	// Finding all tasks IDs from DNS RRs
	taskList, err := getTasksWithDate(writer)
	if err != nil {
		return err
	}
	taskIDs := make([]int, len(taskList))
	for i, task := range taskList {
		taskIDs[i] = task.ID
	}
	/*
		ctx.Log.WithFields(logrus.Fields{
			"tasks": taskList,
		}).Info("Starting query to reset accessible state of all RRs")
		// Set accessible as false for everyone
		res, err := writer.Query(nil, "UPDATE dns_rrs SET accessible = false")
		if err != nil {
			ctx.Log.WithFields(logrus.Fields{
				"error": err,
			}).Error("Error executing query")
			return err
		}
		ctx.Log.WithFields(logrus.Fields{
			"tasks":    taskList,
		}).Info("Reset accessible RRs query executed successfully")
	*/
	for _, task := range taskList {
		date := task.StartDate
		if date.Before(minDate) || date.After(maxDate) {
			ctx.Log.WithFields(logrus.Fields{
				"taskID": task.ID,
				"date": date,
				"minDate": minDate,
				"maxDate": maxDate,
			}).Info("Skipping this task because is outside the limits")
			continue
		}
		ctx.Log.WithFields(logrus.Fields{
			"taskID": task.ID,
			"date": date,
		}).Info("Looking for accessible states in this task")
		accessible, ok := accessibleDate[date.Format(chilean_dns.DayFormat)]
		if !ok {
			ctx.Log.WithFields(logrus.Fields{
				"date":  date,
				"ID":    task.ID,
				"error": err,
			}).Error("Error finding accessible information for this task, skipping...")
			continue
		}
		// Update accessible where IP in set
		for rrType, accessibleRR := range accessible {
			if len(accessibleRR) == 0 {
				ctx.Log.WithFields(logrus.Fields{
					"taskID":            task.ID,
					"date":              date,
					"accessible_length": len(accessibleRR),
				}).Info("Skipping accessible state for RR without IPs")
				continue
			}
			ctx.Log.WithFields(logrus.Fields{
				"taskID":            task.ID,
				"date":              date,
				"RR":                rrType,
				"accessible_length": len(accessibleRR),
			}).Info("Updating accessible state for this RR")
			ips := make([]string, 0, len(accessible))
			for ipStr, _ := range accessibleRR {
				ips = append(ips, ipStr)
			}

			_, err := writer.Query(nil, "UPDATE dns_rrs SET accessible = true"+
				" WHERE scan_type = 1 and derived_type = ? and task_id = ? and ip_value in (?)",
				RRToDerived[rrType],
				task.ID,
				pg.In(ips))
			if err != nil {
				ctx.Log.WithFields(logrus.Fields{
					"date":  date,
					"error": err,
				}).Error("Error executing query, skipping date...")
				continue
			}
			ctx.Log.WithFields(logrus.Fields{
				"taskID": task.ID,
				"RR":     rrType,
				"date":   date,
			}).Info("Update accessible RRs query executed successfully")
		}
	}
	return nil
}

func getTasksWithDate(db *pg.DB) ([]models.Task, error) {
	var taskIDs []int
	var tasks []models.Task

	_, err := db.Query(&taskIDs, "select distinct task_id from dns_rrs")
	if err != nil {
		return nil, err
	}
	err = db.Model(&tasks).
		Where("id in (?)", pg.In(taskIDs)).
		Select()
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, fmt.Errorf("no dnsRR tasks found")
	}
	return tasks, nil
}
