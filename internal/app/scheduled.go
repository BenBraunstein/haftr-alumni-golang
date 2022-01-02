package app

import (
	"github.com/BenBraunstein/haftr-alumni-golang/common/time"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/db"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/email"
	"github.com/BenBraunstein/haftr-alumni-golang/internal/workflow"
)

type ScheduledFunc func() error

func HappyBirthdayEmailScheduled(retrieveAlumnis db.RetrieveAllAlumniFunc,
	provideTime time.EpochProviderFunc,
	getEmailTemplate db.RetrieveEmailTemplateByNameFunc,
	retrieveUserByAlumniId db.RetrieveUserByAlumniIDFunc,
	sendEmail email.SendEmailFunc) ScheduledFunc {
	return func() error {
		happyBirthdayEmail := workflow.HappyBirthdayEmail(retrieveAlumnis, provideTime, getEmailTemplate, retrieveUserByAlumniId, sendEmail)
		if err := happyBirthdayEmail(); err != nil {
			return err
		}
		return nil
	}
}
