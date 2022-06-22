package main_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("Home", func() {
	var page *agouti.Page

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
			Expect(page).To(HaveTitle("Lifelog Home"))
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
