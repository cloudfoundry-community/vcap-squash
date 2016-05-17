package main_test

import (
	. "github.com/joshq00/vcap-squash"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
	var vcap string
	var result []string
	var expected = func(vals ...interface{}) {
		Ω(result).Should(ConsistOf(vals...))
	}
	JustBeforeEach(func() {
		result = Process(vcap)
	})
	It("should be empty", func() {
		expected()
	})
	Context("when `VCAP_SERVICES` is set", func() {
		Context("with a `user-provided` service", func() {
			BeforeEach(func() {
				vcap = `{
					"user-provided": [{
						"name": "hi",
						"credentials": {
							"uri": "http://www.a.b.com",
							"username": "greeter"
						}
					}]
				}`
			})
			It("flattens the service", func() {
				expected(
					`export HI_URI="http://www.a.b.com"`,
					`export HI_USERNAME="greeter"`,
				)
			})

			Context("that has nested objects", func() {
				BeforeEach(func() {
					vcap = `{
						"user-provided": [{
							"name": "NeSt",
							"credentials": {
								"a": {
									"b": "c",
									"e": 1
								}
							}
						}]
					}`
				})
				It("flattens the service", func() {
					expected(
						`export NEST_A_B="c"`,
						`export NEST_A_E=1`,
					)
				})
			})
		})
		Context("with a marketplace service", func() {
			BeforeEach(func() {
				vcap = `{ "p-mysql": [{
					"name": "my-mysql-instance",
					"credentials": {
						"username": "x",
						"password": "y",
						"port": 1234
					}
				}] }`
			})
			It("flattens the service", func() {
				expected(
					`export MY_MYSQL_INSTANCE_USERNAME="x"`,
					`export MY_MYSQL_INSTANCE_PASSWORD="y"`,
					`export MY_MYSQL_INSTANCE_PORT=1234`,
				)
			})
		})
		Context("array fields", func() {
			BeforeEach(func() {
				vcap = `{ "user-provided": [{
					"name": "svc",
					"credentials": {
						"names": [ "abc", 123 ]
					}
				}] }`
			})
			It("flattens with index", func() {
				expected(
					`export SVC_NAMES_0="abc"`,
					`export SVC_NAMES_1=123`,
				)
			})
			Context("that contain non-primitives", func() {
				BeforeEach(func() {
					vcap = `{ "user-provided": [{
						"name": "svc",
						"credentials": {
							"apps": [ {
								"name": "app1"
							}, {
								"name": "app2",
								"nested": { "x": "yz" }
							}, "third" ]
						}
					}] }`
				})
				It("includes index and field", func() {
					expected(
						`export SVC_APPS_0_NAME="app1"`,
						`export SVC_APPS_1_NAME="app2"`,
						`export SVC_APPS_1_NESTED_X="yz"`,
						`export SVC_APPS_2="third"`,
					)
				})
			})
		})
	})

})
