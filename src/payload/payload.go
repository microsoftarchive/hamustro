package payload

// Checks the Collection has payloads or not
func (m *Collection) HasPayloads() bool {
	if m.GetPayloads() == nil {
		return false
	}
	if len(m.GetPayloads()) == 0 {
		return false
	}
	return true
}

// Checks the validity of a Protobuf object
// because not valid identitifed can be parsed via JSON
func (m *Collection) IsValid() bool {
	if m == nil {
		return false
	}
	if m.DeviceId == nil || m.ClientId == nil || m.Session == nil || m.SystemVersion == nil || m.ProductVersion == nil {
		return false
	}
	for _, p := range m.GetPayloads() {
		if p.At == nil || p.Event == nil || p.Nr == nil {
			return false
		}
	}
	return true
}
