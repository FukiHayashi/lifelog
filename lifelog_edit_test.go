package main_test

import (
	"fmt"
	"lifelog/database"
	"lifelog/helpers"
	"lifelog/models"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
	"gorm.io/gorm"
)

var _ = Describe("LifelogEdit", Ordered, func() {
	var page *agouti.Page
	var db *gorm.DB
	BeforeAll(func() {
		// テスト環境のDBに接続
		db = database.DataBaseConnect()
		// DBのマイグレーション
		db.AutoMigrate(&models.User{}, &models.LifeLog{}, &models.Appointment{})
	})
	AfterAll(func() {
		// テストに使用したDBの内容を全て削除する
		db.Migrator().DropTable(&models.User{}, &models.LifeLog{}, &models.Appointment{})
		dbc, _ := db.DB()
		dbc.Close()
	})
	BeforeEach(func() {
		var err error
		page, err = agoutiDriver.NewPage()
		Expect(err).NotTo(HaveOccurred())
		// /からサインインする
		Expect(page.Navigate("http://localhost:3000")).To(Succeed())
		Expect(page.FindByID("signin").Click()).To(Succeed())
		// 画面が切り替わるまで少し待つ
		time.Sleep(1 * time.Second)
		// Auth0のログイン画面が表示されたらログイン
		title, _ := page.Title()
		if title == "Log in | Lifelog App" {
			page.FindByID("username").Fill(os.Getenv("AUTH0_EMAIL"))
			page.FindByID("password").Fill(os.Getenv("AUTH0_PASS"))
			page.FindByButton("Continue").Click()
			time.Sleep(1 * time.Second)
		}
	})

	AfterEach(func() {
		Expect(page.Destroy()).To(Succeed())
	})
	Describe("行動編集画面の表示", Ordered, func() {
		var (
			appointment models.Appointment
			lifelog     models.LifeLog
			user        models.User
		)
		BeforeAll(func() {
			// user情報取得
			db.Where("name = ?", os.Getenv("AUTH0_EMAIL")).First(&user)
			// lifelog情報取得
			db.Where(&models.LifeLog{UserId: user.ID}).Where("name = ?", time.Now().Format("2006/01/02")).First(&lifelog)
			// appointment情報書き込み
			appointment = models.Appointment{
				LifeLogId: lifelog.ID,
				Title:     helpers.GetStringPointer("before"),
			}
			db.Save(&appointment)
			Expect(page.Navigate("http://localhost:3000/lifelog/edit/" + appointment.ID.String())).To(Succeed())
		})
		AfterAll(func() {
			// 使用したappointmentを削除
			db.Delete(&appointment)
		})
		Context("/lifelog/edit/:appointmentIdを表示した時", func() {
			It("行動記録編集画面が表示されること", func() {
				Expect(page).To(HaveTitle("Lifelog | 行動記録編集"))
				Expect(page.FindByName("title")).To(BeFound())
				Expect(page.FindByName("start")).To(BeFound())
				Expect(page.FindByName("end")).To(BeFound())
				Expect(page.FindByName("class")).To(BeFound())
				Expect(page.FindByID("resister-button")).To(BeFound())
				Expect(page.FindByID("cancel-button")).To(BeFound())
				Expect(page.FindByID("js_delete")).To(BeFound())
			})
		})
	})
	Describe("行動記録の編集", func() {
		var (
			appointment models.Appointment
			lifelog     models.LifeLog
			user        models.User
		)
		base_path := "http://localhost:3000/lifelog/edit/"
		BeforeEach(func() {
			// user情報取得
			db.Where("name = ?", os.Getenv("AUTH0_EMAIL")).First(&user)
			// lifelog情報取得
			db.Where(&models.LifeLog{UserId: user.ID}).Where("name = ?", time.Now().Format("2006/01/02")).First(&lifelog)
			// appointment情報書き込み
			appointment = models.Appointment{
				LifeLogId: lifelog.ID,
				Title:     helpers.GetStringPointer("before"),
			}
			db.Save(&appointment)
		})
		AfterEach(func() {
			// 使用したappointmentを削除
			db.Delete(&appointment)
		})
		Context("行動記録のタイトルを編集した時", func() {
			It("タイトルが変更されること", func() {
				appointments := []models.Appointment{}
				Expect(page.Navigate(base_path + appointment.ID.String())).To(Succeed())
				Expect(page.FindByID("js_appointment_title").Attribute("value")).To(Equal("before"))
				Expect(page.FindByName("title").Fill("after")).To(Succeed())
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				db.Where(&models.Appointment{LifeLogId: lifelog.ID}).Find(&appointments)
				for _, app := range appointments {
					if *app.Title == "after" {
						appointment = app
						break
					}
				}
				Expect(page.Navigate(base_path + appointment.ID.String())).To(Succeed())
				Expect(page.FindByID("js_appointment_title").Attribute("value")).To(Equal("after"))
			})
		})
		Context("行動記録のタイトルを空にした時", func() {
			It("タイトルが変更されないこと", func() {
				Expect(page.Navigate(base_path + appointment.ID.String())).To(Succeed())
				Expect(page.FindByID("js_appointment_title").Attribute("value")).To(Equal("before"))
				Expect(page.FindByName("title").Fill("")).To(Succeed())
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				Expect(page.FindByID("js_appointment_title").Attribute("value")).To(Equal("before"))
				Expect(page).To(HaveTitle("Lifelog | 行動記録編集"))
				Expect(page.FindByID("msg-error")).To(HaveText("未入力の項目があります"))
			})
		})
		Context("開始時刻を編集した時", func() {
			It("開始時刻が変更されること", func() {
				before := *lifelog.Name + " " + *appointment.Start
				after := *lifelog.Name + " 00:10"
				appointments := []models.Appointment{}
				Expect(page.Navigate(base_path + appointment.ID.String())).To(Succeed())
				Expect(page.FindByName("start").Attribute("value")).To(Equal(before))
				Expect(page.FindByName("start").Fill(after)).To(Succeed())
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				db.Where(&models.Appointment{LifeLogId: lifelog.ID}).Find(&appointments)
				for _, app := range appointments {
					if *app.Title == "before" {
						appointment = app
						break
					}
				}
				Expect(page.Navigate(base_path + appointment.ID.String())).To(Succeed())
				Expect(page.FindByName("start").Attribute("value")).To(Equal(after))
			})
		})
		Context("終了時刻を編集した時", func() {
			It("終了時刻が変更されること", func() {
				before := *lifelog.Name + " " + *appointment.End
				after := *lifelog.Name + " 00:40"
				appointments := []models.Appointment{}
				Expect(page.Navigate(base_path + appointment.ID.String())).To(Succeed())
				Expect(page.FindByName("end").Attribute("value")).To(Equal(before))
				Expect(page.FindByName("end").Fill(after)).To(Succeed())
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				db.Where(&models.Appointment{LifeLogId: lifelog.ID}).Find(&appointments)
				for _, app := range appointments {
					if *app.Title == "before" {
						appointment = app
						break
					}
				}
				Expect(page.Navigate(base_path + appointment.ID.String())).To(Succeed())
				Expect(page.FindByName("end").Attribute("value")).To(Equal(after))
			})
		})
		Context("分類を編集した時", func() {
			It("分類が変更されること", func() {
				after := "食事"
				appointments := []models.Appointment{}
				Expect(page.Navigate(base_path + appointment.ID.String())).To(Succeed())
				Expect(page.FindByName("class").Select(after)).To(Succeed())
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				db.Where(&models.Appointment{LifeLogId: lifelog.ID}).Find(&appointments)
				for _, app := range appointments {
					if *app.Title == "before" {
						appointment = app
						break
					}
				}
				Expect(page.Navigate(base_path + appointment.ID.String())).To(Succeed())
				fmt.Println(page.FindByName("class").Attribute("value"))
				Expect(page.FindByName("class").Attribute("value")).To(Equal("meal"))
			})
		})

	})
	Describe("キャンセルボタン", func() {
		var (
			appointment models.Appointment
			lifelog     models.LifeLog
			user        models.User
		)
		BeforeAll(func() {
			// user情報取得
			db.Where("name = ?", os.Getenv("AUTH0_EMAIL")).First(&user)
			// lifelog情報取得
			db.Where(&models.LifeLog{UserId: user.ID}).Where("name = ?", time.Now().Format("2006/01/02")).First(&lifelog)
			// appointment情報書き込み
			appointment = models.Appointment{
				LifeLogId: lifelog.ID,
				Title:     helpers.GetStringPointer("before"),
			}
			db.Save(&appointment)
		})
		AfterAll(func() {
			// 使用したappointmentを削除
			db.Delete(&appointment)
		})
		BeforeEach(func() {
			Expect(page.Navigate("http://localhost:3000/lifelog/edit/" + appointment.ID.String())).To(Succeed())
		})
		Context("キャンセルボタンを押した時", func() {
			It("メイン画面に遷移すること", func() {
				page.FindByID("cancel-button").Click()
				Expect(page).To(HaveTitle("Lifelog | メイン"))
			})
		})
	})
	Describe("削除ボタン", func() {
		var (
			appointment models.Appointment
			lifelog     models.LifeLog
			user        models.User
		)
		BeforeAll(func() {
			// user情報取得
			db.Where("name = ?", os.Getenv("AUTH0_EMAIL")).First(&user)
			// lifelog情報取得
			db.Where(&models.LifeLog{UserId: user.ID}).Where("name = ?", time.Now().Format("2006/01/02")).First(&lifelog)
			// appointment情報書き込み
			appointment = models.Appointment{
				LifeLogId: lifelog.ID,
				Title:     helpers.GetStringPointer("before"),
				Class:     "sleep",
			}
			db.Save(&appointment)
		})
		AfterAll(func() {
			// 使用したappointmentを削除
			db.Delete(&appointment)
		})
		Context("削除ボタンを押した時", func() {
			It("行動記録が削除され、メイン画面に遷移すること", func() {
				Expect(page.Navigate("http://localhost:3000/lifelog")).To(Succeed())
				before_count, _ := page.AllByClass("sleep").Count()
				Expect(page.Navigate("http://localhost:3000/lifelog/edit/" + appointment.ID.String())).To(Succeed())
				page.FindByID("js_delete").Click()
				Expect(page).To(HaveTitle("Lifelog | メイン"))
				after_count, _ := page.AllByClass("sleep").Count()
				Expect(after_count).To(Equal(before_count - 1))
			})
		})
	})
})
