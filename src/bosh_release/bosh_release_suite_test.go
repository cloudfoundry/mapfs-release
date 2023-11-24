package bosh_release_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var deployment_name string

func TestBoshReleaseTest(t *testing.T) {
	RegisterFailHandler(Fail)
	deployment_name = os.Getenv("DEPLOYMENT_NAME")
	if deployment_name == "" {
		deployment_name = "bosh_release_test_mapfs"
	}

	RunSpecs(t, "BoshReleaseTest Suite")
}

var _ = BeforeSuite(func() {
	var err error
	Expect(err).ShouldNot(HaveOccurred())
	SetDefaultEventuallyTimeout(10 * time.Minute)
	if !hasStemcell() {
		uploadStemcell()
	}
	deploy()
})

func deploy(opsfiles ...string) {
	stemcell_line := os.Getenv("STEMCELL")
	if stemcell_line == "" {
		stemcell_line = "jammy"
	}
	deployCmd := []string{"deploy",
		"-n",
		"-d",
		deployment_name,
		"./mapfs-manifest.yml",
		"-v", fmt.Sprintf("path_to_mapfs_release=%s", os.Getenv("MAPFS_RELEASE_PATH")),
		"-v", fmt.Sprintf("stemcell_lin=%s", stemcell_line),
	}

	updatedDeployCmd := make([]string, len(deployCmd))
	copy(updatedDeployCmd, deployCmd)
	for _, optFile := range opsfiles {
		updatedDeployCmd = append(updatedDeployCmd, "-o", optFile)
	}

	boshDeployCmd := exec.Command("bosh", updatedDeployCmd...)
	session, err := gexec.Start(boshDeployCmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 60*time.Minute).Should(gexec.Exit(0))
}

func undeploy() {
	deleteDeployCmd := []string{"deld",
		"-n",
		"-d",
		deployment_name,
	}

	boshDeployCmd := exec.Command("bosh", deleteDeployCmd...)
	session, err := gexec.Start(boshDeployCmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 60*time.Minute).Should(gexec.Exit(0))
}

func hasStemcell() bool {
	boshStemcellsCmd := exec.Command("bosh", "stemcells", "--json")
	stemcellOutput := gbytes.NewBuffer()
	session, err := gexec.Start(boshStemcellsCmd, stemcellOutput, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 1*time.Minute).Should(gexec.Exit(0))
	boshStemcellsOutput := &BoshStemcellsOutput{}
	Expect(json.Unmarshal(stemcellOutput.Contents(), boshStemcellsOutput)).Should(Succeed())
	return len(boshStemcellsOutput.Tables[0].Rows) > 0
}

func uploadStemcell() {
	boshUsCmd := exec.Command("bosh", "upload-stemcell", "https://bosh.io/d/stemcells/bosh-warden-boshlite-ubuntu-xenial-go_agent")
	session, err := gexec.Start(boshUsCmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 20*time.Minute).Should(gexec.Exit(0))
}

type BoshStemcellsOutput struct {
	Tables []struct {
		Content string `json:"Content"`
		Header  struct {
			Cid     string `json:"cid"`
			Cpi     string `json:"cpi"`
			Name    string `json:"name"`
			Os      string `json:"os"`
			Version string `json:"version"`
		} `json:"Header"`
		Rows []struct {
			Cid     string `json:"cid"`
			Cpi     string `json:"cpi"`
			Name    string `json:"name"`
			Os      string `json:"os"`
			Version string `json:"version"`
		} `json:"Rows"`
		Notes []string `json:"Notes"`
	} `json:"Tables"`
	Blocks interface{} `json:"Blocks"`
	Lines  []string    `json:"Lines"`
}
