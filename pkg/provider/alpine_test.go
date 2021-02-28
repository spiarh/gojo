package provider

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Alpine Provider", func() {
	Describe("Parse APK Index", func() {
		var (
			apkIndex []byte
			pkgName  string
		)

		Context("with well formed APK Index", func() {
			BeforeEach(func() {
				apkIndex = []byte(`P:mariadb
V:10.4.17-r1
A:x86_64

P:nginx
V:1.18.0-r1
A:x86_64

`)
			})
			It("finds package meta", func() {
				pkgName = "nginx"
				expectedPkgMeta := &AlpinePackageMeta{name: "nginx", version: "1.18.0-r1", arch: "x86_64"}
				actual, err := parseAPKIndex(apkIndex, pkgName)
				Expect(err).To(BeNil())
				Expect(actual).To(Equal(expectedPkgMeta))
			})
			It("fails to find package meta", func() {
				pkgName = "missing"
				actual, err := parseAPKIndex(apkIndex, pkgName)
				Expect(err).To(HaveOccurred())
				Expect(actual).To(BeNil())
			})
		})
		Context("with malformed APK Index", func() {
			BeforeEach(func() {
				pkgName = "nginx"
			})
			It("fails because of incomplete data", func() {
				apkIndex = []byte(`P:mariadb
V:10.4.17-r1
A:x86_64

P:nginx
V:1.18.0-r1

`)
				actual, err := parseAPKIndex(apkIndex, pkgName)
				Expect(err).To(HaveOccurred())
				Expect(actual).To(BeNil())
			})
			It("fails because of bad formating", func() {
				apkIndex = []byte(`P:mariadb
V:10.4.17-r1
\\\
A:x86_64

P:nginx
				V:1.18.0-r1
A:x86_64

`)
				pkgName = "nginx"
				actual, err := parseAPKIndex(apkIndex, pkgName)
				Expect(err).To(HaveOccurred())
				Expect(actual).To(BeNil())
			})
		})
	})
})
