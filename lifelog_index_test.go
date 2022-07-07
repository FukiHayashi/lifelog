package main_test

import (
	"lifelog/database"
	"lifelog/helpers"
	"lifelog/models"
	"net/http"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
	"gorm.io/gorm"
)

var _ = Describe("LifelogIndex", Ordered, func() {
	var page *agouti.Page
	var db *gorm.DB

	BeforeAll(func() {
		db = database.DataBaseConnect()
		// DBのマイグレーション
		db.AutoMigrate(&models.User{}, &models.LifeLog{}, &models.Appointment{}, &models.Remarks{})
	})
	AfterAll(func() {
		// テストに使用したDBの内容を全て削除する
		db.Migrator().DropTable(&models.User{}, &models.LifeLog{}, &models.Appointment{}, &models.Remarks{})
		database.DataBaseClose(db)
	})
	BeforeEach(func() {
		var err error
		page, err = agoutiDriver.NewPage()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(page.Destroy()).To(Succeed())
	})
	Describe("サインイン時", func() {
		BeforeEach(func() {
			// /からサインインする
			Expect(page.Navigate(os.Getenv("SERVER_PATH"))).To(Succeed())
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

		Context("/lifelogを表示した時", func() {
			It("lifelog画面が表示されること", func() {
				Expect(page).To(HaveTitle("Lifelog | メイン"))
			})
			It("行動登録のリンクがあること", func() {
				Expect(page.FindByID("action-resister")).To(BeFound())
			})
			It("備考登録のリンクがあること", func() {
				Expect(page.FindByID("remarks-resister")).To(BeFound())
			})
			It("ログアウトのリンクがあること", func() {
				Expect(page.FindByID("logout")).To(BeFound())
			})
			It("記録表のカレンダーが表示されること", func() {
				t := time.Now()
				tn := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, time.UTC)
				te := tn.AddDate(0, 0, -1)
				sjs := page.AllByClass("sjs-name")
				Expect(sjs.Count()).To(Equal(te.Day()))
			})
		})
		Context("行動記録を押した時", func() {
			It("行動記録登録画面が表示されること", func() {
				page.FindByID("action-resister").Click()
				Expect(page).To(HaveTitle("Lifelog | 行動登録"))
			})
		})
		Context("備考登録を押した時", func() {
			It("備考登録画面が表示されること", func() {
				page.FindByID("remarks-resister").Click()
				Expect(page).To(HaveTitle("Lifelog | 備考登録"))
			})
		})
		Context("ログアウトを押した時", func() {
			It("ログアウトし、ホーム画面に遷移すること", func() {
				page.FindByID("logout").Click()
				Expect(page).To(HaveTitle("Lifelog | ホーム"))
			})
		})
		Context("appointmentを押した時", func() {
			It("行動記録編集画面が表示されること", func() {
				user := models.User{}
				lifelog := models.LifeLog{}
				// user情報取得
				db.Where("name = ?", os.Getenv("AUTH0_EMAIL")).First(&user)
				// lifelog情報取得
				db.Where(&models.LifeLog{UserId: user.ID}).Where("name = ?", time.Now().Format("2006/01/02")).First(&lifelog)
				appointment := models.Appointment{
					LifeLogId: lifelog.ID,
					Title:     helpers.GetStringPointer("appointment"),
					Class:     "sleep",
				}
				db.Save(&appointment)
				Expect(page.Navigate(os.Getenv("SERVER_PATH") + "/lifelog")).To(Succeed())
				Expect(page.FindByClass("sleep").Click()).To(Succeed())
				page.FindByClass("sleep").Click()
				Expect(page.URL()).To(Equal(os.Getenv("SERVER_PATH") + "/lifelog/edit/" + appointment.ID.String()))
			})
		})
		Context("remarksを押した時", func() {
			It("備考編集画面が表示されること", func() {
				user := models.User{}
				lifelog := models.LifeLog{}
				// user情報取得
				db.Where("name = ?", os.Getenv("AUTH0_EMAIL")).First(&user)
				// lifelog情報取得
				db.Where(&models.LifeLog{UserId: user.ID}).Where("name = ?", time.Now().Format("2006/01/02")).First(&lifelog)
				remarks := models.Remarks{
					LifeLogId: lifelog.ID,
					Title:     helpers.GetStringPointer("remarks"),
					Date:      lifelog.Name,
					Class:     "remarks",
				}
				db.Save(&remarks)
				Expect(page.Navigate(os.Getenv("SERVER_PATH") + "/lifelog")).To(Succeed())
				Expect(page.FindByClass("remarks").Click()).To(Succeed())
				page.FindByClass("remarks").Click()
				Expect(page.URL()).To(Equal(os.Getenv("SERVER_PATH") + "/remarks/edit/" + remarks.ID.String()))
			})
		})
		Context("月を切り替えた時", func() {
			It("選択した月のlifelogが表示されること", func() {
				user := models.User{}
				// user情報取得
				db.Where("name = ?", os.Getenv("AUTH0_EMAIL")).First(&user)
				// 次の月のlifelogを作成
				lifelogs := []models.LifeLog{}
				t := time.Now().AddDate(0, 1, 0)
				name_date := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
				for name_date.Month() == t.Month() {
					lifelog_name := name_date.Format("2006/01/02")
					lifelogs = append(lifelogs, models.LifeLog{
						UserId:   user.ID,
						LoggedAt: name_date,
						Name:     &lifelog_name,
					})
					name_date = name_date.AddDate(0, 0, 1)
				}
				db.Create(&lifelogs)
				// ページ更新
				Expect(page.Navigate(os.Getenv("SERVER_PATH") + "/lifelog")).To(Succeed())
				Expect(page.FindByID("js_monthly_selector").Select(t.Format("2006-01"))).To(Succeed())
				Expect(page.URL()).To(Equal(os.Getenv("SERVER_PATH") + "/lifelog/" + t.Format("2006-01")))
			})
		})
	})
	Describe("未サインイン時", func() {
		Context("/lifelogを表示した時", func() {
			It("403が返ること", func() {
				resp, _ := http.Get(os.Getenv("SERVER_PATH") + "/lifelog")
				Expect(resp).To(HaveHTTPStatus(403))
			})
		})
	})
})
