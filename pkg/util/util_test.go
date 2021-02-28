package util_test

import (
	"testing"

	"github.com/lcavajani/gojo/pkg/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "util Test Suite")
}

var _ = Describe("image name parsing", func() {
	It("should parse valid image names", func() {
		tests := []struct {
			fullName, registry, name, tag string
		}{

			{fullName: "registry.fqdn/name:tag", registry: "registry.fqdn", name: "name", tag: "tag"},
			{fullName: "registry.fqdn/project/name:tag", registry: "registry.fqdn/project", name: "name", tag: "tag"},
			{fullName: "registry.fqdn:5000/project/env/name:tag", registry: "registry.fqdn:5000/project/env", name: "name", tag: "tag"},
		}

		for _, tt := range tests {
			reg, name, tag, err := util.ParseImageFullName(tt.fullName)
			Expect(err).To(BeNil())
			Expect(reg).To(Equal(tt.registry))
			Expect(name).To(Equal(tt.name))
			Expect(tag).To(Equal(tt.tag))
		}
	})
	It("should fails with invalid image name", func() {
		tests := []string{
			"name",
			"name:tag",
			"registry.fqdn/project",
			"reg$$%istry.fqdn/project",
		}
		for _, tt := range tests {
			_, _, _, err := util.ParseImageFullName(tt)
			Expect(err).To(HaveOccurred())
		}
	})
})

var _ = Describe("versions", func() {
	It("should sanitize version", func() {
		tests := []struct {
			version, expected string
		}{
			{version: "v1.2.3", expected: "1.2.3"},
			{version: "1.2.3", expected: "1.2.3"},
		}
		for _, tt := range tests {
			version := util.SanitizeVersion(tt.version)
			Expect(version).To(Equal(tt.expected))
		}

	})
})
