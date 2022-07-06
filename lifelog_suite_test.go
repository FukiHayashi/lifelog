package main_test

import (
	"log"
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
	if err := godotenv.Load(".testenv"); err != nil {
		log.Fatalf("テスト用.envファイルの読み込みに失敗しました: %v", err)
	}
	//agoutiDriver = agouti.PhantomJS()
	//agoutiDriver = agouti.Selenium()

	agoutiDriver = agouti.ChromeDriver()

	Expect(agoutiDriver.Start()).To(Succeed())
})

var _ = AfterSuite(func() {
	Expect(agoutiDriver.Stop()).To(Succeed())
})
