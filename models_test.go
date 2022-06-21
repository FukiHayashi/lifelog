package main_test

import (
	"lifelog/models"
	"lifelog/test"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Models", func() {
	Describe("User", func() {
		th := test.TestHelper{}
		BeforeEach(func() {
			th.SetupTest()
		})
		AfterEach(func() {
			th.TearDownTest()
		})
		Context("Userが正しい時", func() {
			aud := time.Now().String()
			user := models.User{
				Aud:  &aud,
				Name: time.Now().String(),
			}
			It("Userが作成される", func() {
				err := th.DB.Create(&user).Error
				Expect(err).Should(BeNil())
			})
		})
		Context("Audがない時", func() {
			erruser := models.User{
				Name: time.Now().String(),
			}
			It("Userが作成されない", func() {
				err := th.DB.Create(&erruser).Error
				Expect(err).ShouldNot(BeNil())
			})
		})
	})
})
