package main_test

import (
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

var _ = Describe("RemarksEdit", Ordered, func() {
	var page *agouti.Page
	var db *gorm.DB
	BeforeAll(func() {
		// テスト環境のDBに接続
		db = database.DataBaseConnect()
		// DBのマイグレーション
		db.AutoMigrate(&models.User{}, &models.LifeLog{}, &models.Appointment{}, &models.Remarks{})
	})
	AfterAll(func() {
		// テストに使用したDBの内容を全て削除する
		db.Migrator().DropTable(&models.User{}, &models.LifeLog{}, &models.Appointment{}, &models.Remarks{})
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

	Describe("備考編集画面の表示", Ordered, func() {
		var (
			remarks models.Remarks
			lifelog models.LifeLog
			user    models.User
		)
		BeforeAll(func() {
			// user情報取得
			db.Where("name = ?", os.Getenv("AUTH0_EMAIL")).First(&user)
			// lifelog情報取得
			db.Where(&models.LifeLog{UserId: user.ID}).Where("name = ?", time.Now().Format("2006/01/02")).First(&lifelog)
			// remarks情報書き込み
			remarks = models.Remarks{
				LifeLogId: lifelog.ID,
				Title:     helpers.GetStringPointer("before"),
				Date:      lifelog.Name,
				Class:     "remarks",
			}
			db.Save(&remarks)
			Expect(page.Navigate("http://localhost:3000/remarks/edit/" + remarks.ID.String())).To(Succeed())
		})
		AfterAll(func() {
			// 使用したremarksを削除
			db.Delete(&remarks)
		})
		Context("/remarks/edit/:remarksIdを表示した時", func() {
			It("備考編集画面が表示されること", func() {
				Expect(page).To(HaveTitle("Lifelog | 備考編集"))
				Expect(page.FindByName("title")).To(BeFound())
				Expect(page.FindByName("date")).To(BeFound())
				Expect(page.FindByID("resister-button")).To(BeFound())
				Expect(page.FindByID("cancel-button")).To(BeFound())
				Expect(page.FindByID("js_delete")).To(BeFound())
			})
		})
	})
	Describe("備考の編集", func() {
		var (
			remarks models.Remarks
			lifelog models.LifeLog
			user    models.User
		)
		base_path := "http://localhost:3000/remarks/edit/"
		BeforeEach(func() {
			// user情報取得
			db.Where("name = ?", os.Getenv("AUTH0_EMAIL")).First(&user)
			// lifelog情報取得
			db.Where(&models.LifeLog{UserId: user.ID}).Where("name = ?", time.Now().Format("2006/01/02")).First(&lifelog)
			// remarks情報書き込み
			remarks = models.Remarks{
				LifeLogId: lifelog.ID,
				Title:     helpers.GetStringPointer("before"),
				Date:      lifelog.Name,
				Class:     "remarks",
			}
			db.Save(&remarks)
		})
		AfterEach(func() {
			// 使用したappointmentを削除
			db.Delete(&remarks)
		})
		Context("備考のタイトルを編集した時", func() {
			It("タイトルが変更されること", func() {
				Expect(page.Navigate(base_path + remarks.ID.String())).To(Succeed())
				Expect(page.FindByID("remarks_title").Attribute("value")).To(Equal(*remarks.Title))
				Expect(page.FindByName("title").Fill("after")).To(Succeed())
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				Expect(page.Navigate(base_path + remarks.ID.String())).To(Succeed())
				Expect(page.FindByID("remarks_title").Attribute("value")).To(Equal("after"))
			})
		})
		Context("備考のタイトルを空にした時", func() {
			It("タイトルが変更されないこと", func() {
				Expect(page.Navigate(base_path + remarks.ID.String())).To(Succeed())
				Expect(page.FindByID("remarks_title").Attribute("value")).To(Equal(*remarks.Title))
				Expect(page.FindByName("title").Fill("")).To(Succeed())
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				Expect(page.FindByID("remarks_title").Attribute("value")).To(Equal("before"))
				Expect(page).To(HaveTitle("Lifelog | 備考編集"))
				Expect(page.FindByID("msg-error")).To(HaveText("未入力の項目があります"))
			})
		})
		Context("日付を変更した時", func() {
			It("日付が変更されること", func() {
				before_date := *remarks.Date
				now := time.Now()
				after_date := now.AddDate(0, 0, 1).Format("2006/01/02")
				Expect(page.Navigate(base_path + remarks.ID.String())).To(Succeed())
				Expect(page.FindByName("date").Attribute("value")).To(Equal(before_date))
				Expect(page.FindByName("date").Fill(after_date)).To(Succeed())
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				Expect(page.Navigate(base_path + remarks.ID.String())).To(Succeed())
				Expect(page.FindByName("date").Attribute("value")).To(Equal(after_date))
			})
		})
		Context("日付を空にした時", func() {
			It("日付が変更されないこと", func() {
				before_date := *remarks.Date
				Expect(page.Navigate(base_path + remarks.ID.String())).To(Succeed())
				Expect(page.FindByName("date").Attribute("value")).To(Equal(before_date))
				Expect(page.FindByName("date").Fill("")).To(Succeed())
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				Expect(page.FindByName("date").Attribute("value")).To(Equal(before_date))
				Expect(page).To(HaveTitle("Lifelog | 備考編集"))
				Expect(page.FindByID("msg-error")).To(HaveText("未入力の項目があります"))
			})
		})
	})
	Describe("キャンセルボタン", func() {
		var (
			remarks models.Remarks
			lifelog models.LifeLog
			user    models.User
		)
		BeforeAll(func() {
			// user情報取得
			db.Where("name = ?", os.Getenv("AUTH0_EMAIL")).First(&user)
			// lifelog情報取得
			db.Where(&models.LifeLog{UserId: user.ID}).Where("name = ?", time.Now().Format("2006/01/02")).First(&lifelog)
			// remarks情報書き込み
			remarks = models.Remarks{
				LifeLogId: lifelog.ID,
				Title:     helpers.GetStringPointer("before"),
				Date:      lifelog.Name,
				Class:     "remarks",
			}
			db.Save(&remarks)
		})
		AfterAll(func() {
			// 使用したremarksを削除
			db.Delete(&remarks)
		})
		BeforeEach(func() {
			Expect(page.Navigate("http://localhost:3000/remarks/edit/" + remarks.ID.String())).To(Succeed())
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
			remarks models.Remarks
			lifelog models.LifeLog
			user    models.User
		)
		BeforeAll(func() {
			// user情報取得
			db.Where("name = ?", os.Getenv("AUTH0_EMAIL")).First(&user)
			// lifelog情報取得
			db.Where(&models.LifeLog{UserId: user.ID}).Where("name = ?", time.Now().Format("2006/01/02")).First(&lifelog)
			// remarks情報書き込み
			remarks = models.Remarks{
				LifeLogId: lifelog.ID,
				Title:     helpers.GetStringPointer("before"),
				Date:      lifelog.Name,
				Class:     "remarks",
			}
			db.Save(&remarks)
		})
		AfterAll(func() {
			// 使用したremarksを削除
			db.Delete(&remarks)
		})
		Context("削除ボタンを押した時", func() {
			It("行動記録が削除され、メイン画面に遷移すること", func() {
				Expect(page.Navigate("http://localhost:3000/lifelog")).To(Succeed())
				before_count, _ := page.AllByClass("remarks").Count()
				Expect(page.Navigate("http://localhost:3000/remarks/edit/" + remarks.ID.String())).To(Succeed())
				page.FindByID("js_delete").Click()
				Expect(page).To(HaveTitle("Lifelog | メイン"))
				after_count, _ := page.AllByClass("remarks").Count()
				Expect(after_count).To(Equal(before_count - 1))
			})
		})
	})
})
