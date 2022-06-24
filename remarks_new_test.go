package main_test

import (
	"lifelog/database"
	"lifelog/models"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("RemarksNew", Ordered, func() {
	var page *agouti.Page
	BeforeAll(func() {
		// テスト環境のDBに接続
		db := database.DataBaseConnect()
		// DBのマイグレーション
		db.AutoMigrate(&models.User{}, &models.LifeLog{}, &models.Appointment{})
	})
	AfterAll(func() {
		// テストに使用したDBの内容を全て削除する
		db := database.DataBaseConnect()
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
	Describe("備考登録画面の表示", func() {
		BeforeEach(func() {
			Expect(page.Navigate("http://localhost:3000/remarks/new")).To(Succeed())
		})
		Context("/remarks/newを表示した時", func() {
			It("備考登録画面が表示されること", func() {
				Expect(page).To(HaveTitle("Lifelog | 備考登録"))
				Expect(page.FindByName("title")).To(BeFound())
				Expect(page.FindByName("date")).To(BeFound())
			})
		})
	})
	Describe("登録ボタン", func() {
		var (
			remarks_count int
		)
		BeforeEach(func() {
			remarks_count, _ = page.AllByClass("remarks").Count()
			Expect(page.Navigate("http://localhost:3000/remarks/new")).To(Succeed())
		})
		Context("タイトルと日付ありで登録した時", func() {
			It("備考が登録されること", func() {
				page.FindByName("title").Fill("腹痛")
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				Expect(page).To(HaveTitle("Lifelog | メイン"))
				Expect(page.AllByClass("remarks").Count()).To(Equal(remarks_count + 1))
			})
		})
		Context("タイトルを空欄にした時", func() {
			It("備考登録されないこと", func() {
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				Expect(page).To(HaveTitle("Lifelog | 備考登録"))
				Expect(page.Navigate("http://localhost:3000/lifelog")).To(Succeed())
				Expect(page.AllByClass("remarks").Count()).To(Equal(remarks_count))
			})
			It("エラーメッセージが表示されること", func() {
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				Expect(page).To(HaveTitle("Lifelog | 備考登録"))
				Expect(page.FindByID("msg-error")).To(HaveText("未入力の項目があります"))
			})
		})
		Context("日付を空欄にした時", func() {
			It("エラーメッセージが表示されること", func() {
				page.FindByName("date").Fill("")
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				Expect(page).To(HaveTitle("Lifelog | 備考登録"))
				Expect(page.FindByID("msg-error")).To(HaveText("未入力の項目があります"))
			})
		})
	})
	Describe("キャンセルボタン", func() {
		BeforeEach(func() {
			Expect(page.Navigate("http://localhost:3000/lifelog/new")).To(Succeed())
		})
		Context("キャンセルボタンを押した時", func() {
			It("メイン画面に遷移すること", func() {
				page.FindByID("cancel-button").Click()
				Expect(page).To(HaveTitle("Lifelog | メイン"))
			})
		})
	})
})
