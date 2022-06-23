package main_test

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("LifelogIndex", Ordered, func() {
	var page *agouti.Page

	BeforeAll(func() {
		if err := godotenv.Load(".testenv"); err != nil {
			log.Fatalf(".envファイルの読み込みに失敗しました: %v", err)
		}
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
})
