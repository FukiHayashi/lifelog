package main_test

import (
	"testing"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
)

func TestLifelog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lifelog Suite")
}

var agoutiDriver *agouti.WebDriver

var _ = BeforeSuite(func() {
	// Choose a WebDriver:
	// 設定ファイル読み込み
	godotenv.Load(".testenv")
	//agoutiDriver = agouti.PhantomJS()
	//agoutiDriver = agouti.Selenium()
	agoutiDriver = agouti.ChromeDriver()

	Expect(agoutiDriver.Start()).To(Succeed())
})

var _ = AfterSuite(func() {
	Expect(agoutiDriver.Stop()).To(Succeed())
})
