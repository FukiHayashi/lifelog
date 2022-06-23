package main_test

import (
	. "lifelog/helpers"
	"lifelog/models"
	"lifelog/test"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Models", func() {
	th := test.TestHelper{}
	BeforeEach(func() {
		th.SetupTest()
	})
	AfterEach(func() {
		th.TearDownTest()
	})
	Describe("User", func() {
		Context("Userが正しい時", func() {
			aud := time.Now().String()
			user := models.User{
				Aud:  &aud,
				Name: time.Now().String(),
			}
			It("Userが作成される", func() {
				err := th.DB.Create(&user).Error
				Expect(err).Should(BeNil())
			})
		})
		Context("Audがない時", func() {
			erruser := models.User{
				Name: time.Now().String(),
			}
			It("Userが作成されない", func() {
				err := th.DB.Create(&erruser).Error
				Expect(err).ShouldNot(BeNil())
			})
		})
	})
	Describe("Lifelog", Ordered, func() {
		var user models.User
		var t time.Time
		var ts string

		BeforeAll(func() {
			aud := time.Now().String()
			user = models.User{
				Aud:  &aud,
				Name: time.Now().String(),
			}
			th.DB.Create(&user)
		})
		BeforeEach(func() {
			t = time.Now()
			ts = t.String()
		})

		Context("Lifelogが正しい時", func() {
			It("Lifelogが作成される", func() {
				lifelog := models.LifeLog{
					UserId:   user.ID,
					Name:     &ts,
					LoggedAt: t,
				}
				err := th.DB.Create(&lifelog).Error
				Expect(err).Should(BeNil())
			})
		})
		Context("UserIdがない時", func() {
			It("Lifelogが作成されない", func() {
				lifelog := models.LifeLog{
					Name:     &ts,
					LoggedAt: t,
				}
				err := th.DB.Create(&lifelog).Error
				Expect(err).ShouldNot(BeNil())
			})
		})
		Context("Nameがない時", func() {
			It("Lifelogが作成されない", func() {
				t := time.Now()
				lifelog := models.LifeLog{
					UserId:   user.ID,
					LoggedAt: t,
				}
				err := th.DB.Create(&lifelog).Error
				Expect(err).ShouldNot(BeNil())
			})
		})
	})
	Describe("Appointment", Ordered, func() {
		var user models.User
		var lifelog models.LifeLog
		BeforeAll(func() {
			t := time.Now()
			ts := t.String()
			user = models.User{
				Aud:  &ts,
				Name: ts,
			}
			th.DB.Create(&user)
			lifelog = models.LifeLog{
				UserId:   user.ID,
				Name:     &ts,
				LoggedAt: t,
			}
			th.DB.Create(&lifelog)
		})
		Context("Appointmentが正しい時", func() {
			It("Appointmentが作成される", func() {
				appointment := models.Appointment{
					LifeLogId: lifelog.ID,
					Title:     GetStringPointer(time.Now().String()),
				}
				err := th.DB.Create(&appointment).Error
				Expect(err).Should(BeNil())
			})
		})
		Context("LifelogIdがない時", func() {
			It("Appointmentが作成されない", func() {
				appointment := models.Appointment{
					Title: GetStringPointer(time.Now().String()),
				}
				err := th.DB.Create(&appointment).Error
				Expect(err).ShouldNot(BeNil())
			})
		})
		Context("Titleがない時", func() {
			It("Appointmentが作成されない", func() {
				appointment := models.Appointment{
					LifeLogId: lifelog.ID,
				}
				err := th.DB.Create(&appointment).Error
				Expect(err).ShouldNot(BeNil())
			})
		})
		Context("Startがない時", func() {
			It("Startが00:00でAppointmentが作成される", func() {
				appointment := models.Appointment{
					LifeLogId: lifelog.ID,
					Title:     GetStringPointer(time.Now().String()),
					End:       GetStringPointer("02:00"),
				}
				err := th.DB.Create(&appointment).Error
				Expect(err).ShouldNot(BeNil())
				Expect(*appointment.Start).To(Equal("00:00"))
			})
		})
		Context("Endがない時", func() {
			It("Endが01:00でAppointmentが作成される", func() {
				appointment := models.Appointment{
					LifeLogId: lifelog.ID,
					Title:     GetStringPointer(time.Now().String()),
					Start:     GetStringPointer("00:30"),
				}
				err := th.DB.Create(&appointment).Error
				Expect(err).ShouldNot(BeNil())
				Expect(*appointment.End).To(Equal("01:00"))
			})
		})
	})
})
