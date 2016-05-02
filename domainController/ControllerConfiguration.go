package domainController




func ForwardingTable() map[string]string {
	table := make(map[string]string)
	table["testDomain"] = "self"
	return table
}