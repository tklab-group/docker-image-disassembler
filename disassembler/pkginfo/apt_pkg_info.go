package pkginfo

import (
	"io"
	"regexp"
	"strings"
)

const AptPkgFilePath string = "/var/lib/dpkg/status"

// AptPkgInfo contains information for the package managed by apt.
// It is minimum information now.
type AptPkgInfo struct {
	Package string
	Version string
}

// ReadAptPkgInfos reads apt packages information and creates AptPkgInfo for each package.
func ReadAptPkgInfos(aptPkgFile io.Reader) ([]*AptPkgInfo, error) {
	b, err := io.ReadAll(aptPkgFile)
	if err != nil {
		return nil, err
	}

	// Each package information is divided by a blank line.
	rawPkgInfoList := strings.Split(string(b), "\n\n")
	reg := regexp.MustCompile(`Package:\s(?P<pkgName>.+)(.|\s)*Version:\s(?P<version>.+)(.|\s)*`)
	list := make([]*AptPkgInfo, 0)
	for _, rawPkgInfo := range rawPkgInfoList {
		if rawPkgInfo == "" {
			continue
		}
		extracted := extractNamedGroup(reg, rawPkgInfo)
		aptPkgInfo := &AptPkgInfo{
			Package: extracted["pkgName"],
			Version: extracted["version"],
		}
		list = append(list, aptPkgInfo)
	}

	return list, nil
}
