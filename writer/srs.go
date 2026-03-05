package writer

import (
	"os"
	"path/filepath"

	"github.com/sagernet/sing-box/common/srs"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

const srsVersion = 2

var _ Writer = (*SRSWriter)(nil)

type SRSWriter struct{}

func (w *SRSWriter) Write(baseDir string, entries []Entry) error {
	dir, err := ensureDir(baseDir, "srs")
	if err != nil {
		return err
	}

	for _, e := range entries {
		if err := writeSRS(filepath.Join(dir, e.Name+".srs"), e.Prefixes); err != nil {
			return err
		}
	}
	return nil
}

func writeSRS(path string, prefixes []string) error {
	var headlessRule option.DefaultHeadlessRule
	headlessRule.IPCIDR = prefixes

	plainRuleSet := option.PlainRuleSet{
		Rules: []option.HeadlessRule{
			{
				Type:           C.RuleTypeDefault,
				DefaultOptions: headlessRule,
			},
		},
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return srs.Write(f, plainRuleSet, srsVersion)
}
