package parsers

import "strings"

// DecryptAllServers attempts to decrypt every provided server entry using
// the Megacloud decryptor (sub/dub/raw). Errors are ignored to keep parity
// with the original JS helper, but the slice omits failed entries.
func DecryptAllServers(servers []ServerData) []DecryptedSources {
	var results []DecryptedSources
	for _, server := range servers {
		if !isDecryptableServer(server.Type) {
			continue
		}
		res, err := DecryptMegacloud(server.ID, server.Name, server.Type)
		if err != nil {
			continue
		}
		results = append(results, res)
	}
	return results
}

func isDecryptableServer(t string) bool {
	switch strings.ToLower(strings.TrimSpace(t)) {
	case "sub", "dub", "raw":
		return true
	default:
		return false
	}
}
