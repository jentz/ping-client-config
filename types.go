package pcc

// Utility functions to make it easier to deal with optional values
func boolPtr(b bool) *bool       { return &b }
func intPtr(i int) *int          { return &i }
func stringPtr(s string) *string { return &s }
