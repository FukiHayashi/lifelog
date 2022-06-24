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

var _ = Describe("LifelogNew", Ordered, func() {
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
	Describe("行動登録画面の表示", func() {
		BeforeEach(func() {
			Expect(page.Navigate("http://localhost:3000/lifelog/new")).To(Succeed())
		})
		Context("/lifelog/newを表示した時", func() {
			It("行動登録画面が表示されること", func() {
				Expect(page).To(HaveTitle("Lifelog | 行動登録"))
				Expect(page.FindByName("title")).To(BeFound())
				Expect(page.FindByName("start")).To(BeFound())
				Expect(page.FindByName("end")).To(BeFound())
				Expect(page.FindByName("class")).To(BeFound())
			})
		})
	})
	Describe("登録ボタン", func() {
		var (
			sleep_count  int
			action_count int
			bath_count   int
			meal_count   int
			other_count  int
		)
		BeforeEach(func() {
			sleep_count, _ = page.AllByClass("sleep").Count()
			action_count, _ = page.AllByClass("action").Count()
			bath_count, _ = page.AllByClass("bath").Count()
			meal_count, _ = page.AllByClass("meal").Count()
			other_count, _ = page.AllByClass("other").Count()
			Expect(page.Navigate("http://localhost:3000/lifelog/new")).To(Succeed())
		})
		Context("分類を睡眠で登録した時", func() {
			It("睡眠が登録されること", func() {
				page.FindByName("title").Fill("すいみん")
				page.FindByName("class").Select("睡眠")
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				Expect(page).To(HaveTitle("Lifelog | メイン"))
				Expect(page.AllByClass("sleep").Count()).To(Equal(sleep_count + 1))
			})
		})
		Context("分類を行動で登録した時", func() {
			It("行動が登録されること", func() {
				page.FindByName("title").Fill("こうどう")
				page.FindByName("class").Select("行動")
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				Expect(page).To(HaveTitle("Lifelog | メイン"))
				Expect(page.AllByClass("action").Count()).To(Equal(action_count + 1))
			})
		})
		Context("分類を風呂で登録した時", func() {
			It("風呂が登録されること", func() {
				page.FindByName("title").Fill("ふろ")
				page.FindByName("class").Select("風呂")
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				Expect(page).To(HaveTitle("Lifelog | メイン"))
				Expect(page.AllByClass("bath").Count()).To(Equal(bath_count + 1))
			})
		})
		Context("分類を食事で登録した時", func() {
			It("食事が登録されること", func() {
				page.FindByName("title").Fill("めし")
				page.FindByName("class").Select("食事")
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				Expect(page).To(HaveTitle("Lifelog | メイン"))
				Expect(page.AllByClass("meal").Count()).To(Equal(meal_count + 1))
			})
		})
		Context("分類をその他で登録した時", func() {
			It("その他が登録されること", func() {
				page.FindByName("title").Fill("そのた")
				page.FindByName("class").Select("その他")
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				Expect(page).To(HaveTitle("Lifelog | メイン"))
				Expect(page.AllByClass("other").Count()).To(Equal(other_count + 1))
			})
		})
		Context("タイトルを空欄にした時", func() {
			It("行動登録されないこと", func() {
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				Expect(page).To(HaveTitle("Lifelog | 行動登録"))
				Expect(page.Navigate("http://localhost:3000/lifelog")).To(Succeed())
				Expect(page.AllByClass("sleep").Count()).To(Equal(sleep_count))
				Expect(page.AllByClass("action").Count()).To(Equal(action_count))
				Expect(page.AllByClass("bath").Count()).To(Equal(bath_count))
				Expect(page.AllByClass("meal").Count()).To(Equal(meal_count))
				Expect(page.AllByClass("other").Count()).To(Equal(other_count))
			})
			It("エラーメッセージが表示されること", func() {
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				Expect(page).To(HaveTitle("Lifelog | 行動登録"))
				Expect(page.FindByID("msg-error")).To(HaveText("未入力の項目があります"))
			})
		})
		Context("開始時刻を空欄にした時", func() {
			It("行動登録されないこと", func() {
				page.FindByName("start").Fill("")
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				Expect(page).To(HaveTitle("Lifelog | 行動登録"))
				Expect(page.Navigate("http://localhost:3000/lifelog")).To(Succeed())
				Expect(page.AllByClass("sleep").Count()).To(Equal(sleep_count))
				Expect(page.AllByClass("action").Count()).To(Equal(action_count))
				Expect(page.AllByClass("bath").Count()).To(Equal(bath_count))
				Expect(page.AllByClass("meal").Count()).To(Equal(meal_count))
				Expect(page.AllByClass("other").Count()).To(Equal(other_count))
			})
		})
		Context("終了時刻を空欄にした時", func() {
			It("行動登録されないこと", func() {
				page.FindByName("end").Fill("")
				Expect(page.FindByID("resister-button").Click()).To(Succeed())
				Expect(page).To(HaveTitle("Lifelog | 行動登録"))
				Expect(page.Navigate("http://localhost:3000/lifelog")).To(Succeed())
				Expect(page.AllByClass("sleep").Count()).To(Equal(sleep_count))
				Expect(page.AllByClass("action").Count()).To(Equal(action_count))
				Expect(page.AllByClass("bath").Count()).To(Equal(bath_count))
				Expect(page.AllByClass("meal").Count()).To(Equal(meal_count))
				Expect(page.AllByClass("other").Count()).To(Equal(other_count))
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
