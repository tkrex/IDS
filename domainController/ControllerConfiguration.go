package domainController


type ForwardingTable map[string]string

func ForwardingTable() ForwardingTable {
	table := make(map[string]string)
	table["testDomain"] = "self"
	return table
}