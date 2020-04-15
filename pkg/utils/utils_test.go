package utils_test

import (
	"strings"

	"github.com/alecthomas/units"
	"github.com/kubevirt/vm-import-operator/pkg/utils"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	k8svalidation "k8s.io/apimachinery/pkg/util/validation"
)

var _ = Describe("Validating Name Normalization", func() {
	table.DescribeTable("should fail for an invalid format", func(tested string) {
		result, err := utils.NormalizeName(tested)
		Expect(result).To(Equal(""))
		Expect(err).To(HaveOccurred())
	},
		table.Entry("Empty string", ""),
		table.Entry("Non-alphanumeric characters", "$!@#$!@#$%"),
		table.Entry("Only dashes without alphanumeric characters", "-----"),
	)
	table.DescribeTable("should normalize given name to expected format", func(tested string, expected string) {
		result, err := utils.NormalizeName(tested)
		Expect(result).To(Equal(expected))
		Expect(err).NotTo(HaveOccurred())
	},
		table.Entry("URL format", "https://my.host.com", "httpsmy-host-com"),
		table.Entry("Leading spaces and dots", " my.host.com", "my-host-com"),
		table.Entry("Uppercase", "MY.HOST.COM", "my-host-com"),
		table.Entry("Leading dash and non alphanumeric last character", "-my-host;", "my-host"),
		table.Entry("Leading dash, mix of letter and digit", "-my-72host;", "my-72host"),
		table.Entry("Alphanumeric characters mixed with illegal symbols", " @#$_#*($%-my-[];.1@##@%2#-host;   ", "my--12-host"),
		table.Entry("A legal name", "my-host", "my-host"),
		table.Entry("A legal name", "12-my-host-123", "12-my-host-123"),
		table.Entry("A legal name", "my-12host", "my-12host"),
		table.Entry("A legal name", "m", "m"),
		table.Entry("A legal name", "0", "0"),
	)
	table.DescribeTable("should normalize long name to 253 length", func(tested string) {
		result, err := utils.NormalizeName(tested)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(result)).To(Equal(k8svalidation.DNS1123SubdomainMaxLength))
	},
		table.Entry("Max length", createStringOfLength(k8svalidation.DNS1123SubdomainMaxLength)),
		table.Entry("Max length exceeded by 1", createStringOfLength(k8svalidation.DNS1123SubdomainMaxLength+1)),
		table.Entry("Double the max length", createStringOfLength(k8svalidation.DNS1123SubdomainMaxLength*2)),
	)
})

var _ = Describe("Converting of bytes", func() {
	table.DescribeTable("should convert to proper suffix", func(bytes int64, expected string) {
		result, _ := utils.FormatBytes(bytes)
		Expect(result).To(Equal(expected))
	},
		table.Entry("To Ki", int64(units.KiB), "1Ki"),
		table.Entry("To Ki", int64(12*units.KiB), "12Ki"),
		table.Entry("To Mi", int64(units.MiB), "1Mi"),
		table.Entry("To Mi", int64(512*units.MiB), "512Mi"),
		table.Entry("To Gi", int64(units.GiB), "1Gi"),
		table.Entry("To Gi", int64(4*units.GiB), "4Gi"),
		table.Entry("To Ti", int64(units.TiB), "1Ti"),
		table.Entry("To Pi", int64(units.PiB), "1Pi"),
		table.Entry("To Ei", int64(units.EiB), "1Ei"),
		table.Entry("To B", int64(1), "1"),
		table.Entry("No conversion", int64(units.GiB-1), "1073741823"),
	)
})

func createStringOfLength(n int) string {
	return strings.Repeat("x", n)
}
