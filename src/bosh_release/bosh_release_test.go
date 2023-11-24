package bosh_release_test

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("BoshReleaseTest", func() {
	BeforeEach(func() {
		deploy()
	})

	It("should have the mapfs binaries", func() {
		expectFileInstalled("/var/vcap/packages/mapfs/bin/mapfs")
		expectFileInstalled("/sbin/mount.fuse3")
		expectFileInstalled("/etc/fuse.conf")
	})

	Context("when mapfs is disabled", func() {

		BeforeEach(func() {
			cmd := exec.Command("bosh", "-d", "bosh_release_test", "delete-deployment", "-n")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))
			deploy("./operations/disable-mapfs.yml")
		})

		It("should not have or configured the fuse package", func() {
			expectFileInstalled("/var/vcap/packages/mapfs/bin/mapfs")
			expectFileNotInstalled("/etc/fuse.conf")
		})
	})
})

func expectFileInstalled(filePath string) {
	cmd := exec.Command("bosh", "-d", "bosh_release_test", "ssh", "-c", fmt.Sprintf("stat %s", filePath))
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session).Should(gexec.Exit(0), fmt.Sprintf("file [%s] was not found", filePath))
}

func expectFileNotInstalled(filePath string) {
	cmd := exec.Command("bosh", "-d", "bosh_release_test", "ssh", "-c", fmt.Sprintf("stat %s", filePath))
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session).Should(gexec.Exit(1), fmt.Sprintf("file [%s] was found when it should not have", filePath))
}
