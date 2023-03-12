package scan

import "Qscan/poc/Confluence"

func CheckAllPoc(url string, flag int) (bool, []map[string]string) {
	var infos []map[string]string

	b, info := Confluence.CVE_2022_26134(url)
	if b {
		infos = append(infos, info)
		return true, infos
	}

	return false, nil
}
