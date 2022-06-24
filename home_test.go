package main_test

import (
	"lifelog/database"
	"lifelog/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("Home", Ordered, func() {
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
		page, err = agoutiDriver.NewPage(agouti.Browser("chrome"))
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(page.Destroy()).To(Succeed())
	})
	Context("/を表示した時", func() {
		BeforeEach(func() {
			Expect(page.Navigate("http://localhost:3000")).To(Succeed())
		})
		It("Home画面が表示されること", func() {
			Expect(page).To(HaveTitle("Lifelog | ホーム"))
		})
		It("サインインのリンクがあること", func() {
			Expect(page.FindByID("signin").Text()).To(Equal("サインイン"))
		})
		It("行動登録のリンクがないこと", func() {
			Expect(page.FindByID("action-resister")).ToNot(BeFound())
		})
		It("備考登録のリンクがないこと", func() {
			Expect(page.FindByID("remarks-resister")).ToNot(BeFound())
		})
		It("ログアウトのリンクがないこと", func() {
			Expect(page.FindByID("logout")).ToNot(BeFound())
		})
	})
})
