package emulator

import (
	"context"

	"github.com/majd/ipatool/v2/internal/sap/assets"
	"github.com/majd/ipatool/v2/internal/sap/unicorn"
	"github.com/majd/ipatool/v2/internal/store/profile"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store error imports", func() {
	DescribeTable("preserves AMSError fields matching the pinned native framework", func(title, message, debug string, underlying bool) {
		ctx := context.Background()
		bundle, err := assets.Load(ctx)
		Expect(err).NotTo(HaveOccurred())
		prepared, err := profile.Load(ctx)
		Expect(err).NotTo(HaveOccurred())
		guest, err := open(ctx, bundle.AppleMediaServices, prepared, true)
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(guest.Close)
		values := []uint64{123, 0, 0, 0}
		if title != "" {
			values[1] = guest.New("NSString", title)
		}
		if message != "" {
			values[2] = guest.New("NSString", message)
		}
		if underlying {
			values[3] = guest.New("NSError", map[string]uint64{
				"domain":   guest.New("NSString", "TestDomain"),
				"code":     guest.New("NSNumber", int64(456)),
				"userInfo": guest.New("NSDictionary", map[string]uint64{}),
			})
		}
		for i, value := range values {
			Expect(guest.engine.RegWrite(unicorn.RegX0+i, value)).To(Succeed())
		}
		guest.symbol("AMSError")
		result := guest.arg(0)
		domain, err := guest.Invoke(result, "domain")
		Expect(err).NotTo(HaveOccurred())
		Expect(guest.String(domain)).To(Equal("AMSErrorDomain"))
		code, err := guest.Invoke(result, "code")
		Expect(err).NotTo(HaveOccurred())
		Expect(code).To(Equal(uint64(123)))
		info, err := guest.Invoke(result, "userInfo")
		Expect(err).NotTo(HaveOccurred())
		fields := guest.objects[info].Value.(map[string]uint64)
		Expect(guest.String(fields["NSDebugDescription"])).To(Equal(debug))
		if title == "" && message == "" {
			title = "An unknown error occurred. Please try again."
		}
		Expect(guest.String(fields["AMSDescription"])).To(Equal(title))
		Expect(guest.String(fields["AMSFailureReason"])).To(Equal(message))
		Expect(fields["NSUnderlyingError"]).To(Equal(values[3]))
	}, Entry("title and message with underlying error", "Title", "Message", "Title Message", true),
		Entry("title only", "Title", "", "Title", false),
		Entry("message only", "", "Message", "Message", false),
		Entry("no description", "", "", "An unknown error occurred. Please try again.", false))
})
